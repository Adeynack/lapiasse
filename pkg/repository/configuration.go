package repository

import (
	"fmt"
	"os"
	"path"

	"adeynack.net/lapiasse/pkg/env"
	"github.com/samber/lo"
)

type Configuration struct {
	// A path to the folder in which all data (e.g. databases, files) is stored.
	BasePath string `json:"base_path"`
}

func ConfigurationDefaults() (Configuration, error) {
	basePath, err := determineDefaultDataDirectory()
	if err != nil {
		return Configuration{}, err
	}

	return Configuration{
		BasePath: basePath,
	}, nil
}

func determineDefaultDataDirectory() (string, error) {
	if env.RunEnv != env.EnvProduction {
		pwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("obtaining working directory: %w", err)
		}

		return path.Join(pwd, "tmp", env.RunEnv.String(), "data"), nil
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
