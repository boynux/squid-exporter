package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/boynux/squid-exporter/collector"
	"github.com/boynux/squid-exporter/config"
	"github.com/prometheus/client_golang/prometheus"
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
	log.Println("Scraping metrics from", fmt.Sprintf("%s:%d", cfg.SquidHostname, cfg.SquidPort))
	e := collector.New(cfg.SquidHostname, cfg.SquidPort, cfg.Login, cfg.Password, cfg.Labels)

	prometheus.MustRegister(e)

	// Serve metrics
	http.Handle(cfg.MetricPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(indexContent))
	})

	log.Println("Listening on", fmt.Sprintf("%s", cfg.ListenAddress))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s", cfg.ListenAddress), nil))
}
