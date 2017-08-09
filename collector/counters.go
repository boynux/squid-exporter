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
	{"server.http", "kbytes_in", "The total number of server http kbytes recevied"},
	{"server.http", "kbytes_out", "The total number of server http kbytes transfered"},

	{"server.all", "requests", "The total number of server all requests"},
	{"server.all", "errors", "The total number of server all errors"},
	{"server.all", "kbytes_in", "The total number of server kbytes recevied"},
	{"server.all", "kbytes_out", "The total number of server kbytes transfered"},

	{"server.ftp", "requests", "The total number of server ftp requests"},
	{"server.ftp", "errors", "The total number of server ftp errors"},
	{"server.ftp", "kbytes_in", "The total number of server ftp kbytes recevied"},
	{"server.ftp", "kbytes_out", "The total number of server ftp kbytes transfered"},

	{"server.other", "requests", "The total number of server other requests"},
	{"server.other", "errors", "The total number of server other errors"},
	{"server.other", "kbytes_in", "The total number of server other kbytes recevied"},
	{"server.other", "kbytes_out", "The total number of server other kbytes transfered"},

	{"swap", "ins", "The total number of server other requests"},
	{"swap", "outs", "The total number of server other errors"},
	{"swap", "files_cleaned", "The total number of server other kbytes recevied"},
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
