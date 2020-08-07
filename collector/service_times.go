package collector

import (
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type squidServiceTimes struct {
	Section     string
	Counter     string
	Suffix      string
	Description string
}

var squidServiceTimess = []squidServiceTimes{
	{"HTTP_Requests", "All", "5", "Service Time Percentiles 5min"},
	{"HTTP_Requests", "All", "10", "Service Time Percentiles 5min"},
	{"HTTP_Requests", "All", "15", "Service Time Percentiles 5min"},
	{"HTTP_Requests", "All", "20", "Service Time Percentiles 5min"},
	{"HTTP_Requests", "All", "25", "Service Time Percentiles 5min"},
	{"HTTP_Requests", "All", "30", "Service Time Percentiles 5min"},
	{"HTTP_Requests", "All", "35", "Service Time Percentiles 5min"},
	{"HTTP_Requests", "All", "40", "Service Time Percentiles 5min"},
	{"HTTP_Requests", "All", "45", "Service Time Percentiles 5min"},
	{"HTTP_Requests", "All", "50", "Service Time Percentiles 5min"},
	{"HTTP_Requests", "All", "55", "Service Time Percentiles 5min"},
	{"HTTP_Requests", "All", "60", "Service Time Percentiles 5min"},
	{"HTTP_Requests", "All", "65", "Service Time Percentiles 5min"},
	{"HTTP_Requests", "All", "70", "Service Time Percentiles 5min"},
	{"HTTP_Requests", "All", "75", "Service Time Percentiles 5min"},
	{"HTTP_Requests", "All", "80", "Service Time Percentiles 5min"},
	{"HTTP_Requests", "All", "85", "Service Time Percentiles 5min"},
	{"HTTP_Requests", "All", "90", "Service Time Percentiles 5min"},
	{"HTTP_Requests", "All", "95", "Service Time Percentiles 5min"},
	{"HTTP_Requests", "All", "100", "Service Time Percentiles 5min"},
	{"Cache_Misses", "", "5", "Service Time Percentiles 5min"},
	{"Cache_Misses", "", "10", "Service Time Percentiles 5min"},
	{"Cache_Misses", "", "15", "Service Time Percentiles 5min"},
	{"Cache_Misses", "", "20", "Service Time Percentiles 5min"},
	{"Cache_Misses", "", "25", "Service Time Percentiles 5min"},
	{"Cache_Misses", "", "30", "Service Time Percentiles 5min"},
	{"Cache_Misses", "", "35", "Service Time Percentiles 5min"},
	{"Cache_Misses", "", "40", "Service Time Percentiles 5min"},
	{"Cache_Misses", "", "45", "Service Time Percentiles 5min"},
	{"Cache_Misses", "", "50", "Service Time Percentiles 5min"},
	{"Cache_Misses", "", "55", "Service Time Percentiles 5min"},
	{"Cache_Misses", "", "60", "Service Time Percentiles 5min"},
	{"Cache_Misses", "", "65", "Service Time Percentiles 5min"},
	{"Cache_Misses", "", "70", "Service Time Percentiles 5min"},
	{"Cache_Misses", "", "75", "Service Time Percentiles 5min"},
	{"Cache_Misses", "", "80", "Service Time Percentiles 5min"},
	{"Cache_Misses", "", "85", "Service Time Percentiles 5min"},
	{"Cache_Misses", "", "90", "Service Time Percentiles 5min"},
	{"Cache_Misses", "", "95", "Service Time Percentiles 5min"},
	{"Cache_Hits", "", "5", "Service Time Percentiles 5min"},
	{"Cache_Hits", "", "10", "Service Time Percentiles 5min"},
	{"Cache_Hits", "", "15", "Service Time Percentiles 5min"},
	{"Cache_Hits", "", "20", "Service Time Percentiles 5min"},
	{"Cache_Hits", "", "25", "Service Time Percentiles 5min"},
	{"Cache_Hits", "", "30", "Service Time Percentiles 5min"},
	{"Cache_Hits", "", "35", "Service Time Percentiles 5min"},
	{"Cache_Hits", "", "40", "Service Time Percentiles 5min"},
	{"Cache_Hits", "", "45", "Service Time Percentiles 5min"},
	{"Cache_Hits", "", "50", "Service Time Percentiles 5min"},
	{"Cache_Hits", "", "55", "Service Time Percentiles 5min"},
	{"Cache_Hits", "", "60", "Service Time Percentiles 5min"},
	{"Cache_Hits", "", "65", "Service Time Percentiles 5min"},
	{"Cache_Hits", "", "70", "Service Time Percentiles 5min"},
	{"Cache_Hits", "", "75", "Service Time Percentiles 5min"},
	{"Cache_Hits", "", "80", "Service Time Percentiles 5min"},
	{"Cache_Hits", "", "85", "Service Time Percentiles 5min"},
	{"Cache_Hits", "", "90", "Service Time Percentiles 5min"},
	{"Cache_Hits", "", "95", "Service Time Percentiles 5min"},
	{"Near_Hits", "", "5", "Service Time Percentiles 5min"},
	{"Near_Hits", "", "10", "Service Time Percentiles 5min"},
	{"Near_Hits", "", "15", "Service Time Percentiles 5min"},
	{"Near_Hits", "", "20", "Service Time Percentiles 5min"},
	{"Near_Hits", "", "25", "Service Time Percentiles 5min"},
	{"Near_Hits", "", "30", "Service Time Percentiles 5min"},
	{"Near_Hits", "", "35", "Service Time Percentiles 5min"},
	{"Near_Hits", "", "40", "Service Time Percentiles 5min"},
	{"Near_Hits", "", "45", "Service Time Percentiles 5min"},
	{"Near_Hits", "", "50", "Service Time Percentiles 5min"},
	{"Near_Hits", "", "55", "Service Time Percentiles 5min"},
	{"Near_Hits", "", "60", "Service Time Percentiles 5min"},
	{"Near_Hits", "", "65", "Service Time Percentiles 5min"},
	{"Near_Hits", "", "70", "Service Time Percentiles 5min"},
	{"Near_Hits", "", "75", "Service Time Percentiles 5min"},
	{"Near_Hits", "", "80", "Service Time Percentiles 5min"},
	{"Near_Hits", "", "85", "Service Time Percentiles 5min"},
	{"Near_Hits", "", "90", "Service Time Percentiles 5min"},
	{"Near_Hits", "", "95", "Service Time Percentiles 5min"},
	{"DNS_Lookups", "", "5", "Service Time Percentiles 5min"},
	{"DNS_Lookups", "", "10", "Service Time Percentiles 5min"},
	{"DNS_Lookups", "", "15", "Service Time Percentiles 5min"},
	{"DNS_Lookups", "", "20", "Service Time Percentiles 5min"},
	{"DNS_Lookups", "", "25", "Service Time Percentiles 5min"},
	{"DNS_Lookups", "", "30", "Service Time Percentiles 5min"},
	{"DNS_Lookups", "", "35", "Service Time Percentiles 5min"},
	{"DNS_Lookups", "", "40", "Service Time Percentiles 5min"},
	{"DNS_Lookups", "", "45", "Service Time Percentiles 5min"},
	{"DNS_Lookups", "", "50", "Service Time Percentiles 5min"},
	{"DNS_Lookups", "", "55", "Service Time Percentiles 5min"},
	{"DNS_Lookups", "", "60", "Service Time Percentiles 5min"},
	{"DNS_Lookups", "", "65", "Service Time Percentiles 5min"},
	{"DNS_Lookups", "", "70", "Service Time Percentiles 5min"},
	{"DNS_Lookups", "", "75", "Service Time Percentiles 5min"},
	{"DNS_Lookups", "", "80", "Service Time Percentiles 5min"},
	{"DNS_Lookups", "", "85", "Service Time Percentiles 5min"},
	{"DNS_Lookups", "", "90", "Service Time Percentiles 5min"},
	{"DNS_Lookups", "", "95", "Service Time Percentiles 5min"},
}

func generateSquidServiceTimes(labels []string) descMap {
	serviceTimes := descMap{}

	for i := range squidServiceTimess {
		serviceTime := squidServiceTimess[i]

		var key string
		var name string

		if serviceTime.Counter != "" {
			key = fmt.Sprintf("%s_%s_%s", serviceTime.Section, serviceTime.Counter, serviceTime.Suffix)
			name = prometheus.BuildFQName(namespace, strings.Replace(serviceTime.Section, ".", "_", -1),
				fmt.Sprintf("%s_%s", serviceTime.Counter, serviceTime.Suffix))
		} else {
			key = fmt.Sprintf("%s_%s", serviceTime.Section, serviceTime.Suffix)
			name = prometheus.BuildFQName(namespace, strings.Replace(serviceTime.Section, ".", "_", -1),
				fmt.Sprintf("%s", serviceTime.Suffix))
		}

		serviceTimes[key] = prometheus.NewDesc(
			name,
			serviceTime.Description,
			labels, nil,
		)
	}

	return serviceTimes
}
