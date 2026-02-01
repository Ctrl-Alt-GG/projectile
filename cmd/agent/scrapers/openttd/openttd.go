package openttd

import (
	"context"
	"fmt"
	"net"

	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers"
	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers/internal"
	"github.com/Ctrl-Alt-GG/projectile/pkg/model"
	"github.com/Ctrl-Alt-GG/projectile/pkg/utils"
	ctxio "github.com/jbenet/go-context/io"
	"go.uber.org/zap"
)

type ScraperConfig struct {
	Address string `mapstructure:"address" default:"127.0.0.1:3979"`
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

	sConfig.Address = utils.WithDefaultPort(sConfig.Address, 3979)

	return Scraper{config: sConfig}, nil
}

func (s Scraper) Scrape(ctx context.Context, logger *zap.Logger) (model.GameServerDynamicData, error) {

	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", s.config.Address)
	if err != nil {
		logger.Error("Failed to dial OpenTTD server", zap.Error(err))
		return model.GameServerDynamicData{}, err
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			logger.Warn("Failed to close connection", zap.Error(err))
		}
	}(conn)

	writer := ctxio.NewWriter(ctx, conn)
	_, err = writer.Write([]byte{0x03, 0x00, 0x07}) // <- I've actually wiresharked that lol
	if err != nil {
		logger.Error("Failed to send the query byte sequence", zap.Error(err))
		return model.GameServerDynamicData{}, err
	}

	// do a check on the context after sending the bytes
	if ctx.Err() != nil {
		return model.GameServerDynamicData{}, err
	}

	var info NetworkServerGameInfo
	reader := ctxio.NewReader(ctx, conn)
	info, err = ParseNetworkGameInfo(reader)
	if err != nil {
		logger.Error("Failed to send the query byte sequence", zap.Error(err))
		return model.GameServerDynamicData{}, err
	}

	companies := "companies"
	if info.CompaniesOn == 1 {
		companies = "company"
	}

	infoStr := fmt.Sprintf("%d %s in %d", info.CompaniesOn, companies, FigureOutYear(info.CalendarDate))

	return model.GameServerDynamicData{
		Info:               infoStr,
		MaxPlayers:         uint32(info.ClientsMax),
		OnlinePlayersCount: utils.Ptr(uint32(info.ClientsOn)),
		OnlinePlayers:      nil, // :(
	}, err
}

func (s Scraper) Capabilities() model.Capabilities {
	return model.Capabilities{
		PlayerCount: true,
		PlayerNames: false,
		PlayerScore: false,
		PlayerTeam:  false,
	}
}
