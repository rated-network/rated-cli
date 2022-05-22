package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/skillz-blockchain/rated-cli/pkg/rated"
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch performances of Ethereum validator keys",
	Run: func(cmd *cobra.Command, args []string) {
		w := rated.NewWatcher("", []string{}, 42)
		err := w.Watch()

		log.WithError(err).Fatal(fmt.Sprintf("unable to watch"))
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
