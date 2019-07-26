/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/paul-d-montag/metro/mapi"
	"github.com/spf13/cobra"
)

var (
	departureTemplate        string
	next                     bool
	defaultDepartureTemplate = `{{ .Description }}: {{ .DepartureText }}{{ if .Actual }}
  warning time not based on actual location{{ end }}
`
)

// departuresCmd represents the departures command
var departuresCmd = &cobra.Command{
	Use:   "departures",
	Short: "List all departures for a route and direction",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("you must supply a route, direction and stop")
		}

		if len(args) < 2 {
			return errors.New("you must supply a direction and stop")
		}

		if len(args) < 3 {
			return errors.New("you must supply a stop")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		c := getMapi()
		route, err := getRoute(c, args[0])
		if err != nil {
			log.Fatal(err)
		}

		direction, err := getDirection(c, route.ID, args[1])
		if err != nil {
			log.Fatal(err)
		}

		stop, err := getStop(c, route.ID, direction.ID, args[2])
		if err != nil {
			log.Fatal(err)
		}

		var filter string
		if len(args) > 4 {
			filter = args[3]
		}

		departures, err := c.FindDepartures(route.ID, direction.ID, stop.ID, filter)
		if err != nil {
			log.Fatalf("Failed to list routes: %s", err)
		}

		for x, l := 0, len(departures); x < l; x++ {
			t, err := template.New("").Parse(departureTemplate)
			if err != nil {
				log.Fatalf("template could not be parsed: %s", err)
			}

			if err := t.Execute(os.Stdout, departures[x]); err != nil {
				log.Fatalf("could not execute template: %s", err)
			}
			// only show the first departure because it will be the newest
			if next {
				break
			}
		}
	},
}

func getStop(c *mapi.Client, routeId, directionId, substr string) (*mapi.Stop, error) {
	stops, err := c.FindStops(routeId, directionId, substr)
	if err != nil {
		log.Fatalf("could not get stops for route: %s", err)
	}

	if len(stops) > 1 {
		var s []string
		for _, v := range stops {
			s = append(s, fmt.Sprintf("%s: %s", v.ID, v.Name))
		}

		return nil, fmt.Errorf("too many matching stops\n%s\n", strings.Join(s, "\n"))
	}

	if len(stops) < 1 {
		return nil, fmt.Errorf("no stops match substring %s", substr)
	}

	return &stops[0], nil
}

// COMMENT: If the api fed me more data it would have been a good idea to implement the templating at a global
// level along with the template flag. I often find myself making bash scripts around my software so any time you
// can specify exactly the data you want instead of having to awk you way to it is a win.
func init() {
	rootCmd.AddCommand(departuresCmd)
	departuresCmd.Flags().BoolVar(&next, "next", false, "Show only the next departure")
	departuresCmd.Flags().StringVar(&departureTemplate, "template", defaultDepartureTemplate, "A golang template to show data instead of the default one")
}
