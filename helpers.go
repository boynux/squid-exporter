package main

import (
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/boynux/squid-exporter/config"

	proxyproto "github.com/pires/go-proxyproto"
)

func createProxyHeader(cfg *config.Config) string {
	la := strings.Split(cfg.ListenAddress, ":")
	if len(la) < 2 {
		log.Printf("Cannot parse listen address (%s). Failed to create proxy header\n", cfg.ListenAddress)
		return ""
	}

	spt, err := strconv.Atoi(la[1])
	if err != nil {
		log.Printf("Failed to create proxy header: %v\n", err.Error())
		return ""
	}

	sip, err := net.LookupIP(la[0])
	if err != nil {
		log.Printf("Failed to create proxy header: %v\n", err.Error())
		return ""
	}

	dip, err := net.LookupIP(cfg.SquidHostname)
	if err != nil {
		log.Printf("Failed to create proxy header: %v\n", err.Error())
		return ""
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
		log.Printf("Failed to create proxy header: %v\n", err.Error())
	}

	// proxyproto adds crlf to the end of the header string, but we will add this later
	// we are triming it here.
	return strings.TrimSuffix(string(phs), "\r\n")
}
