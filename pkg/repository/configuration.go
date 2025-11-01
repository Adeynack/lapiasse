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

func ConfigurationDefaults() (*Configuration, error) {
	basePath, err := determineDefaultDataDirectory()
	if err != nil {
		return nil, err
	}

	return &Configuration{
		BasePath: basePath,
	}, nil
}

func ConfigurationForPath(basePath string) (*Configuration, error) {
	return &Configuration{
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

	documentsPath := lo.FindOrElse(
		[]string{
			"Documents",    // Linux, macOS, ...
			"My Documents", // Windows
		},
		"",
		func(candidate string) bool {
			_, err := os.Stat(path.Join(home, candidate))
			return err == nil
		})

	if documentsPath != "" {
		return path.Join(home, documentsPath, "La Piasse"), nil
	}

	return path.Join(home, ".lapiasse"), nil
}

func (c *Configuration) MainDatabaseFilePath() string {
	return path.Join(c.BasePath, "lapiasse.db")
}
