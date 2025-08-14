package config

import (
	"fmt"
	"github.com/andrerfcsantos/kody/lib/config"
	"sort"

	"github.com/spf13/cobra"
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
	sort.Strings(keys)

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

	paths, err := cfg.ConfigPaths()
	if err != nil {
		fmt.Printf("\n(Error getting config paths: %v)\n", err)
	}
	fmt.Printf("\n(config paths: %v)\n", paths)
}

func GetCmd(config *config.Config) *cobra.Command {
	cfg = config

	cfg.BindFlagConfigToCommand("workshop.dir", configCmd)
	cfg.BindFlagConfigToCommand("workshops.dir", configCmd)

	return configCmd
}
