package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-openapi/errors"
)

type Counter struct {
	Key   string
	Value float64
}

func main() {
	conn, err := net.Dial("tcp", "localhost:3129")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(conn, "GET cache_object://localhost/counters HTTP/1.0\r\nnHost: localhost\r\nUser-Agent: squidclient/3.5.12\r\nAccept: */*\r\n\r\n")
	r, err := http.ReadResponse(bufio.NewReader(conn), nil)

	fmt.Println(r.Status)
	reader := bufio.NewReader(r.Body)

	for {
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		c, err := decodeCounterStrings(line)

		if err == nil {
			fmt.Printf("%s: %f\n", c.Key, c.Value)
		}
	}
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}

func decodeCounterStrings(line string) (*Counter, error) {
	if equal := strings.Index(line, "="); equal >= 0 {
		if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
			value := ""
			if len(line) > equal {
				value = strings.TrimSpace(line[equal+1:])
			}

			if i, err := strconv.ParseFloat(value, 64); err == nil {
				return &Counter{key, i}, nil
			}
		}
	}

	return nil, errors.New(1, "could not parse line")
}
