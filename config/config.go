package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	defaultListenAddress       = "127.0.0.1:9301"
	defaultListenPort          = 9301
	defaultMetricsPath         = "/metrics"
	defaultSquidHostname       = "localhost"
	defaultSquidPort           = 3128
	defaultExtractServiceTimes = true
	defaultUseProxyHeader      = false
)

const (
	squidExporterListenKey      = "SQUID_EXPORTER_LISTEN"
	squidExporterMetricsPathKey = "SQUID_EXPORTER_METRICS_PATH"
	squidHostnameKey            = "SQUID_HOSTNAME"
	squidPortKey                = "SQUID_PORT"
	squidLoginKey               = "SQUID_LOGIN"
	squidPasswordKey            = "SQUID_PASSWORD"
	squidPidfile                = "SQUID_PIDFILE"
	squidExtractServiceTimes    = "SQUID_EXTRACTSERVICETIMES"
	squidUseProxyHeader         = "SQUID_USE_PROXY_HEADER"
)

var (
	VersionFlag *bool
)

type Labels struct {
	Keys   []string
	Values []string
}

/*Config configurations for exporter */
type Config struct {
	ListenAddress       string
	MetricPath          string
	Labels              Labels
	ExtractServiceTimes bool

	SquidHostname string
	SquidPort     int
	Login         string
	Password      string
	Pidfile       string

	UseProxyHeader bool
}

/*NewConfig creates a new config object from command line args */
func NewConfig() *Config {
	c := &Config{}

	flag.StringVar(&c.ListenAddress, "listen",
		loadEnvStringVar(squidExporterListenKey, defaultListenAddress), "Address and Port to bind exporter, in host:port format")
	flag.StringVar(&c.MetricPath, "metrics-path",
		loadEnvStringVar(squidExporterMetricsPathKey, defaultMetricsPath), "Metrics path to expose prometheus metrics")

	flag.BoolVar(&c.ExtractServiceTimes, "extractservicetimes",
		loadEnvBoolVar(squidExtractServiceTimes, defaultExtractServiceTimes), "Extract service times metrics")

	flag.Var(&c.Labels, "label", "Custom metrics to attach to metrics, use -label multiple times for each additional label")

	flag.StringVar(&c.SquidHostname, "squid-hostname",
		loadEnvStringVar(squidHostnameKey, defaultSquidHostname), "Squid hostname")
	flag.IntVar(&c.SquidPort, "squid-port",
		loadEnvIntVar(squidPortKey, defaultSquidPort), "Squid port to read metrics")

	flag.StringVar(&c.Login, "squid-login", loadEnvStringVar(squidLoginKey, ""), "Login to squid service")
	flag.StringVar(&c.Password, "squid-password", loadEnvStringVar(squidPasswordKey, ""), "Password to squid service")

	flag.StringVar(&c.Pidfile, "squid-pidfile", loadEnvStringVar(squidPidfile, ""), "Optional path to the squid PID file for additional metrics")

	flag.BoolVar(&c.UseProxyHeader, "squid-use-proxy-header",
		loadEnvBoolVar(squidUseProxyHeader, defaultUseProxyHeader), "Use proxy headers when fetching metrics")

	VersionFlag = flag.Bool("version", false, "Print the version and exit")

	flag.Parse()

	return c
}

func loadEnvBoolVar(key string, def bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	switch strings.ToLower(val) {
	case "true":
		return true
	case "false":
		return false
	default:
		return def
	}
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

func (l *Labels) String() string {
	var lbls []string
	for i := range l.Keys {
		lbls = append(lbls, l.Keys[i]+"="+l.Values[i])
	}

	return strings.Join(lbls, ", ")
}

func (l *Labels) Set(value string) error {
	args := strings.Split(value, "=")

	if len(args) != 2 || len(args[1]) < 1 {
		return fmt.Errorf("Label must be in 'key=value' format")
	}

	for _, key := range l.Keys {
		if key == args[0] {
			return fmt.Errorf("Labels must be distinct, found duplicated key %s", args[0])
		}
	}
	l.Keys = append(l.Keys, args[0])
	l.Values = append(l.Values, args[1])

	return nil
}
