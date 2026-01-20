package main

import (
	"context"
	"time"

	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/grpc"
	"github.com/Ctrl-Alt-GG/projectile/pkg/framework"
	"github.com/Ctrl-Alt-GG/projectile/pkg/model"
)

func main() {
	logger := framework.SetupLogger()
	defer logger.Sync()
	logger.Info("Starting Projectile agent...")

	cm := grpc.NewClientManager("localhost:50051", "alma", "hello")
	c, err := cm.GetGRPCClient(logger)
	if err != nil {
		panic(err)
	}

	stream, err := c.Updates(context.Background())
	if err != nil {
		panic(err)
	}

	var i uint32
	for {
		i++
		msg := model.GameServer{
			Game:               "",
			Name:               "",
			Addresses:          nil,
			Info:               "",
			Capabilities:       model.Capabilities{},
			MaxPlayers:         i,
			OnlinePlayersCount: nil,
			OnlinePlayers:      nil,
			LastUpdate:         time.Time{},
		}

		err = stream.Send(msg.ToProtobuf())
		if err != nil {
			stream.CloseSend()
			cm.Close(logger)
			time.Sleep(time.Second)
			c, err = cm.GetGRPCClient(logger)
			if err != nil {
				panic(err)
			}
			stream, err = c.Updates(context.Background())
			if err != nil {
				panic(err)
			}
		}

		time.Sleep(time.Second)
	}

	/*
		_, err = c.Withdraw(context.Background(), nil)
		if err != nil {
			logger.Error("nemjolet", zap.Error(err))
		} else {
			logger.Info("jolet")
		}
	*/
}
