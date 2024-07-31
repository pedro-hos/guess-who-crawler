package cmd

import (
	"log"

	"github.com/pedro-hos/guess-who-web/controllers"
	"github.com/pedro-hos/guess-who-web/database"
	"github.com/pedro-hos/guess-who-web/util"
	"github.com/spf13/cobra"
)

var state, city string
var scrapCmd = &cobra.Command{
	Use:     "scrap",
	Aliases: []string{"scrap"},
	Short:   "Scrap information command",
	PreRun: func(cmd *cobra.Command, args []string) {
		if city != "" {
			cmd.MarkFlagRequired("state")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {

		config, errConfig := util.LoadConfig(".")
		if errConfig != nil {
			log.Fatal("cannot load config:", errConfig)
		}

		database.Connect(config, Prod)
		controllers.RunScraper(state, city)
	},
}

func init() {
	scrapCmd.Flags().StringVarP(&state, "state", "s", "", "Use to scrap an specifica brazilian state")
	scrapCmd.Flags().StringVarP(&city, "city", "c", "", "Use to scrap an specifica brazilian city. The state is mandatory when using City!")

	rootCmd.AddCommand(scrapCmd)
}
