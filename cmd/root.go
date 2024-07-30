package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.0.1"
var rootCmd = &cobra.Command{
	Use:     "Guess-Who Crawler",
	Version: version,
	Short:   "Guess-Who Crawler - CLI to scrap information regarding born in Brazil",
	Long: `need to write something better here, explain that this uses LLM, 
	intended to create the card game for the Guess-Who etc`,
	Run: func(cmd *cobra.Command, args []string) {},
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
