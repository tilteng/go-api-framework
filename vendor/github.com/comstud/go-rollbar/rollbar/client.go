package rollbar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	net_url "net/url"
)

// client struct
type client struct {
	httpClient      *http.Client
	apiBaseURL      string
	accessToken     string
	notifierName    string
	notifierVersion string

	ClientOptions
}

func (self *client) decodeBody(http_resp *http.Response, resp interface{}) error {
	if http_resp.StatusCode >= 200 && http_resp.StatusCode < 300 {
		return json.NewDecoder(http_resp.Body).Decode(resp)
	}

	// Hrmph.
	data, err := ioutil.ReadAll(http_resp.Body)
	if err == nil {
		return fmt.Errorf("Got code %d: %s\n", http_resp.StatusCode, string(data))
	} else {
		return fmt.Errorf("Got code %d: <Error reading body: %s>\n", http_resp.StatusCode, err)
	}
}

func (self *client) httpCall(method string, url string, query net_url.Values, data interface{}, resp interface{}) error {
	if query == nil {
		query = make(net_url.Values)
	}

	query.Set("access_token", self.accessToken)

	if url != "" && url[0] != '/' {
		url = self.apiBaseURL + "/" + url
	} else {
		url = self.apiBaseURL + url
	}

	query_str := query.Encode()
	if query_str != "" {
		url += "?" + query_str
	}

	var json_data []byte

	if data != nil {
		var err error

		json_data, err = json.Marshal(data)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(json_data))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	http_resp, err := self.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer http_resp.Body.Close()

	return self.decodeBody(http_resp, resp)

}

func (self *client) httpGet(url string, query net_url.Values, resp interface{}) error {
	return self.httpCall("GET", url, query, nil, resp)
}

func (self *client) httpPatch(url string, data interface{}, resp interface{}) error {
	return self.httpCall("PATCH", url, nil, data, resp)
}

func (self *client) httpPost(url string, data interface{}, resp interface{}) error {
	return self.httpCall("POST", url, nil, data, resp)
}

// Get the client options
func (self *client) Options() *ClientOptions {
	return &self.ClientOptions
}

// Get the base API URL
func (self *client) APIBaseURL() string {
	return self.apiBaseURL
}

// Set the base API URL
func (self *client) SetAPIBaseURL(base_url string) Client {
	self.apiBaseURL = base_url
	return self
}
