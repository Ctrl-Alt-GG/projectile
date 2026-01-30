package scrapers

import (
	"context"
	"encoding/json"
	"errors"
	"os/exec"

	"github.com/Ctrl-Alt-GG/projectile/pkg/model"
	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/zap"
)

type ScriptScraperConfig struct {
	Path         string             `mapstructure:"path"`
	Args         []string           `mapstructure:"args"`
	Env          []string           `mapstructure:"env"`
	Workdir      string             `mapstructure:"workdir"`
	Capabilities model.Capabilities `mapstructure:"capabilities"`
}

type ScriptScraper struct {
	config ScriptScraperConfig
}

func NewScriptScraperFromConfig(cfg map[string]any) (Scraper, error) {
	var sConfig ScriptScraperConfig

	err := mapstructure.Decode(cfg, &sConfig)
	if err != nil {
		return nil, err
	}

	return ScriptScraper{config: sConfig}, nil
}

func (s ScriptScraper) Scrape(ctx context.Context, logger *zap.Logger) (model.GameServerDynamicData, error) {
	// build the command
	x := exec.CommandContext(ctx, s.config.Path, s.config.Args...)
	x.Env = append(x.Environ(), s.config.Env...)
	x.Dir = s.config.Workdir

	logger.Debug("Running command...", zap.String("path", s.config.Path))

	outputBytes, err := x.Output()
	var exitCode int

	if err != nil {
		var exiterr *exec.ExitError
		ok := errors.As(err, &exiterr)
		if ok {
			exitCode = exiterr.ExitCode()

			logger.Error(
				"The script exited with unclean exit code!",
				zap.Error(err),
				zap.Int("exitCode", exitCode),
				zap.ByteString("stderr", exiterr.Stderr),
				zap.ByteString("stdout", outputBytes),
			)
			return model.GameServerDynamicData{}, err
		}

		// some other error
		logger.Error("Failed to run script!", zap.Error(err))
		return model.GameServerDynamicData{}, err
	}

	logger.Debug("Script run successfully, parsing it's output...")

	var output model.GameServerDynamicData
	err = json.Unmarshal(outputBytes, &output)
	if err != nil {
		logger.Error("Failed to parse the output of the script", zap.Error(err))
		return model.GameServerDynamicData{}, err
	}

	logger.Debug("Output parsed successfully!")
	return output, nil
}

func (s ScriptScraper) Capabilities() model.Capabilities {
	return s.config.Capabilities
}
