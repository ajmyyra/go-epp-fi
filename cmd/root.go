package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var configFile string
var configVars = []string{"CLIENT_KEY", "CLIENT_CERT", "USERNAME", "PASSWORD", "SERVER"}
var envPrefix = "FI_EPP"

var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "Command line client for FI EPP asset management (domains, contacts etc)",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	configFile = ".fi-epp.yml"
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFile)

	viper.AutomaticEnv()
	viper.SetEnvPrefix(envPrefix)

	_ = viper.ReadInConfig()

	for _, envVar := range configVars {
		_ = viper.BindEnv(envVar)
	}
}