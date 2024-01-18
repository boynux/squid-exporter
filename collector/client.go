package collector

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/boynux/squid-exporter/types"
)

/*CacheObjectClient holds information about squid manager */
type CacheObjectClient struct {
	ch              connectionHandler
	basicAuthString string
	headers         []string
}

type connectionHandler interface {
	connect() (net.Conn, error)
}

type connectionHandlerImpl struct {
	hostname string
	port     int
}

/*SquidClient provides functionality to fetch squid metrics */
type SquidClient interface {
	GetCounters() (types.Counters, error)
	GetServiceTimes() (types.Counters, error)
	GetInfos() (types.Counters, error)
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

type CacheObjectRequest struct {
	Hostname string
	Port     int
	Login    string
	Password string
	Headers  []string
}

/*NewCacheObjectClient initializes a new cache client */
func NewCacheObjectClient(cor *CacheObjectRequest) *CacheObjectClient {
	return &CacheObjectClient{
		&connectionHandlerImpl{
			cor.Hostname,
			cor.Port,
		},
		buildBasicAuthString(cor.Login, cor.Password),
		cor.Headers,
	}
}

func (c *CacheObjectClient) readFromSquid(endpoint string) (*bufio.Reader, error) {
	conn, err := c.ch.connect()

	if err != nil {
		return nil, err
	}
	r, err := get(conn, endpoint, c.basicAuthString, c.headers)

	if err != nil {
		return nil, err
	}

	if r.StatusCode != 200 {
		return nil, fmt.Errorf("Non success code %d while fetching metrics", r.StatusCode)
	}

	return bufio.NewReader(r.Body), err
}

func readLines(reader *bufio.Reader, lines chan<- string) {
	for {
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("error reading from the bufio.Reader: %v", err)
			break
		}

		lines <- line
	}
	close(lines)
}

/*GetCounters fetches counters from squid cache manager */
func (c *CacheObjectClient) GetCounters() (types.Counters, error) {
	var counters types.Counters

	reader, err := c.readFromSquid("counters")
	if err != nil {
		return nil, fmt.Errorf("error getting counters: %v", err)
	}

	lines := make(chan string)
	go readLines(reader, lines)

	for line := range lines {
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
	var serviceTimes types.Counters

	reader, err := c.readFromSquid("service_times")
	if err != nil {
		return nil, fmt.Errorf("error getting service times: %v", err)
	}

	lines := make(chan string)
	go readLines(reader, lines)

	for line := range lines {
		s, err := decodeServiceTimeStrings(line)
		if err != nil {
			log.Println(err)
		} else {
			if s.Key != "" {
				serviceTimes = append(serviceTimes, s)
			}
		}
	}

	return serviceTimes, err
}

/*GetInfos fetches info from squid cache manager */
func (c *CacheObjectClient) GetInfos() (types.Counters, error) {
	var infos types.Counters

	reader, err := c.readFromSquid("info")
	if err != nil {
		return nil, fmt.Errorf("error getting info: %v", err)
	}

	lines := make(chan string)
	go readLines(reader, lines)

	var infoVarLabels types.Counter
	infoVarLabels.Key = "squid_info"
	infoVarLabels.Value = 1

	for line := range lines {
		i, err := decodeInfoStrings(line)
		if err != nil {
			log.Println(err)
		} else {
			//if i.Key == "Squid_Object_Cache_Version" || i.Key == "Build_Info" || i.Key == "Service_Name" || i.Key == "Start_Time" || i.Key == "Current_Time" {
			if len(i.VarLabels) > 0 {
				infoVarLabels.VarLabels = append(infoVarLabels.VarLabels, i.VarLabels[0])
			} else if i.Key != "" {
				infos = append(infos, i)
			}
		}
	}
	infos = append(infos, infoVarLabels)
	return infos, err
}

func (ch *connectionHandlerImpl) connect() (net.Conn, error) {
	return net.Dial("tcp", fmt.Sprintf("%s:%d", ch.hostname, ch.port))
}

func get(conn net.Conn, path string, basicAuthString string, headers []string) (*http.Response, error) {
	rBody := append(headers, []string{
		fmt.Sprintf(requestProtocol, path),
		"Host: localhost",
		"User-Agent: squidclient/3.5.12",
	}...)

	if len(basicAuthString) > 0 {
		rBody = append(rBody, "Proxy-Authorization: Basic "+basicAuthString)
		rBody = append(rBody, "Authorization: Basic "+basicAuthString)
	}
	rBody = append(rBody, "Accept: */*", "\r\n")
	request := strings.Join(rBody, "\r\n")

	fmt.Fprint(conn, request)

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
	if strings.HasSuffix(line, ":\n") { // A header line isn't a metric
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
	if strings.HasSuffix(line, ":\n") { // A header line isn't a metric
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
			key = strings.Replace(key, ",", "", -1)
			key = strings.Replace(key, "/", "", -1)

			// metrics with string save as label
			if key == "Squid_Object_Cache" || key == "Build_Info" || key == "Service_Name" || key == "Start_Time" || key == "Current_Time" {
				if key == "Squid_Object_Cache" {
					key = key + "_Version"
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
			}

			// Remove additional formating string
			if slices := strings.Split(value, " "); len(slices) > 0 {
				if slices[0] == "5min:" {
					value = slices[1]
				} else {
					value = slices[0]
				}

			}

			value = strings.Replace(value, "%", "", -1)
			value = strings.Replace(value, ",", "", -1)

			if i, err := strconv.ParseFloat(value, 64); err == nil {
				return types.Counter{Key: key, Value: i}, nil
			}
		}
	} else {
		lineTrimed := strings.TrimSpace(line[:])

		if equal := strings.Index(lineTrimed, " "); equal >= 0 {
			key := strings.TrimSpace(lineTrimed[equal+1:])
			key = strings.Replace(key, " ", "_", -1)
			key = strings.Replace(key, "-", "_", -1)

			value := strings.TrimSpace(lineTrimed[:equal])

			if i, err := strconv.ParseFloat(value, 64); err == nil {
				return types.Counter{Key: key, Value: i}, nil
			}
		}
	}

	return types.Counter{}, errors.New("Info - could not parse line: " + line)
}
