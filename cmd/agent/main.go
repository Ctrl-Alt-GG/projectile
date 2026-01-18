package main

import (
	"github.com/Ctrl-Alt-GG/projectile/pkg/framework"
)

func main() {
	logger := framework.SetupLogger()
	defer logger.Sync()
	logger.Info("Starting Projectile agent...")

	logger.Info("Hello world")
}
