package scrapers

import (
	"context"
	"math/rand/v2"

	"github.com/Ctrl-Alt-GG/projectile/pkg/model"
	"github.com/Ctrl-Alt-GG/projectile/pkg/utils"
	"go.uber.org/zap"
)

// DummyScraper is intended for testing purposes only, it just generates some random data each scrape...
type DummyScraper struct{}

func (d DummyScraper) Scrape(ctx context.Context, logger *zap.Logger) (model.GameServerDynamicData, error) {
	dummyNames := []string{
		"ASDfighter", "Yoloman", "Superman",
		"Batman", "Lajos", "UnnamedPlayer",
		"Steve", "Bitterman", "Ayumi",
		"Jenna", "Mememan", "Orang",
		"Leonídász", "RandomPistike", "Fighter",
		"Goku", "Vegeta", "Micimack",
	}
	dummyTeams := []string{
		"alpha", "bravo", "charlie", "sierra",
	}
	dummyInfo := []string{
		"dead", "flag",
	}

	maxPly := len(dummyNames)
	onlinePly := rand.IntN(maxPly)

	ply := make([]model.Player, onlinePly)
	for i := range onlinePly {
		info := ""
		if rand.IntN(10) == 0 {
			info = dummyInfo[rand.IntN(len(dummyInfo))]
		}

		ply[i] = model.Player{
			Name:  dummyNames[i],
			Score: utils.Ptr(int32(rand.IntN(110) - 10)),
			Team:  utils.Ptr(dummyTeams[i%len(dummyTeams)]),
			Info:  info,
		}
	}

	return model.GameServerDynamicData{
		Info:               "Currently doing dummy things",
		MaxPlayers:         uint32(maxPly),
		OnlinePlayersCount: utils.Ptr(uint32(onlinePly)),
		OnlinePlayers:      &ply,
	}, nil
}

func (d DummyScraper) Capabilities() model.Capabilities {
	return model.Capabilities{
		PlayerCount: true,
		PlayerNames: true,
		PlayerScore: true,
		PlayerTeam:  true,
	}
}
