package mapi

import (
	"fmt"
	"strings"
)

type Direction struct {
	ID   string `json:"Value"`
	Name string `json:"Text"`
}

func (c *Client) Directions(routeId string) ([]Direction, error) {
	var directions []Direction
	err := c.get("nextrip/directions/%s?format=json", &directions, routeId)
	return directions, err
}

func (c *Client) FindDirections(routeId, substr string) ([]Direction, error) {
	substr = strings.ToLower(substr)
	directions, err := c.Directions(routeId)
	if err != nil {
		return nil, fmt.Errorf("could not get directions %s", err)
	}
	var matches []Direction

	for x, l := 0, len(directions); x < l; x++ {
		if strings.Contains(strings.ToLower(directions[x].Name), substr) {
			matches = append(matches, directions[x])
		}
	}

	return matches, nil
}
