package app

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Configuration struct {
	// A path to the folder in which all data (e.g. databases, files) is stored.
	DataFolder string
}

func InitializeConfiguration(cmd *cobra.Command) (*viper.Viper, error) {
	v := viper.New()
	if err := setConfigurationDefaults(v); err != nil {
		return nil, fmt.Errorf("setting configuration defaults: %w", err)
	}

	configFilePath, err := cmd.Flags().GetString("config")
	if err != nil {
		configFilePath = ""
	}

	if configFilePath == "" {
		if configFilePath, err = setupDefaultConfigurationEnvironment(v); err != nil {
			return nil, fmt.Errorf("setting up default configuration environment: %w", err)
		}
	} else {
		if _, err := os.Stat(configFilePath); err != nil {
			return nil, fmt.Errorf("accessing configuration file %q: %w", configFilePath, err)
		}
		v.SetConfigFile(configFilePath)
	}

	// Load from configuration file
	if err := v.ReadInConfig(); err != nil {
		// It's OK if the config file does not exist.
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			slog.Info(fmt.Sprintf("configuration file not found, creating with defaults at %s", configFilePath))
			if err := v.WriteConfigAs(configFilePath); err != nil {
				return nil, fmt.Errorf("writing configuration to %q: %w", configFilePath, err)
			}
		} else {
			return nil, fmt.Errorf("reading config: %w", err)
		}
	}

	// Load from ENV
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(
		"-", "_", // --foo-bar=123  ==>  FOO_BAR=123
		".", "_", // --web.port-admin=8081  ==>  WEB_PORT_ADMIN=8081
	))

	// Bind CLI flags to configuration.
	for config, flag := range map[string]string{
		"web.expose": "serve-web",
		"data.path":  "data",
	} {
		if err := v.BindPFlag(config, cmd.Flags().Lookup(flag)); err != nil {
			return nil, fmt.Errorf("binging configuration %q to CLI flag %q: %w", config, flag, err)
		}
	}

	return v, nil
}

func setupDefaultConfigurationEnvironment(v *viper.Viper) (string, error) {
	const configType = "json"
	configName := "lapiasse"
	var appConfigDir string

	if RunEnv == EnvProduction {
		userConfigDir, err := os.UserConfigDir()
		cobra.CheckErr(err)

		appConfigDir = path.Join(userConfigDir, "LaPiasse")
	} else {
		pwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("obtaining working directory: %w", err)
		}

		appConfigDir = path.Join(pwd, "tmp", RunEnv.String(), "configuration")
	}

	if err := os.MkdirAll(appConfigDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("creating application configuration folder at %q: %w", appConfigDir, err)
	}

	v.AddConfigPath(appConfigDir)
	v.SetConfigName(configName)
	v.SetConfigType(configType)

	return path.Join(appConfigDir, configName+"."+configType), nil
}

func setConfigurationDefaults(v *viper.Viper) error {
	v.SetDefault("web.expose", false)
	v.SetDefault("web.port", 8080)

	if defaultDataPath, err := determineDefaultDataDirectory(); err == nil {
		v.SetDefault("data.path", defaultDataPath)
	} else {
		return fmt.Errorf("determining default data directory: %w", err)
	}

	return nil
}

func determineDefaultDataDirectory() (string, error) {
	if RunEnv != EnvProduction {
		pwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("obtaining working directory: %w", err)
		}

		return path.Join(pwd, "tmp", RunEnv.String(), "data"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	lo.FindOrElse(
		[]string{
			"Documents",
			"My Documents",
		},
		"",
		func(candidate string) bool {
			_, err := os.Stat(path.Join(home, "Documents"))
			return err == nil
		})

	return "", nil
}
