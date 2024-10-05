package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var (
	workshopPath string
	outputDir    string
)

func init() {
	flag.StringVar(&workshopPath, "workshop", ".", "Path to the current workshop")
	flag.StringVar(&outputDir, "output", ".", "Path to the output directory")
}

func checkFlags() error {
	if workshopPath == "" {
		return errors.New("please provide a path to the workshop folder using the --workshop flag")
	}

	if outputDir == "" {
		return errors.New("please provide a path to the output directory using the --output flag")
	}

	return nil
}

func main() {
	flag.Parse()

	if err := checkFlags(); err != nil {
		fmt.Printf("Flag error: %s\n", err)
		os.Exit(2)
	}

	isWorkshop := IsWorkshopFolder(workshopPath)
	if !isWorkshop {
		fmt.Printf("'%s' does not look like an Epic React Dev workshop folder\n", workshopPath)
		os.Exit(1)
	}

	hasPlayground := HasPlayground(workshopPath)
	if !hasPlayground {
		fmt.Printf("'%s' does not have a playground folder\n", workshopPath)
		os.Exit(1)
	}

	playgroundPath := filepath.Join(workshopPath, "playground")
	playgroundHash, err := HashFromPath(playgroundPath)
	if err != nil {
		fmt.Printf("Error getting hash for '%s': %s\n", playgroundPath, err)
		os.Exit(2)
	}

	exercise, err := LookupExerciseFromHash(workshopPath, playgroundHash)
	if err != nil {
		fmt.Printf("Error looking up exercise from playground hash %s: %s\n", playgroundHash, err)
		os.Exit(2)
	}

	if exercise == nil {
		fmt.Printf("No exercise found for playground hash %s\n", playgroundHash)
		os.Exit(1)
	}
	workshopSlug := filepath.Base(workshopPath)

	fmt.Printf("Looks like you are doing exercise %s\n", exercise.BreadCrumbsWithWorkshop(workshopSlug))

	exerciseDir := filepath.Join(outputDir, workshopSlug, exercise.SectionFolderName(), exercise.FolderName())
	err = CopyExercise(playgroundPath, exerciseDir)
	if err != nil {
		fmt.Printf("Error copying exercise %s > %s: %s\n", playgroundPath, outputDir, err)
		os.Exit(2)
	}

	fmt.Printf("Copied exercise from playground '%s' > '%s'\n", playgroundPath, exerciseDir)

}
