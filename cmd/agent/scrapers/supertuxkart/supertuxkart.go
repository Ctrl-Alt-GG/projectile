package supertuxkart

import (
	"context"

	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers"
	"github.com/Ctrl-Alt-GG/projectile/pkg/model"
	"github.com/Ctrl-Alt-GG/projectile/pkg/utils"
	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/zap"

	_ "github.com/mattn/go-sqlite3"
)

type ScraperConfig struct {
	DBPath     string `mapstructure:"db_path"`
	PSGrep     string `mapstructure:"psgrep"`      // ensure that the server is running by looking for the executable
	MaxPlayers uint32 `mapstructure:"max_players"` // I'm just gonna hard-code this here...
}

type Scraper struct {
	config ScraperConfig
}

func New(cfg map[string]any) (scrapers.Scraper, error) {
	var sConfig ScraperConfig

	err := mapstructure.Decode(cfg, &sConfig)
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
			return model.GameServerDynamicData{}, ErrServerProcessNotRunning
		}
		logger.Debug("The expected process is running")
	}

	players, err := QueryPlayers(ctx, logger, s.config.DBPath)
	if err != nil {
		logger.Error("Failed to query players", zap.Error(err))
		return model.GameServerDynamicData{}, err
	}

	return model.GameServerDynamicData{
		MaxPlayers:         s.config.MaxPlayers,
		OnlinePlayersCount: utils.Ptr(uint32(len(players))),
		OnlinePlayers:      &players,
	}, nil
}

func (s Scraper) Capabilities() model.Capabilities {
	return model.Capabilities{
		PlayerCount: true,
		PlayerNames: true,
		PlayerScore: false,
		PlayerTeam:  false,
	}
}
