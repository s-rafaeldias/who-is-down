package pkg

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Client interface {
	getEndpointData(url *url.URL) (map[string]interface{}, error)
}

type HTTPClient struct{}

func (h *HTTPClient) getEndpointData(url *url.URL) (map[string]interface{}, error) {
	// Make HTTP request
	resp, err := http.Get(url.String())
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// get response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// parse response body
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return data, nil
}
