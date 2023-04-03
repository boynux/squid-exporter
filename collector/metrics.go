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

type CollectorConfig struct {
	Hostname string
	Port     int
	Login    string
	Password string
	Labels   config.Labels
	Headers  []string
}

/*New initializes a new exporter */
func New(c *CollectorConfig) *Exporter {
	counters = generateSquidCounters(c.Labels.Keys)
	if ExtractServiceTimes {
		serviceTimes = generateSquidServiceTimes(c.Labels.Keys)
	}

	return &Exporter{
		NewCacheObjectClient(&CacheObjectRequest{
			c.Hostname,
			c.Port,
			c.Login,
			c.Password,
			c.Headers,
		}),

		c.Hostname,
		c.Port,

		c.Labels,
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
					c <- prometheus.MustNewConstMetric(d, prometheus.GaugeValue, insts[i].Value, e.labels.Values...)
				}
			}
		} else {
			log.Println("Could not fetch service times metrics from squid instance: ", err)
		}
	}

	e.up.Collect(c)
}
