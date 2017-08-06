package main

import (
	"log"
	"net/http"
	"squid-exporter/collector"

	"github.com/prometheus/client_golang/prometheus"
)

const IndexContent = `<html>
             <head><title>Squid Exporter</title></head>
             <body>
             <h1>Squid Exporter</h1>
             <p><a href='` + "/metrics" + `'>Metrics</a></p>
             </body>
             </html>`

func main() {
	e := collector.New("localhost", 3129)

	prometheus.MustRegister(e)

	// Serve metrics
	http.Handle("/metrics", prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(IndexContent))
	})

	log.Println("Listening on", "0.0.0.0:8088")
	log.Fatal(http.ListenAndServe("0.0.0.0:8088", nil))
}
