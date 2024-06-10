package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/boynux/squid-exporter/collector"
	"github.com/boynux/squid-exporter/config"
	kitlog "github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	versioncollector "github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
)

func init() {
	prometheus.MustRegister(versioncollector.NewCollector("squid_exporter"))
}

func main() {
	cfg := config.NewConfig()
	if *config.VersionFlag {
		log.Println(version.Print("squid_exporter"))
		os.Exit(0)
	}
	collector.ExtractServiceTimes = cfg.ExtractServiceTimes

	headers := []string{}

	if cfg.UseProxyHeader {
		headers = append(headers, createProxyHeader(cfg))
	}

	log.Println("Scraping metrics from", fmt.Sprintf("%s:%d", cfg.SquidHostname, cfg.SquidPort))
	e := collector.New(&collector.CollectorConfig{
		Hostname: cfg.SquidHostname,
		Port:     cfg.SquidPort,
		Login:    cfg.Login,
		Password: cfg.Password,
		Labels:   cfg.Labels,
		Headers:  headers,
	})
	prometheus.MustRegister(e)

	if cfg.Pidfile != "" {
		procExporter := collectors.NewProcessCollector(collectors.ProcessCollectorOpts{
			PidFn: func() (int, error) {
				content, err := os.ReadFile(cfg.Pidfile)
				if err != nil {
					return 0, fmt.Errorf("can't read pid file %q: %s", cfg.Pidfile, err)
				}
				value, err := strconv.Atoi(strings.TrimSpace(string(content)))
				if err != nil {
					return 0, fmt.Errorf("can't parse pid file %q: %s", cfg.Pidfile, err)
				}
				return value, nil
			},
			Namespace: "squid",
		})
		prometheus.MustRegister(procExporter)
	}

	// Serve metrics
	http.Handle(cfg.MetricPath, promhttp.Handler())

	if cfg.MetricPath != "/" {
		landingConfig := web.LandingConfig{
			Name:        "Squid Exporter",
			Description: "Prometheus exporter for Squid caching proxy servers",
			HeaderColor: "#15a5be",
			Version:     version.Info(),
			Links: []web.LandingLinks{
				{
					Address: cfg.MetricPath,
					Text:    "Metrics",
				},
			},
		}
		landingPage, err := web.NewLandingPage(landingConfig)
		if err != nil {
			log.Fatal(err)
		}
		http.Handle("/", landingPage)
	}

	systemdSocket := false
	toolkitFlags := &web.FlagConfig{
		WebListenAddresses: &[]string{cfg.ListenAddress},
		WebSystemdSocket:   &systemdSocket,
		WebConfigFile:      &cfg.WebConfigFile,
	}
	logger := kitlog.NewLogfmtLogger(kitlog.StdlibWriter{})

	server := &http.Server{}
	log.Fatal(web.ListenAndServe(server, toolkitFlags, logger))
}
