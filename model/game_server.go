package model

import (
	"fmt"
	"time"
)

type Player struct {
	Name  string `json:"name"`
	Score int    `json:"score"` // optional, may not be shown...
	Info  string `json:"info"`  // optional short text
}

// Copy creates a deep copy
func (p Player) Copy() Player {
	return Player{
		Name:  p.Name,
		Score: p.Score,
		Info:  p.Info,
	}
}

// GameServerMeta lets us know what dynamic info the agent is capable of reporting
type GameServerMeta struct {
	KnowOnlinePlayerCount bool `json:"knowOnlinePlayerCount"` // the playerCount field is usable
	KnowPlayers           bool `json:"knowPlayers"`           // The onlinePlayers field is usable
	PlayersHasScore       bool `json:"playersHasScore"`       // the Score field for the players are usable
}

func (m GameServerMeta) Copy() GameServerMeta {
	return GameServerMeta{
		KnowOnlinePlayerCount: m.KnowOnlinePlayerCount,
		KnowPlayers:           m.KnowPlayers,
		PlayersHasScore:       m.PlayersHasScore,
	}
}

type GameServer struct {
	Address  string `json:"address"`  // this will be our key for simplicity
	Game     string `json:"game"`     // stuff like, cs2, minecraft, etc.
	LongName string `json:"longName"` // "Counter-Strike 2" for example

	Info string `json:"info"` // optional short text, kinda like motd.

	// meta
	PlayerCount int `json:"playerCount"`
	MaxPlayers  int `json:"maxPlayers"`

	Meta       GameServerMeta `json:"meta"`
	LastUpdate time.Time      `json:"lastUpdate"`

	OnlinePlayers []Player `json:"onlinePlayers"`
}

func (g GameServer) Validate() error {

	if g.MaxPlayers < 1 {
		return fmt.Errorf("max players can not be lass then 1")
	}

	if g.Address == "" {
		return fmt.Errorf("address can not be empty")
	}

	if g.Game == "" {
		return fmt.Errorf("game can not be empty")
	}

	return nil
}

// Copy creates a deep copy (also copying the player slice)
func (g GameServer) Copy() GameServer {

	playersCopy := make([]Player, len(g.OnlinePlayers))
	for i, val := range g.OnlinePlayers {
		playersCopy[i] = val.Copy()
	}

	return GameServer{
		Address:       g.Address,
		Game:          g.Game,
		LongName:      g.LongName,
		Info:          g.Info,
		PlayerCount:   g.PlayerCount,
		MaxPlayers:    g.MaxPlayers,
		Meta:          g.Meta.Copy(),
		LastUpdate:    g.LastUpdate,
		OnlinePlayers: playersCopy,
	}
}
