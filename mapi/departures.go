package mapi

import (
	"fmt"
	"strings"
)

type Departure struct {
	Actual        bool   `json:"Actual"`
	BlockNumber   int    `json:"BlockNumber"`
	DepartureText string `json:"DepartureText"`
	// Comment: The date that was passed was in a format which I couldn't find an existing library for.
	// although implementing a lexor using the bufio package and parser could be kind of fun.
	DepartureTime    string  `json:"DepartureTime"`
	Description      string  `json:"Description"`
	Gate             string  `json:"Gate"`
	Route            string  `json:"Route"`
	RouteDirection   string  `json:"RouteDirection"`
	Terminal         string  `json:"Terminal"`
	VehicleHeading   int     `json:"VehicleHeading"`
	VehicleLatitude  float32 `json:"VehicleLatitude"`
	VehicleLongitude float32 `json:"VehicleLongitude"`
}

// COMMENT: If given more time I would have implemented a caching mechanism here of some
// kind. Possibly tiedot, so that request data could have been stored in a database locally
// then allowing the user the searching capabilities of the library. Tiedot used to consume a
// massive amount of memory but I implemented a PR to fix the issue. All the hard work of
// figuring out how to do it was done by others, but I implemented it https://github.com/HouzuoGuo/tiedot/pull/157
func (c *Client) Departures(routeId, directionId, stopId string) ([]Departure, error) {
	var departures []Departure
	err := c.get("nextrip/%s/%s/%s?format=json", &departures, routeId, directionId, stopId)
	return departures, err
}

func (c *Client) FindDepartures(routeId, directionId, stopId, substr string) ([]Departure, error) {
	substr = strings.ToLower(substr)
	departures, err := c.Departures(routeId, directionId, stopId)
	if err != nil {
		return nil, fmt.Errorf("could not get departures %s", err)
	}
	var matches []Departure

	for x, l := 0, len(departures); x < l; x++ {
		if strings.Contains(strings.ToLower(departures[x].Description), substr) {
			matches = append(matches, departures[x])
		}
	}

	return matches, nil
}
