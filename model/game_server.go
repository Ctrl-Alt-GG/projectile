package model

import "fmt"

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

type GameServer struct {
	Address string `json:"address"` // this will be our key for simplicity
	Game    string `json:"game"`    // stuff like, cs2, minecraft, etc.

	Info string `json:"info"` // optional short text, kinda like motd.

	// meta
	MaxPlayers      int  `json:"maxPlayers"`
	PlayersHasScore bool `json:"playersHasScore"`

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
		Address:         g.Address,
		Game:            g.Game,
		Info:            g.Info,
		MaxPlayers:      g.MaxPlayers,
		PlayersHasScore: g.PlayersHasScore,
		OnlinePlayers:   playersCopy,
	}
}
