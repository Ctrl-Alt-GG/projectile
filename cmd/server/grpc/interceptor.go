package grpc

import (
	"context"

	"github.com/Ctrl-Alt-GG/projectile/cmd/server/authn"
	"github.com/Ctrl-Alt-GG/projectile/pkg/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	authnKey  = "a"
	loggerKey = "l"
)

type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

func streamingLoggerInterceptor(logger *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// We are overriding context here...
		wrapped := &wrappedServerStream{
			ServerStream: ss,
			ctx:          context.WithValue(ss.Context(), loggerKey, logger),
		}
		return handler(srv, wrapped)
	}
}

func unaryLoggerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		newCtx := context.WithValue(ctx, loggerKey, logger)
		return handler(newCtx, req)
	}
}

func interceptorDoAuth(ctx context.Context, provider *authn.BasicAuthProvider) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing metadata")
	}

	id, key, ok := auth.FromMetadata(md)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "invalid auth")
	}

	if !provider.ValidateUsernamePassword(id, key) {
		return "", status.Error(codes.Unauthenticated, "invalid auth")
	}

	// success
	return id, nil
}

func streamingAuthnInterceptor(provider *authn.BasicAuthProvider) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		id, err := interceptorDoAuth(ctx, provider)
		if err != nil {
			return err
		}

		// We are overriding context here...
		wrapped := &wrappedServerStream{
			ServerStream: ss,
			ctx:          context.WithValue(ctx, authnKey, id),
		}

		return handler(srv, wrapped)
	}
}

func unaryAuthnInterceptor(provider *authn.BasicAuthProvider) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		id, err := interceptorDoAuth(ctx, provider)
		if err != nil {
			return nil, err
		}

		newCtx := context.WithValue(ctx, authnKey, id)
		return handler(newCtx, req)
	}
}

func getIDFromCtx(ctx context.Context) string {
	val := ctx.Value(authnKey)
	if val == nil {
		panic("authn is nil")
	}
	return val.(string)
}

func getLoggerFromCtx(ctx context.Context) *zap.Logger {
	val := ctx.Value(loggerKey)
	if val == nil {
		panic("logger is nil")
	}
	return val.(*zap.Logger)
}
