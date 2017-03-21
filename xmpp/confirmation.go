package xmpp

import (
	"git.kingpenguin.tk/chteufleur/go-xmpp.git/src/xmpp"

	"log"
	"strconv"
	"strings"
)

const (
	REPLY_UNREACHABLE = "reply_unreachable"
	REPLY_DENY        = "reply_deny"
	REPLY_OK          = "reply_ok"

	TYPE_SEND_MESSAGE = "type_send_message"
	TYPE_SEND_IQ      = "type_send_iq"

	TEMPLATE_DOMAIN          = "_DOMAIN_"
	TEMPLATE_METHOD          = "_METHOD_"
	TEMPLATE_VALIDATION_CODE = "_VALIDE_CODE_"
	DEFAULT_MESSAGE          = "_DOMAIN_ (with method _METHOD_) need to validate your identity, do you agree ?\nValidation code : _VALIDE_CODE_\nPlease check that this code is the same as on _DOMAIN_.\n\nIf your client doesn't support that functionnality, please send back the validation code to confirm the request."
)

var (
	MapLangs = make(map[string]string)
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
	confirmation.setBodies(&m)
	m.Confirm = &xmpp.Confirm{Id: confirmation.Transaction, Method: confirmation.Method, URL: confirmation.Domain}

	log.Printf("%sSend message %v", LogInfo, m)
	WaitMessageAnswers[confirmation.Transaction] = confirmation
	comp.Out <- m

	confirmation.TypeSend = TYPE_SEND_MESSAGE
	confirmation.IdMap = confirmation.Transaction
}

func (confirmation *Confirmation) setBodies(message *xmpp.Message) {
	msg := DEFAULT_MESSAGE
	if len(MapLangs) == 0 {
		msg = strings.Replace(msg, TEMPLATE_DOMAIN, confirmation.Domain, -1)
		msg = strings.Replace(msg, TEMPLATE_METHOD, confirmation.Method, -1)
		msg = strings.Replace(msg, TEMPLATE_VALIDATION_CODE, confirmation.Transaction, -1)
		message.Body = append(message.Body, xmpp.MessageBody{Lang: "en", Value: msg})
	} else {
		for key, val := range MapLangs {
			msg = val
			msg = strings.Replace(msg, TEMPLATE_DOMAIN, confirmation.Domain, -1)
			msg = strings.Replace(msg, TEMPLATE_METHOD, confirmation.Method, -1)
			msg = strings.Replace(msg, TEMPLATE_VALIDATION_CODE, confirmation.Transaction, -1)
			msg = strings.Replace(msg, "\\n", "\n", -1)
			msg = strings.Replace(msg, "\\r", "\r", -1)
			msg = strings.Replace(msg, "\\t", "\t", -1)
			message.Body = append(message.Body, xmpp.MessageBody{Lang: key, Value: msg})
		}
	}
}
