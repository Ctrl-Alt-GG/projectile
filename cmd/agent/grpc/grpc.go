package grpc

import (
	"crypto/tls"
	"sync"

	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/config"
	"github.com/Ctrl-Alt-GG/projectile/pkg/agentmsg"
	"gitlab.com/MikeTTh/env"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ClientManager struct {
	addr string
	auth IDKeyAuth
	conn *grpc.ClientConn
	mu   sync.Mutex
}

func NewClientManager(addr, id, key string) *ClientManager {
	return &ClientManager{
		addr: addr,
		auth: IDKeyAuth{
			ID:  id,
			Key: key,
		},
		conn: nil,
		mu:   sync.Mutex{},
	}
}

func NewClientManagerFromConfig(server config.Server) *ClientManager {
	return NewClientManager(server.Address, server.ID, server.Key)
}

func (cm *ClientManager) GetGRPCClient(logger *zap.Logger) (agentmsg.GameServersClient, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.conn == nil {
		// Set up a connection to the server.
		var err error

		var creds credentials.TransportCredentials
		customCACert := env.String("GRPC_CUSTOM_CA_CERT", "")
		if customCACert != "" {
			creds, err = credentials.NewClientTLSFromFile(customCACert, "")
			if err != nil {
				logger.Error("Failed to load custom CA cert", zap.Error(err))
				return nil, err
			}
		} else {
			creds = credentials.NewTLS(&tls.Config{})
		}

		cm.conn, err = grpc.NewClient(cm.addr,
			grpc.WithTransportCredentials(creds),
			grpc.WithPerRPCCredentials(cm.auth),
		)
		if err != nil {
			logger.Error("Could not connect to GRPC server", zap.Error(err))
			return nil, err
		}
	}

	return agentmsg.NewGameServersClient(cm.conn), nil
}

func (cm *ClientManager) Close(logger *zap.Logger) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if cm.conn != nil {
		err := cm.conn.Close()
		if err != nil {
			logger.Error("Failed to close connection", zap.Error(err))
		}
		cm.conn = nil
	}
}
