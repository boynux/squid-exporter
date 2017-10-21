package main

import (
	"flag"
)

const (
	defaultListenAddress = "127.0.0.1"
	defaultListenPort    = 9301
	defaultMetricsPath   = "/metrics"
	defaultSquidHostname = "localhost"
	defaultSquidPort     = 3128
)

type config struct {
	ListenAddress string
	ListenPort    int
	MetricPath    string

	SquidHostname string
	SquidPort     int
}

func NewConfig() *config {
	c := &config{}

	flag.StringVar(&c.ListenAddress, "listen-address", defaultListenAddress, "Address to bind exporter")
	flag.IntVar(&c.ListenPort, "listen-port", defaultListenPort, "Port to bind exporter")
	flag.StringVar(&c.MetricPath, "metrics-path", defaultMetricsPath, "Metrics path to expose prometheus metrics")

	flag.StringVar(&c.SquidHostname, "squid-hostname", defaultSquidHostname, "Squid hostname")
	flag.IntVar(&c.SquidPort, "squid-port", defaultSquidPort, "Squid port to read metrics")

	flag.Parse()

	return c
}
