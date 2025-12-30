package app

import (
	"context"
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

func (c *ConfigurationHolder) Save() (err error) {
	file, err := os.OpenFile(c.Path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("opening configuration file %q: %w", c.Path, err)
	}
	defer loex.OnErrJoin(&err, file.Close)

	err = c.WriteTo(file)
	if err != nil {
		return fmt.Errorf("writing configuration to %q: %w", c.Path, err)
	}

	return nil
}

type Configuration struct {
	Data *repository.Configuration `json:"data"`
	Web  *web.Configuration        `json:"web"`
	Log  *logConfiguration         `json:"log"`

	DryStart bool `json:"-"`
}

func ConfigurationDefaults(ctx context.Context) (*Configuration, error) {
	d, err := repository.ConfigurationDefaults(ctx)
	if err != nil {
		return nil, err
	}

	w, err := web.ConfigurationDefaults()
	if err != nil {
		return nil, err
	}

	return &Configuration{
		Data: d,
		Web:  w,
	}, nil
}

type CliFlags struct {
	Config   string
	Data     string
	ServeWeb bool
	DryStart bool
}

func InitializeConfiguration(ctx context.Context, flags CliFlags) (*ConfigurationHolder, error) {
	defaultConfiguration, err := ConfigurationDefaults(ctx)
	if err != nil {
		return nil, fmt.Errorf("obtaining default configuration: %w", err)
	}

	configHolder := ConfigurationHolder{
		Path:          flags.Config,
		Configuration: defaultConfiguration,
	}

	if configHolder.Path == "" {
		if configHolder.Path, err = setupDefaultConfigurationEnvironment(ctx); err != nil {
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

func setupDefaultConfigurationEnvironment(ctx context.Context) (string, error) {
	const configName = "lapiasse"
	var appConfigDir string

	runEnv := env.GetRunEnv(ctx)
	switch runEnv {
	case env.EnvProduction:
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("obtaining user configuration directory: %w", err)
		}

		appConfigDir = path.Join(userConfigDir, "LaPiasse")
	case env.EnvTest:
		panic(`Use "CreateTestAppCtx" for tests, not a real one`)
	default:
		workspace_root := os.Getenv("WORKSPACE_ROOT")
		if workspace_root == "" {
			panic("WORKSPACE_ROOT has to be set during development, look at `Makefile` and at `.vscode/settings.json`")
		}

		appConfigDir = path.Join(workspace_root, "tmp", runEnv.String(), "configuration")
	}

	slog.Debug("Ensure app config dir exists", "appConfigDir", appConfigDir)
	if err := os.MkdirAll(appConfigDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("creating application configuration folder at %q: %w", appConfigDir, err)
	}

	return path.Join(appConfigDir, configName+".json"), nil
}

// Apply CLI flags to configuration.
// Perform this after saving the loaded + defaults configuration, since the CLI flags
// configurations are valid only for this run of the application, and should not be persisted.
func applyCliFlagsToConfiguration(cfg *Configuration, flags CliFlags) {
	// Override data directory if specified via CLI flag.
	if flags.Data != "" {
		cfg.Data.BasePath = flags.Data
	}

	// Enable web server if requested. Will be on if the configuration file says so,
	// but the CLI flag has precedence.
	if flags.ServeWeb {
		cfg.Web.Expose = true
	}

	// Enable dry-start.
	cfg.DryStart = flags.DryStart
}
