package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/boynux/squid-exporter/collector"
	"github.com/boynux/squid-exporter/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
)

const indexContent = `<html>
             <head><title>Squid Exporter</title></head>
             <body>
             <h1>Squid Exporter</h1>
             <p><a href='` + "/metrics" + `'>Metrics</a></p>
             </body>
             </html>`

func init() {
	prometheus.MustRegister(version.NewCollector("squid_exporter"))
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
		procExporter := prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{
			PidFn: func() (int, error) {
				content, err := ioutil.ReadFile(cfg.Pidfile)
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(indexContent))
	})

	log.Println("Listening on", fmt.Sprintf("%s", cfg.ListenAddress))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s", cfg.ListenAddress), nil))
}
