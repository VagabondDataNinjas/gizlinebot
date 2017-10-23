package cmd

import (
	"errors"

	"github.com/VagabondDataNinjas/gizlinebot/domain"
	"github.com/VagabondDataNinjas/gizlinebot/line"
	"github.com/VagabondDataNinjas/gizlinebot/storage"
	"github.com/VagabondDataNinjas/gizlinebot/survey"
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

		// @TODO move outside code
		qs := domain.NewQuestions()
		questions := [][]string{
			{"Q1", `Thank you for following us!
If you'd like to complete the survey online please go to https://google.com
Otherwise you can complete the form here.
You can start by replying back with your location (area or island name):`},
			{"Q2", "How much do you pay for diesel in your area?"},
			{"Q3", "What is your occupation?"},
			{"Q4", "What is your line id?"},
			{"Q5", "What did you do today?"},
			{"Q6", "Thank you for all your help! We might ask you more questions in the future"},
		}

		for _, question := range questions {
			err = qs.Add(question[0], question[1])
			checkErr(err)
		}
		surv := survey.NewSurvey(s, qs)

		port := cfgStr("SERVER_PORT")
		server, err := line.NewLineServer(port, surv, s, viper.GetString("GIZLB_LINE_SECRET"), viper.GetString("GIZLB_LINE_TOKEN"))
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
