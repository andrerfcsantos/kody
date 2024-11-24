package status

import (
	"fmt"
	"github.com/spf13/cobra"
	"kody/lib/config"
	"kody/lib/workshop"
)

var (
	workshopPath string
)

func init() {
	statusCmd.PersistentFlags().StringVarP(&workshopPath, "workshop", "w", ".", "Path to the current workshop")
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Information about the current exercise",
	Long:  `This command gives information about the current exercise based on the current playground.`,
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

		return nil
	},
}

func GetCmd(config *config.Config) *cobra.Command {
	return statusCmd
}
