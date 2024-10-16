package workshop

import (
	"fmt"
	"kody/lib/directory"
	"path/filepath"
)

func IsWorkshopFolder(workshopPath string) bool {
	if dirExists := directory.Exists(workshopPath); !dirExists {
		return false
	}

	if dirExists := directory.Exists(filepath.Join(workshopPath, "exercises")); !dirExists {
		return false
	}

	if dirExists := directory.Exists(filepath.Join(workshopPath, "epicshop")); !dirExists {
		return false
	}

	return true
}

func WorkshopName(workshopFolder string) string {
	return filepath.Base(workshopFolder)
}

func HasPlayground(workshopFolder string) bool {
	return directory.Exists(filepath.Join(workshopFolder, "playground"))
}

func LookupExerciseFromHash(workshopFolder string, targetHash string) (*Exercise, error) {
	exercisePaths, err := filepath.Glob(filepath.Join(workshopFolder, "exercises", "*", "*.problem.*"))
	if err != nil {
		return nil, fmt.Errorf("getting exercise paths: %w\n", err)
	}
	for _, path := range exercisePaths {
		hash, err := HashFromPath(path)
		if err != nil {
			return nil, fmt.Errorf("getting hash for exercise at '%s': %w\n", path, err)
		}

		if hash == targetHash {
			exercise, err := ExerciseFromPath(path)
			if err != nil {
				return nil, fmt.Errorf("getting exercise from path '%s': %w\n", path, err)
			}
			return exercise, nil
		}
	}

	return nil, nil
}
