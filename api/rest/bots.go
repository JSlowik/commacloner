package rest

import (
	"fmt"
	"github.com/jslowik/commacloner/config"
	"github.com/jslowik/commacloner/log"
	"strings"
)

const (
	pairParameter = "pair"
)

const (
	ShowBot          = "/ver1/bots/%d/show"
	StartNewBotDeal  = "/ver1/bots/%d/start_new_deal"
	CancelBotDeal    = "/ver1/deals/%d/cancel"
	PanicSellBotDeal = "/ver1/deals/%d/panic_sell"
)

// StartNewDeal invokes the API to start a new deal based on the bot mapping for the given pair
func StartNewDeal(apiConfig config.API, bot config.BotMapping, pair string) error {
	logger := log.NewLogger("bots")
	route := fmt.Sprintf(StartNewBotDeal, bot.Destination.ID)
	path := apiConfig.RestURL + route

	// override the pair if set
	pairSplit := strings.Split(pair, "_")
	if bot.Overrides.QuoteCurrency != "" {
		pairSplit[0] = bot.Overrides.QuoteCurrency
	}
	if bot.Overrides.BaseCurrency != "" {
		pairSplit[1] = bot.Overrides.BaseCurrency
	}

	params := make(map[string]string)
	params[pairParameter] = fmt.Sprintf("%s_%s", pairSplit[0], pairSplit[1])

	query := generateQuery(path, params)

	logger.Infof("generating new deal: %s", query.String())
	_, err := makeRequest("POST", query, apiConfig)

	return err
}

// CancelDeal cancels an existing deal
func CancelDeal(apiConfig config.API, dealID int, panicSell bool) error {
	logger := log.NewLogger("CancelDeal")
	route := fmt.Sprintf(CancelBotDeal, dealID)

	if panicSell {
		route = fmt.Sprintf(PanicSellBotDeal, dealID)
	}

	path := apiConfig.RestURL + route

	query := generateQuery(path, nil)

	logger.Infof("cancelling deal: %s", query.String())
	_, err := makeRequest("POST", query, apiConfig)

	return err
}
