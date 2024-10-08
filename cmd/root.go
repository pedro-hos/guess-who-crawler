package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.0.1"
var Prod bool = false
var rootCmd = &cobra.Command{
	Use:     "Guess-Who Crawler",
	Version: version,
	Short:   "Guess-Who Crawler: A CLI tool to scrape information about famous people born in Brazil from Wikipedia.",
	Long: `This CLI scrapes data from Wikipedia about famous Brazilians to create cards for the Guess-Who 
	game. You can specify city and state 
	parameters to get data for specific locations, or leave them blank to retrieve data for all Brazilian 
	states and cities.`,
	Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&Prod, "start-prod", false, "Start the production mode. (default: false)")
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
