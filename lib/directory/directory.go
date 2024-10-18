package directory

import (
	"os"
	"path/filepath"
)

func Exists(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	return f.IsDir()
}

func IsGitRepo(path string) bool {
	_, err := os.Stat(filepath.Join(path, ".git"))
	if err != nil {
		return false
	}

	return true
}
