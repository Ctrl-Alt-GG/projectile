package http

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/MikeTTh/env"
	"go.uber.org/zap"
)

func RunServer(logger *zap.Logger) error {
	r := gin.New()
	r.Use(goodLoggerMiddleware(logger), gin.Recovery())

	key := env.StringOrPanic("ADMIN_KEY")

	// API root
	apiGroup := r.Group("api")
	apiGroup.Use(goodCORSMiddleware)

	// Admin stuff
	adminGroup := apiGroup.Group("admin")
	adminGroup.Use(goodKeyAuthMiddleware(key))

	serversGroup := adminGroup.Group("servers")
	serversGroup.PUT(":id", updateGameServer)
	serversGroup.DELETE(":id", deleteGameServer)
	serversGroup.GET(":id", getGameServer)
	serversGroup.GET("", getAllGameServersWithID)

	announcementGroup := adminGroup.Group("announcement")
	announcementGroup.GET("", getAnnouncement)
	announcementGroup.PUT("", setAnnouncement)
	announcementGroup.DELETE("", clearAnnouncement)

	// no-auth stuff

	apiGroup.GET("servers", getAllGameServersNoID)
	apiGroup.GET("announcement", getAnnouncement)
	apiGroup.GET("bundle", getBundle) // get all data in one request
	apiGroup.GET("ping")

	// start stuff
	tlsCert := env.String("HTTPS_TLS_CERT", "")
	tlsKey := env.String("HTTPS_TLS_KEY", "")
	tlsEnabled := tlsCert != "" && tlsKey != ""

	address := env.String("HTTP_BIND_ADDRESS", ":8080")

	// Run gin
	logger.Info("Starting HTTP Server...", zap.Bool("tlsEnabled", tlsEnabled), zap.String("address", address))
	if tlsEnabled {
		return r.RunTLS(address, tlsCert, tlsKey)
	} else {
		logger.Warn("Running HTTP mode!")
		return r.Run(address)
	}
}
