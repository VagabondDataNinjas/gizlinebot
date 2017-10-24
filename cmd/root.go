package cmd

import (
	"errors"
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
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
		fmt.Println(err)
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
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".gizlinebot" (without extension).

		viper.SetConfigType("toml")
		viper.AddConfigPath(home)
		viper.SetConfigName(".gizlinebot")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	checkErr(err)
	fmt.Println("Using config file:", viper.ConfigFileUsed())
}

func checkErr(err error) {
	if err != nil {
		fmt.Printf("\nError: %s\n", err)
		os.Exit(1)
	}
}

func cfgStr(key string) string {
	if key == "PORT" {
		return viper.GetString(key)
	}
	return viper.GetString("GIZLB_" + key)
}

func validateEnv() {
	reqEnv := []string{"LINE_SECRET", "LINE_TOKEN", "SQL_DB", "SQL_USER", "SQL_PASS", "SQL_HOST", "SQL_PORT"}
	for _, v := range reqEnv {
		if val := cfgStr(v); val == "" {
			checkErr(errors.New("GIZLB_" + v + " is not defined in config file or variable"))
		}
	}

	if val := viper.GetString("PORT"); val == "" {
		checkErr(errors.New("PORT is not defined in config file or env variable"))
	}
}
