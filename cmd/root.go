package cmd

import (
	"errors"
	"os"
	"strconv"

	logrus_papertrail "github.com/polds/logrus-papertrail-hook"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "list",
	Short: "List commands",
	Long:  `List commands.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gizlinebot.toml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigType("toml")
		viper.AddConfigPath(".")
		viper.SetConfigName(".gizlinebot")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	checkErr(err)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	log.Infof("Using conf file: %s", viper.ConfigFileUsed())

	err = setupPapertrailLogging(cfgStr("PTRAIL_PORT"), cfgStr("PTRAIL_HOST"), cfgStr("PTRAIL_APP"))
	checkErr(err)
}

func setupPapertrailLogging(portStr, hostname, name string) error {
	if portStr == "" || hostname == "" || name == "" {
		log.Info("Skipping papertrail setup (missing PTRAIL_* vars)")
		return nil
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return err
	}
	hook, err := logrus_papertrail.NewPapertrailHook(&logrus_papertrail.Hook{
		Host:     "logs6.papertrailapp.com",
		Port:     port,
		Hostname: hostname,
		Appname:  name,
	})
	if err != nil {
		return err
	}

	hook.SetLevels([]log.Level{log.ErrorLevel, log.WarnLevel, log.InfoLevel, log.DebugLevel})

	log.AddHook(hook)
	return nil
}

func checkErr(err error) {
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func cfgPort() int {
	return cfgInt("PORT")
}
func cfgInt(key string) int {
	val, err := strconv.Atoi(cfgStr(key))
	checkErr(err)
	return val
}

func cfgStr(key string) string {
	if key == "PORT" {
		return viper.GetString(key)
	}
	return viper.GetString("GIZLB_" + key)
}

func validateEnv() {
	reqEnv := []string{"HOSTNAME", "LINE_SECRET", "LINE_TOKEN", "SQL_DB", "SQL_USER", "SQL_PASS", "SQL_HOST", "SQL_PORT"}
	for _, v := range reqEnv {
		if val := cfgStr(v); val == "" {
			checkErr(errors.New("GIZLB_" + v + " is not defined in config file or variable"))
		}
	}

	if val := viper.GetString("PORT"); val == "" {
		checkErr(errors.New("PORT is not defined in config file or env variable"))
	}
}
