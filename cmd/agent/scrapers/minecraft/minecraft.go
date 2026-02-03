package minecraft

import (
	"context"
	"fmt"

	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers"
	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers/internal"
	"github.com/Ctrl-Alt-GG/projectile/pkg/model"
	"github.com/Ctrl-Alt-GG/projectile/pkg/utils"
	"github.com/dreamscached/minequery/v2"
	"go.uber.org/zap"
)

type ScraperConfig struct {
	Address string `mapstructure:"address" default:"127.0.0.1:25565"`
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

	sConfig.Address = utils.WithDefaultPort(sConfig.Address, 25565)

	return Scraper{config: sConfig}, nil
}

func (m Scraper) Scrape(ctx context.Context, logger *zap.Logger) (model.GameServerDynamicData, error) {

	res, err := minequery.QueryFull("localhost", 25565)
	if err != nil {
		logger.Error("Failed query the Minecraft server", zap.Error(err), zap.String("addr", m.config.Address))
		return model.GameServerDynamicData{}, err
	}

	plyList := make([]model.Player, len(res.SamplePlayers))
	for i, playerName := range res.SamplePlayers {
		plyList[i] = model.Player{
			Name: playerName,
		}
	}

	fmt.Printf("%+v", res)

	return model.GameServerDynamicData{
		Info:               res.MOTD,
		MaxPlayers:         uint32(res.MaxPlayers),
		OnlinePlayersCount: utils.Ptr(uint32(res.OnlinePlayers)),
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
