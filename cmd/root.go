package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"kody/cmd/save"
	"os"
)

func init() {
	rootCmd.AddCommand(save.GetCmd())
}

var rootCmd = &cobra.Command{
	Use:   "kody",
	Short: "CLI tool to help manage Epic React Dev workshops and exercises.",
	Long:  `Management of Epic React Dev workshops and exercises and other automation tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
