package xmpp

import (
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
  client.ChanReply <- false
}
