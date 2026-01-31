package registry

import (
	"sync"

	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers"
)

var reg sync.Map

type ScraperNewFunc func(map[string]any) (scrapers.Scraper, error)

func RegisterScraper(name string, newFn ScraperNewFunc) {
	_, exists := reg.LoadOrStore(name, newFn)
	if exists {
		panic("Scraper already registered")
	}
}

func instantiateScraper(name string, config map[string]any) (scrapers.Scraper, error) {
	v, ok := reg.Load(name)
	if !ok {
		return nil, ErrInvalidScraper
	}

	newFn := v.(ScraperNewFunc)
	return newFn(config)
}

func ListScrapers() []string {
	s := make([]string, 0)
	reg.Range(func(key, _ any) bool {
		s = append(s, key.(string))
		return true
	})
	return s
}
