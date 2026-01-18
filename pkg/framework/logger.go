package framework

import (
	"os"
	"slices"

	"github.com/gin-gonic/gin"
	"gitlab.com/MikeTTh/env"
	"go.uber.org/zap"
)

func SetupLogger() *zap.Logger {
	debugMode := env.Bool("DEBUG", false) || slices.Contains(os.Args, "--debug") || slices.Contains(os.Args, "-debug")

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

	return logger
}
