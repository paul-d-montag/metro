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

	// COMMENT: in hindsight, there should be a request backoff here. If a user
	// waiting is truely an issue you could run the request loop in a goroutine
	// and send the response and error back through a channel. While in the do
	// function you could set up a select that listens on that channel as well
	// as on a time.After() with the amount of time set to a acceptable time to
	// wait for a failure. you could then have a buffered channel that would
	// collect errors as they occured so on timeout you could fetch the first error
	// to append to the timeout error. Otherwise the cause of the timeout would be
	// hard to track down. You would also want the timeout of the http client to be less
	// than the timeout within the select, otherwise you could create a situation
	// where the client returns successfully after the select timeout, and no error would
	// be present in the error buffer to return to the user. That would result in an
	// error like
	// timed out waiting for request: nil
	// which could be hard to debug
	// granted all of this may be over engineering for such a simple client
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
