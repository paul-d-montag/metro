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
	"fmt"
	"log"
	"strings"

	"github.com/paul-d-montag/metro/mapi"
	"github.com/spf13/cobra"
)

var direction bool

// routesCmd represents the routes command
var routesCmd = &cobra.Command{
	Use:   "routes",
	Short: "list all the routes",
	Long: `list all the routes that exist in the Metro Transit API
usage:

	metro routes <search term>

example:

	metro routes north
`,
	Run: func(cmd *cobra.Command, args []string) {
		c := getMapi()
		var filter string
		if len(args) > 0 {
			filter = args[0]
		}

		routes, err := c.FindRoutes(filter)
		if err != nil {
			log.Fatalf("Failed to list routes: %s", err)
		}

		for x, l := 0, len(routes); x < l; x++ {
			fmt.Print(routes[x].Description)
			if direction {
				d, err := listDirections(c, routes[x].ID)

				if err != nil {
					log.Printf("WARNING: %s", err)
					d = "directions unavailable"
				}

				fmt.Printf(" %s", d)
			}
			fmt.Print("\n")
		}
	},
}

func init() {
	routesCmd.Flags().BoolVar(&direction, "direction", false, "show the directions this route can take")
	rootCmd.AddCommand(routesCmd)
}

func listDirections(c *mapi.Client, routeID string) (string, error) {
	directions, err := c.Directions(routeID)

	if err != nil {
		return "", fmt.Errorf("listing route directions: %s", err)
	}

	d := make([]string, len(directions))
	for x, l := 0, len(directions); x < l; x++ {
		d[x] = directions[x].Name
	}

	return strings.Join(d, ", "), nil
}
