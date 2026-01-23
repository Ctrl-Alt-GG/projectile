package scrapers

import (
	"context"

	"github.com/Ctrl-Alt-GG/projectile/pkg/model"
	"github.com/Ctrl-Alt-GG/projectile/pkg/utils"
	"github.com/SpencerSharkey/gomc/query"
	"go.uber.org/zap"
)

type MinecraftScraper struct {
	addr string
}

func NewMinecraftScraperFromConfig(cfg map[string]any) (Scraper, error) {
	return MinecraftScraper{addr: "TODO"}, nil // TODO
}

func (m MinecraftScraper) Scrape(ctx context.Context, logger *zap.Logger) (model.GameServerDynamicData, error) {
	req := query.NewRequest()
	err := req.Connect(m.addr)
	if err != nil {
		logger.Error("Failed to connect to the Minecraft server", zap.Error(err), zap.String("addr", m.addr))
		return model.GameServerDynamicData{}, err
	}

	if ctx.Err() != nil {
		return model.GameServerDynamicData{}, ctx.Err()
	}

	res, err := req.Full()
	if err != nil {
		logger.Error("Failed query the Minecraft server", zap.Error(err), zap.String("addr", m.addr))
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

func (m MinecraftScraper) Capabilities() model.Capabilities {
	return model.Capabilities{
		PlayerCount: true,
		PlayerNames: true,
		PlayerScore: false,
		PlayerTeam:  false,
	}
}
