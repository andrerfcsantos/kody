package version

import (
	"fmt"
	"github.com/andrerfcsantos/kody/lib/config"

	"github.com/spf13/cobra"
)

var (
	cfg *config.Config
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Long:  `Display the version information for kody. Use -v flag for verbose output including commit and date.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		buildInfo := cfg.GetBuildInfo()

		if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
			fmt.Printf("Version: %s\n", buildInfo.Version)
			fmt.Printf("Commit:  %s\n", buildInfo.Commit)
			fmt.Printf("Date:    %s\n", buildInfo.Date)
		} else {
			fmt.Println(buildInfo.Version)
		}

		return nil
	},
}

func GetCmd(configuration *config.Config) *cobra.Command {
	cfg = configuration

	// Add verbose flag
	versionCmd.Flags().BoolP("verbose", "v", false, "Show detailed version information")

	return versionCmd
}
