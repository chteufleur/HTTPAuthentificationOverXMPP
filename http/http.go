package http

import (
  "fmt"
  "log"
  "net/http"
  "strconv"
  "strings"
  "time"
)

const (
	LogInfo  = "\t[HTTP INFO]\t"
	LogError = "\t[HTTP ERROR]\t"
	LogDebug = "\t[HTTP DEBUG]\t"

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
  TimeoutSec = 60
)


func indexHandler(w http.ResponseWriter, r *http.Request) {
  // TODO
  fmt.Fprintf(w, "Welcome to HTTP authentification over XMPP")
}

func authHandler(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  jid := strings.Join(r.Form[PARAM_JID], "")
  method := strings.Join(r.Form[METHOD_ACCESS], "")
  domain := strings.Join(r.Form[DOMAIN_ACCESS], "")
  transaction := strings.Join(r.Form[TRANSACTION_ID], "")
  log.Printf("%sAuth %s", LogDebug, jid)

  chanAnswer := make(chan bool)

  ChanRequest <- jid
  ChanRequest <- method
  ChanRequest <- domain
  ChanRequest <- transaction
  ChanRequest <- chanAnswer

  ret := RETURN_VALUE_NOK
  answer := false
  select {
  case answer = <- chanAnswer:
  case <- time.After(time.Duration(TimeoutSec) * time.Second):
    answer = false
  }
  if answer {
    ret = RETURN_VALUE_OK
  }
  fmt.Fprintf(w, ret)
}



func Run() {
  log.Printf("%sRunning", LogInfo)
  http.HandleFunc("/", indexHandler) // set router
  http.HandleFunc("/toto", authHandler)

  port := strconv.Itoa(HttpPortBind)
  log.Printf("%sListenning on port %s", LogInfo, port)
  err := http.ListenAndServe(":"+port, nil) // set listen port
  if err != nil {
    log.Fatal("%sListenAndServe: ", LogError, err)
  }
}
