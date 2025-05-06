package collector

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/boynux/squid-exporter/types"
)

// CacheObjectClient holds information about Squid manager.
type CacheObjectClient struct {
	baseURL     string
	username    string
	password    string
	proxyHeader string
}

// SquidClient provides functionality to fetch Squid metrics.
type SquidClient interface {
	GetCounters() (types.Counters, error)
	GetServiceTimes() (types.Counters, error)
	GetInfos() (types.Counters, error)
}

// NewCacheObjectClient initializes a new cache client.
func NewCacheObjectClient(baseURL, username, password, proxyHeader string) *CacheObjectClient {
	return &CacheObjectClient{
		baseURL,
		username,
		password,
		proxyHeader,
	}
}

func (c *CacheObjectClient) readFromSquid(endpoint string) (*bufio.Reader, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-success code %d while fetching metrics", resp.StatusCode)
	}

	return bufio.NewReader(resp.Body), err
}

// GetCounters fetches counters from Squid cache manager.
func (c *CacheObjectClient) GetCounters() (types.Counters, error) {
	var counters types.Counters

	reader, err := c.readFromSquid("counters")
	if err != nil {
		return nil, fmt.Errorf("error getting counters: %w", err)
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		if c, err := decodeCounterStrings(scanner.Text()); err != nil {
			log.Println(err)
		} else {
			counters = append(counters, c)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return counters, err
}

// GetServiceTimes fetches service times from Squid cache manager.
func (c *CacheObjectClient) GetServiceTimes() (types.Counters, error) {
	var serviceTimes types.Counters

	reader, err := c.readFromSquid("service_times")
	if err != nil {
		return nil, fmt.Errorf("error getting service times: %w", err)
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		if s, err := decodeServiceTimeStrings(scanner.Text()); err != nil {
			log.Println(err)
		} else if s.Key != "" {
			serviceTimes = append(serviceTimes, s)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return serviceTimes, err
}

// GetInfos fetches info from Squid cache manager.
func (c *CacheObjectClient) GetInfos() (types.Counters, error) {
	var infos types.Counters

	reader, err := c.readFromSquid("info")
	if err != nil {
		return nil, fmt.Errorf("error getting info: %w", err)
	}

	infoVarLabels := types.Counter{Key: "squid_info", Value: 1}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		if dis, err := decodeInfoStrings(scanner.Text()); err != nil {
			log.Println(err)
		} else {
			if len(dis.VarLabels) > 0 {
				if dis.VarLabels[0].Key == "5min" {
					var infoAvg5, infoAvg60 types.Counter

					infoAvg5.Key = dis.Key + "_" + dis.VarLabels[0].Key
					infoAvg60.Key = dis.Key + "_" + dis.VarLabels[1].Key

					if value, err := strconv.ParseFloat(dis.VarLabels[0].Value, 64); err == nil {
						infoAvg5.Value = value
						infos = append(infos, infoAvg5)
					}
					if value, err := strconv.ParseFloat(dis.VarLabels[1].Value, 64); err == nil {
						infoAvg60.Value = value
						infos = append(infos, infoAvg60)
					}
				} else {
					infoVarLabels.VarLabels = append(infoVarLabels.VarLabels, dis.VarLabels[0])
				}
			} else if dis.Key != "" {
				infos = append(infos, dis)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	infos = append(infos, infoVarLabels)
	return infos, err
}

func decodeCounterStrings(line string) (types.Counter, error) {
	if equal := strings.Index(line, "="); equal >= 0 {
		if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
			value := ""
			if len(line) > equal {
				value = strings.TrimSpace(line[equal+1:])
			}

			// Remove additional formating string from `sample_time`
			if slices := strings.Split(value, " "); len(slices) > 0 {
				value = slices[0]
			}

			if i, err := strconv.ParseFloat(value, 64); err == nil {
				return types.Counter{Key: key, Value: i}, nil
			}
		}
	}

	return types.Counter{}, errors.New("counter - could not parse line: " + line)
}

func decodeServiceTimeStrings(line string) (types.Counter, error) {
	if strings.HasSuffix(line, ":") { // A header line isn't a metric
		return types.Counter{}, nil
	}
	if equal := strings.Index(line, ":"); equal >= 0 {
		if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
			value := ""
			if len(line) > equal {
				value = strings.TrimSpace(line[equal+1:])
			}
			key = strings.Replace(key, " ", "_", -1)
			key = strings.Replace(key, "(", "", -1)
			key = strings.Replace(key, ")", "", -1)

			if equalTwo := strings.Index(value, "%"); equalTwo >= 0 {
				if keyTwo := strings.TrimSpace(value[:equalTwo]); len(keyTwo) > 0 {
					if len(value) > equalTwo {
						value = strings.Split(strings.TrimSpace(value[equalTwo+1:]), " ")[0]
					}
					key = key + "_" + keyTwo
				}
			}

			if value, err := strconv.ParseFloat(value, 64); err == nil {
				return types.Counter{Key: key, Value: value}, nil
			}
		}
	}

	return types.Counter{}, errors.New("service times - could not parse line: " + line)
}

func decodeInfoStrings(line string) (types.Counter, error) {
	if strings.HasSuffix(line, ":") { // A header line isn't a metric
		return types.Counter{}, nil
	}

	if idx := strings.Index(line, ":"); idx >= 0 { // detect if line contain metric format like "metricName: value"
		if key := strings.TrimSpace(line[:idx]); len(key) > 0 {
			value := ""
			if len(line) > idx {
				value = strings.TrimSpace(line[idx+1:])
			}
			key = strings.Replace(key, " ", "_", -1)
			key = strings.Replace(key, "(", "", -1)
			key = strings.Replace(key, ")", "", -1)
			key = strings.Replace(key, ",", "", -1)
			key = strings.Replace(key, "/", "", -1)

			// metrics with value as string need to save as label, format like "Squid Object Cache: Version 6.1" (the 3 first metrics)
			if key == "Squid_Object_Cache" || key == "Build_Info" || key == "Service_Name" {
				if key == "Squid_Object_Cache" { // To clarify that the value is the squid version.
					key += "_Version"
					if slices := strings.Split(value, " "); len(slices) > 0 {
						value = slices[1]
					}
				}
				var infoVarLabel types.VarLabel
				infoVarLabel.Key = key
				infoVarLabel.Value = value

				var infoCounter types.Counter
				infoCounter.Key = key
				infoCounter.VarLabels = append(infoCounter.VarLabels, infoVarLabel)
				return infoCounter, nil
			} else if key == "Start_Time" || key == "Current_Time" { // discart this metrics
				return types.Counter{}, nil
			}

			// Remove additional information in value metric
			if slices := strings.Split(value, " "); len(slices) > 0 {
				if slices[0] == "5min:" && slices[2] == "60min:" { // catch metrics with avg in 5min and 60min format like "Hits as % of bytes sent: 5min: -0.0%, 60min: -0.0%"
					var infoAvg5mVarLabel types.VarLabel
					infoAvg5mVarLabel.Key = slices[0]
					infoAvg5mVarLabel.Value = slices[1]

					infoAvg5mVarLabel.Key = strings.Replace(infoAvg5mVarLabel.Key, ":", "", -1)
					infoAvg5mVarLabel.Value = strings.Replace(infoAvg5mVarLabel.Value, "%", "", -1)
					infoAvg5mVarLabel.Value = strings.Replace(infoAvg5mVarLabel.Value, ",", "", -1)

					var infoAvg60mVarLabel types.VarLabel
					infoAvg60mVarLabel.Key = slices[2]
					infoAvg60mVarLabel.Value = slices[3]

					infoAvg60mVarLabel.Key = strings.Replace(infoAvg60mVarLabel.Key, ":", "", -1)
					infoAvg60mVarLabel.Value = strings.Replace(infoAvg60mVarLabel.Value, "%", "", -1)
					infoAvg60mVarLabel.Value = strings.Replace(infoAvg60mVarLabel.Value, ",", "", -1)

					var infoAvgCounter types.Counter
					infoAvgCounter.Key = key
					infoAvgCounter.VarLabels = append(infoAvgCounter.VarLabels, infoAvg5mVarLabel, infoAvg60mVarLabel)

					return infoAvgCounter, nil
				}

				value = slices[0]
			}

			value = strings.Replace(value, "%", "", -1)
			value = strings.Replace(value, ",", "", -1)

			if i, err := strconv.ParseFloat(value, 64); err == nil {
				return types.Counter{Key: key, Value: i}, nil
			}
		}
	} else {
		// this catch the last 4 metrics format like "value metricName"
		lineTrimed := strings.TrimSpace(line)

		if idx := strings.Index(lineTrimed, " "); idx >= 0 {
			key := strings.TrimSpace(lineTrimed[idx+1:])
			key = strings.Replace(key, " ", "_", -1)
			key = strings.Replace(key, "-", "_", -1)

			value := strings.TrimSpace(lineTrimed[:idx])

			if i, err := strconv.ParseFloat(value, 64); err == nil {
				return types.Counter{Key: key, Value: i}, nil
			}
		}
	}

	return types.Counter{}, errors.New("Info - could not parse line: " + line)
}
