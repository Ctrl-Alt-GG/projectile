package valve

import (
	"context"

	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers"
	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers/internal"
	"github.com/Ctrl-Alt-GG/projectile/pkg/model"
	"github.com/Ctrl-Alt-GG/projectile/pkg/utils"
	"github.com/rumblefrog/go-a2s"
	"go.uber.org/zap"
)

type ScraperConfig struct {
	Address       string `mapstructure:"address" default:"127.0.0.1"`
	WorkaroundCS2 bool   `mapstructure:"workaround_cs2" default:"false"` // THIS ALSO NEEDS THIS TO BE INSTALLED ON THE SERVER: https://github.com/Source2ZE/ServerListPlayersFix !!!
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

	if s.config.WorkaroundCS2 { // CS2 has this weird thing, that upon startup it's showing two players with no name (if name fix is installed)
		if info.Players == 2 && players != nil && len(players.Players) == 2 && players.Count == 2 &&
			players.Players[0].Name == "" && players.Players[1].Name == "" && players.Players[0].Score == 0 &&
			players.Players[1].Score == 0 {
			// Yup we are in that weird state again...
			// just alter the response we got...
			players = nil
			info.Players = 0
		}
	}

	// Convert data, then send it back

	var onlinePlayers []model.Player
	if players != nil {
		onlinePlayers = make([]model.Player, len(players.Players))
		for i, ply := range players.Players {
			onlinePlayers[i] = model.Player{
				Name:  ply.Name,
				Score: utils.Ptr(int32(ply.Score)), // this is presented as uint32 but that's wrong...
			}
		}
	} else {
		onlinePlayers = make([]model.Player, 0) // make sure it's always presented as an empty list instead of null
	}

	return model.GameServerDynamicData{
		Info:               info.Map,
		MaxPlayers:         uint32(info.MaxPlayers),
		OnlinePlayersCount: utils.Ptr(uint32(info.Players)),
		OnlinePlayers:      &onlinePlayers,
	}, nil
}

func (s Scraper) Capabilities() model.Capabilities {
	return model.Capabilities{
		PlayerCount: true,
		PlayerNames: true,
		PlayerScore: true,
		PlayerTeam:  false,
	}
}
