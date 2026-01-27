package scrapers

import (
	"context"

	"github.com/Ctrl-Alt-GG/projectile/pkg/model"
	"github.com/Ctrl-Alt-GG/projectile/pkg/utils"
	"github.com/go-viper/mapstructure/v2"
	"github.com/rumblefrog/go-a2s"
	"go.uber.org/zap"
)

type ValveScraperConfig struct {
	Address string `mapstructure:"address"`
}

type ValveScraper struct {
	config ValveScraperConfig
}

func NewValveScraperFromConfig(cfg map[string]any) (Scraper, error) {
	var sConfig ValveScraperConfig

	err := mapstructure.Decode(cfg, &sConfig)
	if err != nil {
		return nil, err
	}

	return ValveScraper{config: sConfig}, nil
}

func (s ValveScraper) Scrape(ctx context.Context, logger *zap.Logger) (model.GameServerDynamicData, error) {
	client, err := a2s.NewClient(s.config.Address)
	if err != nil {
		logger.Error("Failed to create new client to query the server", zap.Error(err))
		return model.GameServerDynamicData{}, err
	}
	defer func(client *a2s.Client) {
		err := client.Close()
		if err != nil {
			logger.Warn("Failed to close query client", zap.Error(err))
		}
	}(client)

	info, err := client.QueryInfo() // QueryInfo, QueryPlayer, QueryRules
	if err != nil {
		logger.Error("Failed to query the server", zap.Error(err))
		return model.GameServerDynamicData{}, err
	}

	players, err := client.QueryPlayer()
	if err != nil {
		logger.Warn("Failed to query players. Continuing...", zap.Error(err))
		players = nil
	}

	var onlinePlayers []model.Player
	if players != nil {
		onlinePlayers = make([]model.Player, len(players.Players))
		for i, ply := range players.Players {
			onlinePlayers[i] = model.Player{
				Name:  ply.Name,
				Score: utils.Ptr(int32(ply.Score)), // this is presented as uint32 but that's wrong...
			}
		}
	}

	return model.GameServerDynamicData{
		Info:               info.Map,
		MaxPlayers:         uint32(info.MaxPlayers),
		OnlinePlayersCount: utils.Ptr(uint32(info.Players)),
		OnlinePlayers:      &onlinePlayers,
	}, nil
}

func (s ValveScraper) Capabilities() model.Capabilities {
	return model.Capabilities{
		PlayerCount: true,
		PlayerNames: true,
		PlayerScore: true,
		PlayerTeam:  false,
	}
}
