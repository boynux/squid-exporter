package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/boynux/squid-exporter/collector"
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
	config := NewConfig()
	if *versionFlag {
		log.Println(version.Print("squid_exporter"))
		os.Exit(0)
	}
	log.Println("Scraping metrics from", fmt.Sprintf("%s:%d", config.SquidHostname, config.SquidPort))
	e := collector.New(config.SquidHostname, config.SquidPort, config.Login, config.Password)

	prometheus.MustRegister(e)

	// Serve metrics
	http.Handle(config.MetricPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(indexContent))
	})

	log.Println("Listening on", fmt.Sprintf("%s", config.ListenAddress))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s", config.ListenAddress), nil))
}
