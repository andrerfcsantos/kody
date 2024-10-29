package workshop

import (
	"fmt"
	"kody/lib/directory"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

type Workshop struct {
	Path   string
	config *PackageConfig
}

func (w *Workshop) Slug() string {
	return w.config.Epicshop.Product.Slug
}

func (w *Workshop) Title() string {
	return w.config.Epicshop.Title
}

func (w *Workshop) AsciiTitle() string {
	cleanString := strings.Map(func(r rune) rune {
		if utf8.RuneLen(r) > 1 {
			return -1
		}
		return r
	}, w.Title())

	return strings.TrimSpace(cleanString)
}

func (w *Workshop) PlaygroundPath() string {
	return filepath.Join(w.Path, "playground")
}

func (w *Workshop) HasPlayground() bool {
	return directory.Exists(w.PlaygroundPath())
}

func (w *Workshop) PlaygroundHash() (string, error) {
	if !w.HasPlayground() {
		return "", fmt.Errorf("workshop '%s' does not have a playground folder\n", w.Path)
	}
	return HashFromPath(filepath.Join(w.Path, "playground"))
}

func (w *Workshop) LookupExerciseFromHash(targetHash string) (*Exercise, error) {
	exercisePaths, err := filepath.Glob(filepath.Join(w.Path, "exercises", "*", "*.problem.*"))
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

func (w *Workshop) PlaygroundExercise() (*Exercise, error) {
	playgroundHash, err := w.PlaygroundHash()
	if err != nil {
		return nil, fmt.Errorf("getting hash for '%s': %w", w.Path, err)
	}

	exercise, err := w.LookupExerciseFromHash(playgroundHash)
	if err != nil {
		return nil, fmt.Errorf("looking up exercise from playground hash '%s': %w", playgroundHash, err)
	}

	if exercise == nil {
		return nil, fmt.Errorf("no exercise found for playground hash '%s'", playgroundHash)
	}

	return exercise, nil
}

func WorkshopFromPath(workshopPath string) (*Workshop, error) {

	if !isWorkshopFolder(workshopPath) {
		return nil, fmt.Errorf("'%s' does not look like an Epic React Dev workshop folder\n", workshopPath)
	}

	config, err := LoadPackageConfig(workshopPath)
	if err != nil {
		return nil, fmt.Errorf("loading package config: %w\n", err)
	}

	return &Workshop{
		Path:   workshopPath,
		config: config,
	}, nil
}

func isWorkshopFolder(workshopPath string) bool {
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
