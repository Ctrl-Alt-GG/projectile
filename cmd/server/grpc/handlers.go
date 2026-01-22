package grpc

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/Ctrl-Alt-GG/projectile/cmd/server/db"
	"github.com/Ctrl-Alt-GG/projectile/pkg/agentmsg"
	"github.com/Ctrl-Alt-GG/projectile/pkg/model"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	ErrInvalidUpdate = errors.New("invalid update")
)

type GameServersHandler struct {
	agentmsg.UnimplementedGameServersServer
}

func (s *GameServersHandler) Updates(stream grpc.ClientStreamingServer[agentmsg.GameServer, emptypb.Empty]) error {
	id := getIDFromCtx(stream.Context())
	l := getLoggerFromCtx(stream.Context()).With(zap.String("id", id))

	startTime := time.Now()
	for {
		updateMsg, err := stream.Recv()
		if err == io.EOF {
			l.Info("Stream closed", zap.Duration("openTime", time.Since(startTime)))
			return nil
		}
		if err != nil {
			l.Error("Error in reading from stream", zap.Error(err))
			return err
		}

		update, ok := model.GameServerDataFromProtobuf(updateMsg)
		if !ok {
			l.Warn("Client sent an invalid error")
			return ErrInvalidUpdate
		}

		db.StoreGameServerData(id, update)
		l.Debug("Update stored")
	}
}

func (s *GameServersHandler) Withdraw(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	l := getLoggerFromCtx(ctx)
	id := getIDFromCtx(ctx)
	cnt := db.DeleteServer(id)
	if cnt > 0 {
		l.Info("Withdrawn server info", zap.String("id", id))
	}
	return nil, nil
}
