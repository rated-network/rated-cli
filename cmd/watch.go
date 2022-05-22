package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/skillz-blockchain/rated-cli/pkg/rated"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch performances of Ethereum validator keys",
	Run: func(cmd *cobra.Command, args []string) {
		w := rated.NewWatcher()
		err := w.Watch()

		log.Fatal(fmt.Sprintf("rater watcher exited with err=%v", err))
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
