package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func DirectoryExists(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	return f.IsDir()
}

func HasPlayground(workshopFolder string) bool {
	return DirectoryExists(filepath.Join(workshopFolder, "playground"))
}

func WorkshopName(workshopFolder string) string {
	return filepath.Base(workshopFolder)
}

func ExerciseFromPath(exercisePath string) (*Exercise, error) {
	parts := strings.Split(exercisePath, string(filepath.Separator))
	if len(parts) < 2 {
		return nil, fmt.Errorf("exercise path '%s' does not seem to contain [...]/<section>/<exercise>", exercisePath)
	}
	parts = parts[len(parts)-2:]

	sectionParts := strings.Split(parts[0], ".")
	if len(sectionParts) < 2 {
		return nil, fmt.Errorf("section path '%s' does not seem to be in the format [...]/<section>.<number>/<exercise>", parts[0])
	}

	sectionNumber, err := strconv.Atoi(sectionParts[0])
	if err != nil {
		return nil, fmt.Errorf("section number '%s' is not a number", sectionParts[1])
	}

	section := Section{
		Number: sectionNumber,
		Slug:   sectionParts[1],
	}

	exerciseParts := strings.Split(parts[1], ".")
	if len(exerciseParts) < 3 {
		return nil, fmt.Errorf("exercise path '%s' does not seem to be in the format [...]/<section>/<exercise>.problem.<number>", exercisePath)
	}

	exerciseNumber, err := strconv.Atoi(exerciseParts[0])
	if err != nil {
		return nil, fmt.Errorf("exercise number '%s' is not a number", exerciseParts[1])
	}

	exercise := Exercise{
		Number:  exerciseNumber,
		Slug:    exerciseParts[2],
		Section: section,
	}

	return &exercise, nil
}

func LookupExerciseFromHash(workshopFolder string, targetHash string) (*Exercise, error) {
	exercisePaths, err := filepath.Glob(filepath.Join(workshopPath, "exercises", "*", "*.problem.*"))
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

func IsWorkshopFolder(workshopPath string) bool {
	if dirExists := DirectoryExists(workshopPath); !dirExists {
		return false
	}

	if dirExists := DirectoryExists(filepath.Join(workshopPath, "exercises")); !dirExists {
		return false
	}

	if dirExists := DirectoryExists(filepath.Join(workshopPath, "epicshop")); !dirExists {
		return false
	}

	return true
}
