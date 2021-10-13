package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jslowik/commacloner/config"
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

	if customFunc != nil {
		rtr.HandleFunc(customPath, customFunc)
	}
	return httptest.NewServer(rtr), nil
}

func TestStartNewDeal(t *testing.T) {
	type customHandlerFields struct {
		handlerPath string
		handler     func(w http.ResponseWriter, r *http.Request)
	}

	tests := []struct {
		name      string
		apiConfig config.API
		bot       config.BotMapping
		handler   customHandlerFields
		pair      string
		wantErr   bool
	}{
		{
			name: "clean path, successful deal setting",
			apiConfig: config.API{
				Key:    "abcd1234",
				Secret: "zyxw9876",
			},
			bot: config.BotMapping{
				ID: "standard_mapping",
				Source: config.BotConfig{
					ID: 1,
				},
				Destination: config.BotConfig{
					ID: 2,
				},
				Overrides: config.BotOverrides{
					QuoteCurrency: "",
					BaseCurrency:  "",
				},
			},
			pair:    "BTC_USDT",
			wantErr: false,
		},
		{
			name: "clean path, override quote",
			apiConfig: config.API{
				Key:    "abcd1234",
				Secret: "zyxw9876",
			},
			bot: config.BotMapping{
				ID: "standard_mapping",
				Source: config.BotConfig{
					ID: 1,
				},
				Destination: config.BotConfig{
					ID: 2,
				},
				Overrides: config.BotOverrides{
					QuoteCurrency: "",
					BaseCurrency:  "USDC",
				},
			},
			pair:    "BTC_USDT",
			wantErr: false,
		},
		{
			name: "clean path, override base",
			apiConfig: config.API{
				Key:    "abcd1234",
				Secret: "zyxw9876",
			},
			bot: config.BotMapping{
				ID: "standard_mapping",
				Source: config.BotConfig{
					ID: 1,
				},
				Destination: config.BotConfig{
					ID: 2,
				},
				Overrides: config.BotOverrides{
					QuoteCurrency: "ETH",
					BaseCurrency:  "",
				},
			},
			pair:    "BTC_USDT",
			wantErr: false,
		},
		{
			name: "max deals reached",
			apiConfig: config.API{
				Key:    "abcd1234",
				Secret: "zyxw9876",
			},
			bot: config.BotMapping{
				ID: "standard_mapping",
				Source: config.BotConfig{
					ID: 1,
				},
				Destination: config.BotConfig{
					ID: 2,
				},
				Overrides: config.BotOverrides{
					QuoteCurrency: "",
					BaseCurrency:  "",
				},
			},
			handler: customHandlerFields{
				handlerPath: StartNewDealPath,
				handler: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusUnprocessableEntity)
				},
			},
			pair:    "BTC_USDT",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test3CServer, _ := newTest3CServer(tt.handler.handlerPath, tt.handler.handler)
			tt.apiConfig.RestURL = test3CServer.URL

			if err := StartNewDeal(tt.apiConfig, tt.bot, tt.pair); (err != nil) != tt.wantErr {
				t.Errorf("StartNewDeal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCancelDeal(t *testing.T) {
	type customHandlerFields struct {
		handlerPath string
		handler     func(w http.ResponseWriter, r *http.Request)
	}

	tests := []struct {
		name      string
		apiConfig config.API
		handler   customHandlerFields
		panicSell bool
		wantErr   bool
	}{
		{
			name: "cancel unavailable deal",
			apiConfig: config.API{
				Key:    "abcd1234",
				Secret: "zyxw9876",
			},
			handler: customHandlerFields{
				handlerPath: CancelDealPath,
				handler: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusCreated)
				},
			},
			panicSell: false,
			wantErr:   false,
		},
		{
			name: "panic sell unavailable deal",
			apiConfig: config.API{
				Key:    "abcd1234",
				Secret: "zyxw9876",
			},
			handler: customHandlerFields{
				handlerPath: PanicSellDealPath,
				handler: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusCreated)
				},
			},
			panicSell: true,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test3CServer, _ := newTest3CServer(tt.handler.handlerPath, tt.handler.handler)
			tt.apiConfig.RestURL = test3CServer.URL

			if err := CancelDeal(tt.apiConfig, 1234, tt.panicSell); (err != nil) != tt.wantErr {
				t.Errorf("StartNewDeal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
