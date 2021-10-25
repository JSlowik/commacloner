package config

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Config is the top level of the configuration yaml
type Config struct {
	API           API           `json:"api"`
	LunarcrushAPI LunarCrushAPI `json:"lunarcrush"`
	Bots          []BotMapping  `json:"bots"`
	Logging       Logger        `json:"logging"`
}

// Logger holds configuration required to customize logging
type Logger struct {
	// Level sets logging level severity.
	Level string `json:"level"`

	// Format specifies the format to be used for logging.
	Format string `json:"format"`

	//Destination
	Destination string `json:"destination"`
}

// API contains the configuration elements for an API
type API struct {
	Key          string `json:"key"`
	Secret       string `json:"secret"`
	WebsocketURL string `json:"websocket_url"`
	RestURL      string `json:"rest_url"`
}

// LunarCrushAPI contains configuration elements for the LunarCrush API
type LunarCrushAPI struct {
	Enabled   bool       `json:"enabled"`
	Key       string     `json:"key"`
	RestURL   string     `json:"rest_url"`
	Blacklist []string   `json:"blacklist"`
	Cache     LunarCache `json:"cache"`
}

// LunarCache contains caching configurations
type LunarCache struct {
	Enabled bool   `json:"enabled"`
	Every   string `json:"every"`
}

// BotMapping contains both a source bot id to look for deals from the websockets api, and a destination bot to generate
// a matching deal for.  Additional overrides are allowed for conversion between stablecoins/currencies. (ie converting
// USDT to USD)
type BotMapping struct {
	ID          string             `json:"id"`
	Source      BotConfig          `json:"source"`
	Destination BotConfig          `json:"dest"`
	Overrides   BotOverrides       `json:"overrides"`
	Pairs       PairsConfiguration `json:"pairs"`
}

// BotConfig contains configuration elements for the 3commas bots.
type BotConfig struct {
	ID int `json:"bot_id"`
}

// BotOverrides contains bot-specific overrides when translating a deal from the source bot to the destination bot
// This includes both overriding a base/quote currency, or cancelling a deal on the source bot if it's not available
// on the destination bot's exchange.
type BotOverrides struct {
	QuoteCurrency             string `json:"quote_currency"`
	BaseCurrency              string `json:"base_currency"`
	CancelUnavailableDeals    bool   `json:"cancelUnavailableDeals"`
	PanicSellUnavailableDeals bool   `json:"panicSellUnavailableDeals"`
}

// PairsConfiguration contains how a bot mapping should or should not update its pairs.  if manual, it will be left alone
type PairsConfiguration struct {
	Mode   string             `json:"mode"`
	Config LunarConfiguration `json:"config"`
}

// LunarConfiguration contains LunarCrush configurations such as currencies, refresh time, sort category, and max pairs
// to return
type LunarConfiguration struct {
	QuoteCurrency struct {
		Source string `json:"source"`
		Dest   string `json:"dest"`
	} `json:"quote_currency"`
	Refresh  string `json:"refresh"`
	Category string `json:"category"`
	MaxPairs int    `json:"max"`
}

// Validate the configuration
func (c Config) Validate() error {
	// Fast checks. Perform these first for a more responsive CLI.
	checks := []struct {
		bad    bool
		errMsg string
	}{
		{len(c.Bots) == 0, "no bot mappings defined"},
	}

	var checkErrors []string

	for _, check := range checks {
		if check.bad {
			checkErrors = append(checkErrors, check.errMsg)
		}
	}

	// Validate the logging configs
	checkErrors = append(checkErrors, c.Logging.validate()...)

	// Validate the 3Commas API configs
	checkErrors = append(checkErrors, c.API.validate()...)

	// Validate the Lunarcrush API configs
	checkErrors = append(checkErrors, c.LunarcrushAPI.validate()...)

	// Validate the bot mappings
	for _, mapping := range c.Bots {
		checkErrors = append(checkErrors, mapping.validate()...)
	}

	if len(checkErrors) != 0 {
		return fmt.Errorf("invalid Config:\n\t-\t%s", strings.Join(checkErrors, "\n\t-\t"))
	}
	return nil
}

func (c Logger) validate() []string {
	// Fast checks. Perform these first for a more responsive CLI.
	checks := []struct {
		bad    bool
		errMsg string
	}{
		{c.Format == "" || c.Format != "json" && c.Format != "console", "log format must be \"json\" or \"console\""},

		{c.Destination == "" || c.Destination != "file" && c.Destination != "console", "log destination must be \"file\" or \"console\""},
	}

	var checkErrors []string
	for _, check := range checks {
		if check.bad {
			checkErrors = append(checkErrors, check.errMsg)
		}
	}
	var lvl zap.AtomicLevel
	err := lvl.UnmarshalText([]byte(c.Level))
	if err != nil {
		checkErrors = append(checkErrors, "invalid log level: "+c.Level)
	}

	return checkErrors
}

func (c API) validate() []string {
	// Fast checks. Perform these first for a more responsive CLI.
	checks := []struct {
		bad    bool
		errMsg string
	}{
		{c.Key == "", "no api key specified in config file"},
		{c.Secret == "", "no api secret in config file"},
		{c.WebsocketURL == "", "no websocket url defined"},
		{c.RestURL == "", "no rest url defined"},
	}

	var checkErrors []string

	for _, check := range checks {
		if check.bad {
			checkErrors = append(checkErrors, check.errMsg)
		}
	}
	return checkErrors
}

func (c LunarCrushAPI) validate() []string {
	// Fast checks. Perform these first for a more responsive CLI.
	checks := []struct {
		bad    bool
		errMsg string
	}{
		{c.Enabled && c.Key == "", "no api key specified in config file"},
		{c.Enabled && c.RestURL == "", "no rest url defined"},
	}

	var checkErrors []string

	for _, check := range checks {
		if check.bad {
			checkErrors = append(checkErrors, check.errMsg)
		}
	}

	// Validate the Lunarcrush Cache configs
	checkErrors = append(checkErrors, c.Cache.validate()...)

	return checkErrors
}

func (c LunarCache) validate() []string {
	// Fast checks. Perform these first for a more responsive CLI.
	checks := []struct {
		bad    bool
		errMsg string
	}{
		{c.Enabled && c.Every == "", "no cache refresh defined"},
	}

	var checkErrors []string

	for _, check := range checks {
		if check.bad {
			checkErrors = append(checkErrors, check.errMsg)
		}
	}

	//validate duration
	if c.Every != "" {
		d, err := time.ParseDuration(c.Every)
		if err != nil {
			checkErrors = append(checkErrors, fmt.Sprintf("invalid duriation for cache refresh: %v", err))
		} else if d < 0 {
			checkErrors = append(checkErrors, fmt.Sprintf("cache refresh time cannot be less than 0: %s", c.Every))
		}
	}

	return checkErrors
}

func (m BotMapping) validate() []string {
	var checkErrors []string

	checks := []struct {
		bad    bool
		errMsg string
	}{
		{m.ID == "", "no bot mapping id defined"},
	}
	for _, check := range checks {
		if check.bad {
			checkErrors = append(checkErrors, check.errMsg)
		}
	}

	// Validate BotConfigs
	checkErrors = append(checkErrors, m.Source.validate()...)
	checkErrors = append(checkErrors, m.Destination.validate()...)
	checkErrors = append(checkErrors, m.Pairs.validate()...)
	return checkErrors
}

func (m BotConfig) validate() []string {
	var checkErrors []string

	checks := []struct {
		bad    bool
		errMsg string
	}{
		{m.ID == 0, "no bot mapping id defined"},
	}

	for _, check := range checks {
		if check.bad {
			checkErrors = append(checkErrors, check.errMsg)
		}
	}
	return checkErrors
}

func (p PairsConfiguration) validate() []string {
	var checkErrors []string

	checks := []struct {
		bad    bool
		errMsg string
	}{
		{p.Mode != "manual" && p.Mode != "lunarcrush", "invalid mode, use 'manual' or 'lunarcrush'"},
	}

	for _, check := range checks {
		if check.bad {
			checkErrors = append(checkErrors, check.errMsg)
		}
	}
	if len(checkErrors) == 0 && p.Mode == "lunarcrush" {
		c := p.Config
		lcChecks := []struct {
			bad    bool
			errMsg string
		}{
			{c.Category != "galaxyscore" && c.Category != "altrank", "invalid category, use 'galaxyscore' or 'altrank'"},
			{c.MaxPairs <= 0, "Max Pairs must be greater than 0"},
			{c.QuoteCurrency.Source == "", "source quote currency cannot be blank"},
			{c.QuoteCurrency.Dest == "", "destination quote currency cannot be blank"},
		}
		for _, lcCheck := range lcChecks {
			if lcCheck.bad {
				checkErrors = append(checkErrors, lcCheck.errMsg)
			}
		}
		//validate duration
		d, err := time.ParseDuration(c.Refresh)
		if err != nil {
			checkErrors = append(checkErrors, fmt.Sprintf("invalid duriation for pair refresh: %v", err))
		} else if d <= 0 {
			checkErrors = append(checkErrors, fmt.Sprintf("pair refresh time cannot be less than 0: %s", c.Refresh))
		}
	}

	return checkErrors
}
