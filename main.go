package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:3129")
	if err != nil {
		log.Fatal(err)
		// handle error
	}

	fmt.Fprintf(conn, "GET cache_object://localhost/counters HTTP/1.0\r\nnHost: localhost\r\nUser-Agent: squidclient/3.5.12\r\nAccept: */*\r\n\r\n")
	r, err := http.ReadResponse(bufio.NewReader(conn), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(r.Status)
}
