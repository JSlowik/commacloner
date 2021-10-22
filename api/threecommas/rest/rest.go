package rest

import (
	"fmt"
	"github.com/jslowik/commacloner/api/threecommas"
	"github.com/jslowik/commacloner/config"
	"io/ioutil"
	"net/http"
	"net/url"
)

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
func makeRequest(method string, query *url.URL, apiConfig config.API) ([]byte, error) {
	// Generate Signature
	sig := threecommas.ComputeSignature(fmt.Sprintf("%s?%s", query.Path, query.RawQuery), apiConfig.Secret)

	req, err := http.NewRequest(method, query.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("APIKEY", apiConfig.Key)
	req.Header.Set("Signature", sig)

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
	case http.StatusUnprocessableEntity:
		e = fmt.Errorf("%d - Unprocessable Entity - %v ", resp.StatusCode, string(responseBody))
	default:
		e = fmt.Errorf("%d - Unexpected Response Code - %v", resp.StatusCode, string(responseBody))
	}
	return responseBody, e
}
