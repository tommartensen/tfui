package client

import (
	"bytes"
	"fmt"
	"net/http"
)

// TODO: use config for address
// TODO: use application token from config
var terraformUIAddr = "http://localhost:8080"

func RequestWithoutBody(method string, url string) (*http.Request, error) {
	return http.NewRequest(method, url, nil)
}

func RequestWithBody(method string, url string, body []byte) (*http.Request, error) {
	return http.NewRequest(method, url, bytes.NewReader(body))
}

func Request(method string, uri string, body []byte) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/%s", terraformUIAddr, uri)
	var req *http.Request
	var err error

	if body == nil {
		req, err = RequestWithoutBody(method, url)
	} else {
		req, err = RequestWithBody(method, url, body)
	}
	if err != nil {
		return &http.Response{}, err
	}
	return http.DefaultClient.Do(req)

}
