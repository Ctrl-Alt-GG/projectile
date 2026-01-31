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

func GetOutboundIP(probeAddr string) (net.IP, error) {
	if probeAddr == "" {
		probeAddr = "8.8.8.8:53"
	}
	conn, err := net.Dial("udp", probeAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}
