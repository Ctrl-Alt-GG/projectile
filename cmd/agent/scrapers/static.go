package scrapers

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/Ctrl-Alt-GG/projectile/pkg/model"
	"github.com/go-viper/mapstructure/v2"
	"github.com/mitchellh/go-ps"
	"go.uber.org/zap"
)

var ErrStaticScraperProcessNotRunning = errors.New("process not running")

type StaticScraperConfig struct {
	PSGrep     string `mapstructure:"psgrep"`
	Info       string `mapstructure:"info"`
	MaxPlayers uint32 `mapstructure:"max_players"`
}

type StaticScraper struct {
	config StaticScraperConfig
}

func NewStaticScraperFromConfig(cfg map[string]any) (Scraper, error) {
	var sConfig StaticScraperConfig

	err := mapstructure.Decode(cfg, &sConfig)
	if err != nil {
		return nil, err
	}

	return StaticScraper{config: sConfig}, nil
}

func (s StaticScraper) Scrape(ctx context.Context, logger *zap.Logger) (model.GameServerDynamicData, error) {

	if s.config.PSGrep != "" {
		processes, err := ps.Processes()
		if err != nil {
			logger.Error("Failed to list processes running on the system...", zap.Error(err))
			return model.GameServerDynamicData{}, err
		}
		found := slices.ContainsFunc(processes, func(proc ps.Process) bool {
			return strings.Contains(proc.Executable(), s.config.PSGrep)
		})
		if !found {
			return model.GameServerDynamicData{}, ErrStaticScraperProcessNotRunning
		}
		logger.Debug("The expected process is running")
	}

	return model.GameServerDynamicData{
		Info:       s.config.Info,
		MaxPlayers: s.config.MaxPlayers,
	}, nil
}

func (s StaticScraper) Capabilities() model.Capabilities {
	return model.Capabilities{
		PlayerCount: false,
		PlayerNames: false,
		PlayerScore: false,
		PlayerTeam:  false,
	}
}
