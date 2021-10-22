package rest

import (
	"encoding/json"
	"github.com/jslowik/commacloner/api/threecommas/rest/dobjs"
	"github.com/jslowik/commacloner/config"
	"net/http"
	"reflect"
	"testing"
)

func TestGetExchangeAccounts(t *testing.T) {
	testExchanges := []dobjs.Exchange{
		{
			ID:           123456,
			ExchangeName: "Coinbase Pro (GDAX)",
			MarketCode:   "gdax",
		},
		{
			ID:           56789,
			ExchangeName: "Kucoin",
			MarketCode:   "kucoin",
		},
	}

	type customHandlerFields struct {
		handlerPath string
		handler     func(w http.ResponseWriter, r *http.Request)
	}

	tests := []struct {
		name    string
		handler customHandlerFields
		wantErr bool
	}{
		{
			name: "Clean Path",
			handler: customHandlerFields{
				handlerPath: GetExchanges,
				handler: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					data, _ := json.Marshal(testExchanges)
					_, e := w.Write(data)
					if e != nil {
						t.Fatalf("could not marshal test data: %v", e)
					}
				},
			},
			wantErr: false,
		},
		{
			name: "Unauthorized User",
			handler: customHandlerFields{
				handlerPath: GetExchanges,
				handler: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusUnauthorized)
					msg := struct {
						Error       string `json:"error"`
						Description string `json:"error_description"`
					}{
						Error:       "api_key_invalid_or_expired",
						Description: "Unauthorized. Invalid or expired api key.",
					}
					data, _ := json.Marshal(msg)
					_, e := w.Write(data)
					if e != nil {
						t.Fatalf("could not marshal test data: %v", e)
					}
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test3CServer, _ := newTest3CServer(tt.handler.handlerPath, tt.handler.handler)
			apiConfig := config.API{
				RestURL: test3CServer.URL,
			}

			if accounts, err := GetExchangeAccounts(apiConfig); (err != nil) != tt.wantErr {
				t.Errorf("GetExchangeAccounts() error = %v, wantErr %v", err, tt.wantErr)
			} else if reflect.DeepEqual(testExchanges, accounts) == false && !tt.wantErr {
				t.Errorf("accounts not equal.")
			}
		})
	}
}

func TestGetExchangePairs(t *testing.T) {
	testPairs := []string{
		"USDT_BTC",
		"USDT_LTC",
		"USD_BTC",
		"USD_LTC",
		"BTC_ETH",
		"BTC_LTC",
	}

	type customHandlerFields struct {
		handlerPath string
		handler     func(w http.ResponseWriter, r *http.Request)
	}

	tests := []struct {
		name       string
		marketCode string
		handler    customHandlerFields
		wantErr    bool
	}{
		{
			name: "Clean Path",
			handler: customHandlerFields{
				handlerPath: GetMarketPairs,
				handler: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					data, _ := json.Marshal(testPairs)
					_, e := w.Write(data)
					if e != nil {
						t.Fatalf("could not marshal test data: %v", e)
					}
				},
			},
			wantErr:    false,
			marketCode: "gdax",
		},
		{
			name: "Unauthorized User",
			handler: customHandlerFields{
				handlerPath: GetMarketPairs,
				handler: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusUnauthorized)
					msg := struct {
						Error       string `json:"error"`
						Description string `json:"error_description"`
					}{
						Error:       "api_key_invalid_or_expired",
						Description: "Unauthorized. Invalid or expired api key.",
					}
					data, _ := json.Marshal(msg)
					_, e := w.Write(data)
					if e != nil {
						t.Fatalf("could not marshal test data: %v", e)
					}
				},
			},
			wantErr:    true,
			marketCode: "gdax",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test3CServer, _ := newTest3CServer(tt.handler.handlerPath, tt.handler.handler)
			apiConfig := config.API{
				RestURL: test3CServer.URL,
			}

			if pairs, err := GetExchangePairs(apiConfig, tt.marketCode); (err != nil) != tt.wantErr {
				t.Errorf("GetExchangePairs() error = %v, wantErr %v", err, tt.wantErr)
			} else if reflect.DeepEqual(testPairs, pairs) == false && !tt.wantErr {
				t.Errorf("accounts not equal.")
			}
		})
	}
}
