package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/boynux/squid-exporter/config"
)

func TestCreateProxyHelper(t *testing.T) {
	cfg := &config.Config{
		ListenAddress: "192.0.2.1:3192",
		SquidHostname: "127.0.0.1",
		SquidPort:     3128,
	}

	expectedHProxyString := "PROXY TCP4 192.0.2.1 127.0.0.1 3192 3128"

	p, err := createProxyHeader(cfg)
	require.NoError(t, err)
	assert.Equal(t, expectedHProxyString, p, "Proxy headers do not match!")
}
