package http

import (
	"git.kingpenguin.tk/chteufleur/HTTPAuthentificationOverXMPP.git/xmpp"

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
	TIMEOUTE       = "timeout"

	ROUTE_ROOT = "/"
	ROUTE_AUTH = "/auth"

	RETURN_VALUE_OK  = "OK"
	RETURN_VALUE_NOK = "NOK"

	StatusUnknownError = 520
	StatusUnreachable  = 523
)

var (
	HttpPortBind  = 9090
	HttpsPortBind = 9093
	CertPath      = "./cert.pem"
	KeyPath       = "./key.pem"

	ChanRequest = make(chan interface{}, 5)
	TimeoutSec  = 60  // 1 min
	MaxTimeout  = 300 // 5 min
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

	if jid == "" || method == "" || domain == "" {
		// If mandatory params is missing
		log.Printf("%sMandatory params is missing", LogInfo)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	transaction := strings.Join(r.Form[TRANSACTION_ID], "")
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
		go runHttp()
	}
	if HttpsPortBind > 0 {
		go runHttps()
	}
}

func runHttp() {
	port := strconv.Itoa(HttpPortBind)
	log.Printf("%sHTTP listenning on port %s", LogInfo, port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("%sListenAndServe: ", LogError, err)
	}
}

func runHttps() {
	port := strconv.Itoa(HttpsPortBind)
	log.Printf("%sHTTPS listenning on port %s", LogInfo, port)
	err := http.ListenAndServeTLS(":"+port, CertPath, KeyPath, nil)
	if err != nil {
		log.Fatal("%sListenAndServe: ", LogError, err)
	}
}
