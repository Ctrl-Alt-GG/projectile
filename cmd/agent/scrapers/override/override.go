package override

import (
	"context"

	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/config"
	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers"
	"github.com/Ctrl-Alt-GG/projectile/pkg/model"
	"go.uber.org/zap"
)

// ScraperOverrideWrapper is a wrapper that overrides some data gathered by the scraper if needed
type ScraperOverrideWrapper struct {
	Scraper scrapers.Scraper
	Config  config.ScraperOverrides
}

func (sow ScraperOverrideWrapper) Scrape(ctx context.Context, logger *zap.Logger) (model.GameServerDynamicData, error) {
	data, err := sow.Scraper.Scrape(ctx, logger)
	if err != nil {
		return data, err
	}

	// Override if needed
	if sow.Config.MaxPlayers != nil {
		data.MaxPlayers = *sow.Config.MaxPlayers
	}

	if sow.Config.Info != nil {
		data.Info = *sow.Config.Info
	}

	return data, nil
}

func (sow ScraperOverrideWrapper) Capabilities() model.Capabilities {
	return sow.Scraper.Capabilities()
}
