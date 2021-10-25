package rest

import (
	"encoding/json"
	"github.com/jslowik/commacloner/api/threecommas/rest/dobjs"
	"github.com/jslowik/commacloner/config"
)

const (
	marketCodeParameter = "market_code"
)

const (
	GetMarketPairs = "/ver1/accounts/market_pairs"
	GetExchanges   = "/ver1/accounts"
)

//GetExchangeAccounts gets the exchanges associated with an API key
func GetExchangeAccounts(apiConfig config.API) (map[int]dobjs.Exchange, error) {
	route := GetExchanges

	path := apiConfig.RestURL + route

	query := generateQuery(path, nil)

	resp, err := makeRequest("GET", query, apiConfig,nil, nil)
	if err != nil {
		return nil, err
	}

	var accounts []dobjs.Exchange
	err = json.Unmarshal(resp, &accounts)
	if err != nil {
		return nil, err
	}

	exchangeMap := make(map[int]dobjs.Exchange)
	for _, account := range accounts {
		exchangeMap[account.ID] = account
	}
	return exchangeMap, nil
}

//GetExchangePairs gets the pairs that are allowed on a given exchange
func GetExchangePairs(apiConfig config.API, marketCode string) ([]string, error) {
	route := GetMarketPairs

	path := apiConfig.RestURL + route

	params := make(map[string]string)
	params[marketCodeParameter] = marketCode

	query := generateQuery(path, params)
	resp, err := makeRequest("GET", query, apiConfig,nil, nil)
	if err != nil {
		return nil, err
	}

	var pairs []string
	err = json.Unmarshal(resp, &pairs)
	if err != nil {
		return nil, err
	}
	return pairs, nil

}
