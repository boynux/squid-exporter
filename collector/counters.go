package collector

import (
	"fmt"
	"strings"

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
	{"server.http", "requests", "The total number of server http requests"},
	{"server.http", "errors", "The total number of server http errors"},
	{"server.http", "kbytes_in", "The total number of server kbytes recevied"},
	{"server.http", "kbytes_out", "The total number of server kbytes transfered"},
}

func generateSquidCounters() DescMap {
	counters := DescMap{}

	for i := range squidCounters {
		counter := squidCounters[i]

		counters[fmt.Sprintf("%s.%s", counter.Section, counter.Counter)] = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, strings.Replace(counter.Section, ".", "_", -1), counter.Counter),
			counter.Description,
			[]string{}, nil,
		)
	}

	return counters
}
