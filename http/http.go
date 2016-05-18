package http

import (
  "fmt"
  "net/http"
  "strconv"
  "strings"
  "log"
)

const (
  PARAM_JID      = "jid"
  METHOD_ACCESS  = "method"
  DOMAIN_ACCESS  = "domain"
  TRANSACTION_ID = "transaction_id"

  RETURN_VALUE_OK  = "OK"
  RETURN_VALUE_NOK = "NOK"
)

var (
  HttpPortBind = 9090

  ChanRequest = make(chan interface{}, 5)
)


func indexHandler(w http.ResponseWriter, r *http.Request) {
  // TODO
  fmt.Fprintf(w, "Welcome to HTTP authentification over XMPP")
}

func authHandler(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  chanAnswer := make(chan bool)

  ChanRequest <- strings.Join(r.Form[PARAM_JID], "")
  ChanRequest <- strings.Join(r.Form[METHOD_ACCESS], "")
  ChanRequest <- strings.Join(r.Form[DOMAIN_ACCESS], "")
  ChanRequest <- strings.Join(r.Form[TRANSACTION_ID], "")
  ChanRequest <- chanAnswer

  answer := <- chanAnswer
  ret := RETURN_VALUE_NOK
  if answer {
    ret = RETURN_VALUE_OK
  }
  fmt.Fprintf(w, ret)
}



func Run() {
  http.HandleFunc("/", indexHandler) // set router
  http.HandleFunc("/toto", authHandler)

  port := strconv.Itoa(HttpPortBind)
  log.Println("Listenning on port "+port)
  err := http.ListenAndServe(":"+port, nil) // set listen port
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
