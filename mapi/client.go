package mapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

// Create an interface for the Http client for easier testing
type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	// The http client to use
	C HttpClient
	// Endpoint is the endpoint the client uses when making requests
	Endpoint string

	i sync.Once
}

func (c *Client) init() {
	c.i.Do(func() {
		if c.C == nil {
			c.C = http.DefaultClient
		}

		if c.Endpoint == "" {
			c.Endpoint = "http://svc.metrotransit.org"
		}

		// strip any trailing slashes
		if c.Endpoint[len(c.Endpoint)-1] == '/' {
			c.Endpoint = c.Endpoint[0 : len(c.Endpoint)-1]
		}
	})
}

func (c *Client) url(uri string, v ...interface{}) string {
	if len(v) > 0 {
		uri = fmt.Sprintf(uri, v...)
	}

	return fmt.Sprintf("%s/%s", c.Endpoint, uri)
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	c.init()

	// if this were a server application instead of a terminal client
	// I would add some exponential backoff code here that tried
	// to run the request a set number of times before giving up.
	// however, because this is a terminal client it makes more sense to
	// fail quickly to let the user know something is wrong, instead of
	// making them wait for multiple requests to fail.
	res, err := c.C.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= http.StatusMultipleChoices {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		res.Body.Close()

		return nil, fmt.Errorf("[%s] :: %s", res.Status, b)
	}

	return res, nil
}

// get will reach out to the endpoint at the uri and unmarshal the json into what is passed into v
func (c *Client) get(uri string, v interface{}, vs ...interface{}) error {
	c.init()

	req, err := http.NewRequest(http.MethodGet, c.url(uri, vs...), nil)
	if err != nil {
		return fmt.Errorf("request %s::%s could not be created %s", http.MethodGet, uri, err)
	}

	res, err := c.do(req)
	if err != nil {
		return fmt.Errorf("request %s::%s failed: %s", http.MethodGet, uri, err)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading response %s::%s failed: %s", http.MethodGet, uri, err)
	}

	if err := json.Unmarshal(b, v); err != nil {
		return fmt.Errorf("could not unmarshal response: %s\n\n%s", err, b)
	}

	return nil
}
