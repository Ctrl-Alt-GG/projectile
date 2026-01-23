package scrapers

import (
	"context"

	"github.com/Ctrl-Alt-GG/projectile/pkg/model"
	"go.uber.org/zap"
)

type Scraper interface {
	Scrape(ctx context.Context, logger *zap.Logger) (model.GameServerDynamicData, error)
	Capabilities() model.Capabilities
}

type ScraperMock struct {
	ScrapeFn       func(ctx context.Context, logger *zap.Logger) (model.GameServerDynamicData, error)
	CapabilitiesFn func() model.Capabilities
}

func (sm ScraperMock) Scrape(ctx context.Context, logger *zap.Logger) (model.GameServerDynamicData, error) {
	if sm.ScrapeFn != nil {
		return sm.ScrapeFn(ctx, logger)
	}
	return model.GameServerDynamicData{}, nil
}

func (sm ScraperMock) Capabilities() model.Capabilities {
	if sm.CapabilitiesFn != nil {
		return sm.CapabilitiesFn()
	}
	return model.Capabilities{}
}
