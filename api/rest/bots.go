package rest

import (
	"fmt"
	"github.com/jslowik/commacloner/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/jslowik/commacloner/api"
	"github.com/jslowik/commacloner/config"
)

const (
	pairParameter = "pair"
	// skipSignalChecks    = "skip_signal_checks"
	// skipOpenDealsChecks = "skip_open_deals_checks"
	// botID               = "bot_id"
)

const (
	StartNewBotDeal  = "/ver1/bots/%d/start_new_deal"
	CancelBotDeal    = "/ver1/deals/%d/cancel"
	PanicSellBotDeal = "/ver1/deals/%d/panic_sell"
)

func generateQuery(path string, queryParameters map[string]string) *url.URL {
	u, _ := url.Parse(path)
	q, _ := url.ParseQuery(u.RawQuery)

	for key, element := range queryParameters {
		q.Add(key, element)
	}
	u.RawQuery = q.Encode()

	return u
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

	// Generate Signature
	sig := api.ComputeSignature(fmt.Sprintf("%s?%s", query.Path, query.RawQuery), apiConfig.Secret)

	req, err := http.NewRequest("POST", query.String(), nil)
	if err != nil {
		return fmt.Errorf("could not generate new deal request: %v", err)
	}

	req.Header.Set("APIKEY", apiConfig.Key)
	req.Header.Set("Signature", sig)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not send new deal request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusCreated:
		break
	case 422:
		return fmt.Errorf("cannot create new deal: %s", string(responseBody))
	default:
		return fmt.Errorf("bad status %d - %s", resp.StatusCode, string(responseBody))
	}
	return nil
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

	// Generate Signature
	sig := api.ComputeSignature(fmt.Sprintf("%s?%s", query.Path, query.RawQuery), apiConfig.Secret)

	req, err := http.NewRequest("POST", query.String(), nil)
	if err != nil {
		return fmt.Errorf("could not generate new deal request: %v", err)
	}

	req.Header.Set("APIKEY", apiConfig.Key)
	req.Header.Set("Signature", sig)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not send new deal request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		break
	case http.StatusUnprocessableEntity:
		logger.Warnf("cannot cancel deal: %s", string(responseBody))
	default:
		return fmt.Errorf("bad status %d - %s", resp.StatusCode, string(responseBody))
	}
	return nil
}
