package xmpp

import (
	"git.kingpenguin.tk/chteufleur/go-xmpp.git/src/xmpp"

	"log"
	"strings"
)

const (
	LogInfo  = "\t[XMPP INFO]\t"
	LogError = "\t[XMPP ERROR]\t"
	LogDebug = "\t[XMPP DEBUG]\t"

	DEFAULT_SERVER_ADDRESS = "127.0.0.1"
	DEFAULT_SERVER_PORT    = "5347"
)

var (
	Addr   = ""
	Port   = ""
	JidStr = ""
	Secret = ""

	SoftVersion = ""

	stanzaID = 0

	jid    xmpp.JID
	stream = new(xmpp.Stream)
	comp   = new(xmpp.XMPP)

	ChanAction = make(chan string)

	WaitMessageAnswers = make(map[string]*Confirmation)
	WaitIqMessages     = make(map[string]*Confirmation)

	Debug              = true
	VerifyCertValidity = true
)

func Run() {
	var addr string
	var isComponent bool

	log.Printf("%sRunning", LogInfo)
	// Create stream and configure it as a component connection.
	jid = must(xmpp.ParseJID(JidStr)).(xmpp.JID)
	isComponent = jid.Node == ""

	if isComponent {
		// component
		if Addr == "" {
			Addr = DEFAULT_SERVER_ADDRESS
		}
		if Port == "" {
			Port = DEFAULT_SERVER_PORT
		}
		addr = Addr + ":" + Port
	} else {
		// client
		if Addr == "" {
			addrs := must(xmpp.HomeServerAddrs(jid)).([]string)
			addr = addrs[0]
		} else {
			if Port == "" {
				Port = DEFAULT_SERVER_PORT
			}
			addr = Addr + ":" + Port
		}
	}

	log.Printf("%sConnecting to %s", LogInfo, addr)
	stream = must(xmpp.NewStream(addr, &xmpp.StreamConfig{LogStanzas: Debug})).(*xmpp.Stream)

	if isComponent {
		comp = must(xmpp.NewComponentXMPP(stream, jid, Secret)).(*xmpp.XMPP)
	} else {
		comp = must(xmpp.NewClientXMPP(stream, jid, Secret, &xmpp.ClientConfig{InsecureSkipVerify: !VerifyCertValidity})).(*xmpp.XMPP)
		comp.Out <- xmpp.Presence{}
	}

	mainXMPP()
	log.Printf("%sReach main method's end", LogInfo)
	go Run()
}

func mainXMPP() {
	for x := range comp.In {
		switch v := x.(type) {
		case *xmpp.Presence:

		case *xmpp.Message:
			confirm := v.Confir
			if confirm != nil {
				confirmation := WaitMessageAnswers[confirm.Id]
				processConfirm(v, confirmation)
			} else {
				// If body is the confirmation id, it will be considerated as accepted.
				// In order to be compatible with all confirmations.
				confirmation := WaitMessageAnswers[v.Body]
				jidFrom, _ := xmpp.ParseJID(v.From)
				if confirmation != nil && confirmation.JID == jidFrom.Bare() {
					processConfirm(v, confirmation)
				}
			}

		case *xmpp.Iq:
			switch v.PayloadName().Space {
			case xmpp.NSDiscoItems:
				execDiscoCommand(v)

			case xmpp.NSVCardTemp:
				reply := v.Response(xmpp.IQTypeResult)
				vcard := &xmpp.VCard{}
				reply.PayloadEncode(vcard)
				comp.Out <- reply

			case xmpp.NSJabberClient:
				reply := v.Response(xmpp.IQTypeResult)
				reply.PayloadEncode(&xmpp.SoftwareVersion{Name: "HTTP authentification component", Version: SoftVersion})
				comp.Out <- reply

			case xmpp.NSHTTPAuth:
				confirm := &xmpp.Confirm{}
				v.PayloadDecode(confirm)
				confirmation := WaitIqMessages[v.Id]
				processConfirm(v, confirmation)

			default:
				// Handle reply iq that doesn't contain HTTP-Auth namespace
				confirmation := WaitIqMessages[v.Id]
				processConfirm(v, confirmation)

				if confirmation == nil {
					reply := v.Response(xmpp.IQTypeError)
					reply.PayloadEncode(xmpp.NewError("cancel", xmpp.FeatureNotImplemented, ""))
					comp.Out <- reply
				}
			}

		default:
			log.Printf("%srecv: %v", LogDebug, x)
		}
	}
}

func processConfirm(x interface{}, confirmation *Confirmation) {
	mes, mesOK := x.(*xmpp.Message)
	iq, iqOK := x.(*xmpp.Iq)

	if confirmation != nil {
		if mesOK && mes.Error != nil {
			// Message error
			errCondition := mes.Error.Condition()
			if errCondition == xmpp.ServiceUnavailable {
				// unreachable
				confirmation.ChanReply <- REPLY_UNREACHABLE
			} else {
				confirmation.ChanReply <- REPLY_DENY
			}

		} else if iqOK && iq.Error != nil {
			// IQ error
			errCondition := iq.Error.Condition()
			if errCondition == xmpp.ServiceUnavailable || errCondition == xmpp.FeatureNotImplemented {
				// send by message if client doesn't implemente it
				confirmation.JID = strings.SplitN(confirmation.JID, "/", 2)[0]
				go confirmation.SendConfirmation()
			} else if errCondition == xmpp.RemoteServerNotFound {
				// unreachable
				confirmation.ChanReply <- REPLY_UNREACHABLE
			} else {
				confirmation.ChanReply <- REPLY_DENY
			}

		} else {
			// No error
			confirmation.ChanReply <- REPLY_OK
		}
	}
}

func must(v interface{}, err error) interface{} {
	if err != nil {
		log.Fatal(LogError, err)
	}
	return v
}

func execDiscoCommand(iq *xmpp.Iq) {
	log.Printf("%sDiscovery item iq received", LogInfo)
	reply := iq.Response(xmpp.IQTypeResult)
	discoItem := &xmpp.DiscoItems{Node: xmpp.NodeAdHocCommand}
	reply.PayloadEncode(discoItem)
	comp.Out <- reply
}
