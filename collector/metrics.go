package collector

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/version"
)

const (
	namespace = "squid"
	timeout   = 10 * time.Second
)

var (
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last query of squid successful?",
		[]string{"region"}, nil,
	)

	client_http_requests = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "client_http", "requests"),
		"The total number of client requests",
		[]string{}, nil,
	)

	client_http_hits = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "client_http", "hits"),
		"The total number of client cache hits",
		[]string{}, nil,
	)

	client_http_errors = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "client_http", "errors"),
		"The total number of client http errors",
		[]string{}, nil,
	)

	client_http_received = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "client_http", "received"),
		"The total number of client kbytes recevied",
		[]string{}, nil,
	)

	client_http_transfered = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "client_http", "transfered"),
		"The total number of client kbytes transfered",
		[]string{}, nil,
	)

	client_http_kbytes_hit = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "client_http", "kbytes_hit"),
		"The total number of client kbytes cache hit",
		[]string{}, nil,
	)
)

type Exporter struct {
	client SquidClient

	hostname string
	port     int
}

func New(hostname string, port int) *Exporter {
	c := NewCacheObjectClient(hostname, port)

	return &Exporter{
		c,

		hostname,
		port,
	}
}

// Describe describes all the metrics ever exported by the ECS exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- client_http_requests
}

func (e *Exporter) Collect(c chan<- prometheus.Metric) {
	prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0, e.hostname)

	counters, err := e.client.GetCounters()

	if err == nil {
		for i := range counters {
			if counters[i].Key == "client_http.requests" {
				c <- prometheus.MustNewConstMetric(client_http_requests, prometheus.CounterValue, counters[i].Value)
			}
			if counters[i].Key == "client_http.hits" {
				c <- prometheus.MustNewConstMetric(client_http_hits, prometheus.CounterValue, counters[i].Value)
			}
			if counters[i].Key == "client_http.errors" {
				c <- prometheus.MustNewConstMetric(client_http_errors, prometheus.CounterValue, counters[i].Value)
			}
			if counters[i].Key == "client_http.kbytes_in" {
				c <- prometheus.MustNewConstMetric(client_http_received, prometheus.CounterValue, counters[i].Value)
			}
			if counters[i].Key == "client_http.kbytes_out" {
				c <- prometheus.MustNewConstMetric(client_http_transfered, prometheus.CounterValue, counters[i].Value)
			}
			if counters[i].Key == "client_http.hit_kbytes_out" {
				c <- prometheus.MustNewConstMetric(client_http_kbytes_hit, prometheus.CounterValue, counters[i].Value)
			}
		}

	}
}

func init() {
	prometheus.MustRegister(version.NewCollector("squid_exporter"))
}
