package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/boynux/squid-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
)

const indexContent = `<html>
             <head><title>Squid Exporter</title></head>
             <body>
             <h1>Squid Exporter</h1>
             <p><a href='` + "/metrics" + `'>Metrics</a></p>
             </body>
             </html>`

func main() {
	config := NewConfig()

	log.Println("Scraping metrics from", fmt.Sprintf("%s:%d", config.SquidHostname, config.SquidPort))
	e := collector.New(config.SquidHostname, config.SquidPort)

	prometheus.MustRegister(e)

	// Serve metrics
	http.Handle(config.MetricPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(indexContent))
	})

	log.Println("Listening on", fmt.Sprintf("%s:%d", config.ListenAddress, config.ListenPort))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", config.ListenAddress, config.ListenPort), nil))
}
