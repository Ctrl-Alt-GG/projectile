package config

import (
	"fmt"
	"time"
)

type Server struct {
	Address  string `yaml:"address"`
	ID       string `yaml:"id"`
	Key      string `yaml:"key"`
	MaxRetry int    `yaml:"maxRetry" default:"6"`
}

type GameData struct {
	Game              string   `yaml:"game"`              // Unified name for the game that is understood by the frontend e.g.: `openarena`
	Name              string   `yaml:"name"`              // Friendly name with some extra info e.g.: `OpenArena - CTF`
	AddressesOverride []string `yaml:"addressesOverride"` // Optional: skips auto address lookup
}

func (gd GameData) Validate() error {
	if gd.Game == "" {
		return fmt.Errorf("game is empty")
	}
	if gd.Name == "" {
		return fmt.Errorf("name is empty")
	}
	return nil
}

type ScraperOverrides struct {
	MaxPlayers *uint32 `yaml:"maxPlayers"`
	Info       *string `yaml:"info"`
}

func (so ScraperOverrides) Any() bool {
	return so.MaxPlayers != nil || so.Info != nil
}

type Scraper struct {
	Module    string           `yaml:"module"`
	Config    map[string]any   `yaml:"config"`
	Interval  time.Duration    `yaml:"interval" default:"15s"`
	Timeout   time.Duration    `yaml:"timeout" default:"10s"`
	Overrides ScraperOverrides `yaml:"overrides"`
}

func (s Scraper) Validate() error {
	if s.Timeout > s.Interval {
		return fmt.Errorf("timeout is longer than interval")
	}
	return nil
}

type AgentConfig struct {
	GameData GameData `yaml:"gameData"`
	Server   Server   `yaml:"server"`
	Scraper  Scraper  `yaml:"scraper"`
}

func (ac AgentConfig) Validate() error {
	err := ac.Scraper.Validate()
	if err != nil {
		return err
	}
	err = ac.GameData.Validate()
	if err != nil {
		return err
	}
	return nil
}
