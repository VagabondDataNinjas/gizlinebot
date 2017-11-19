package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/VagabondDataNinjas/gizlinebot/storage"
	"github.com/VagabondDataNinjas/gizlinebot/survey"
	"github.com/spf13/cobra"
)

// normaliseCmd represents the lineBot command
var normaliseCmd = &cobra.Command{
	Use:   "normalise",
	Short: "Normalise data!",
	Long:  `Continuously monitors SQL tables, parses the answer data and normalises it`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		validateEnv()

		s, err := storage.NewSql(cfgStr("SQL_USER") + ":" + cfgStr("SQL_PASS") + "@(" + cfgStr("SQL_HOST") + ":" + cfgStr("SQL_PORT") + ")/" + cfgStr("SQL_DB"))
		checkErr(err)

		normaliser := survey.NewNormaliser(s)
		errc := make(chan error)
		go func() {
			normaliser.Start(errc)
		}()

		for err = range errc {
			log.Errorf("[normaliser] Error: %s", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(normaliseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// normaliseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// normaliseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
