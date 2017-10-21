package collector

import (
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/version"
)

type descMap map[string]*prometheus.Desc

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

	counters descMap
)

/*Exporter entry point to squid exporter */
type Exporter struct {
	client SquidClient

	hostname string
	port     int
}

/*New initializes a new exporter */
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

	for _, v := range counters {
		ch <- v
	}

}

/*Collect fetches metrics from squid manager and pushes them to promethus */
func (e *Exporter) Collect(c chan<- prometheus.Metric) {
	prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0, e.hostname)

	insts, err := e.client.GetCounters()

	if err == nil {
		for i := range insts {
			if d, ok := counters[insts[i].Key]; ok {
				c <- prometheus.MustNewConstMetric(d, prometheus.CounterValue, insts[i].Value)
			}
		}
	} else {
		log.Println("Could not fetch metrics from squid instance: ", err)
	}
}

func init() {
	prometheus.MustRegister(version.NewCollector("squid_exporter"))

	counters = generateSquidCounters()
}
