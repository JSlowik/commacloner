package rest

import (
	"fmt"
	"github.com/jslowik/commacloner/api"
	"github.com/jslowik/commacloner/config"
	"github.com/jslowik/commacloner/log"
	"io/ioutil"
	"net/http"
)

const (
	marketCodeParameter = "market_code"
)

const (
	GetMarketPairs = "/ver1/accounts/market_pairs"
	GetExchanges   = "/ver1/accounts"
)

//func generateQuery(path string, queryParameters map[string]string) *url.URL {
//	u, _ := url.Parse(path)
//	q, _ := url.ParseQuery(u.RawQuery)
//
//	for key, element := range queryParameters {
//		q.Add(key, element)
//	}
//	u.RawQuery = q.Encode()
//
//	return u
//}

func GetExchangeAccounts(apiConfig config.API) error {
	logger := log.NewLogger("GetExchanges")
	route := GetExchanges

	path := apiConfig.RestURL + route

	query := generateQuery(path, nil)

	// Generate Signature
	sig := api.ComputeSignature(fmt.Sprintf("%s?%s", query.Path, query.RawQuery), apiConfig.Secret)

	req, err := http.NewRequest("GET", query.String(), nil)
	if err != nil {
		return fmt.Errorf("could not get user accounts: %v", err)
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
		break
	case http.StatusUnprocessableEntity:
		logger.Warnf("cannot cancel deal: %s", string(responseBody))
	default:
		return fmt.Errorf("bad status %d - %s", resp.StatusCode, string(responseBody))
	}
	return nil
}

func GetExchangePairs(apiConfig config.API, marketCode string) error {
	logger := log.NewLogger("GetExchangePairs")
	route := GetMarketPairs

	path := apiConfig.RestURL + route

	params := make(map[string]string)
	params[marketCodeParameter] = marketCode

	query := generateQuery(path, params)

	// Generate Signature
	sig := api.ComputeSignature(fmt.Sprintf("%s?%s", query.Path, query.RawQuery), apiConfig.Secret)

	req, err := http.NewRequest("GET", query.String(), nil)
	if err != nil {
		return fmt.Errorf("could not get user accounts: %v", err)
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
		break
	case http.StatusUnprocessableEntity:
		logger.Warnf("cannot cancel deal: %s", string(responseBody))
	default:
		return fmt.Errorf("bad status %d - %s", resp.StatusCode, string(responseBody))
	}
	return nil
}
