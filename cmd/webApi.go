package cmd

import (
	"github.com/VagabondDataNinjas/gizlinebot/http"
	"github.com/VagabondDataNinjas/gizlinebot/storage"
	"github.com/spf13/cobra"
)

// webApi
var webApiCmd = &cobra.Command{
	Use:   "webApi",
	Short: "Start the webApi service",
	Long:  `webApi service`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		validateEnv()

		s, err := storage.NewSql(cfgStr("SQL_USER") + ":" + cfgStr("SQL_PASS") + "@(" + cfgStr("SQL_HOST") + ":" + cfgStr("SQL_PORT") + ")/" + cfgStr("SQL_DB"))
		checkErr(err)

		// cfgStr("PORT")
		api := http.NewApi(cfgPort(), s)
		checkErr(api.Serve())
	},
}

func init() {
	RootCmd.AddCommand(webApiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// webApiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// webApiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
