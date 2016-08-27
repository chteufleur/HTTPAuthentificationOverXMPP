package http

import (
	"git.kingpenguin.tk/chteufleur/HTTPAuthentificationOverXMPP.git/xmpp"

	"fmt"
	"log"
	"math/rand"
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
	TIMEOUTE       = "timeout"

	ROUTE_ROOT = "/"
	ROUTE_AUTH = "/auth"

	RETURN_VALUE_OK  = "OK"
	RETURN_VALUE_NOK = "NOK"

	MAX_PORT_VAL = 65535

	StatusUnknownError = 520
	StatusUnreachable  = 523
)

var (
	HttpPortBind  = -1
	HttpsPortBind = -1
	CertPath      = "./cert.pem"
	KeyPath       = "./key.pem"

	ChanRequest = make(chan interface{}, 5)
	TimeoutSec  = 60  // 1 min
	MaxTimeout  = 300 // 5 min

	BindAddressIPv4 = "127.0.0.1"
	BindAddressIPv6 = "[::1]"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to HTTP authentification over XMPP")
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	jid := strings.Join(r.Form[PARAM_JID], "")
	method := strings.Join(r.Form[METHOD_ACCESS], "")
	domain := strings.Join(r.Form[DOMAIN_ACCESS], "")
	transaction := strings.Join(r.Form[TRANSACTION_ID], "")

	if jid == "" || method == "" || domain == "" || transaction == "" {
		// If mandatory params is missing
		log.Printf("%sMandatory params is missing", LogInfo)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	timeoutStr := strings.Join(r.Form[TIMEOUTE], "")
	log.Printf("%sAuth %s", LogInfo, jid)
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil || timeout <= 0 {
		timeout = TimeoutSec
	}
	if timeout > MaxTimeout {
		timeout = MaxTimeout
	}

	chanAnswer := make(chan string)

	confirmation := new(xmpp.Confirmation)
	confirmation.JID = jid
	confirmation.Method = method
	confirmation.Domain = domain
	confirmation.Transaction = transaction
	confirmation.ChanReply = chanAnswer
	confirmation.SendConfirmation()

	select {
	case answer := <-chanAnswer:
		switch answer {
		case xmpp.REPLY_OK:
			w.WriteHeader(http.StatusOK)

		case xmpp.REPLY_DENY:
			w.WriteHeader(http.StatusUnauthorized)

		case xmpp.REPLY_UNREACHABLE:
			w.WriteHeader(StatusUnreachable)

		default:
			w.WriteHeader(StatusUnknownError)
		}
	case <-time.After(time.Duration(timeout) * time.Second):
		w.WriteHeader(http.StatusUnauthorized)
	}

	switch confirmation.TypeSend {
	case xmpp.TYPE_SEND_IQ:
		log.Printf("%sDelete IQ", LogDebug)
		delete(xmpp.WaitIqMessages, confirmation.IdMap)

	case xmpp.TYPE_SEND_MESSAGE:
		log.Printf("%sDelete Message", LogDebug)
		delete(xmpp.WaitMessageAnswers, confirmation.IdMap)
	}
}

func Run() {
	log.Printf("%sRunning", LogInfo)

	http.HandleFunc(ROUTE_ROOT, indexHandler)
	http.HandleFunc(ROUTE_AUTH, authHandler)

	if HttpPortBind > 0 {
		go runHttp(BindAddressIPv4)
		if BindAddressIPv4 != "0.0.0.0" {
			go runHttp(BindAddressIPv6)
		}
	} else if HttpPortBind == 0 {
		HttpPortBind = rand.Intn(MAX_PORT_VAL)
		go runHttp(BindAddressIPv4)
		if BindAddressIPv4 != "0.0.0.0" {
			go runHttp(BindAddressIPv6)
		}
	}
	if HttpsPortBind > 0 {
		go runHttps(BindAddressIPv4)
		if BindAddressIPv6 != "0.0.0.0" {
			go runHttps(BindAddressIPv6)
		}
	} else if HttpsPortBind == 0 {
		HttpsPortBind = rand.Intn(MAX_PORT_VAL)
		go runHttps(BindAddressIPv4)
		if BindAddressIPv6 != "0.0.0.0" {
			go runHttps(BindAddressIPv6)
		}
	}
}

func runHttp(bindAddress string) {
	port := strconv.Itoa(HttpPortBind)
	log.Printf("%sHTTP listenning on %s:%s", LogInfo, bindAddress, port)
	err := http.ListenAndServe(bindAddress+":"+port, nil)
	if err != nil {
		log.Fatal("%sListenAndServe: ", LogError, err)
	}
}

func runHttps(bindAddress string) {
	port := strconv.Itoa(HttpsPortBind)
	log.Printf("%sHTTPS listenning on %s:%s", LogInfo, bindAddress, port)
	err := http.ListenAndServeTLS(bindAddress+":"+port, CertPath, KeyPath, nil)
	if err != nil {
		log.Fatal("%sListenAndServe: ", LogError, err)
	}
}
