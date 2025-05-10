package main

import (
	"github.com/Ctrl-Alt-GG/projectile/db"
	"github.com/Ctrl-Alt-GG/projectile/http"
	"github.com/gin-gonic/gin"
	"gitlab.com/MikeTTh/env"
	"go.uber.org/zap"
	"time"
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

	go func() {
		for {
			time.Sleep(time.Second * 30)
			cnt, err := db.Cleanup()
			if err != nil {
				logger.Error("Error while running cleanup job", zap.Error(err))
			} else {
				if cnt > 0 {
					logger.Info("Cleanup cleaned up some records", zap.Int("count", cnt))
				}
			}
		}
	}()

	err = http.RunHTTP(logger)

	if err != nil {
		panic(err)
	}
}
