package ts3

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers"
	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers/internal"
	"github.com/Ctrl-Alt-GG/projectile/pkg/model"
	"github.com/Ctrl-Alt-GG/projectile/pkg/utils"
	"go.uber.org/zap"
)

type ScraperConfig struct {
	APIURL          string `mapstructure:"apiURL" default:"http://127.0.0.1:10080"`
	APIKey          string `mapstructure:"apiKey"`
	VirtualServerID int    `mapstructure:"virtualServerID" default:"1"`
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

var client = http.Client{}

func (s Scraper) makeRequest(ctx context.Context, command string, out any) error {
	fullURL, err := url.JoinPath(s.config.APIURL, command)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return err
	}
	req.Header.Add("x-api-key", s.config.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	d := json.NewDecoder(resp.Body)
	return d.Decode(&out)
}

func (s Scraper) Scrape(ctx context.Context, logger *zap.Logger) (model.GameServerDynamicData, error) {
	logger.Debug("Collecting server info")

	var serverInfo Envelope[VServerInfo]
	err := s.makeRequest(ctx, "serverinfo", &serverInfo)
	if err != nil {
		logger.Error("Failed to query serverinfo", zap.Error(err))
		return model.GameServerDynamicData{}, err
	}

	if serverInfo.Status.Code != 0 {
		logger.Error("Unexpected status in response for serverinfo", zap.Int("code", serverInfo.Status.Code), zap.String("message", serverInfo.Status.Message), zap.String("failed_permission", serverInfo.Status.FailedPermission), zap.String("extra_message", serverInfo.Status.ExtraMessage))
		return model.GameServerDynamicData{}, fmt.Errorf("unexpected status: %d", serverInfo.Status.Code)
	}

	if len(serverInfo.Body) != 1 {
		logger.Error("Unexpected amount of body entries for serverinfo query", zap.Int("len", len(serverInfo.Body)))
		return model.GameServerDynamicData{}, fmt.Errorf("unexpected amount of body fields: %d", len(serverInfo.Body))
	}

	logger.Debug("Collecting client list")

	var clientList Envelope[Client]
	err = s.makeRequest(ctx, "clientlist", &clientList)
	if err != nil {
		logger.Error("Failed to query client list", zap.Error(err))
		return model.GameServerDynamicData{}, err
	}

	if serverInfo.Status.Code != 0 {
		logger.Error("Unexpected status in response for client list", zap.Int("code", serverInfo.Status.Code), zap.String("message", serverInfo.Status.Message), zap.String("failed_permission", serverInfo.Status.FailedPermission), zap.String("extra_message", serverInfo.Status.ExtraMessage))
		return model.GameServerDynamicData{}, fmt.Errorf("unexpected status: %d", serverInfo.Status.Code)
	}

	// now parse the results
	var maxClients uint64
	maxClients, err = strconv.ParseUint(serverInfo.Body[0].VirtualserverMaxclients, 10, 32)
	if err != nil {
		logger.Error("Failed to parse virtualserver_maxclients field", zap.Error(err))
		return model.GameServerDynamicData{}, err
	}

	var players []model.Player

	for _, ply := range clientList.Body {
		if ply.ClientType == "0" { // TODO: maybe de-duplicate based on database id?
			players = append(players, model.Player{
				Name: ply.ClientNickname,
			})
		}
	}

	return model.GameServerDynamicData{
		Info:               "",
		MaxPlayers:         uint32(maxClients),
		OnlinePlayersCount: utils.Ptr(uint32(len(players))), // ts counts all sorts of things as clients, so this is more reliable
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
