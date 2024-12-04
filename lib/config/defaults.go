package config

import "path/filepath"

func DefaultSaveDir(cfg *Config) string {
	dataDir, err := cfg.DataDir()
	if err != nil {
		dataDir = "."
	}

	return filepath.Join(dataDir, "save")

}
