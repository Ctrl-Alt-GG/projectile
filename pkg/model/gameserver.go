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

// GameServerDynamicData is data, that have the probability to change every scrape.
type GameServerDynamicData struct {
	Info               string    `json:"info"`
	MaxPlayers         uint32    `json:"max_players"`
	OnlinePlayersCount *uint32   `json:"online_players,omitempty"`
	OnlinePlayers      *[]Player `json:"players,omitempty"`
}

// GameServerStaticData is data, that probably does not change between each scrape.
// That's because it probably comes from the config or is a characteristic of the scraper.
// That data can still change on the server side, because the agent can be restarted with a different config.
type GameServerStaticData struct {
	Game         string       `json:"game"`
	Name         string       `json:"name"`
	Addresses    []string     `json:"addresses"`
	Capabilities Capabilities `json:"capabilities"`
}

func (gssd GameServerStaticData) Validate() error {
	if gssd.Game == "" {
		return fmt.Errorf("game is empty")
	}
	if gssd.Name == "" {
		return fmt.Errorf("name is empty")
	}
	if len(gssd.Addresses) == 0 {
		return fmt.Errorf("addresses is empty")
	}
	for _, addr := range gssd.Addresses {
		if addr == "" {
			return fmt.Errorf("one or more addresses are empty")
		}
	}

	if gssd.Capabilities.IsValid() {
		return fmt.Errorf("capabilities are invalid")
	}

	return nil
}

// GameServerData is all the info about a game server we want to present
type GameServerData struct {
	GameServerStaticData
	GameServerDynamicData
}

func (gsd GameServerData) ToProtobuf() *agentmsg.GameServer {
	pb := agentmsg.GameServer{
		Game:      gsd.Game,
		Name:      gsd.Name,
		Addresses: gsd.Addresses,
		Info:      utils.NilStrPtr(gsd.Info),
		Capabilities: &agentmsg.GameServer_Capabilities{
			PlayerCount: gsd.Capabilities.PlayerCount,
			PlayerNames: gsd.Capabilities.PlayerNames,
			PlayerScore: gsd.Capabilities.PlayerScore,
			PlayerTeam:  gsd.Capabilities.PlayerTeam,
		},
		MaxPlayers:         gsd.MaxPlayers,
		OnlinePlayersCount: nil,
		OnlinePlayers:      nil,
	}
	if gsd.Capabilities.PlayerCount {
		pb.OnlinePlayersCount = gsd.OnlinePlayersCount
	}
	if gsd.Capabilities.PlayerNames && gsd.OnlinePlayers != nil {
		pb.OnlinePlayers = make([]*agentmsg.GameServer_Player, len(*gsd.OnlinePlayers))
		for i, ply := range *gsd.OnlinePlayers {
			pb.OnlinePlayers[i] = &agentmsg.GameServer_Player{
				Name:  ply.Name,
				Score: nil,
				Team:  nil,
				Info:  utils.NilStrPtr(ply.Info),
			}

			if gsd.Capabilities.PlayerScore {
				pb.OnlinePlayers[i].Score = ply.Score
			}

			if gsd.Capabilities.PlayerTeam {
				pb.OnlinePlayers[i].Team = ply.Team
			}
		}
	}

	return &pb
}

func (gsd GameServerData) Copy() GameServerData {

	var playersCopy *[]Player

	if gsd.OnlinePlayers != nil {
		p := make([]Player, len(*gsd.OnlinePlayers))

		for i, ply := range *gsd.OnlinePlayers {
			p[i] = Player{
				Name:  ply.Name,
				Score: utils.ValCopy(ply.Score),
				Team:  utils.ValCopy(ply.Team),
				Info:  ply.Info,
			}
		}

		playersCopy = &p
	}

	return GameServerData{
		GameServerStaticData: GameServerStaticData{
			Game:         gsd.Game,
			Name:         gsd.Name,
			Addresses:    slices.Clone(gsd.Addresses),
			Capabilities: gsd.Capabilities,
		},
		GameServerDynamicData: GameServerDynamicData{
			Info:               gsd.Info,
			MaxPlayers:         gsd.MaxPlayers,
			OnlinePlayersCount: utils.ValCopy(gsd.OnlinePlayersCount), // Ptr actually makes a copy, while simply dereferencing doesn't. See https://goplay.tools/snippet/ipMDVGHhgOU
			OnlinePlayers:      playersCopy,
		},
	}
}

func GameServerDataFromProtobuf(server *agentmsg.GameServer) (GameServerData, bool) {
	serverCaps := server.GetCapabilities()
	if serverCaps == nil {
		return GameServerData{}, false
	}
	translatedCaps := Capabilities{
		PlayerCount: serverCaps.GetPlayerCount(),
		PlayerNames: serverCaps.GetPlayerNames(),
		PlayerScore: serverCaps.GetPlayerScore(),
		PlayerTeam:  serverCaps.GetPlayerTeam(),
	}

	gsd := GameServerData{
		GameServerStaticData: GameServerStaticData{
			Game:         server.GetName(),
			Name:         server.GetName(),
			Addresses:    server.GetAddresses(),
			Capabilities: translatedCaps,
		},
		GameServerDynamicData: GameServerDynamicData{
			Info:               server.GetInfo(),
			MaxPlayers:         server.GetMaxPlayers(),
			OnlinePlayersCount: nil,
			OnlinePlayers:      nil,
		},
	}
	if translatedCaps.PlayerCount {
		gsd.OnlinePlayersCount = utils.Ptr(server.GetOnlinePlayersCount())
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

		gsd.OnlinePlayers = &players
	}

	return gsd, true
}

// StoredGameServerData is just GameServerData with a timestamp
// We store these on the server side
type StoredGameServerData struct {
	GameServerData
	LastUpdate time.Time `json:"last_update"`
}

func (sgsd StoredGameServerData) Copy() StoredGameServerData {
	return StoredGameServerData{
		GameServerData: sgsd.GameServerData.Copy(),
		LastUpdate:     sgsd.LastUpdate,
	}
}
