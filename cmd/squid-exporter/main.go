package main

import (
	"fmt"
	"log"
	"squid-exporter/collector"
)

func main() {
	c := collector.NewCacheObjectClient("localhost", 3129)

	counters, err := c.GetCounters()

	if err != nil {
		log.Fatal(err)
	}

	for i := range counters {
		fmt.Printf("%s: %f\n", counters[i].Key, counters[i].Value)
	}

	return
}
