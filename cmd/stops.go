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
	"strings"

	"github.com/paul-d-montag/metro/mapi"
	"github.com/spf13/cobra"
)

// stopsCmd represents the stops command
var stopsCmd = &cobra.Command{
	Use:   "stops",
	Short: "List all of the stops on a route",
	Long: `List all of the stops on a route and direction.
usage:

metro stops <route> <direction> [stop filter]`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("you must supply a route and direction")
		}

		if len(args) < 2 {
			return errors.New("you must supply a direction")
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

		var filter string
		if len(args) > 3 {
			filter = args[2]
		}

		stops, err := c.FindStops(route.ID, direction.ID, filter)
		if err != nil {
			log.Fatalf("Failed to list routes: %s", err)
		}

		for x, l := 0, len(stops); x < l; x++ {
			fmt.Printf("%s: %s\n", stops[x].ID, stops[x].Name)
		}
	},
}

func getDirection(c *mapi.Client, routeId string, substr string) (*mapi.Direction, error) {
	directions, err := c.FindDirections(routeId, substr)
	if err != nil {
		log.Fatalf("could not get directions for route: %s", err)
	}

	if len(directions) > 1 {
		var s []string
		for _, v := range directions {
			s = append(s, v.Name)
		}
		return nil, fmt.Errorf("too many matching directions\n%s\n", strings.Join(s, ", "))
	}

	if len(directions) < 1 {
		return nil, fmt.Errorf("no directions match substring %s", substr)
	}

	return &directions[0], nil
}

func getRoute(c *mapi.Client, substr string) (*mapi.Route, error) {
	routes, err := c.FindRoutes(substr)
	if err != nil {
		return nil, fmt.Errorf("failed getting routes: %s", substr)
	}

	if len(routes) > 1 {
		var s []string
		for _, v := range routes {
			s = append(s, v.Description)
		}
		return nil, fmt.Errorf("too many matching routes\n%s\n", strings.Join(s, "\n"))
	}

	if len(routes) < 1 {
		return nil, fmt.Errorf("no routes match substring %s", substr)
	}

	return &routes[0], nil
}

func init() {
	rootCmd.AddCommand(stopsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stopsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
