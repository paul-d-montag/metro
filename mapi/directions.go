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

// COMMENT: All of these Find functions have a specific smell to them, but because
// go doesn't allow interfaces to be applied to slice items without making a new
// slice I have yet to find a way to make this readable and generic (obviously an issue in go :) )
// There is likely a solution utilizing both closures and interfaces, and if you could
// pass a slice of []Direction to a function that wanted a slice of []interface{} it would be
// trivial, but you cant. Many have argued that this implementation would be harder to read
// but I would argue the extra code to instatiate a new slice with a new type
// jumbles up the code more and is logic for type handling sake and not to solve
// the actuall problem at hand. So these end up being more copy paste and search and
// replace functions to maintain readability. If given more time I might find a
// proper solution
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
