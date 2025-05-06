package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	proxyproto "github.com/pires/go-proxyproto"

	"github.com/boynux/squid-exporter/config"
)

func createProxyHeader(cfg *config.Config) (string, error) {
	la := strings.Split(cfg.ListenAddress, ":")
	if len(la) < 2 {
		return "", fmt.Errorf("cannot parse listen address %q", cfg.ListenAddress)
	}

	spt, err := strconv.Atoi(la[1])
	if err != nil {
		return "", fmt.Errorf("cannot parse listen port: %w", err)
	}

	sip, err := net.LookupIP(la[0])
	if err != nil {
		return "", fmt.Errorf("cannot resolve listen IP address: %w", err)
	}

	dip, err := net.LookupIP(cfg.SquidHostname)
	if err != nil {
		return "", fmt.Errorf("cannot resolve Squid host: %w", err)
	}

	ph := &proxyproto.Header{
		Version:           1,
		Command:           proxyproto.PROXY,
		TransportProtocol: proxyproto.TCPv4,
		SourceAddr: &net.TCPAddr{
			IP:   sip[0],
			Port: spt,
		},
		DestinationAddr: &net.TCPAddr{
			IP:   dip[0],
			Port: cfg.SquidPort,
		},
	}

	phs, err := ph.Format()
	if err != nil {
		return "", fmt.Errorf("cannot create proxy header: %w", err)
	}

	// proxyproto adds crlf to the end of the header string, but we will add this later
	// we are triming it here.
	return strings.TrimSuffix(string(phs), "\r\n"), nil
}
