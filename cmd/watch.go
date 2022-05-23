package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/skillz-blockchain/rated-cli/pkg/core"
	"github.com/skillz-blockchain/rated-cli/pkg/watcher"
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch performances of Ethereum validator keys",
	Run: func(cmd *cobra.Command, args []string) {
		go core.ListenAndServe(&cfg)

		w, err := watcher.NewWatcher(&cfg)
		if err != nil {
			log.WithError(err).Fatal("unable to initialize watcher")
		}

		err = w.Watch()

		log.WithError(err).Fatal("unable to watch")
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
