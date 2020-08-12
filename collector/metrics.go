package collector

import (
	"log"
	"time"

	"github.com/boynux/squid-exporter/config"
	"github.com/prometheus/client_golang/prometheus"
)

type descMap map[string]*prometheus.Desc

const (
	namespace = "squid"
	timeout   = 10 * time.Second
)

var (
	counters            descMap
	serviceTimes        descMap // ExtractServiceTimes decides if we want to extract service times
	ExtractServiceTimes bool
)

/*Exporter entry point to squid exporter */
type Exporter struct {
	client SquidClient

	hostname string
	port     int

	labels config.Labels
	up     *prometheus.GaugeVec
}

/*New initializes a new exporter */
func New(hostname string, port int, login string, password string, labels config.Labels) *Exporter {
	counters = generateSquidCounters(labels.Keys)
	if ExtractServiceTimes {
		serviceTimes = generateSquidServiceTimes(labels.Keys)
	}
	c := NewCacheObjectClient(hostname, port, login, password)

	return &Exporter{
		c,

		hostname,
		port,

		labels,
		prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "up",
			Help:      "Was the last query of squid successful?",
		}, []string{"host"}),
	}
}

// Describe describes all the metrics ever exported by the ECS exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.up.Describe(ch)

	for _, v := range counters {
		ch <- v
	}

	if ExtractServiceTimes {
		for _, v := range serviceTimes {
			ch <- v
		}
	}

}

/*Collect fetches metrics from squid manager and pushes them to promethus */
func (e *Exporter) Collect(c chan<- prometheus.Metric) {
	insts, err := e.client.GetCounters()

	if err == nil {
		e.up.With(prometheus.Labels{"host": e.hostname}).Set(1)
		for i := range insts {
			if d, ok := counters[insts[i].Key]; ok {
				c <- prometheus.MustNewConstMetric(d, prometheus.CounterValue, insts[i].Value, e.labels.Values...)
			}
		}
	} else {
		e.up.With(prometheus.Labels{"host": e.hostname}).Set(0)
		log.Println("Could not fetch counter metrics from squid instance: ", err)
	}

	if ExtractServiceTimes {
		insts, err = e.client.GetServiceTimes()

		if err == nil {
			for i := range insts {
				if d, ok := serviceTimes[insts[i].Key]; ok {
					c <- prometheus.MustNewConstMetric(d, prometheus.CounterValue, insts[i].Value, e.labels.Values...)
				}
			}
		} else {
			log.Println("Could not fetch service times metrics from squid instance: ", err)
		}
	}

	e.up.Collect(c)
}
