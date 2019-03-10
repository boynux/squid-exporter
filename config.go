package main

import (
	"flag"
	"log"
	"strconv"
	"os"
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

	flag.StringVar(&c.ListenAddress, "listen",
		loadEnvStringVar("SQUID_EXPORTER_LISTEN", defaultListenAddress), "Address and Port to bind exporter, in host:port format")
	flag.StringVar(&c.MetricPath, "metrics-path",
		loadEnvStringVar("SQUID_EXPORTER_METRICS_PATH", defaultMetricsPath), "Metrics path to expose prometheus metrics")

	flag.StringVar(&c.SquidHostname, "squid-hostname",
		loadEnvStringVar("SQUID_HOSTNAME", defaultSquidHostname), "Squid hostname")
	flag.IntVar(&c.SquidPort, "squid-port",
		loadEnvIntVar("SQUID_PORT", defaultSquidPort), "Squid port to read metrics")

	flag.StringVar(&c.Login, "squid-login", loadEnvStringVar("SQUID_LOGIN", ""), "Login to squid service")
	flag.StringVar(&c.Password, "squid-password", loadEnvStringVar("SQUID_PASSWORD", ""), "Password to squid service")

	versionFlag = flag.Bool("version", false, "Print the version and exit")

	flag.Parse()

	return c
}

func loadEnvStringVar(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}

	return val
}

func loadEnvIntVar(key string, def int) int {
	valStr := os.Getenv(key)
	if valStr != "" {
		val, err := strconv.ParseInt(valStr, 0, 32)
		if err == nil {
			return int(val)
		}

		log.Printf("Error parsing  %s='%s'. Integer value expected", key, valStr)
	}

	return def
}
