package factorio

import (
	"context"
	"strconv"
	"strings"

	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers"
	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers/internal"
	"github.com/Ctrl-Alt-GG/projectile/pkg/model"
	"github.com/Ctrl-Alt-GG/projectile/pkg/utils"
	"github.com/gtaylor/factorio-rcon"
	"go.uber.org/zap"
)

type ScraperConfig struct {
	Address  string `mapstructure:"address" default:"127.0.0.1:34198"`
	Password string `mapstructure:"password" validate:"required"`
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

	sConfig.Address = utils.WithDefaultPort(sConfig.Address, 34198)

	return Scraper{config: sConfig}, nil
}

func (s Scraper) Scrape(ctx context.Context, logger *zap.Logger) (model.GameServerDynamicData, error) {
	r, err := rcon.Dial(s.config.Address)
	if err != nil {
		logger.Error("Failed to dial RCON for Factorio server", zap.Error(err))
		return model.GameServerDynamicData{}, err
	}
	defer func(r *rcon.RCON) {
		err := r.Close()
		if err != nil {
			logger.Warn("Failed to close RCON connection")
		}
	}(r)

	err = r.Authenticate(s.config.Password)
	if err != nil {
		logger.Error("RCON authentication failure to the Factorio server", zap.Error(err))
		return model.GameServerDynamicData{}, err
	}

	// Get max players first
	response, err := r.Execute("/config get max-players")
	if err != nil {
		logger.Error("RCON failed to execute config command", zap.Error(err))
		return model.GameServerDynamicData{}, err
	}

	var maxPlayers int
	maxPlayers, err = strconv.Atoi(strings.TrimSpace(response.Body))
	if err != nil {
		logger.Error("Could not parse max player count")
		return model.GameServerDynamicData{}, err
	}

	response, err = r.Execute("/players")
	if err != nil {
		logger.Error("RCON failed to execute players command", zap.Error(err))
		return model.GameServerDynamicData{}, err
	}

	playerNames, err := parseRCONPlayersList(response.Body)
	if err != nil {
		logger.Error("Error while parsing RCON response", zap.Error(err))
		return model.GameServerDynamicData{}, err
	}

	players := make([]model.Player, len(playerNames))
	for i, name := range playerNames {
		players[i] = model.Player{
			Name: name,
		}
	}

	return model.GameServerDynamicData{
		MaxPlayers:         uint32(maxPlayers),
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
