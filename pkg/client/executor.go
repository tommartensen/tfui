package client

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/tommartensen/tfui/pkg/config"
)

func RequestWithoutBody(method string, url string) (*http.Request, error) {
	return http.NewRequest(method, url, nil)
}

func RequestWithBody(method string, url string, body []byte) (*http.Request, error) {
	return http.NewRequest(method, url, bytes.NewReader(body))
}

func Request(method string, uri string, body []byte) (*http.Response, error) {
	configuration := config.New()
	url := fmt.Sprintf("%s/api/%s", configuration.Addr, uri)
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
	if len(configuration.ClientToken) > 0 {
		b64EncodedClientToken := base64.StdEncoding.EncodeToString([]byte(configuration.ClientToken))
		bearer := fmt.Sprintf("Bearer %s", b64EncodedClientToken)
		req.Header.Add("Authorization", bearer)
	}
	return http.DefaultClient.Do(req)
}
