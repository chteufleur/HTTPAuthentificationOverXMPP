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
)

var (
	HttpPortBind = 9090

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
	transaction := strings.Join(r.Form[TRANSACTION_ID], "")
	timeoutStr := strings.Join(r.Form[TIMEOUTE], "")
	log.Printf("%sAuth %s", LogDebug, jid)

	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		timeout = TimeoutSec
	}
	if timeout > MaxTimeout {
		timeout = MaxTimeout
	}

	chanAnswer := make(chan bool)

	ChanRequest <- jid
	ChanRequest <- method
	ChanRequest <- domain
	ChanRequest <- transaction
	ChanRequest <- chanAnswer

	select {
	case answer := <-chanAnswer:
		if answer {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	case <-time.After(time.Duration(timeout) * time.Second):
		w.WriteHeader(http.StatusUnauthorized)
		delete(xmpp.WaitMessageAnswers, transaction)
	}
}

func Run() {
	log.Printf("%sRunning", LogInfo)

	http.HandleFunc(ROUTE_ROOT, indexHandler)
	http.HandleFunc(ROUTE_AUTH, authHandler)

	port := strconv.Itoa(HttpPortBind)
	log.Printf("%sListenning on port %s", LogInfo, port)
	err := http.ListenAndServe(":"+port, nil) // set listen port
	if err != nil {
		log.Fatal("%sListenAndServe: ", LogError, err)
	}
}
