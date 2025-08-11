package test

import (
	"fmt"
	"kody/lib/config"
	"kody/lib/workshop"

	"github.com/spf13/cobra"
)

var (
	cfg *config.Config
)

var (
	workshopPath    string
	workshopsDir    string
	currentWorkshop *workshop.Workshop
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Testing command",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkAndSetupConfigs(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("workshop path: %v\n", workshopPath)
		fmt.Printf("workshops dir: %v\n", workshopsDir)
		return nil
	},
}

func checkAndSetupConfigs(cmd *cobra.Command) error {
	workshopPath = cfg.GetString("workshop.path")
	workshopsDir = cfg.GetString("workshops.dir")

	// Check if flags were passed directly
	if workshopPathFlag := cmd.Flags().Lookup("workshop"); workshopPathFlag != nil && workshopPathFlag.Changed {
		workshopPath = workshopPathFlag.Value.String()
	}
	if workshopsDirFlag := cmd.Flags().Lookup("workshops-dir"); workshopsDirFlag != nil && workshopsDirFlag.Changed {
		workshopsDir = workshopsDirFlag.Value.String()
	}

	// If workshopPath is not provided but workshopsDir is, auto-detect the current workshop
	if workshopPath == "" && workshopsDir != "" {
		var err error
		currentWorkshop, err = workshop.DetectCurrentWorkshop(workshopsDir)
		if err != nil {
			return fmt.Errorf("auto-detecting workshop from workshopsDir '%s': %w", workshopsDir, err)
		}
		workshopPath = currentWorkshop.Path
	}

	return nil
}

func GetCmd(config *config.Config) *cobra.Command {
	cfg = config

	testCmd.PersistentFlags().StringP("workshop", "w", "", "Path to the current workshop")
	err := cfg.BindPFlag("workshop.path", testCmd.PersistentFlags().Lookup("workshop"))

	if err != nil {
		fmt.Printf("Error setting up the workshop flag: %s\n", err)
	}

	testCmd.PersistentFlags().StringP("workshops-dir", "p", "", "Directory containing workshop sub-directories for auto-detection")
	err = cfg.BindPFlag("workshops.dir", testCmd.PersistentFlags().Lookup("workshops-dir"))

	if err != nil {
		fmt.Printf("Error setting up the workshop-paths flag: %s\n", err)
	}

	return testCmd
}
