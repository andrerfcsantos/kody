package restore

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"kody/lib/config"
	"kody/lib/directory"
	"kody/lib/workshop"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	cfg *config.Config
)

var (
	workshopPath string
	outputDir    string
	sectionNo    int
	exerciseNo   int
)

func checkAndSetupConfigs() error {
	workshopPath = cfg.GetString("workshop.path")
	outputDir = cfg.GetString("save.output.directory")

	if workshopPath == "" {
		return errors.New("please provide a path to the workshop folder using the --workshop flag or the workshop.path configuration")
	}

	if outputDir == "" {
		return errors.New("please provide a path to the output directory using the --output flag or the save.output.directory configuration")
	}

	return nil
}

var restoreCmd = &cobra.Command{
	Use:    "restore",
	Hidden: true,
	Short:  "Testing command",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := checkAndSetupConfigs(); err != nil {
			return fmt.Errorf("flag error: %w", err)
		}
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("restore requires an argument with the exercise to restore, in the format <section_number>.<exercise_number>, e.g. \"01.02\"")
		}

		exerciseStr := args[0]
		exerciseSplit := strings.Split(exerciseStr, ".")
		if len(exerciseSplit) != 2 {
			return errors.New("restore requires an argument with the exercise to restore, in the format <section_number>.<exercise_number>, e.g. \"01.02\"")
		}
		var err error
		sectionNo, err = strconv.Atoi(exerciseSplit[0])
		if err != nil {
			return fmt.Errorf("parsing section number: %w", err)
		}

		exerciseNo, err = strconv.Atoi(exerciseSplit[1])
		if err != nil {
			return fmt.Errorf("parsing exercise number: %w", err)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		w, err := workshop.WorkshopFromPath(workshopPath)
		if err != nil {
			return fmt.Errorf("getting workshop from path '%s': %w", workshopPath, err)
		}

		sectionGlob := fmt.Sprintf("%02d.*", sectionNo)
		exerciseGlob := fmt.Sprintf("%02d.*", exerciseNo)
		searchGlob := filepath.Join(outputDir, w.Slug(), sectionGlob, exerciseGlob)

		matches, err := filepath.Glob(searchGlob)
		if err != nil {
			return fmt.Errorf("globbing files: %w", err)
		}

		if len(matches) == 0 {
			return fmt.Errorf("no candidate directory in the output directory found for the exercise was found (search glob: %s)", searchGlob)
		}

		if len(matches) != 1 {
			return fmt.Errorf("more than one candidate directory in the output directory found for the exercise was found")
		}

		restorePath := matches[0]
		err = directory.CopyFS(w.PlaygroundPath(), os.DirFS(restorePath))
		if err != nil {
			return fmt.Errorf("restoring files: %w", err)
		}
		return nil
	},
}

func GetCmd(configuration *config.Config) *cobra.Command {
	cfg = configuration

	restoreCmd.PersistentFlags().StringP("workshop", "w", ".", "Path to the current workshop")
	cfg.BindPFlag("workshop.path", restoreCmd.PersistentFlags().Lookup("workshop"))

	restoreCmd.PersistentFlags().StringP("source", "s", config.DefaultSaveDir(cfg), "Source directory from where to get the exercises. This is usually the same directory the save command uses to save the exercises.")
	cfg.BindPFlag("save.output.directory", restoreCmd.PersistentFlags().Lookup("source"))

	return restoreCmd
}
