package collector

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type SquidCounter struct {
	Section     string
	Counter     string
	Description string
}

var squidCounters = []SquidCounter{
	{"client_http", "requests", "The total number of client requests"},
	{"client_http", "hits", "The total number of client cache hits"},
	{"client_http", "errors", "The total number of client http errors"},
	{"client_http", "kbytes_in", "The total number of client kbytes recevied"},
	{"client_http", "kbytes_out", "The total number of client kbytes transfered"},
	{"client_http", "hit_kbytes_out", "The total number of client kbytes cache hit"},
}

func generateSquidCounters() DescMap {
	counters := DescMap{}

	for i := range squidCounters {
		counter := squidCounters[i]

		counters[fmt.Sprintf("%s.%s", counter.Section, counter.Counter)] = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, counter.Section, counter.Counter),
			counter.Description,
			[]string{}, nil,
		)
	}

	return counters
}
