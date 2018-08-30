package main

import (
	"flag"
)

const (
	defaultListenAddress = "127.0.0.1:9301"
	defaultListenPort    = 9301
	defaultMetricsPath   = "/metrics"
	defaultSquidHostname = "localhost"
	defaultSquidPort     = 3128
)

/*Config configurations for exporter */
type Config struct {
	ListenAddress string
	ListenPort    int
	MetricPath    string

	SquidHostname string
	SquidPort     int
}

/*NewConfig creates a new config object from command line args */
func NewConfig() *Config {
	c := &Config{}

	flag.StringVar(&c.ListenAddress, "listen", defaultListenAddress, "Address and Port to bind exporter, in host:port format")
	flag.StringVar(&c.MetricPath, "metrics-path", defaultMetricsPath, "Metrics path to expose prometheus metrics")

	flag.StringVar(&c.SquidHostname, "squid-hostname", defaultSquidHostname, "Squid hostname")
	flag.IntVar(&c.SquidPort, "squid-port", defaultSquidPort, "Squid port to read metrics")

	flag.Parse()

	return c
}
