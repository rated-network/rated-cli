package cmd

import (
	"math/rand"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/skillz-blockchain/rated-cli/pkg/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfg core.Config

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
	rand.Seed(time.Now().UTC().UnixNano())
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

	viper.SetEnvPrefix("rated")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.WithError(err).Fatal("unable to read configuration file")
	}

	cfg.ApiEndpoint = viper.GetString("rated.apiEndpoint")

	if viper.InConfig("rated.network") {
		cfg.Network = viper.GetString("rated.network")
	} else {
		cfg.Network = "mainnet"
	}

	cfg.WatcherValidationKeys = viper.GetStringSlice("rated.watcher.validationKeys")
	cfg.WatcherRefreshRate = time.Second * time.Duration(viper.GetInt64("rated.watcher.refreshRateInSeconds"))
	cfg.ListenOn = viper.GetString("rated.listenOn")

	log.WithFields(log.Fields{
		"config": viper.ConfigFileUsed(),
	}).Info("successfully read configuration file")
}
