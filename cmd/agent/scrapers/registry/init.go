package registry

import (
	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers/dummy"
	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers/minecraft"
	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers/openttd"
	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers/script"
	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers/static"
	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers/supertuxkart"
	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers/valve"
)

// Eventually, this didn't go as I planned :(
func init() {
	RegisterScraper("dummy", dummy.New)
	RegisterScraper("minecraft", minecraft.New)
	RegisterScraper("openttd", openttd.New)
	RegisterScraper("script", script.New)
	RegisterScraper("static", static.New)
	RegisterScraper("valve", valve.New)
	RegisterScraper("supertuxkart", supertuxkart.New)
}
