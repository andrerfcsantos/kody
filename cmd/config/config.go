package config

import (
	"fmt"
	"github.com/spf13/cobra"
	"kody/lib/config"
)

var (
	cfg *config.Config
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage kody configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			getAllConfig()
			return nil
		}

		if len(args) == 1 {
			value := cfg.Get(args[0])
			if value != nil {
				fmt.Printf("%s: %v\n", args[0], value)
				return nil
			}
		}

		if len(args) == 2 {
			cfg.Set(args[0], args[1])
			err := cfg.Write()
			if err != nil {
				return fmt.Errorf("writing config: %w", err)
			}
		}

		return nil
	},
}

func getAllConfig() {
	keys := cfg.AllKeys()

	if len(keys) == 0 {
		fmt.Println("No configurations defined yet.")
		return
	}

	for _, key := range keys {
		value := cfg.Get(key)
		if value != nil {
			fmt.Printf("%s: %v\n", key, value)
		}
	}
}

func GetCmd(config *config.Config) *cobra.Command {
	cfg = config

	configCmd.PersistentFlags().StringP("workshop", "w", ".", "Path to the current workshop")
	cfg.BindPFlag("workshop.path", configCmd.PersistentFlags().Lookup("workshop"))

	return configCmd
}
