package collector

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/boynux/squid-exporter/types"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

/*CacheObjectClient holds information about squid manager */
type CacheObjectClient struct {
	hostname        string
	port            int
	basicAuthString string
	headers         map[string]string
}

/*SquidClient provides functionality to fetch squid metrics */
type SquidClient interface {
	GetCounters() (types.Counters, error)
	GetServiceTimes() (types.Counters, error)
}

const (
	requestProtocol = "GET cache_object://localhost/%s HTTP/1.0"
)

func buildBasicAuthString(login string, password string) string {
	if len(login) == 0 {
		return ""
	} else {
		return base64.StdEncoding.EncodeToString([]byte(login + ":" + password))
	}
}

/*NewCacheObjectClient initializes a new cache client */
func NewCacheObjectClient(hostname string, port int, login string, password string) *CacheObjectClient {
	return &CacheObjectClient{
		hostname,
		port,
		buildBasicAuthString(login, password),
		map[string]string{},
	}
}

/*GetCounters fetches counters from squid cache manager */
func (c *CacheObjectClient) GetCounters() (types.Counters, error) {
	conn, err := connect(c.hostname, c.port)

	if err != nil {
		return types.Counters{}, err
	}

	r, err := get(conn, "counters", c.basicAuthString)

	if err != nil {
		return nil, err
	}

	if r.StatusCode != 200 {
		return nil, fmt.Errorf("Non success code %d while fetching metrics", r.StatusCode)
	}

	var counters types.Counters

	// TODO: Move to another func
	reader := bufio.NewReader(r.Body)

	for {
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		c, err := decodeCounterStrings(line)
		if err != nil {
			log.Println(err)
		} else {
			counters = append(counters, c)
		}
	}

	return counters, err
}

/*GetServiceTimes fetches service times from squid cache manager */
func (c *CacheObjectClient) GetServiceTimes() (types.Counters, error) {
	conn, err := connect(c.hostname, c.port)

	if err != nil {
		return types.Counters{}, err
	}

	r, err := get(conn, "service_times", c.basicAuthString)

	if err != nil {
		return nil, err
	}

	if r.StatusCode != 200 {
		return nil, fmt.Errorf("Non success code %d while fetching service_times metrics", r.StatusCode)
	}

	var serviceTimes types.Counters

	reader := bufio.NewReader(r.Body)

	for {
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		c, err := decodeServiceTimeStrings(line)
		if err != nil {
			log.Println(err)
		} else {
			serviceTimes = append(serviceTimes, c)
		}
	}

	return serviceTimes, err
}

func connect(hostname string, port int) (net.Conn, error) {
	return net.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port))
}

func get(conn net.Conn, path string, basicAuthString string) (*http.Response, error) {
	rBody := []string{
		fmt.Sprintf(requestProtocol, path),
		"Host: localhost",
		"User-Agent: squidclient/3.5.12",
	}
	if len(basicAuthString) > 0 {
		rBody = append(rBody, "Proxy-Authorization: Basic "+basicAuthString)
	}
	rBody = append(rBody, "Accept: */*", "\r\n")
	request := strings.Join(rBody, "\r\n")

	fmt.Fprintf(conn, request)
	return http.ReadResponse(bufio.NewReader(conn), nil)
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

	return types.Counter{}, errors.New("servicec times - could not parse line: " + line)
}
