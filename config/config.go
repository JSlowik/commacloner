package config

import (
	"fmt"
	"strings"

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

type LunarCrushAPI struct {
	Key     string `json:"key"`
	RestURL string `json:"rest_url"`
}

// BotMapping contains both a source bot id to look for deals from the websockets api, and a destination bot to generate
// a matching deal for.  Additional overrides are allowed for conversion between stablecoins/currencies. (ie converting
// USDT to USD)
type BotMapping struct {
	ID          string       `json:"id"`
	Source      BotConfig    `json:"source"`
	Destination BotConfig    `json:"dest"`
	Overrides   BotOverrides `json:"overrides"`
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
		{c.Key == "", "no api key specified in config file"},
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
