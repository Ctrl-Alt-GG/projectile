package grpc

import (
	"context"
	"net"

	"github.com/Ctrl-Alt-GG/projectile/pkg/agentmsg"
	"gitlab.com/MikeTTh/env"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GameServersServerImpl struct {
	agentmsg.UnimplementedGameServersServer
}

func (s *GameServersServerImpl) Updates(grpc.ClientStreamingServer[agentmsg.GameServer, agentmsg.Error]) error {
	// TODO
	return nil
}

func (s *GameServersServerImpl) Withdraw(context.Context, *agentmsg.Identifier) (*agentmsg.Error, error) {
	// TODO
	return &agentmsg.Error{Type: agentmsg.ErrorType_HAPPY}, nil
}

func RunServer(logger *zap.Logger) error {
	address := env.String("GRPC_BIND_ADDRESS", ":50051")

	lis, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error("failed to open GRPC listen socket", zap.Error(err))
		return err
	}

	grpcServer := grpc.NewServer()
	agentmsg.RegisterGameServersServer(grpcServer, &GameServersServerImpl{})

	logger.Info("Starting GRPC Server...", zap.String("address", address)) // TODO: TLS // TODO: Auth
	return grpcServer.Serve(lis)
}
