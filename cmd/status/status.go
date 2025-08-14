package status

import (
	"fmt"
	"github.com/andrerfcsantos/kody/lib/config"
	"github.com/andrerfcsantos/kody/lib/workshop"

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

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Information about the current exercise",
	Long:  `This command gives information about the current exercise based on the current playground.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkAndSetupConfigs(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var w *workshop.Workshop
		var err error

		// Use the already loaded workshop if available, otherwise load it from path
		if currentWorkshop != nil {
			w = currentWorkshop
		} else {
			w, err = workshop.WorkshopFromPath(workshopPath)
			if err != nil {
				return fmt.Errorf("getting workshop from path '%s': %w", workshopPath, err)
			}
		}

		exercise, err := w.PlaygroundExercise()
		if err != nil {
			return fmt.Errorf("getting playground exercise: %w", err)
		}

		fmt.Printf("Looks like you are doing exercise %s\n", exercise.BreadCrumbsWithWorkshop(w.Slug()))

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

	if workshopPath == "" {
		return fmt.Errorf("please provide a path to the workshop folder using the --workshop flag or the workshop.path configuration, or use --workshops-dir to auto-detect")
	}

	return nil
}

func GetCmd(config *config.Config) *cobra.Command {
	cfg = config

	cfg.BindFlagConfigToCommand("workshop.dir", statusCmd)
	cfg.BindFlagConfigToCommand("workshops.dir", statusCmd)

	return statusCmd
}
