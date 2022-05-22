package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "rated-cli",
	Short: "Beacon Chain validator ratings from the command line.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rated-cli.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".rated-cli")
	}

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.WithError(err).Fatal("unable to read configuration file")
	}

	log.WithFields(log.Fields{
		"config": viper.ConfigFileUsed(),
	}).Info("successfully read configuration file")
}
