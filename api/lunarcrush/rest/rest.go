package rest

import (
	"encoding/json"
	"fmt"
	"github.com/jslowik/commacloner/api/lunarcrush/rest/dobjs"
	"io/ioutil"
	"net/http"
	"net/url"
)

//func sortByAltRank(pairs [] dobjs.PairData) [] dobjs.PairData {
//	sort.Slice(pairs, func(i, j int) bool {
//		return pairs[i].Acr < pairs[j].Acr
//	})
//	return pairs
//}
//
//func sortByGalaxyScore(pairs [] dobjs.PairData) [] dobjs.PairData {
//	sort.Slice(pairs, func(i, j int) bool {
//		return pairs[i].Gs > pairs[j].Gs
//	})
//	return pairs
//}

//func GetNPairs(apiKey string, sort string, validPairs map[string][]string, max int) ([]string, error) {
//	allPairs, err := GetPairs(apiKey,sort)
//	if err != nil {
//		return nil, err
//	}
//	sorted := sortByGalaxyScore(allPairs)
//
//	pairs := make([]string,0)
//	for _, pair := range sorted {
//		if validPairs[pair.S] != nil {
//			pairs = append(pairs,pair.S)
//		}
//		if len(pairs) == max {
//			break
//		}
//	}
//	return pairs, nil
//}

func GetPairs(apiKey string, sort string)([]dobjs.PairData, error) {
	route := "/v2?data=market&key=ytoqgk0erazk0wkjfzr3h&type=fast"

	path := "https://api.lunarcrush.com" + route

	q := generateQuery(path,nil)
	res, e := makeRequest("GET",q)
	if e != nil {
		return nil, e
	}

	var respData dobjs.Response
	e = json.Unmarshal(res, &respData)
	if e != nil {
		return nil, e
	}
	return respData.Data, nil

}

//generateQuery generates a query with the given map of query parameters
func generateQuery(path string, queryParameters map[string]string) *url.URL {
	u, _ := url.Parse(path)
	q, _ := url.ParseQuery(u.RawQuery)

	for key, element := range queryParameters {
		q.Add(key, element)
	}
	u.RawQuery = q.Encode()

	return u
}

//makeRequest makes and signs an http request, and returns the response
func makeRequest(method string, query *url.URL) ([]byte, error) {
	req, err := http.NewRequest(method, query.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var e error
	switch resp.StatusCode {
	case http.StatusCreated, http.StatusOK:
		e = nil
	//case http.StatusUnprocessableEntity:
	//	e = fmt.Errorf("%d - Unprocessable Entity - %v ", resp.StatusCode, string(responseBody))
	default:
		e = fmt.Errorf("%d - Unexpected Response Code - %v", resp.StatusCode, string(responseBody))
	}
	return responseBody, e
}
