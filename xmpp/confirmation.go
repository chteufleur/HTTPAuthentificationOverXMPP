package xmpp

import (
	"git.kingpenguin.tk/chteufleur/go-xmpp.git/src/xmpp"

	"log"
	"strconv"
)

const (
	REPLY_UNREACHABLE = "reply_unreachable"
	REPLY_DENY        = "reply_deny"
	REPLY_OK          = "reply_ok"

	TYPE_SEND_MESSAGE = "type_send_message"
	TYPE_SEND_IQ      = "type_send_iq"
)

type Confirmation struct {
	JID         string
	Method      string
	Domain      string
	Transaction string

	TypeSend string
	IdMap    string

	ChanReply chan string
}

func (confirmation *Confirmation) SendConfirmation() {
	log.Printf("%sQuery JID %s", LogInfo, confirmation.JID)
	clientJID, _ := xmpp.ParseJID(confirmation.JID)
	if clientJID.Resource == "" {
		confirmation.askViaMessage()
	} else {
		confirmation.askViaIQ()
	}
}

func (confirmation *Confirmation) askViaIQ() {
	stanzaID++
	stanzaIDstr := strconv.Itoa(stanzaID)
	m := xmpp.Iq{Type: xmpp.IQTypeGet, To: confirmation.JID, From: jid.Full(), Id: stanzaIDstr}
	confirm := &xmpp.Confirm{Id: confirmation.Transaction, Method: confirmation.Method, URL: confirmation.Domain}
	m.PayloadEncode(confirm)
	WaitIqMessages[stanzaIDstr] = confirmation
	comp.Out <- m

	confirmation.TypeSend = TYPE_SEND_IQ
	confirmation.IdMap = stanzaIDstr
}

func (confirmation *Confirmation) askViaMessage() {
	m := xmpp.Message{From: jid.Full(), To: confirmation.JID, Type: xmpp.MessageTypeNormal}

	m.Thread = xmpp.SessionID()
	m.Body = confirmation.Domain + " (with method " + confirmation.Method + ") need to validate your identity, do you agree ?"
	m.Body += "\nValidation code : " + confirmation.Transaction
	m.Body += "\nPlease check that this code is the same as on " + confirmation.Domain
	m.Body += "\n\nIf your client doesn't support that functionnality, please send back the validation code to confirm the request."
	m.Confir = &xmpp.Confirm{Id: confirmation.Transaction, Method: confirmation.Method, URL: confirmation.Domain}

	log.Printf("%sSend message %v", LogInfo, m)
	WaitMessageAnswers[confirmation.Transaction] = confirmation
	comp.Out <- m

	confirmation.TypeSend = TYPE_SEND_MESSAGE
	confirmation.IdMap = confirmation.Transaction
}
