package test

import (
	"fmt"
	"github.com/spf13/cobra"
	"kody/lib/config"
)

var (
	cfg *config.Config
)

func init() {
	testCmd.PersistentFlags().StringP("workshop", "w", ".", "Path to the current workshop")
	cfg.BindPFlag("workshop.path", testCmd.PersistentFlags().Lookup("workshop"))
}

var testCmd = &cobra.Command{
	Use:    "test",
	Hidden: true,
	Short:  "Testing command",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("workshop path: %v\n", cfg.Get("workshop.path"))
		return nil
	},
}

func GetCmd(config *config.Config) *cobra.Command {
	cfg = config
	return testCmd
}
