package restore

import (
	"errors"
	"fmt"
	"kody/lib/config"
	"kody/lib/directory"
	"kody/lib/workshop"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	cfg *config.Config
)

var (
	workshopPath    string
	workshopsDir    string
	currentWorkshop *workshop.Workshop
	outputDir       string
	sectionNo       int
	exerciseNo      int
)

func checkAndSetupConfigs(cmd *cobra.Command) error {
	workshopPath = cfg.GetString("workshop.path")
	workshopsDir = cfg.GetString("workshops.dir")
	outputDir = cfg.GetString("save.output.directory")

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
		return errors.New("please provide a path to the workshop folder using the --workshop flag or the workshop.path configuration, or use --workshops-dir to auto-detect")
	}

	if outputDir == "" {
		return errors.New("please provide a path to the output directory using the --output flag or the save.output.directory configuration")
	}

	return nil
}

var restoreCmd = &cobra.Command{
	Use:    "restore [exercise]",
	Hidden: true,
	Short:  "Restore an exercise to the playground",
	Long:   `Restore an exercise to the playground. If no exercise is specified, automatically detects the current exercise from the playground.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := checkAndSetupConfigs(cmd); err != nil {
			return fmt.Errorf("flag error: %w", err)
		}
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		// If no exercise is provided, we'll auto-detect from the playground
		if len(args) == 0 {
			return nil
		}

		// If an exercise is provided, parse it
		if len(args) != 1 {
			return errors.New("restore accepts at most one argument with the exercise to restore, in the format <section_number>.<exercise_number>, e.g. \"01.02\"")
		}

		exerciseStr := args[0]
		exerciseSplit := strings.Split(exerciseStr, ".")
		if len(exerciseSplit) != 2 {
			return errors.New("exercise argument must be in the format <section_number>.<exercise_number>, e.g. \"01.02\"")
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

		// If no exercise was specified, auto-detect from the playground
		if len(args) == 0 {
			playgroundExercise, err := w.PlaygroundExercise()
			if err != nil {
				return fmt.Errorf("auto-detecting exercise from playground: %w", err)
			}

			// Use the section and exercise numbers from the detected exercise
			sectionNo = playgroundExercise.Section.Number
			exerciseNo = playgroundExercise.Number

			fmt.Printf("Auto-detected exercise: %s > %s\n", playgroundExercise.BreadCrumbsWithWorkshop(w.Slug()), playgroundExercise.Descriptor())
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

	cfg.BindFlagConfigToCommand("workshop.dir", restoreCmd)
	cfg.BindFlagConfigToCommand("workshops.dir", restoreCmd)
	cfg.BindFlagConfigToCommand("save.output.directory", restoreCmd)

	return restoreCmd
}
