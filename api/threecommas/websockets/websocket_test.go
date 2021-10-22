package websockets

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
)

const (
	StartNewDealPath  = "/ver1/bots/{id:[a-zA-Z0-9]+}/start_new_deal"
	CancelDealPath    = "/ver1/deals/{id:[a-zA-Z0-9]+}/cancel"
	PanicSellDealPath = "/ver1/deals/{id:[a-zA-Z0-9]+}/panic_sell"
)

// NewTest3CServer mocks the 3Commas API Server.  pass in a func to set a custom request handler
func newTest3CServer(customPath string, customFunc func(w http.ResponseWriter, r *http.Request)) (*httptest.Server, error) {
	rtr := mux.NewRouter()

	dealFunc := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}

	if customFunc != nil {
		dealFunc = customFunc
	}

	if customPath != StartNewDealPath {
		rtr.HandleFunc(StartNewDealPath, dealFunc)
	}

	if customPath != CancelDealPath {
		rtr.HandleFunc(CancelDealPath, dealFunc)
	}

	if customFunc != nil {
		rtr.HandleFunc(customPath, customFunc)
	}
	return httptest.NewServer(rtr), nil
}
