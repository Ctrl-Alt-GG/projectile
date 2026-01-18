package main

import (
	"time"

	"github.com/Ctrl-Alt-GG/projectile/cmd/server/db"
	"github.com/Ctrl-Alt-GG/projectile/cmd/server/grpc"
	"github.com/Ctrl-Alt-GG/projectile/cmd/server/http"
	"github.com/Ctrl-Alt-GG/projectile/pkg/framework"
	"go.uber.org/zap"
)

func main() {
	logger := framework.SetupLogger()
	defer logger.Sync()
	logger.Info("Starting Projectile server...")

	// db cleanup job
	go func() {
		for range time.Tick(time.Second * 30) {
			cnt, err := db.CleanupJob()
			if err != nil {
				logger.Error("Error while running cleanup job", zap.Error(err))
			} else {
				if cnt > 0 {
					logger.Info("CleanupJob cleaned up some records", zap.Int("count", cnt))
				}
			}
		}
	}()

	// GRPC server
	go func() {
		err := grpc.RunServer(logger)
		if err != nil {
			panic(err)
		}
	}()

	// HTTP server
	err := http.RunServer(logger)
	if err != nil {
		panic(err)
	}
}
