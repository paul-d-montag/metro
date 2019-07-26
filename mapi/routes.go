package mapi

import (
	"fmt"
	"strings"
)

type Route struct {
	Description string `json:"Description"`
	ProviderID  string `json:"ProviderID"`
	ID          string `json:"Route"`
}

func (c *Client) Routes() ([]Route, error) {
	var routes []Route
	err := c.get("nextrip/routes?format=json", &routes)
	return routes, err
}

func (c *Client) FindRoutes(substr string) ([]Route, error) {
	substr = strings.ToLower(substr)
	routes, err := c.Routes()
	if err != nil {
		return nil, fmt.Errorf("could not get routes %s", err)
	}
	var matches []Route

	for x, l := 0, len(routes); x < l; x++ {
		if strings.Contains(strings.ToLower(routes[x].Description), substr) {
			matches = append(matches, routes[x])
		}
	}

	return matches, nil
}
