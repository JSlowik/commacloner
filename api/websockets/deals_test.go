package websockets

import (
	"reflect"
	"testing"

	"github.com/jslowik/commacloner/config"
)

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

//func TestDealsStream_HandleDeal(t *testing.T) {
//	tests := []struct {
//		name    string
//		config config.API
//		botMaps map[int][]config.BotMapping
//		rawMessage *Message
//		wantErr bool
//	}{
//		{
//			name: "non deal message",
//			config: config.API{},
//			botMaps: map[int][]config.BotMapping{
//				1234:{{
//					ID:          "example",
//					Source:      config.BotConfig{},
//					Destination: config.BotConfig{},
//					Overrides:   config.BotOverrides{},
//				}},
//			},
//			rawMessage: &Message{
//					Type:       "ping",
//			},
//		},
//		//{
//		//	name: "deal message, ignored status",
//		//	config: config.API{},
//		//	botMaps: map[int][]config.BotMapping{
//		//		1234:{{
//		//			ID:          "example",
//		//			Source:      config.BotConfig{},
//		//			Destination: config.BotConfig{},
//		//			Overrides:   config.BotOverrides{},
//		//		}},
//		//	},
//		//	rawMessage: &api.Deal{
//		//		Identifier: "ping",
//		//	},
//		//},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			logger := zap.NewExample()
//
//			d := DealsStream{
//				APIConfig: tt.config,
//				Bots:      tt.botMaps,
//			}
//			m, e := json.Marshal(tt.rawMessage)
//			if e != nil {
//				t.Fatalf("cannot marshall test message: %v", e)
//			}
//			if err := d.HandleDeal(m, logger); (err != nil) != tt.wantErr {
//				t.Errorf("HandleDeal() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
