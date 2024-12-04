package status

import (
	"fmt"
	"github.com/spf13/cobra"
	"kody/lib/config"
	"kody/lib/workshop"
)

var (
	cfg *config.Config
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Information about the current exercise",
	Long:  `This command gives information about the current exercise based on the current playground.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		workshopPath := cfg.GetString("workshop.path")
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
	cfg = config
	statusCmd.PersistentFlags().StringP("workshop", "w", ".", "Path to the current workshop")
	err := cfg.BindPFlag("workshop.path", statusCmd.PersistentFlags().Lookup("workshop"))

	if err != nil {
		fmt.Printf("Error setting up the workshop flag: %s\n", err)
	}
	return statusCmd
}
