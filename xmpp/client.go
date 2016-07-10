package xmpp

import (
	"git.kingpenguin.tk/chteufleur/go-xmpp.git/src/xmpp"

	"log"
	"strconv"
)

type Client struct {
	JID         string
	Method      string
	Domain      string
	Transaction string

	ChanReply chan bool
}

func (client *Client) QueryClient() {
	log.Printf("%sQuery JID %s", LogInfo, client.JID)
	clientJID, _ := xmpp.ParseJID(client.JID)
	if clientJID.Resource == "" {
		client.askViaMessage()
	} else {
		client.askViaIQ()
	}
}

func (client *Client) askViaIQ() {
	stanzaID++
	stanzaIDstr := strconv.Itoa(stanzaID)
	m := xmpp.Iq{Type: xmpp.IQTypeGet, To: client.JID, From: jid.Domain, Id: stanzaIDstr}
	confirm := &xmpp.Confirm{Id: client.Transaction, Method: client.Method, URL: client.Domain}
	m.PayloadEncode(confirm)
	WaitIqMessages[stanzaIDstr] = client
	comp.Out <- m
}

func (client *Client) askViaMessage() {
	m := xmpp.Message{From: jid.Domain, To: client.JID, Type: "normal"}

	m.Thread = xmpp.SessionID()
	m.Body = "Auth request for " + client.Domain + ".\nTransaction identifier is: " + client.Transaction + "\nReply to this message to confirm the request."
	m.Confir = &xmpp.Confirm{Id: client.Transaction, Method: client.Method, URL: client.Domain}

	log.Printf("%sSenp message %v", LogInfo, m)
	WaitMessageAnswers[client.Transaction] = client
	comp.Out <- m
}
