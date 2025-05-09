package main

import (
	"github.com/Ctrl-Alt-GG/projectile/http"
	"github.com/gin-gonic/gin"
	"gitlab.com/MikeTTh/env"
	"go.uber.org/zap"
)

func main() {
	debugMode := env.Bool("DEBUG", false)

	// Setup logger
	var err error
	var logger *zap.Logger
	if debugMode {
		gin.SetMode(gin.DebugMode)
		logger, err = zap.NewDevelopment()
	} else {
		gin.SetMode(gin.ReleaseMode)
		logger, err = zap.NewProduction()
	}
	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	logger.Info("Starting Projectile server...")

	if debugMode {
		logger.Warn("Running in DEBUG mode!")
	}

	err = http.RunHTTP(logger)

	if err != nil {
		panic(err)
	}
}
