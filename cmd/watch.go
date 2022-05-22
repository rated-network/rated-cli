package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var keywatchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch performances of Ethereum validator keys",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("keywatch called")
	},
}

func init() {
	rootCmd.AddCommand(keywatchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// keywatchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// keywatchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
