package config

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	gap "github.com/muesli/go-app-paths"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configName = "config"
	configType = "yml"
	appName    = "kody"
)

type Config struct {
	AppName       string
	gapScope      *gap.Scope
	viper         *viper.Viper
	configFlagMap map[string]interface{}
}

func NewConfig(appName string) *Config {
	viperCfg := viper.NewWithOptions(
		viper.EnvKeyReplacer(strings.NewReplacer("_", ".")),
	)

	viperCfg.SetEnvPrefix(appName)
	viperCfg.AutomaticEnv()

	viperCfg.SetConfigName(configName)
	viperCfg.SetConfigType(configType)

	return &Config{
		AppName:       appName,
		viper:         viperCfg,
		gapScope:      gap.NewScope(gap.User, appName),
		configFlagMap: make(map[string]interface{}),
	}
}

func (c *Config) BindFlagConfigToCommand(key string, cmd *cobra.Command) {
	if fc, ok := c.configFlagMap[key]; ok {
		switch v := fc.(type) {
		case FlagConfig[string]:
			cmd.PersistentFlags().StringP(v.FlagName, v.FlagShortHand, v.Default, v.Description)
			c.viper.BindPFlag(key, cmd.PersistentFlags().Lookup(v.FlagName))
		case FlagConfig[int]:
			cmd.PersistentFlags().IntP(v.FlagName, v.FlagShortHand, v.Default, v.Description)
			c.viper.BindPFlag(key, cmd.PersistentFlags().Lookup(v.FlagName))
		case FlagConfig[bool]:
			cmd.PersistentFlags().BoolP(v.FlagName, v.FlagShortHand, v.Default, v.Description)
			c.viper.BindPFlag(key, cmd.PersistentFlags().Lookup(v.FlagName))
		default:
			panic(fmt.Sprintf("unsupported type: %T", fc))
		}
	}
}

func AddFlagConfig[T string | int | bool](c *Config, flagConfig FlagConfig[T]) {
	c.configFlagMap[flagConfig.Key] = flagConfig
}

func (c *Config) Read() error {
	paths, err := c.gapScope.LookupConfig(configName + "." + configType)
	if err != nil {
		return fmt.Errorf("getting config path: %w", err)
	}

	for _, path := range slices.Backward(paths) {
		c.viper.SetConfigFile(path)
		err := c.viper.MergeInConfig()
		if err != nil {
			return fmt.Errorf("merging config from '%s': %w", path, err)
		}
	}

	return nil
}

func (c *Config) Write() error {
	paths, err := c.ConfigPaths()
	if err != nil {
		return err
	}

	if len(paths) == 0 {
		configDirs, err := c.gapScope.ConfigDirs()
		if err != nil {
			return fmt.Errorf("getting config dir alternatives to create a new config file: %w", err)
		}
		configDir := configDirs[0]
		err = os.MkdirAll(configDir, 0750)
		if err != nil {
			return fmt.Errorf("creating config dir: %w", err)
		}

		configPath := filepath.Join(configDir, configName+"."+configType)
		c.viper.SetConfigFile(configPath)
	} else {
		c.viper.SetConfigFile(paths[0])
	}

	err = c.viper.WriteConfig()
	if err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	return nil
}

func (c *Config) ConfigPaths() ([]string, error) {
	paths, err := c.gapScope.LookupConfig(configName + "." + configType)
	if err != nil {
		return nil, fmt.Errorf("getting config path: %w", err)
	}
	return paths, nil
}

func (c *Config) DataDir() (string, error) {
	dataPaths, err := c.gapScope.DataDirs()
	if err != nil {
		return "", fmt.Errorf("getting data dirs: %w", err)
	}

	if len(dataPaths) == 0 {
		return "", fmt.Errorf("no suitable data dirs found")
	}

	return dataPaths[0], nil
}

func (c *Config) Get(key string) any {
	return c.viper.Get(key)
}

func (c *Config) GetString(key string) string {
	return c.viper.GetString(key)
}

func (c *Config) GetBool(key string) bool {
	return c.viper.GetBool(key)
}

func (c *Config) GetInt(key string) int {
	return c.viper.GetInt(key)
}

func (c *Config) Set(key, value string) {
	c.viper.Set(key, value)
}

func (c *Config) AllKeys() []string {
	return c.viper.AllKeys()
}
