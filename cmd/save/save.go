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

		w, err := workshop.WorkshopFromPath(workshopPath)
		if err != nil {
			return fmt.Errorf("getting workshop from path '%s': %w", workshopPath, err)
		}

		exercise, err := w.PlaygroundExercise()
		if err != nil {
			return fmt.Errorf("getting playground exercise: %w", err)
		}

		fmt.Printf("Looks like you are doing exercise %s\n", exercise.BreadCrumbsWithWorkshop(w.Slug()))

		exerciseDir := filepath.Join(outputDir, w.Slug(), exercise.SectionFolderName(), exercise.FolderName())
		err = workshop.CopyExercise(w.PlaygroundPath(), exerciseDir)
		if err != nil {
			return fmt.Errorf("error copying exercise %s > %s: %w", w.PlaygroundPath(), outputDir, err)
		}

		if shouldCommit {
			err = HandleCommit(outputDir, exerciseDir, fmt.Sprintf("Add exercise %s", exercise.Descriptor()))
			if err != nil {
				return fmt.Errorf("committing exercise '%s': %w", exerciseDir, err)
			}

		}

		fmt.Printf("Copied exercise from playground '%s' > '%s'\n", w.PlaygroundPath(), exerciseDir)
		return nil
	},
}

func HandleCommit(repoPath string, exercisePath string, message string) error {
	if !directory.IsGitRepo(repoPath) {
		return fmt.Errorf("output directory '%s' is not a git repository", repoPath)
	}

	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current working directory: %w", err)
	}

	err = os.Chdir(repoPath)
	if err != nil {
		return fmt.Errorf("changing directory to '%s': %w", repoPath, err)
	}

	output, err := cmder.ExecuteCommand("git", "add", "-A", exercisePath)
	if err != nil {
		fmt.Print(output)
		return fmt.Errorf("adding exercise to git repository: %w", err)
	}

	fmt.Println(output)

	output, err = cmder.ExecuteCommand("git", "commit", "-m", message)
	if err != nil {
		return fmt.Errorf("committing exercise to git repository: %w", err)
	}

	fmt.Println(output)

	err = os.Chdir(dir)
	if err != nil {
		return fmt.Errorf("changing directory back to '%s': %w", dir, err)
	}

	return nil
}

func GetCmd() *cobra.Command {
	return saveCmd
}
