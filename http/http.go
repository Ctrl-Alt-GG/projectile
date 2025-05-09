package http

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/MikeTTh/env"
	"go.uber.org/zap"
)

func RunHTTP(logger *zap.Logger) error {
	r := gin.New()
	r.Use(goodLoggerMiddleware(logger), gin.Recovery())

	key := env.StringOrPanic("KEY")

	// register stuff
	r.StaticFile("", "frontend/index.html")
	r.Static("/static", "frontend/static")

	apiGroup := r.Group("api")
	apiGroup.Use(goodCORSMiddleware)
	serversGroup := apiGroup.Group("servers")
	serversGroup.PUT("", goodKeyAuthMiddleware(key), updateGameServer)
	serversGroup.DELETE(":addr", goodKeyAuthMiddleware(key), deleteGameServer)
	serversGroup.GET("", getAll)

	// start stuff
	tlsCert := env.String("TLS_CERT", "")
	tlsKey := env.String("TLS_KEY", "")
	tlsEnabled := tlsCert != "" && tlsKey != ""

	address := env.String("BIND_ADDRESS", ":8080")

	// Run gin
	logger.Info("Starting HTTP Server...", zap.Bool("tlsEnabled", tlsEnabled), zap.String("address", address))
	if tlsEnabled {
		return r.RunTLS(address, tlsCert, tlsKey)
	} else {
		logger.Warn("Running HTTP mode!")
		return r.Run(address)
	}

}
