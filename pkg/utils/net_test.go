package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithDefaultPort(t *testing.T) {
	testCases := []struct {
		name string

		addr string
		port uint16

		expectedResult string
	}{
		{
			name: "ipv4_has_port",

			addr: "127.0.0.1:1234",
			port: 4567,

			expectedResult: "127.0.0.1:1234",
		},
		{
			name: "ipv4_has_no_port",

			addr: "127.0.0.1",
			port: 4567,

			expectedResult: "127.0.0.1:4567",
		},
		{
			name: "ipv6_has_port",

			addr: "[::1]:1234",
			port: 4567,

			expectedResult: "[::1]:1234",
		},
		{
			name: "ipv6_has_no_port",

			addr: "::1",
			port: 4567,

			expectedResult: "[::1]:4567",
		},
		{
			name: "domain_has_port",

			addr: "example.com:1234",
			port: 4567,

			expectedResult: "example.com:1234",
		},
		{
			name: "domain_has_no_port",

			addr: "example.com",
			port: 4567,

			expectedResult: "example.com:4567",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := WithDefaultPort(tc.addr, tc.port)
			assert.Equal(t, tc.expectedResult, got)
		})
	}
}
