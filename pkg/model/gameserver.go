package model

import (
	"fmt"
	"slices"
	"time"

	"github.com/Ctrl-Alt-GG/projectile/pkg/agentmsg"
	"github.com/Ctrl-Alt-GG/projectile/pkg/utils"
)

// Currently it does not seem to me a good idea to store protobuf messages directly
// They seem to contain some state information

type Identifier struct {
	Address string `json:"address"`
	Game    string `json:"game"`
}

func (i Identifier) Equals(otherID Identifier) bool {
	return i.Address == otherID.Address && i.Game == otherID.Game
}

func (i Identifier) IsValid() bool {
	return i.Address != "" && i.Game != ""
}

type Capabilities struct {
	PlayerCount bool `json:"player_count"`
	PlayerNames bool `json:"player_names"`
	PlayerScore bool `json:"player_score"`
	PlayerTeam  bool `json:"player_team"`
}

func (c Capabilities) IsValid() bool {
	if !c.PlayerCount && (c.PlayerNames || c.PlayerScore || c.PlayerTeam) {
		// If count can not be not reported, then no player data could be reported either
		return false
	}
	if !c.PlayerNames && (c.PlayerScore || c.PlayerTeam) {
		// If player names can not be not reported, then no player data could be reported either
		return false
	}

	return true
}

type Player struct {
	Name  string  `json:"name"`
	Score *int32  `json:"score,omitempty"`
	Team  *string `json:"team,omitempty"`
	Info  string  `json:"info"`
}

type GameServer struct {
	ID                 Identifier   `json:"id"`
	AlternateAddresses []string     `json:"alternate_addresses"`
	Name               string       `json:"name"`
	Info               string       `json:"info"`
	Capabilities       Capabilities `json:"capabilities"`
	MaxPlayers         uint32       `json:"max_players"`
	OnlinePlayersCount *uint32      `json:"online_players,omitempty"`
	OnlinePlayers      *[]Player    `json:"players,omitempty"`

	// I'm not a fan of this
	LastUpdate time.Time `json:"last_update"`
}

func (gs GameServer) ToProtobuf() *agentmsg.GameServer {
	pb := agentmsg.GameServer{
		Id: &agentmsg.Identifier{
			Address: gs.ID.Address,
			Game:    gs.ID.Game,
		},
		Name:               gs.Name,
		AlternateAddresses: gs.AlternateAddresses,
		Info:               utils.NilStrPtr(gs.Info),
		Capabilities: &agentmsg.GameServer_Capabilities{
			PlayerCount: gs.Capabilities.PlayerCount,
			PlayerNames: gs.Capabilities.PlayerNames,
			PlayerScore: gs.Capabilities.PlayerScore,
			PlayerTeam:  gs.Capabilities.PlayerTeam,
		},
		MaxPlayers:         gs.MaxPlayers,
		OnlinePlayersCount: nil,
		OnlinePlayers:      nil,
	}
	if gs.Capabilities.PlayerCount {
		pb.OnlinePlayersCount = gs.OnlinePlayersCount
	}
	if gs.Capabilities.PlayerNames && gs.OnlinePlayers != nil {
		pb.OnlinePlayers = make([]*agentmsg.GameServer_Player, len(*gs.OnlinePlayers))
		for i, ply := range *gs.OnlinePlayers {
			pb.OnlinePlayers[i] = &agentmsg.GameServer_Player{
				Name:  ply.Name,
				Score: nil,
				Team:  nil,
				Info:  utils.NilStrPtr(ply.Info),
			}

			if gs.Capabilities.PlayerScore {
				pb.OnlinePlayers[i].Score = ply.Score
			}

			if gs.Capabilities.PlayerTeam {
				pb.OnlinePlayers[i].Team = ply.Team
			}
		}
	}

	return &pb
}

func (gs GameServer) Copy() GameServer {

	var playersCopy *[]Player

	if gs.OnlinePlayers != nil {
		p := make([]Player, len(*gs.OnlinePlayers))

		for i, ply := range *gs.OnlinePlayers {
			p[i] = Player{
				Name:  ply.Name,
				Score: utils.ValCopy(ply.Score),
				Team:  utils.ValCopy(ply.Team),
				Info:  ply.Info,
			}
		}

		playersCopy = &p
	}

	return GameServer{
		ID:                 gs.ID,
		AlternateAddresses: slices.Clone(gs.AlternateAddresses),
		Name:               gs.Name,
		Info:               gs.Info,
		Capabilities:       gs.Capabilities,
		MaxPlayers:         gs.MaxPlayers,
		OnlinePlayersCount: utils.ValCopy(gs.OnlinePlayersCount), // Ptr actually makes a copy, while simply dereferencing doesn't. See https://goplay.tools/snippet/ipMDVGHhgOU
		OnlinePlayers:      playersCopy,
		LastUpdate:         gs.LastUpdate,
	}
}

func (gs GameServer) Validate() error {
	if gs.ID.IsValid() {
		return fmt.Errorf("ID is invalid")
	}

	if gs.Capabilities.IsValid() {
		return fmt.Errorf("capabilities are invalid")
	}

	return nil
}

func GameServerFromProtobuf(server *agentmsg.GameServer) (GameServer, bool) {
	serverID := server.GetId()
	if serverID == nil {
		return GameServer{}, false
	}
	serverCaps := server.GetCapabilities()
	if serverCaps == nil {
		return GameServer{}, false
	}
	translatedCaps := Capabilities{
		PlayerCount: serverCaps.GetPlayerCount(),
		PlayerNames: serverCaps.GetPlayerNames(),
		PlayerScore: serverCaps.GetPlayerScore(),
		PlayerTeam:  serverCaps.GetPlayerTeam(),
	}
	if !translatedCaps.IsValid() {
		return GameServer{}, false
	}

	gs := GameServer{
		ID: Identifier{
			Address: serverID.GetAddress(),
			Game:    serverID.GetGame(),
		},
		AlternateAddresses: server.GetAlternateAddresses(),
		Name:               server.GetName(),
		Info:               server.GetInfo(),
		Capabilities:       translatedCaps,
		MaxPlayers:         server.GetMaxPlayers(),
		OnlinePlayersCount: nil,
		OnlinePlayers:      nil,
	}
	if translatedCaps.PlayerCount {
		gs.OnlinePlayersCount = utils.Ptr(server.GetOnlinePlayersCount())
	}

	serverOnlinePlayers := server.GetOnlinePlayers()

	if translatedCaps.PlayerNames && serverOnlinePlayers != nil { // not using len(x) on purpose, I want an empty array if names available
		players := make([]Player, len(serverOnlinePlayers))

		for i, ply := range serverOnlinePlayers {
			players[i] = Player{
				Name:  ply.GetName(),
				Score: nil,
				Team:  nil,
				Info:  ply.GetInfo(),
			}

			if translatedCaps.PlayerScore {
				players[i].Score = utils.Ptr(ply.GetScore())
			}
			if translatedCaps.PlayerTeam {
				players[i].Team = utils.Ptr(ply.GetTeam())
			}

		}

		gs.OnlinePlayers = &players
	}

	return gs, true
}
