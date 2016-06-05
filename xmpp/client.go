package xmpp

import (
  "git.kingpenguin.tk/chteufleur/go-xmpp.git/src/xmpp"

  "log"
)

type Client struct {
  JID         string
  Method      string
  Domain      string
  Transaction string

  ChanReply   chan bool
}


func (client *Client) QueryClient() {
  log.Printf("%sQuery JID %s", LogInfo, client.JID)
  client.askViaMessage()
}

func (client *Client) askViaIQ() {

}

func (client *Client) askViaMessage() {
  m := xmpp.Message{From: jid.Domain, To: client.JID, Type: "normal"}

  m.Thread = xmpp.SessionID()
  m.Body = "Auth request for "+client.Domain+".\nTransaction identifier is: "+client.Transaction+"\nReply to this message to confirm the request."
  m.Confir = &xmpp.Confirm{ID: client.Transaction, Method: client.Method, URL: client.Domain}

  log.Printf("%sSenp message %v", LogInfo, m)
  comp.Out <- m

  WaitMessageAnswers[client.Transaction] = client
}
