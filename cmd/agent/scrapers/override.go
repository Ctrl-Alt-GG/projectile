package scrapers

import (
	"context"

	"github.com/Ctrl-Alt-GG/projectile/pkg/model"
	"go.uber.org/zap"
)

// ScraperOverrideWrapper is a wrapper that overrides some data gathered by the scraper if needed
type ScraperOverrideWrapper struct {
	scraper    Scraper
	MaxPlayers *uint32
	Info       *string
}

func (sow ScraperOverrideWrapper) Scrape(ctx context.Context, logger *zap.Logger) (model.GameServerDynamicData, error) {
	data, err := sow.scraper.Scrape(ctx, logger)
	if err != nil {
		return data, err
	}

	// Override if needed
	if sow.MaxPlayers != nil {
		data.MaxPlayers = *sow.MaxPlayers
	}

	if sow.Info != nil {
		data.Info = *sow.Info
	}

	return data, nil
}

func (sow ScraperOverrideWrapper) Capabilities() model.Capabilities {
	return sow.scraper.Capabilities()
}
