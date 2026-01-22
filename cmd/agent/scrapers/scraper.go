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
