package collector

import (
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type squidCounter struct {
	Section     string
	Counter     string
	Suffix      string
	Description string
}

var squidCounters = []squidCounter{
	{"client_http", "requests", "total", "The total number of client requests"},
	{"client_http", "hits", "total", "The total number of client cache hits"},
	{"client_http", "errors", "total", "The total number of client http errors"},
	{"client_http", "kbytes_in", "kbytes_total", "The total number of client kbytes received"},
	{"client_http", "kbytes_out", "kbytes_total", "The total number of client kbytes transferred"},
	{"client_http", "hit_kbytes_out", "bytes_total", "The total number of client kbytes cache hit"},

	{"server.http", "requests", "total", "The total number of server http requests"},
	{"server.http", "errors", "total", "The total number of server http errors"},
	{"server.http", "kbytes_in", "kbytes_total", "The total number of server http kbytes received"},
	{"server.http", "kbytes_out", "kbytes_total", "The total number of server http kbytes transferred"},

	{"server.all", "requests", "total", "The total number of server all requests"},
	{"server.all", "errors", "total", "The total number of server all errors"},
	{"server.all", "kbytes_in", "kbytes_total", "The total number of server kbytes received"},
	{"server.all", "kbytes_out", "kbytes_total", "The total number of server kbytes transferred"},

	{"server.ftp", "requests", "total", "The total number of server ftp requests"},
	{"server.ftp", "errors", "total", "The total number of server ftp errors"},
	{"server.ftp", "kbytes_in", "kbytes_total", "The total number of server ftp kbytes received"},
	{"server.ftp", "kbytes_out", "kbytes_total", "The total number of server ftp kbytes transferred"},

	{"server.other", "requests", "total", "The total number of server other requests"},
	{"server.other", "errors", "total", "The total number of server other errors"},
	{"server.other", "kbytes_in", "kbytes_total", "The total number of server other kbytes received"},
	{"server.other", "kbytes_out", "kbytes_total", "The total number of server other kbytes transferred"},

	{"swap", "ins", "total", "The number of objects read from disk"},
	{"swap", "outs", "total", "The number of objects saved to disk"},
	{"swap", "files_cleaned", "total", "The number of orphaned cache files removed by the periodic cleanup procedure"},
}

func generateSquidCounters(labels []string) descMap {
	counters := descMap{}

	for i := range squidCounters {
		counter := squidCounters[i]

		counters[fmt.Sprintf("%s.%s", counter.Section, counter.Counter)] = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, strings.Replace(counter.Section, ".", "_", -1),
				fmt.Sprintf("%s_%s", counter.Counter, counter.Suffix)),
			counter.Description,
			labels, nil,
		)
	}

	return counters
}
