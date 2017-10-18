package cmd

import (
	"errors"

	"github.com/VagabondDataNinjas/gizlinebot/line"

	"github.com/VagabondDataNinjas/gizlinebot/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// lineBotCmd represents the lineBot command
var lineBotCmd = &cobra.Command{
	Use:   "lineBot",
	Short: "Start the linebot server",
	Long:  `Linebot server`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		validateEnv()

		s, err := storage.NewSql(cfgStr("SQL_USER") + ":" + cfgStr("SQL_PASS") + "@(" + cfgStr("SQL_HOST") + ":" + cfgStr("SQL_PORT") + ")/" + cfgStr("SQL_DB"))
		checkErr(err)

		port := cfgStr("SERVER_PORT")
		server, err := line.NewLineServer(port, s, viper.GetString("GIZLB_LINE_SECRET"), viper.GetString("GIZLB_LINE_TOKEN"))
		checkErr(err)

		err = server.Serve()
		checkErr(err)
	},
}

func cfgStr(key string) string {
	return viper.GetString("GIZLB_" + key)
}

func validateEnv() {
	reqEnv := []string{"LINE_SECRET", "LINE_TOKEN", "SERVER_PORT", "SQL_DB", "SQL_USER", "SQL_PASS", "SQL_HOST", "SQL_PORT"}
	for _, v := range reqEnv {
		if val := cfgStr(v); val == "" {
			checkErr(errors.New("GIZLB_" + v + " is not defined in config file"))
		}
	}
}

func init() {
	RootCmd.AddCommand(lineBotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lineBotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lineBotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
