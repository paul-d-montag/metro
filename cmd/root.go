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
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/paul-d-montag/metro/mapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "metro",
	Short: "a command line tool for interacting with the Metro API",
	Long: `Interact with the Metro Transit api so you can:

list Routes:
	metro routes
`,
	Run: func(cmd *cobra.Command, args []string) {
		routesCmd.Run(cmd, args)
	},
}

// COMMENT: The overall usage of this command line client is clumsy while being
// able to slowly discover all the things needed to render the final command, I
// found myself constantly going back to the begining of the line and editing
// the subcommand. In truth the command should flow, where as you offer it more data,
// the context of what it shows you changes. The issue with doing this currently is the
// filter argument at the end of each subcommand. This should be changed to a flag instead
// so the flow can function. It could be made as a persistent flag since each subcommand uses it
// I felt it was good to leave out a directions sub command as the directions are so closely
// married to a route. The only reason it wasn't a default action is because it does demand a
// new request to the api for each item

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.metro.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".metro" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName("metro")
	}

	viper.SetDefault("endpoint", "http://svc.metrotransit.org/")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func getMapi() *mapi.Client {
	return &mapi.Client{
		Endpoint: viper.GetString("endpoint"),
	}

}
