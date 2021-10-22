package rest

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	StartNewDealPath  = "/ver1/bots/{id:[a-zA-Z0-9]+}/start_new_deal"
	CancelDealPath    = "/ver1/deals/{id:[a-zA-Z0-9]+}/cancel"
	PanicSellDealPath = "/ver1/deals/{id:[a-zA-Z0-9]+}/panic_sell"
)

// newTest3CServer mocks the 3Commas API Server.  pass in a func to set a custom request handler
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

func Test_generateQuery(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		queryParameters map[string]string
		want            string
	}{
		{
			name: "generate new deal test",
			path: "https://api.3commas.io/public/api/ver1/bots/1234/start_new_deal",
			queryParameters: map[string]string{
				"pair":               "BTC_USD",
				"skip_signal_checks": "true",
			},
			want: "https://api.3commas.io/public/api/ver1/bots/1234/start_new_deal?pair=BTC_USD&skip_signal_checks=true",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateQuery(tt.path, tt.queryParameters); got.String() != tt.want {
				t.Errorf("generateQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
