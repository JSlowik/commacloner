package websockets

import (
	"github.com/jslowik/commacloner/api/threecommas/dobjs"
	dobjs2 "github.com/jslowik/commacloner/api/threecommas/websockets/dobjs"
	"net/http"
	"reflect"
	"testing"

	"github.com/jslowik/commacloner/config"
)

func TestDealsStream_Build(t *testing.T) {
	tests := []struct {
		name      string
		APIKey    string
		APISecret string
		want      *dobjs2.IdentifierMessage
		wantErr   bool
	}{
		{
			name:      "clean path",
			APIKey:    "myapikey",
			APISecret: "s0m3s3cr3t!!",
			want: &dobjs2.IdentifierMessage{
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
		want      dobjs2.Identifier
	}{
		{
			name:      "clean path",
			APIKey:    "myapikey",
			APISecret: "s0m3s3cr3t!!",
			want: dobjs2.Identifier{
				Channel: "DealsChannel",
				Users: []dobjs2.User{
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
		deal    dobjs2.DealsMessage
		handler customHandlerFields
		wantErr bool
	}{
		{
			name:    "safety trade should ignore",
			botMaps: nil,
			deal: dobjs2.DealsMessage{
				Details: dobjs.Deal{
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
			deal: dobjs2.DealsMessage{
				Details: dobjs.Deal{
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
			deal: dobjs2.DealsMessage{
				Details: dobjs.Deal{
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
			deal: dobjs2.DealsMessage{
				Details: dobjs.Deal{
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
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test3CServer, _ := newTest3CServer(tt.handler.handlerPath, tt.handler.handler)

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
