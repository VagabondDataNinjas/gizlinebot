package cmd

import (
	"github.com/VagabondDataNinjas/gizlinebot/domain"
	"github.com/VagabondDataNinjas/gizlinebot/http"
	"github.com/VagabondDataNinjas/gizlinebot/storage"
	"github.com/VagabondDataNinjas/gizlinebot/survey"
	"github.com/line/line-bot-sdk-go/linebot"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// startCmd represents the lineBot command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start your engines!",
	Long:  `Start the linebot server, the api and serve static files`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		validateEnv()

		s, err := storage.NewSql(cfgStr("SQL_USER") + ":" + cfgStr("SQL_PASS") + "@(" + cfgStr("SQL_HOST") + ":" + cfgStr("SQL_PORT") + ")/" + cfgStr("SQL_DB"))
		checkErr(err)

		bot, err := linebot.New(viper.GetString("GIZLB_LINE_SECRET"), viper.GetString("GIZLB_LINE_TOKEN"))
		checkErr(err)

		qs, err := s.GetQuestions()
		checkErr(err)
		surv := survey.NewSurvey(s, qs)

		globalVars := &domain.GlobalTplVars{
			Hostname: cfgStr("HOSTNAME"),
		}

		initiatorDelay := int64(cfgInt("INITIATOR_DELAY_SEC"))
		errc := make(chan error)
		initiator := survey.NewInitiator(surv, s, bot, globalVars)
		go func() {
			initiator.Monitor(initiatorDelay, errc)
		}()

		apiConf := &http.ApiConf{
			Port:       cfgPort(),
			GlobalVars: globalVars,
		}
		api := http.NewApi(s, bot, surv, apiConf)
		log.Info("Starting the API...")
		checkErr(api.Serve())
	},
}

func init() {
	RootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
