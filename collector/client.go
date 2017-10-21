package collector

import (
	"bufio"
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
	hostname string
	port     int
	headers  map[string]string
}

/*SquidClient provides functionality to fetch squid metrics */
type SquidClient interface {
	GetCounters() (types.Counters, error)
}

const (
	requestProtocol = "GET cache_object://localhost/%s HTTP/1.0"
)

/*NewCacheObjectClient initializes a new cache client */
func NewCacheObjectClient(hostname string, port int) *CacheObjectClient {
	return &CacheObjectClient{
		hostname,
		port,
		map[string]string{},
	}
}

/*GetCounters fetches counters from squid cache manager */
func (c *CacheObjectClient) GetCounters() (types.Counters, error) {
	conn, err := connect(c.hostname, c.port)

	if err != nil {
		return types.Counters{}, err
	}

	r, err := get(conn, "counters")

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

func connect(hostname string, port int) (net.Conn, error) {
	return net.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port))
}

func get(conn net.Conn, path string) (*http.Response, error) {
	rBody := []string{
		fmt.Sprintf(requestProtocol, path),
		"Host: localhost",
		"User-Agent: squidclient/3.5.12",
		"Accept: */*",
		"\r\n",
	}

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

			if i, err := strconv.ParseFloat(value, 64); err == nil {
				return types.Counter{key, i}, nil
			}
		}
	}

	return types.Counter{}, errors.New("could not parse line")
}
