package utils

import (
	"fmt"
	"net"
)

func WithDefaultPort(addr string, defaultPort uint16) string {
	_, port, err := net.SplitHostPort(addr)
	if err == nil && port != "" {
		return addr
	}
	return net.JoinHostPort(addr, fmt.Sprintf("%d", defaultPort))
}
