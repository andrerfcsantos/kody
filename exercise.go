package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Section struct {
	Number int
	Slug   string
}

type Exercise struct {
	Number  int
	Slug    string
	Section Section
}

func (e *Exercise) BreadCrumbs() string {
	return fmt.Sprintf("[%0.2d] %s > [%0.2d] %s", e.Section.Number, e.Section.Slug, e.Number, e.Slug)
}

func (e *Exercise) BreadCrumbsWithWorkshop(workshop string) string {
	return fmt.Sprintf("%s > %s", workshop, e.BreadCrumbs())
}

func (e *Exercise) FolderName() string {
	return fmt.Sprintf("%0.2d.%s", e.Number, e.Slug)
}

func (e *Exercise) SectionFolderName() string {
	return fmt.Sprintf("%0.2d.%s", e.Section.Number, e.Section.Slug)
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

	return GetMD5Hash(readmeData), nil
}

func GetMD5Hash(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

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
