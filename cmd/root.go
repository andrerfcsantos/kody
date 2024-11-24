package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	configCmd "kody/cmd/config"
	"kody/cmd/save"
	"kody/cmd/status"
	"kody/cmd/test"
	"kody/lib/config"
	"os"
)

var cfg *config.Config

func init() {
	cfg = config.NewConfig("kody")

	rootCmd.AddCommand(save.GetCmd(cfg))
	rootCmd.AddCommand(status.GetCmd(cfg))
	rootCmd.AddCommand(configCmd.GetCmd(cfg))
	rootCmd.AddCommand(test.GetCmd(cfg))
}

var rootCmd = &cobra.Command{
	Use:   "kody",
	Short: "CLI tool to help manage Epic React Dev workshops and exercises.",
	Long:  `Management of Epic React Dev workshops and exercises and other automation tasks.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return cfg.Read()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
