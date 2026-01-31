package registry

import (
	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/config"
	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers"
	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers/override"
)

func NewScraperFromConfig(scraperCfg config.Scraper) (scrapers.Scraper, error) {
	scraper, err := instantiateScraper(scraperCfg.Module, scraperCfg.Config)
	if err != nil {
		return nil, err
	}

	if scraperCfg.Overrides.Any() {
		scraper = override.ScraperOverrideWrapper{
			Scraper: scraper,
			Config:  scraperCfg.Overrides,
		}
	}

	return scraper, nil
}
