package workshop

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type PackageConfig struct {
	Name     string `json:"name"`
	Epicshop struct {
		Title      string `json:"title"`
		GithubRoot string `json:"githubRoot"`
		GithubRepo string `json:"githubRepo"`
		Subtitle   string `json:"subtitle"`
		Product    struct {
			Host             string   `json:"host"`
			Slug             string   `json:"slug"`
			DisplayName      string   `json:"displayName"`
			DisplayNameShort string   `json:"displayNameShort"`
			Logo             string   `json:"logo"`
			DiscordChannelId string   `json:"discordChannelId"`
			DiscordTags      []string `json:"discordTags"`
		} `json:"product"`
		OnboardingVideo string `json:"onboardingVideo"`
	} `json:"epicshop"`
	Type       string            `json:"type"`
	Scripts    map[string]string `json:"scripts"`
	Keywords   []string          `json:"keywords"`
	Author     string            `json:"author"`
	License    string            `json:"license"`
	Workspaces []string          `json:"workspaces"`
}

func LoadPackageConfig(workshopPath string) (*PackageConfig, error) {

	configPath := filepath.Join(workshopPath, "package.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading package config: %w\n", err)
	}

	config := PackageConfig{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling package config: %w\n", err)
	}

	return &config, nil
}
