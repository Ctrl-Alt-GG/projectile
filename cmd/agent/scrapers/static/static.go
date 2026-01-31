package static

import (
	"context"

	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers"
	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers/internal"
	"github.com/Ctrl-Alt-GG/projectile/pkg/model"
	"github.com/Ctrl-Alt-GG/projectile/pkg/utils"
	"go.uber.org/zap"
)

type ScraperConfig struct {
	PSGrep     string `mapstructure:"psgrep"`
	Info       string `mapstructure:"info"`
	MaxPlayers uint32 `mapstructure:"max_players" validate:"required"`
}

type Scraper struct {
	config ScraperConfig
}

func New(cfg map[string]any) (scrapers.Scraper, error) {
	var sConfig ScraperConfig
	err := internal.LoadScraperConfig(cfg, &sConfig)
	if err != nil {
		return nil, err
	}

	return Scraper{config: sConfig}, nil
}

func (s Scraper) Scrape(ctx context.Context, logger *zap.Logger) (model.GameServerDynamicData, error) {

	if s.config.PSGrep != "" {
		running, err := utils.PSGrep(s.config.PSGrep)
		if err != nil {
			logger.Error("Error while trying to check for running processes on the system...", zap.Error(err))
			return model.GameServerDynamicData{}, err
		}
		if !running {
			return model.GameServerDynamicData{}, ErrStaticScraperProcessNotRunning
		}
		logger.Debug("The expected process is running")
	}

	return model.GameServerDynamicData{
		Info:       s.config.Info,
		MaxPlayers: s.config.MaxPlayers,
	}, nil
}

func (s Scraper) Capabilities() model.Capabilities {
	return model.Capabilities{
		PlayerCount: false,
		PlayerNames: false,
		PlayerScore: false,
		PlayerTeam:  false,
	}
}
