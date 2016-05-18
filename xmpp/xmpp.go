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

	jid    xmpp.JID
	stream = new(xmpp.Stream)
	comp   = new(xmpp.XMPP)

	ChanAction = make(chan string)

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

func must(v interface{}, err error) interface{} {
	if err != nil {
		log.Fatal(LogError, err)
	}
	return v
}


func SendMessage(to, subject, message string) {
	m := xmpp.Message{From: jid.Domain, To: to, Body: message, Type: "chat"}

	if subject != "" {
		m.Subject = subject
	}

	log.Printf("%sSenp message %v", LogInfo, m)
	comp.Out <- m
}


func execDiscoCommand(iq *xmpp.Iq) {
	log.Printf("%sDiscovery item iq received", LogInfo)
	reply := iq.Response(xmpp.IQTypeResult)
	discoItem := &xmpp.DiscoItems{Node: xmpp.NodeAdHocCommand}
	reply.PayloadEncode(discoItem)
	comp.Out <- reply
}
