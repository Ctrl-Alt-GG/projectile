package scrapers

import "github.com/Ctrl-Alt-GG/projectile/cmd/agent/config"

func FromConfig(scraperCfg config.Scraper) (Scraper, error) {
	var scraper Scraper
	var err error

	switch scraperCfg.Module {
	case "dummy":
		scraper = DummyScraper{}
	case "minecraft":
		scraper, err = NewMinecraftScraperFromConfig(scraperCfg.Config)
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
