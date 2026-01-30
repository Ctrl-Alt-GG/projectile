package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/config"
	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/grpc"
	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers"
	"github.com/Ctrl-Alt-GG/projectile/pkg/framework"
	"github.com/Ctrl-Alt-GG/projectile/pkg/model"
	"github.com/Ctrl-Alt-GG/projectile/pkg/utils"
	"go.uber.org/zap"
)

func printStruct(data any) {
	fmt.Println()
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err := enc.Encode(data)
	if err != nil {
		fmt.Println("  FAILED!", err)
		return
	}
	fmt.Println()
}

func testPingServer(logger *zap.Logger, cm *grpc.ClientManager) {
	fmt.Println("Pinging the server...")

	c, err := cm.GetGRPCClient(logger)
	if err != nil {
		fmt.Println("  FAILED: Could not create client:", err)
		return
	}

	startTime := time.Now()
	_, err = c.Ping(context.Background(), nil)
	delay := time.Since(startTime)
	if err != nil {
		fmt.Println("  FAILED: Error while pinging:", err)
		return
	}
	fmt.Println("  Success! latency:", delay)
}

func testScrapeServer(logger *zap.Logger, scraper scrapers.Scraper) (model.GameServerDynamicData, error) {
	fmt.Println("Running a test scrape...")

	startTime := time.Now()
	data, err := scraper.Scrape(context.Background(), logger)
	delay := time.Since(startTime)
	if err != nil {
		fmt.Println("  FAILED:", err)
		return model.GameServerDynamicData{}, err
	}

	fmt.Println("  Info:", data.Info)
	fmt.Println("  MaxPlayers:", data.MaxPlayers)
	fmt.Println("  OnlinePlayersCount:", utils.FormatPtr(data.OnlinePlayersCount))
	if data.OnlinePlayers != nil {
		fmt.Println("  Players:")
		for _, ply := range *data.OnlinePlayers {
			fmt.Println("  - Name:", ply.Name)
			fmt.Println("    Score:", utils.FormatPtr(ply.Score))
			fmt.Println("    Team:", utils.FormatPtr(ply.Team))
			fmt.Println("    Info:", ply.Info)
		}
	} else {
		fmt.Println("  Players: <nil>")
	}

	fmt.Println("Scrape took", delay)
	return data, nil
}

func testScraperStaticInfo(scraper scrapers.Scraper) {
	caps := scraper.Capabilities()
	fmt.Println("Scraper capabilities:")
	fmt.Println("  PlayerCount:", caps.PlayerCount)
	fmt.Println("  PlayerNames:", caps.PlayerNames)
	fmt.Println("  PlayerScore:", caps.PlayerScore)
	fmt.Println("  PlayerTeam:", caps.PlayerTeam)
}

func testMsg(logger *zap.Logger, cfg config.GameData, scraper scrapers.Scraper, previouslyScrapedData model.GameServerDynamicData) {
	data := initGameServerData(logger, cfg, scraper)
	data.GameServerDynamicData = previouslyScrapedData

	fmt.Println("Would send over this message:")
	printStruct(data)
}

func test() {
	fmt.Println("Running agent test...")
	logger := framework.SetupLogger(true)

	cfg, err := config.LoadConfig(logger, "")
	if err != nil {
		fmt.Println("Failed to load config:", err)
		return
	}
	fmt.Println("Loaded config!")
	printStruct(cfg)

	scraper, err := scrapers.FromConfig(cfg.Scraper)
	if err != nil {
		fmt.Println("Failed to instantiate scraper from config!", err)
		return
	}
	fmt.Println("Loaded scraper module:", cfg.Scraper.Module)

	fmt.Println("---")
	testScraperStaticInfo(scraper)

	fmt.Println("---")
	data, scrapeErr := testScrapeServer(logger, scraper)

	if scrapeErr == nil {
		fmt.Println("---")
		testMsg(logger, cfg.GameData, scraper, data)
	}

	fmt.Println("---")
	cm := grpc.NewClientManagerFromConfig(cfg.Server)
	testPingServer(logger, cm)
}
