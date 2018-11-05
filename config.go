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

var (
	versionFlag *bool
)

/*Config configurations for exporter */
type Config struct {
	ListenAddress string
	ListenPort    int
	MetricPath    string

	SquidHostname string
	SquidPort     int
	Login         string
	Password      string
}

/*NewConfig creates a new config object from command line args */
func NewConfig() *Config {
	c := &Config{}

	flag.StringVar(&c.ListenAddress, "listen", defaultListenAddress, "Address and Port to bind exporter, in host:port format")
	flag.StringVar(&c.MetricPath, "metrics-path", defaultMetricsPath, "Metrics path to expose prometheus metrics")

	flag.StringVar(&c.SquidHostname, "squid-hostname", defaultSquidHostname, "Squid hostname")
	flag.IntVar(&c.SquidPort, "squid-port", defaultSquidPort, "Squid port to read metrics")

	flag.StringVar(&c.Login, "squid-login", "", "Login to squid service")
	flag.StringVar(&c.Password, "squid-password", "", "Password to squid service")
	versionFlag = flag.Bool("version", false, "Print the version and exit")

	flag.Parse()

	return c
}
