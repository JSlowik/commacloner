package rest

import (
	"encoding/json"
	"fmt"
	"github.com/jslowik/commacloner/api/threecommas/dobjs"
	"github.com/jslowik/commacloner/config"
	"github.com/jslowik/commacloner/log"
	"net/url"
	"strconv"
	"strings"
)

const (
	pairParameter                   = "pair"
	pairsParameter                  = "pairs"
	takeProfitParameter             = "take_profit"
	nameParameter                   = "name"
	baseOrderVolumeParameter        = "base_order_volume"
	safetyOrderVolumeParameter      = "safety_order_volume"
	martingaleVolumeParameter       = "martingale_volume_coefficient"
	martingaleStepParameter         = "martingale_step_coefficient"
	maxSafetyOrderParameter         = "max_safety_orders"
	activeSafetyOrderCountParameter = "active_safety_orders_count"
	safetyOrderStepParameter        = "safety_order_step_percentage"
	takeProfitTypeParameter         = "take_profit_type"
	strategyListParameter           = "strategy_list"
)

const (
	ShowBot          = "/ver1/bots/%d/show"
	StartNewBotDeal  = "/ver1/bots/%d/start_new_deal"
	CancelBotDeal    = "/ver1/deals/%d/cancel"
	PanicSellBotDeal = "/ver1/deals/%d/panic_sell"
	UpdateBotPairs   = "/ver1/bots/%d/update"
)

//GetBot returns the JSON description of a bot based on its bot id.
func GetBot(apiConfig config.API, id int) (dobjs.Bot, error) {
	logger := log.NewLogger("bots")
	route := fmt.Sprintf(ShowBot, id)
	path := apiConfig.RestURL + route

	query := generateQuery(path, nil)

	logger.Infof("getting bot info: %s", query.String())
	info, err := makeRequest("GET", query, apiConfig, nil, nil)
	if err != nil {
		return dobjs.Bot{}, err
	}

	var bot dobjs.Bot
	err = json.Unmarshal(info, &bot)
	return bot, err
}

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
	_, err := makeRequest("POST", query, apiConfig, nil, nil)

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
	_, err := makeRequest("POST", query, apiConfig, nil, nil)

	return err
}

func UpdatePairs(apiConfig config.API, bot dobjs.Bot, quoteCurrency string, pairs []string) error {
	logger := log.NewLogger("UpdatePairs")
	route := fmt.Sprintf(UpdateBotPairs, bot.ID)

	newPairs := make([]string, 0)
	for _, pair := range pairs {
		newPairs = append(newPairs, fmt.Sprintf("%s_%s", quoteCurrency, pair))
	}

	params := make(map[string]string)

	pString, err := json.Marshal(newPairs)
	if err != nil {
		return err
	}

	data := url.Values{}
	data.Set(nameParameter, bot.Name)
	data.Set(pairsParameter, string(pString))
	data.Set(baseOrderVolumeParameter, bot.BaseOrderVolume)
	data.Set(takeProfitParameter, bot.TakeProfit)
	data.Set(safetyOrderVolumeParameter, bot.SafetyOrderVolume)
	data.Set(martingaleVolumeParameter, bot.MartingaleVolumeCoefficient)
	data.Set(martingaleStepParameter, bot.MartingaleStepCoefficient)
	data.Set(maxSafetyOrderParameter, strconv.Itoa(bot.MaxSafetyOrders))
	data.Set(activeSafetyOrderCountParameter, strconv.Itoa(bot.ActiveSafetyOrdersCount))
	data.Set(safetyOrderStepParameter, bot.SafetyOrderStepPercentage)
	data.Set(takeProfitTypeParameter, bot.TakeProfitType)
	strategy, _ := json.Marshal(bot.StrategyList)
	data.Set(strategyListParameter, string(strategy))

	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	path := apiConfig.RestURL + route
	query := generateQuery(path, params)

	logger.Infof("updating bot pairs: %s", query.String())
	_, err = makeRequest("PATCH", query, apiConfig, data, headers)

	return err
}
