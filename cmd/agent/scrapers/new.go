package scrapers

import (
	"errors"

	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/config"
)

var ErrInvalidScraper = errors.New("invalid scraper")

func FromConfig(scraperCfg config.Scraper) (Scraper, error) {
	var scraper Scraper
	var err error

	switch scraperCfg.Module {
	case "dummy":
		scraper = DummyScraper{}
	case "minecraft":
		scraper, err = NewMinecraftScraperFromConfig(scraperCfg.Config)
	case "static":
		scraper, err = NewStaticScraperFromConfig(scraperCfg.Config)
	case "valve":
		scraper, err = NewValveScraperFromConfig(scraperCfg.Config)
	default:
		return nil, ErrInvalidScraper
	}
	if err != nil {
		return nil, err
	}

	if scraperCfg.Overrides.Any() {
		scraper = ScraperOverrideWrapper{
			scraper:    scraper,
			MaxPlayers: scraperCfg.Overrides.MaxPlayers,
			Info:       scraperCfg.Overrides.Info,
		}
	}

	return scraper, nil
}
