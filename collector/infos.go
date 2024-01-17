package collector

import (
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type squidInfos struct {
	Section     string
	Description string
	Unit        string
}

var squidInfoss = []squidInfos{
	{"Number_of_clients_accessing_cache", "", "number"},
	{"Number_of_HTTP_requests_received", "", "number"},
	{"Number_of_ICP_messages_received", "", "number"},
	{"Number_of_ICP_messages_sent", "", "number"},
	{"Number_of_queued_ICP_replies", "", "number"},
	{"Number_of_HTCP_messages_received", "", "number"},
	{"Number_of_HTCP_messages_sent", "", "number"},
	{"Request_failure_ratio", "", "%"},
	{"Average_HTTP_requests_per_minute_since_start", "", "%"},
	{"Average_ICP_messages_per_minute_since_start", "", "%"},
	{"Select_loop_called", "", "number"},
	{"Hits_as_%_of_all_requests", "", "5min %"},
	{"Hits_as_%_of_bytes_sent", "", "5min %"},
	{"Memory_hits_as_%_of_hit_requests", "", "5min %"},
	{"Disk_hits_as_%_of_hit_requests", "", "5min %"},
	{"Storage_Swap_size", "", "KB"},
	{"Storage_Swap_capacity", "", "% use"},
	{"Storage_Mem_size", "", "KB"},
	{"Storage_Mem_capacity", "", "% used"},
	{"Mean_Object_Size", "", "KB"},
	{"Requests_given_to_unlinkd", "", "number"},
	{"UP_Time", "time squid is up", "seconds"},
	{"CPU_Time", "", "seconds"},
	{"CPU_Usage", "of cpu usage", "%"},
	{"CPU_Usage_5_minute_avg", "of cpu usage", "%"},
	{"CPU_Usage_60_minute_avg", "of cpu usage", "%"},
	{"Maximum_Resident_Size", "", "KB"},
	{"Page_faults_with_physical_i_o", "", "number"},
	{"Total_accounted", "", "KB"},
	{"memPoolAlloc_calls", "", "number"},
	{"memPoolFree_calls", "", "number"},
	{"Maximum_number_of_file_descriptors", "", "number"},
	{"Largest_file_desc_currently_in_use", "", "number"},
	{"Number_of_file_desc_currently_in_use", "", "number"},
	{"Files_queued_for_open", "", "number"},
	{"Available_number_of_file_descriptors", "", "number"},
	{"Reserved_number_of_file_descriptors", "", "number"},
	{"Store_Disk_files_open", "", "number"},
	{"StoreEntries", "", "number"},
	{"StoreEntries_with_MemObjects", "", "number"},
	{"Hot_Object_Cache_Items", "", "number"},
	{"on_disk_objects", "", "number"},
}

func generateSquidInfos(labels []string) descMap {
	infos := descMap{}

	for i := range squidInfoss {
		info := squidInfoss[i]

		var key string
		var name string
		var description string

		key = fmt.Sprintf("%s", info.Section)
		name = prometheus.BuildFQName(namespace, "info", strings.Replace(info.Section, "%", "pct", -1))

		if info.Description == "" {
			description = strings.Replace(info.Section, "_", " ", -1)
		} else {
			description = info.Description
		}

		description = description + " in " + info.Unit

		infos[key] = prometheus.NewDesc(
			name,
			description,
			labels, nil,
		)
	}

	return infos
}
