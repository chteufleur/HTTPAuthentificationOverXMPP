package main

import (
	"git.kingpenguin.tk/chteufleur/HTTPAuthentificationOverXMPP.git/http"
	"git.kingpenguin.tk/chteufleur/HTTPAuthentificationOverXMPP.git/xmpp"

	"github.com/jimlawless/cfg"

	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	Version               = "v0.3.1"
	configurationFilePath = "httpAuth.cfg"

	default_xmpp_server_address = "127.0.0.1"
	default_xmpp_server_port    = "5347"
)

var (
	mapConfig = make(map[string]string)
)

func init() {
	log.Printf("Running HTTP-Auth %v", Version)

	err := cfg.Load(configurationFilePath, mapConfig)
	if err != nil {
		log.Fatal("Failed to load configuration file.", err)
	}

	// HTTP config
	httpTimeout, err := strconv.Atoi(mapConfig["http_timeoute_sec"])
	if err == nil && httpTimeout > 0 && httpTimeout < http.MaxTimeout {
		log.Println("Define HTTP timeout to " + strconv.Itoa(httpTimeout) + " second")
		http.TimeoutSec = httpTimeout
	}
	httpPort, err := strconv.Atoi(mapConfig["http_port"])
	if err == nil {
		log.Println("Define HTTP port to " + strconv.Itoa(httpPort))
		http.HttpPortBind = httpPort
	}
	httpsPort, err := strconv.Atoi(mapConfig["https_port"])
	if err == nil {
		log.Println("Define HTTPS port to " + strconv.Itoa(httpsPort))
		http.HttpsPortBind = httpsPort
		http.CertPath = mapConfig["https_cert_path"]
		http.KeyPath = mapConfig["https_key_path"]
	}

	// XMPP config
	xmpp_server_address := mapConfig["xmpp_server_address"]
	if xmpp_server_address == "" {
		xmpp_server_address = default_xmpp_server_address
	}
	xmpp_server_port := mapConfig["xmpp_server_port"]
	if xmpp_server_port == "" {
		xmpp_server_port = default_xmpp_server_port
	}

	xmpp.Addr = xmpp_server_address + ":" + xmpp_server_port
	xmpp.JidStr = mapConfig["xmpp_hostname"]
	xmpp.Secret = mapConfig["xmpp_secret"]
	xmpp.Debug = mapConfig["xmpp_debug"] == "true"
}

func getChanString(c chan interface{}) string {
	ret := ""
	i := <-c
	if v, ok := i.(string); ok {
		ret = v
	}
	return ret
}

func main() {

	go http.Run()
	go xmpp.Run()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	signal.Notify(sigchan, syscall.SIGTERM)
	signal.Notify(sigchan, os.Kill)
	<-sigchan

	log.Println("Exit main()")
	time.Sleep(1 * time.Second)
}
