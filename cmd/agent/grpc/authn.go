package grpc

import (
	"context"

	"github.com/Ctrl-Alt-GG/projectile/pkg/auth"
)

type IDKeyAuth struct {
	ID  string
	Key string
}

func (t IDKeyAuth) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	return auth.ToMetadata(t.ID, t.Key), nil
}

func (t IDKeyAuth) RequireTransportSecurity() bool {
	return true
}
