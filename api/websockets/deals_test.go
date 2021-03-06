package websockets

import (
	"github.com/gorilla/mux"
	"github.com/jslowik/commacloner/api"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/jslowik/commacloner/config"
)

const (
	StartNewDealPath  = "/ver1/bots/{id:[a-zA-Z0-9]+}/start_new_deal"
	CancelDealPath    = "/ver1/deals/{id:[a-zA-Z0-9]+}/cancel"
	PanicSellDealPath = "/ver1/deals/{id:[a-zA-Z0-9]+}/panic_sell"
)

// NewTest3CServer mocks the 3Commas API Server.  pass in a func to set a custom request handler
func NewTest3CServer(customPath string, customFunc func(w http.ResponseWriter, r *http.Request)) (*httptest.Server, error) {
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

func TestDealsStream_Build(t *testing.T) {
	tests := []struct {
		name      string
		APIKey    string
		APISecret string
		want      *Message
		wantErr   bool
	}{
		{
			name:      "clean path",
			APIKey:    "myapikey",
			APISecret: "s0m3s3cr3t!!",
			want: &Message{
				Identifier: "{\"channel\":\"DealsChannel\",\"users\":[{\"api_key\":\"myapikey\",\"signature\":\"0a77586521ce9d268f87e6d3bcf5a3c0995481c37dce4502914d07f61562f57f\"}]}",
				Command:    "subscribe",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DealsStream{
				APIConfig: config.API{
					Key:    tt.APIKey,
					Secret: tt.APISecret,
				},
			}
			got, err := d.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Build() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDealsStream_BuildSignature(t *testing.T) {
	tests := []struct {
		name      string
		APIKey    string
		APISecret string
		want      string
	}{
		{
			name:      "clean path",
			APIKey:    "myapikey",
			APISecret: "s0m3s3cr3t!!",
			want:      "0a77586521ce9d268f87e6d3bcf5a3c0995481c37dce4502914d07f61562f57f",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DealsStream{
				APIConfig: config.API{
					Key:    tt.APIKey,
					Secret: tt.APISecret,
				},
			}
			if got := d.BuildSignature(dealsEndpoint); got != tt.want {
				t.Errorf("BuildSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDealsStream_buildIdentifier(t *testing.T) {
	tests := []struct {
		name      string
		APIKey    string
		APISecret string
		want      Identifier
	}{
		{
			name:      "clean path",
			APIKey:    "myapikey",
			APISecret: "s0m3s3cr3t!!",
			want: Identifier{
				Channel: "DealsChannel",
				Users: []User{
					{
						APIKey:    "myapikey",
						Signature: "0a77586521ce9d268f87e6d3bcf5a3c0995481c37dce4502914d07f61562f57f",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DealsStream{
				APIConfig: config.API{
					Key:    tt.APIKey,
					Secret: tt.APISecret,
				},
			}
			if got := d.buildIdentifier(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildIdentifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDealsStream_HandleDeal(t *testing.T) {
	type customHandlerFields struct {
		handlerPath string
		handler     func(w http.ResponseWriter, r *http.Request)
	}

	tests := []struct {
		name    string
		config  config.API
		botMaps map[int][]config.BotMapping
		deal    api.DealsMessage
		handler customHandlerFields
		wantErr bool
	}{
		{
			name:    "safety trade should ignore",
			botMaps: nil,
			deal: api.DealsMessage{
				Details: api.DealDetails{
					Status:                     "bought",
					CompletedSafetyOrdersCount: 1,
				},
			},
			wantErr: false,
		},
		{
			name:    "no defined bots",
			config:  config.API{},
			botMaps: nil,
			deal: api.DealsMessage{
				Details: api.DealDetails{
					Status: "bought",
				},
			},
			wantErr: true,
		},
		{
			name:   "success path",
			config: config.API{},
			botMaps: map[int][]config.BotMapping{
				1234: {{
					ID: "example",
					Source: config.BotConfig{
						ID: 1234,
					},
					Destination: config.BotConfig{
						ID: 5678,
					},
				}},
			},
			deal: api.DealsMessage{
				Details: api.DealDetails{
					BotID:                      1234,
					Status:                     "bought",
					CompletedSafetyOrdersCount: 0,
					Pair:                       "BTC_USD",
				},
			},
			wantErr: false,
		},
		{
			name: "max deals reached",
			botMaps: map[int][]config.BotMapping{
				1234: {{
					ID: "example",
					Source: config.BotConfig{
						ID: 1234,
					},
					Destination: config.BotConfig{
						ID: 5678,
					},
					Overrides: config.BotOverrides{
						CancelUnavailableDeals: true,
					},
				}},
			},
			deal: api.DealsMessage{
				Details: api.DealDetails{
					BotID:                      1234,
					Status:                     "bought",
					CompletedSafetyOrdersCount: 0,
					Pair:                       "BTC_USD",
				},
			},
			handler: customHandlerFields{
				handlerPath: StartNewDealPath,
				handler: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusUnprocessableEntity)
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//logger := zap.NewExample()
			test3CServer, _ := NewTest3CServer(tt.handler.handlerPath, tt.handler.handler)

			d := DealsStream{
				APIConfig: tt.config,
				Bots:      tt.botMaps,
			}
			d.APIConfig.RestURL = test3CServer.URL
			if err := d.HandleDeal(tt.deal); (err != nil) != tt.wantErr {
				t.Errorf("HandleDeal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
