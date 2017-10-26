package cmd

import (
	"github.com/VagabondDataNinjas/gizlinebot/line"
	"github.com/VagabondDataNinjas/gizlinebot/storage"
	"github.com/VagabondDataNinjas/gizlinebot/survey"
	"github.com/line/line-bot-sdk-go/linebot"
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

		store, err := storage.NewSql(cfgStr("SQL_USER") + ":" + cfgStr("SQL_PASS") + "@(" + cfgStr("SQL_HOST") + ":" + cfgStr("SQL_PORT") + ")/" + cfgStr("SQL_DB"))
		checkErr(err)

		qs, err := store.GetQuestions()
		checkErr(err)
		surv := survey.NewSurvey(store, qs)

		bot, err := linebot.New(viper.GetString("GIZLB_LINE_SECRET"), viper.GetString("GIZLB_LINE_TOKEN"))
		checkErr(err)

		errc := make(chan error)
		initiator := line.NewInitiator(surv, store, bot)
		go func() {
			initiator.Monitor(120, errc)
		}()

		port := cfgStr("PORT")
		server, err := line.NewLineServer(port, surv, store, bot)
		checkErr(err)

		err = server.Serve()
		checkErr(err)
	},
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
