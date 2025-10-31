package app

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"

	"adeynack.net/lapiasse/pkg/env"
	"adeynack.net/lapiasse/pkg/platform/loex"
	"adeynack.net/lapiasse/pkg/repository"
	"adeynack.net/lapiasse/pkg/web"
	"github.com/samber/lo"
)

type ConfigurationHolder struct {
	Path          string
	Configuration *Configuration
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
	Data *repository.Configuration `json:"data"`
	Web  *web.Configuration        `json:"web"`
}

func ConfigurationDefaults() (*Configuration, error) {
	d, w, err := loex.GetAllOrErr2(
		repository.ConfigurationDefaults,
		web.ConfigurationDefaults,
	)
	if err != nil {
		return nil, err
	}

	return &Configuration{
		Data: d,
		Web:  w,
	}, nil
}

type CliFlags struct {
	Config   *string
	Data     *string
	ServeWeb *bool
}

func InitializeConfiguration(flags CliFlags) (*ConfigurationHolder, error) {
	defaultConfiguration, err := ConfigurationDefaults()
	if err != nil {
		return nil, fmt.Errorf("obtaining default configuration: %w", err)
	}

	configHolder := ConfigurationHolder{
		Path:          lo.FromPtrOr(flags.Config, ""),
		Configuration: defaultConfiguration,
	}

	if configHolder.Path == "" {
		if configHolder.Path, err = setupDefaultConfigurationEnvironment(); err != nil {
			return nil, fmt.Errorf("setting up default configuration environment: %w", err)
		}
	}

	if err := configHolder.Load(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			slog.Info(fmt.Sprintf("configuration file not found, creating with defaults at %s", configHolder.Path))
		} else {
			return nil, fmt.Errorf("reading config: %w", err)
		}
	}

	// Saving back to configuration file.
	if err := configHolder.Save(); err != nil {
		return nil, fmt.Errorf("saving configuration to %q: %w", configHolder.Path, err)
	}

	applyCliFlagsToConfiguration(configHolder.Configuration, flags)

	return &configHolder, nil
}

func setupDefaultConfigurationEnvironment() (string, error) {
	const configName = "lapiasse"
	var appConfigDir string

	switch env.RunEnv {
	case env.EnvProduction:
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("obtaining user configuration directory: %w", err)
		}

		appConfigDir = path.Join(userConfigDir, "LaPiasse")
	case env.EnvTest:
		appConfigDir = path.Join(os.TempDir(), "lapiasse", "test-configuration")
	default:
		pwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("obtaining working directory: %w", err)
		}

		appConfigDir = path.Join(pwd, "tmp", env.RunEnv.String(), "configuration")
	}

	slog.Debug("Ensure app config dir exists", slog.String("appConfigDir", appConfigDir))
	if err := os.MkdirAll(appConfigDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("creating application configuration folder at %q: %w", appConfigDir, err)
	}

	return path.Join(appConfigDir, configName+".json"), nil
}

// Apply CLI flags to configuration.
// Perform this after saving the loaded + defaults configuration, since the CLI flags
// configurations are valid only for this run of the application, and should not be persisted.
func applyCliFlagsToConfiguration(cfg *Configuration, flags CliFlags) {
	if lo.FromPtrOr(flags.Data, "") != "" {
		cfg.Data.BasePath = *flags.Data
	}

	if flags.ServeWeb != nil {
		cfg.Web.Expose = *flags.ServeWeb
	}
}
