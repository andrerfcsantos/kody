package save

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"kody/lib/cmder"
	"kody/lib/directory"
	"kody/lib/workshop"
	"os"
	"path/filepath"
)

var (
	workshopPath string
	outputDir    string
	shouldCommit bool
)

func init() {
	saveCmd.PersistentFlags().StringVarP(&workshopPath, "workshop", "w", ".", "Path to the current workshop")
	saveCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", ".", "Path to the output directory")
	saveCmd.PersistentFlags().BoolVarP(&shouldCommit, "commit", "c", false, "After adding the exercise to the output directory, commit the changes. This requires the output directory to be a git repository.")
}

func checkSaveFlags() error {
	if workshopPath == "" {
		return errors.New("please provide a path to the workshop folder using the --workshop flag")
	}

	if outputDir == "" {
		return errors.New("please provide a path to the output directory using the --output flag")
	}

	return nil
}

var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save current playground to a more permanent location",
	Long:  `This command allows to save the current contents of a playground to a more permanent location.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := checkSaveFlags(); err != nil {
			return fmt.Errorf("flag error: %w", err)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		isWorkshop := workshop.IsWorkshopFolder(workshopPath)
		if !isWorkshop {
			return fmt.Errorf("'%s' does not look like an Epic React Dev workshop folder", workshopPath)
		}

		hasPlayground := workshop.HasPlayground(workshopPath)
		if !hasPlayground {
			return fmt.Errorf("'%s' does not have a playground folder", workshopPath)
		}

		playgroundPath := filepath.Join(workshopPath, "playground")
		playgroundHash, err := workshop.HashFromPath(playgroundPath)
		if err != nil {
			return fmt.Errorf("error getting hash for '%s': %w", playgroundPath, err)
		}

		exercise, err := workshop.LookupExerciseFromHash(workshopPath, playgroundHash)
		if err != nil {
			return fmt.Errorf("error looking up exercise from playground hash %s: %w", playgroundHash, err)
		}

		if exercise == nil {
			return fmt.Errorf("no exercise found for playground hash %s", playgroundHash)
		}
		workshopSlug := filepath.Base(workshopPath)

		fmt.Printf("Looks like you are doing exercise %s\n", exercise.BreadCrumbsWithWorkshop(workshopSlug))

		exerciseDir := filepath.Join(outputDir, workshopSlug, exercise.SectionFolderName(), exercise.FolderName())
		err = workshop.CopyExercise(playgroundPath, exerciseDir)
		if err != nil {
			return fmt.Errorf("error copying exercise %s > %s: %w", playgroundPath, outputDir, err)
		}

		if shouldCommit {
			if !directory.IsGitRepo(outputDir) {
				return fmt.Errorf("output directory '%s' is not a git repository", outputDir)
			}

			dir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("getting current working directory: %w", err)
			}

			err = os.Chdir(outputDir)
			if err != nil {
				return fmt.Errorf("changing directory to '%s': %w", outputDir, err)
			}

			output, err := cmder.ExecuteCommand("git", "add", "-A", exerciseDir)
			if err != nil {
				fmt.Print(output)
				return fmt.Errorf("adding exercise to git repository: %w", err)
			}

			fmt.Println(output)

			output, err = cmder.ExecuteCommand("git", "commit", "-m", fmt.Sprintf("Add exercise %s", exercise.Descriptor()))
			if err != nil {
				return fmt.Errorf("committing exercise to git repository: %w", err)
			}

			fmt.Println(output)

			err = os.Chdir(dir)
			if err != nil {
				return fmt.Errorf("changing directory back to '%s': %w", dir, err)
			}
		}

		fmt.Printf("Copied exercise from playground '%s' > '%s'\n", playgroundPath, exerciseDir)
		return nil
	},
}

func GetCmd() *cobra.Command {
	return saveCmd
}
