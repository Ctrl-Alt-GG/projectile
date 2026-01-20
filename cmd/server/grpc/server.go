package grpc

import (
	"net"

	"github.com/Ctrl-Alt-GG/projectile/cmd/server/authn"
	"github.com/Ctrl-Alt-GG/projectile/pkg/agentmsg"
	"gitlab.com/MikeTTh/env"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func RunServer(logger *zap.Logger) error {
	address := env.String("GRPC_BIND_ADDRESS", ":50051")

	lis, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error("failed to open GRPC listen socket", zap.Error(err))
		return err
	}
	defer func(lis net.Listener) {
		err := lis.Close()
		if err != nil {
			logger.Warn("Failed to close listener", zap.Error(err))
		}
	}(lis)

	// setup some options
	var opts []grpc.ServerOption

	// first, the TLS
	tlsCert := env.String("GRPC_TLS_CERT", "")
	tlsKey := env.String("GRPC_TLS_KEY", "")
	tlsEnabled := tlsCert != "" && tlsKey != ""

	if tlsEnabled {
		var tlsCreds credentials.TransportCredentials
		tlsCreds, err = credentials.NewServerTLSFromFile(tlsCert, tlsKey)
		if err != nil {
			logger.Error("Failed to load TLS certs", zap.Error(err))
			return err
		}
		opts = append(opts, grpc.Creds(tlsCreds))
	} else {
		logger.Warn("GRPC server running without TLS!")
	}

	// Second, the interceptors

	basichAuthn, err := authn.NewBasicAuthProvider(env.StringOrPanic("AGENT_HTPASSWD_PATH"))
	if err != nil {
		logger.Error("Failed to load basic authentication credentials", zap.Error(err))
		return err
	}

	opts = append(opts,
		grpc.ChainUnaryInterceptor(
			unaryAuthnInterceptor(basichAuthn),
			unaryLoggerInterceptor(logger),
		),
		grpc.ChainStreamInterceptor(
			streamingAuthnInterceptor(basichAuthn),
			streamingLoggerInterceptor(logger),
		),
	)

	grpcServer := grpc.NewServer(opts...)
	agentmsg.RegisterGameServersServer(grpcServer, &GameServersHandler{})

	logger.Info("Starting GRPC Server...", zap.String("address", address), zap.Bool("tlsEnabled", tlsEnabled))
	return grpcServer.Serve(lis)
}
