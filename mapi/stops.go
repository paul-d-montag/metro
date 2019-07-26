package mapi

import (
	"fmt"
	"strings"
)

type Stop struct {
	ID   string `json:"Value"`
	Name string `json:"Text"`
}

func (c *Client) Stops(routeId, directionId string) ([]Stop, error) {
	var stops []Stop
	err := c.get("nextrip/stops/%s/%s?format=json", &stops, routeId, directionId)
	return stops, err
}

func (c *Client) FindStops(routeId, directionId, substr string) ([]Stop, error) {
	substr = strings.ToLower(substr)
	stops, err := c.Stops(routeId, directionId)
	if err != nil {
		return nil, fmt.Errorf("could not get stops %s", err)
	}
	var matches []Stop

	for x, l := 0, len(stops); x < l; x++ {
		if strings.Contains(strings.ToLower(fmt.Sprintf("%s %s", stops[x].Name, stops[x].ID)), substr) {
			matches = append(matches, stops[x])
		}
	}

	return matches, nil
}
