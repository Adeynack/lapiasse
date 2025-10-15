package app

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/creasty/defaults"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type ConfigurationHolder struct {
	Path          string
	Configuration Configuration
}

func (c *ConfigurationHolder) WriteTo(out *os.File) error {
	err := json.MarshalWrite(out, c, jsontext.WithIndent("  "))
	if err != nil {
		return fmt.Errorf("writing configuration to %q: %w", c.Path, err)
	}

	return nil
}

func (c *ConfigurationHolder) Load() error {
	b, err := os.ReadFile(c.Path)
	if err != nil {
		return fmt.Errorf("loading configuration from %q: %w", c.Path, err)
	}

	err = json.Unmarshal(b, c)
	if err != nil {
		return fmt.Errorf("unmarshalling configuration from %q: %w", c.Path, err)
	}

	return nil
}

func (c *ConfigurationHolder) Save() error {
	file, err := os.OpenFile(c.Path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("opening configuration file %q: %w", c.Path, err)
	}
	defer file.Close()

	err = c.WriteTo(file)
	if err != nil {
		return fmt.Errorf("writing configuration to %q: %w", c.Path, err)
	}

	return nil
}

type Configuration struct {
	Data DataConfiguration `json:"data"`
	Web  WebConfiguration  `json:"web"`
}

func (c *Configuration) ApplyDefaults() error {
	return defaults.Set(c)
}

type BasePathConfiguration string

type DataConfiguration struct {
	// A path to the folder in which all data (e.g. databases, files) is stored.
	BasePath string `json:"base_path"`
}

// SetDefaults implements the [defaults.Setter] interface
func (c *DataConfiguration) SetDefaults() {
	var err error

	if c.BasePath, err = determineDefaultDataDirectory(); err != nil {
		panic(fmt.Errorf("determining default data directory: %w", err))
	}
}

type WebConfiguration struct {
	Expose bool `json:"expose" default:"false"`
	Port   int  `json:"port" default:"8080"`
}

type CliFlags struct {
	Config   *string
	Data     *string
	ServeWeb *bool
}

func InitializeConfiguration(flags CliFlags) (*ConfigurationHolder, error) {
	var err error
	configHolder := ConfigurationHolder{Path: lo.FromPtrOr(flags.Config, "")}

	if configHolder.Path == "" {
		if configHolder.Path, err = setupDefaultConfigurationEnvironment(); err != nil {
			return nil, fmt.Errorf("setting up default configuration environment: %w", err)
		}
	} else {
		if _, err := os.Stat(configHolder.Path); err != nil {
			return nil, fmt.Errorf("accessing configuration file %q: %w", configHolder.Path, err)
		}
	}

	// Load from configuration file
	if err := configHolder.Load(); err != nil {
		// It's OK if the config file does not exist.
		if errors.Is(err, os.ErrNotExist) {
			slog.Info(fmt.Sprintf("configuration file not found, creating with defaults at %s", configHolder.Path))
		} else {
			return nil, fmt.Errorf("reading config: %w", err)
		}
	}

	// Applying defaults & saving back to configuration file.
	if err := configHolder.Configuration.ApplyDefaults(); err != nil {
		return nil, fmt.Errorf("applying defaults to configuration: %w", err)
	}
	if err := configHolder.Save(); err != nil {
		return nil, fmt.Errorf("saving configuration to %q: %w", configHolder.Path, err)
	}

	// Bind CLI flags to configuration.
	// Perform this after saving the loaded + defaults configuration, since the CLI flags
	// configurations are valid only for this run of the application, and should not be persisted.
	if lo.FromPtrOr(flags.Data, "") != "" {
		configHolder.Configuration.Data.BasePath = *flags.Data
	}
	if flags.ServeWeb != nil {
		configHolder.Configuration.Web.Expose = *flags.ServeWeb
	}

	return &configHolder, nil
}

func setupDefaultConfigurationEnvironment() (string, error) {
	const configName = "lapiasse"
	var appConfigDir string

	switch RunEnv {
	case EnvProduction:
		userConfigDir, err := os.UserConfigDir()
		cobra.CheckErr(err)

		appConfigDir = path.Join(userConfigDir, "LaPiasse")
	case EnvTest:
		appConfigDir = path.Join(os.TempDir(), "lapiasse", "test-configuration")
	default:
		pwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("obtaining working directory: %w", err)
		}

		appConfigDir = path.Join(pwd, "tmp", RunEnv.String(), "configuration")
	}

	if err := os.MkdirAll(appConfigDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("creating application configuration folder at %q: %w", appConfigDir, err)
	}

	return path.Join(appConfigDir, configName+".json"), nil
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
