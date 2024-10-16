package workshop

import (
	"fmt"
	"os"
)

func CopyExercise(playgroundPath string, outputDir string) error {

	err := os.RemoveAll(outputDir)
	if err != nil {
		return fmt.Errorf("removing '%s' folder: %w\n", outputDir, err)
	}

	err = os.CopyFS(outputDir, os.DirFS(playgroundPath))
	if err != nil {
		return fmt.Errorf("copying '%s' folder to '%s': %w\n", playgroundPath, outputDir, err)
	}

	return nil
}
