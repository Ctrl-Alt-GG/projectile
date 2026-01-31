package satisfactory

import (
	"context"
	"fmt"

	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers"
	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers/internal"
	"github.com/Ctrl-Alt-GG/projectile/pkg/model"
	"github.com/Ctrl-Alt-GG/projectile/pkg/utils"
	"go.uber.org/zap"
)

type ScraperConfig struct {
	URL   string `mapstructure:"url" default:"https://127.0.0.1:7777/api/v1"`
	Token string `mapstructure:"token"`
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
	resp, err := doQuery(ctx, logger, s.config.URL, s.config.Token)
	if err != nil {
		logger.Error("Query failed", zap.Error(err))
		return model.GameServerDynamicData{}, err
	}

	return model.GameServerDynamicData{
		Info:               fmt.Sprintf("Tier %d", resp.Data.ServerGameState.TechTier),
		MaxPlayers:         uint32(resp.Data.ServerGameState.PlayerLimit),
		OnlinePlayersCount: utils.Ptr(uint32(resp.Data.ServerGameState.NumConnectedPlayers)),
	}, nil
}

func (s Scraper) Capabilities() model.Capabilities {
	return model.Capabilities{
		PlayerCount: true,
		PlayerNames: false,
		PlayerScore: false,
		PlayerTeam:  false,
	}
}
