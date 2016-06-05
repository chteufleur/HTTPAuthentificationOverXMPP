package xmpp

import (
	"git.kingpenguin.tk/chteufleur/go-xmpp.git/src/xmpp"

	"log"
)

const (
	LogInfo  = "\t[XMPP COMPONENT INFO]\t"
	LogError = "\t[XMPP COMPONENT ERROR]\t"
	LogDebug = "\t[XMPP COMPONENT DEBUG]\t"
)

var (
	Addr   = "127.0.0.1:5347"
	JidStr = ""
	Secret = ""

	SoftVersion = ""

	stanzaID = 0

	jid    xmpp.JID
	stream = new(xmpp.Stream)
	comp   = new(xmpp.XMPP)

	ChanAction = make(chan string)

  WaitMessageAnswers = make(map[string]*Client)

	Debug = true
)



func Run() {
	log.Printf("%sRunning", LogInfo)
	// Create stream and configure it as a component connection.
	jid = must(xmpp.ParseJID(JidStr)).(xmpp.JID)
	stream = must(xmpp.NewStream(Addr, &xmpp.StreamConfig{LogStanzas: Debug})).(*xmpp.Stream)
	comp = must(xmpp.NewComponentXMPP(stream, jid, Secret)).(*xmpp.XMPP)

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
				client := WaitMessageAnswers[confirm.Id]
				delete(WaitMessageAnswers, confirm.Id)
				processConfirm(v, client)
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
				client := WaitMessageAnswers[v.Id]
				delete(WaitMessageAnswers, v.Id)
				processConfirm(v, client)

			default:
				reply := v.Response(xmpp.IQTypeError)
				reply.PayloadEncode(xmpp.NewError("cancel", xmpp.FeatureNotImplemented, ""))
				comp.Out <- reply
			}

		default:
			log.Printf("%srecv: %v", LogDebug, x)
		}
	}
}

func processConfirm(x interface{}, client *Client) {
	mes, mesOK := x.(*xmpp.Message)
	iq, iqOK := x.(*xmpp.Iq)

	if client != nil {
		if mesOK && mes.Error != nil {
			client.ChanReply <- false
		} else if iqOK && iq.Error != nil {
			client.ChanReply <- false
		} else {
			client.ChanReply <- true
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
