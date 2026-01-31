package override

import (
	"context"
	"testing"

	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers"
	"github.com/Ctrl-Alt-GG/projectile/pkg/model"
	"github.com/Ctrl-Alt-GG/projectile/pkg/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestOverride(t *testing.T) {
	testCases := []struct {
		name         string
		Override     ScraperOverrideWrapper
		expectedData model.GameServerDynamicData
	}{
		{
			name: "happy__no_override",
			Override: ScraperOverrideWrapper{
				scraper: scrapers.ScraperMock{
					ScrapeFn: func(ctx context.Context, logger *zap.Logger) (model.GameServerDynamicData, error) {
						return model.GameServerDynamicData{
							Info:               "hello",
							MaxPlayers:         12,
							OnlinePlayersCount: utils.Ptr(uint32(12)),
							OnlinePlayers:      nil,
						}, nil
					},
				},
				MaxPlayers: nil,
				Info:       nil,
			},
			expectedData: model.GameServerDynamicData{
				Info:               "hello",
				MaxPlayers:         12,
				OnlinePlayersCount: utils.Ptr(uint32(12)),
				OnlinePlayers:      nil,
			},
		},
		{
			name: "happy__override_max_players",
			Override: ScraperOverrideWrapper{
				scraper: scrapers.ScraperMock{
					ScrapeFn: func(ctx context.Context, logger *zap.Logger) (model.GameServerDynamicData, error) {
						return model.GameServerDynamicData{
							Info:               "hello",
							MaxPlayers:         12,
							OnlinePlayersCount: utils.Ptr(uint32(12)),
							OnlinePlayers:      nil,
						}, nil
					},
				},
				MaxPlayers: utils.Ptr(uint32(14)),
				Info:       nil,
			},
			expectedData: model.GameServerDynamicData{
				Info:               "hello",
				MaxPlayers:         14,
				OnlinePlayersCount: utils.Ptr(uint32(12)),
				OnlinePlayers:      nil,
			},
		},
		{
			name: "happy__override_info",
			Override: ScraperOverrideWrapper{
				scraper: scrapers.ScraperMock{
					ScrapeFn: func(ctx context.Context, logger *zap.Logger) (model.GameServerDynamicData, error) {
						return model.GameServerDynamicData{
							Info:               "hello",
							MaxPlayers:         12,
							OnlinePlayersCount: utils.Ptr(uint32(12)),
							OnlinePlayers:      nil,
						}, nil
					},
				},
				MaxPlayers: nil,
				Info:       utils.Ptr("good bye"),
			},
			expectedData: model.GameServerDynamicData{
				Info:               "good bye",
				MaxPlayers:         12,
				OnlinePlayersCount: utils.Ptr(uint32(12)),
				OnlinePlayers:      nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger := zaptest.NewLogger(t)
			gotData, err := tc.Override.Scrape(context.Background(), logger)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedData, gotData)
		})
	}
}
