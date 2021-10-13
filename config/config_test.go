package config

import (
	"testing"

	"go.uber.org/zap"
)

func TestConfig_Validate(t *testing.T) {
	baselineConfig := Config{
		Logging: Logger{
			Level:       zap.DebugLevel.String(),
			Format:      "json",
			Destination: "console",
		},
		API: API{
			Key:          "abcd1234",
			Secret:       "a1b2c3d4e5",
			WebsocketURL: "wss://ws.3commas.io/websocket",
			RestURL:      "https://api.3commas.io/public/api",
		},
		LunarcrushAPI: LunarCrushAPI{
			Enabled: true,
			Key:     "qwerty",
			RestURL: "https://api.lunarcrush.com/v2",
		},
		Bots: []BotMapping{
			{
				ID: "longbot",
				Source: BotConfig{
					ID: 1234,
				},
				Destination: BotConfig{
					ID: 5678,
				},
				Overrides: BotOverrides{
					QuoteCurrency:             "USD",
					CancelUnavailableDeals:    false,
					PanicSellUnavailableDeals: false,
				},
			},
		},
	}

	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name:   "success case",
			config: baselineConfig,
		},
		{
			name: "no api key",
			config: Config{
				Logging: baselineConfig.Logging,
				API: API{
					Key:          "",
					Secret:       baselineConfig.API.Secret,
					WebsocketURL: baselineConfig.API.WebsocketURL,
					RestURL:      baselineConfig.API.RestURL,
				},
				Bots: baselineConfig.Bots,
			},
			wantErr: true,
		},
		{
			name: "no api secret",
			config: Config{
				Logging: baselineConfig.Logging,
				API: API{
					Key:          baselineConfig.API.Key,
					Secret:       "",
					WebsocketURL: baselineConfig.API.WebsocketURL,
					RestURL:      baselineConfig.API.RestURL,
				},
				Bots: baselineConfig.Bots,
			},
			wantErr: true,
		},
		{
			name: "no websocket",
			config: Config{
				Logging: baselineConfig.Logging,
				API: API{
					Key:          baselineConfig.API.Key,
					Secret:       baselineConfig.API.Secret,
					WebsocketURL: "",
					RestURL:      baselineConfig.API.RestURL,
				},
				Bots: baselineConfig.Bots,
			},
			wantErr: true,
		},
		{
			name: "no rest api",
			config: Config{
				Logging: baselineConfig.Logging,
				API: API{
					Key:          baselineConfig.API.Key,
					Secret:       baselineConfig.API.Secret,
					WebsocketURL: baselineConfig.API.WebsocketURL,
					RestURL:      "",
				},
				Bots: baselineConfig.Bots,
			},
			wantErr: true,
		},
		{
			name: "no bot mappings",
			config: Config{
				Logging: baselineConfig.Logging,
				API:     baselineConfig.API,
				Bots:    nil,
			},
			wantErr: true,
		},
		{
			name: "no bot mapping id",
			config: Config{
				Logging: baselineConfig.Logging,
				API:     baselineConfig.API,
				Bots: []BotMapping{
					{
						ID:          "",
						Source:      baselineConfig.Bots[0].Source,
						Destination: baselineConfig.Bots[0].Destination,
						Overrides:   baselineConfig.Bots[0].Overrides,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no bot source id",
			config: Config{
				Logging: baselineConfig.Logging,
				API:     baselineConfig.API,
				Bots: []BotMapping{
					{
						ID:          "",
						Source:      BotConfig{},
						Destination: baselineConfig.Bots[0].Destination,
						Overrides:   baselineConfig.Bots[0].Overrides,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no bot destination id",
			config: Config{
				Logging: baselineConfig.Logging,
				API:     baselineConfig.API,
				Bots: []BotMapping{
					{
						ID:          "",
						Source:      baselineConfig.Bots[0].Source,
						Destination: BotConfig{},
						Overrides:   baselineConfig.Bots[0].Overrides,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "bad log format",
			config: Config{
				Logging: Logger{
					Level:       "debug",
					Format:      "html",
					Destination: "console",
				},
				API:  baselineConfig.API,
				Bots: baselineConfig.Bots,
			},
			wantErr: true,
		},
		{
			name: "invalid log level",
			config: Config{
				Logging: Logger{
					Level:       "fakelevel",
					Format:      "json",
					Destination: "console",
				},
				API:  baselineConfig.API,
				Bots: baselineConfig.Bots,
			},
			wantErr: true,
		},
		{
			name: "invalid log destination",
			config: Config{
				Logging: Logger{
					Level:       "debug",
					Format:      "json",
					Destination: "/my/fake/path",
				},
				API:  baselineConfig.API,
				Bots: baselineConfig.Bots,
			},
			wantErr: true,
		},
		{
			name: "lunarcrush enabled but no api key",
			config: Config{
				Logging: baselineConfig.Logging,
				API:     baselineConfig.API,
				Bots:    baselineConfig.Bots,
				LunarcrushAPI: LunarCrushAPI{
					Enabled: true,
					Key:     "",
					RestURL: baselineConfig.LunarcrushAPI.RestURL,
				},
			},
			wantErr: true,
		},
		{
			name: "lunarcrush enabled but no url",
			config: Config{
				Logging: baselineConfig.Logging,
				API:     baselineConfig.API,
				Bots:    baselineConfig.Bots,
				LunarcrushAPI: LunarCrushAPI{
					Enabled: true,
					Key:     baselineConfig.LunarcrushAPI.Key,
					RestURL: "",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.config
			if err := c.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
