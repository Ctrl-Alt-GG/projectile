package minecraft

import (
	"context"

	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers"
	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers/internal"
	"github.com/Ctrl-Alt-GG/projectile/pkg/model"
	"github.com/Ctrl-Alt-GG/projectile/pkg/utils"
	"github.com/SpencerSharkey/gomc/query"
	"go.uber.org/zap"
)

type ScraperConfig struct {
	Address string `mapstructure:"address" default:"127.0.0.1"`
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

func (m Scraper) Scrape(ctx context.Context, logger *zap.Logger) (model.GameServerDynamicData, error) {
	req := query.NewRequest()
	err := req.Connect(m.config.Address)
	if err != nil {
		logger.Error("Failed to connect to the Minecraft server", zap.Error(err), zap.String("addr", m.config.Address))
		return model.GameServerDynamicData{}, err
	}

	if ctx.Err() != nil {
		return model.GameServerDynamicData{}, ctx.Err()
	}

	res, err := req.Full()
	if err != nil {
		logger.Error("Failed query the Minecraft server", zap.Error(err), zap.String("addr", m.config.Address))
		return model.GameServerDynamicData{}, err
	}

	plyList := make([]model.Player, len(res.Players))
	for i, playerName := range res.Players {
		plyList[i] = model.Player{
			Name: playerName,
		}
	}

	return model.GameServerDynamicData{
		Info:               "",
		MaxPlayers:         uint32(res.MaxPlayers),
		OnlinePlayersCount: utils.Ptr(uint32(res.NumPlayers)),
		OnlinePlayers:      &plyList,
	}, nil
}

func (m Scraper) Capabilities() model.Capabilities {
	return model.Capabilities{
		PlayerCount: true,
		PlayerNames: true,
		PlayerScore: false,
		PlayerTeam:  false,
	}
}
