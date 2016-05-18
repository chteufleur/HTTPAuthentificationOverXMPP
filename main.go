package main

import (
  "git.kingpenguin.tk/chteufleur/HTTPAuthentificationOverXMPP.git/http"

  "github.com/jimlawless/cfg"

  "log"
  "os"
	"os/signal"
	"syscall"
	"time"
)

const (
	Version               = "v0.1.0"
	configurationFilePath = "httpAuth.cfg"
)

var (
	mapConfig = make(map[string]string)
)

func init() {
  err := cfg.Load(configurationFilePath, mapConfig)
	if err != nil {
		log.Fatal("Failed to load configuration file.", err)
	}
  // TODO make config
}


func request() {
  for {
    jid := <- http.ChanRequest
    log.Println(jid)
    method := <- http.ChanRequest
    log.Println(method)
    domain := <- http.ChanRequest
    log.Println(domain)
    transaction := <- http.ChanRequest
    log.Println(transaction)

    chanResult := <- http.ChanRequest
    // TODO make the XMPP request
    if v, ok := chanResult.(chan bool); ok {
      v <- false
    }
  }
}

func main() {
  // TODO start ressources
  go http.Run()
  go request()

  sigchan := make(chan os.Signal, 1)
  signal.Notify(sigchan, os.Interrupt)
  signal.Notify(sigchan, syscall.SIGTERM)
  signal.Notify(sigchan, os.Kill)
  <-sigchan

  // TODO close all ressources

  log.Println("Exit main()")
  time.Sleep(1 * time.Second)
}
