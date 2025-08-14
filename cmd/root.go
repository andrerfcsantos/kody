package cmd

import (
	"fmt"
	configCmd "github.com/andrerfcsantos/kody/cmd/config"
	"github.com/andrerfcsantos/kody/cmd/restore"
	"github.com/andrerfcsantos/kody/cmd/save"
	"github.com/andrerfcsantos/kody/cmd/status"
	"github.com/andrerfcsantos/kody/cmd/test"
	"github.com/andrerfcsantos/kody/cmd/version"
	"github.com/andrerfcsantos/kody/lib/config"
	"os"

	"github.com/spf13/cobra"
)

var cfg *config.Config

func init() {
	cfg = config.NewConfig("kody")

	config.AddFlagConfig(cfg, config.FlagConfig[string]{
		Key:           "workshop.dir",
		FlagName:      "workshop",
		FlagShortHand: "w",
		Default:       "",
		Description:   "Path to the current workshop directory. If this is specificed, Kody will only consider save/restore operations on this workshop. Tipically you want to pass the path to the workshop you are currently working on. For automatic detection of the current workshop, use the --workshops flag instead. [config key: workshop.dir]",
	})

	config.AddFlagConfig(cfg, config.FlagConfig[string]{
		Key:           "workshops.dir",
		FlagName:      "workshops",
		FlagShortHand: "d",
		Default:       "",
		Description:   "Path to the workshops directory, where all the workshops sub-directories are located. If this is provided, the current workshop will be automatically calculated to be the one with the most recent playground modification time. Use the --workshop flag if don't want to use this automatic workshop detection. [config key: workshops.dir]",
	})

	config.AddFlagConfig(cfg, config.FlagConfig[string]{
		Key:           "save.output.directory",
		FlagName:      "output-dir",
		FlagShortHand: "o",
		Default:       config.DefaultSaveDir(cfg),
		Description:   "Directory to save the exercises to. This is usually the same directory the save command uses to save the exercises.",
	})

	config.AddFlagConfig(cfg, config.FlagConfig[bool]{
		Key:           "save.shouldCommit",
		FlagName:      "commit",
		FlagShortHand: "c",
		Default:       false,
		Description:   "After adding the exercise to the output directory, commit the changes. This requires the output directory to be a git repository.",
	})

	config.AddFlagConfig(cfg, config.FlagConfig[string]{
		Key:           "save.commit.message",
		FlagName:      "commitMessage",
		FlagShortHand: "m",
		Default:       "[{{ .Workshop.Slug }}] Add exercise {{ .Exercise.BreadCrumbs }}",
		Description:   "Commit message to use, in case the --commit flag is set or the save.shouldCommit configuration is set to true. The template is rendered using Go's text/template package.",
	})

	rootCmd.AddCommand(save.GetCmd(cfg))
	rootCmd.AddCommand(restore.GetCmd(cfg))
	rootCmd.AddCommand(status.GetCmd(cfg))
	rootCmd.AddCommand(configCmd.GetCmd(cfg))
	rootCmd.AddCommand(test.GetCmd(cfg))
	rootCmd.AddCommand(version.GetCmd(cfg))
}

var rootCmd = &cobra.Command{
	Use:   "kody",
	Short: "CLI tool to help manage Epic React Dev workshops and exercises.",
	Long:  `Management of Epic React Dev workshops and exercises and other automation tasks.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return cfg.Read()
	},
}

func Execute(buildInfo config.BuildInfo) {
	cfg.SetBuildInfo(buildInfo)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
