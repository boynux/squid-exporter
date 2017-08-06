package main

import (
	"fmt"
	"log"
	"squid-exporter/collector"
	"squid-exporter/types"
	"strconv"
	"strings"

	"github.com/go-openapi/errors"
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

func decodeCounterStrings(line string) (*types.Counter, error) {
	if equal := strings.Index(line, "="); equal >= 0 {
		if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
			value := ""
			if len(line) > equal {
				value = strings.TrimSpace(line[equal+1:])
			}

			if i, err := strconv.ParseFloat(value, 64); err == nil {
				return &types.Counter{key, i}, nil
			}
		}
	}

	return nil, errors.New(1, "could not parse line")
}
