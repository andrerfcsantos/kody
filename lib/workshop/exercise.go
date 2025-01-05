package workshop

import (
	"fmt"
	"io"
	"kody/lib/hash"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Exercise struct {
	Number  int
	Slug    string
	Section Section
	path    string
}

func (e *Exercise) BreadCrumbs() string {
	return fmt.Sprintf("[%0.2d] %s > [%0.2d] %s", e.Section.Number, e.Section.Slug, e.Number, e.Slug)
}

func (e *Exercise) BreadCrumbsWithWorkshop(workshop string) string {
	return fmt.Sprintf("%s > %s", workshop, e.BreadCrumbs())
}

func (e *Exercise) Descriptor() string {
	return fmt.Sprintf("%0.2d-%s", e.Number, e.Slug)
}

func (e *Exercise) FolderName() string {
	return fmt.Sprintf("%0.2d.%s", e.Number, e.Slug)
}

func (e *Exercise) SectionFolderName() string {
	return fmt.Sprintf("%0.2d.%s", e.Section.Number, e.Section.Slug)
}

func (e *Exercise) Hash() (string, error) {
	return HashFromPath(e.path)
}

func (e *Exercise) Path() string {
	return e.path
}

func HashFromPath(exerciseDir string) (s string, err error) {
	readmePath := filepath.Join(exerciseDir, "README.mdx")

	readmeFile, err := os.Open(readmePath)
	defer func(readmeFile *os.File) {
		err = readmeFile.Close()
	}(readmeFile)

	if err != nil {
		return "", fmt.Errorf("opening '%s' file: %w\n", readmePath, err)
	}
	readmeData, err := io.ReadAll(readmeFile)
	if err != nil {
		return "", fmt.Errorf("reading '%s' file: %w\n", readmePath, err)
	}

	return hash.MD5Hex(readmeData), nil
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
		path:    exercisePath,
	}

	return &exercise, nil
}

func DefaultExerciseDir(outputDir string, w *Workshop, exercise *Exercise) string {
	return filepath.Join(outputDir, w.Slug(), exercise.SectionFolderName(), exercise.FolderName())
}
