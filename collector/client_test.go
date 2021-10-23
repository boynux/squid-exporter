package collector

import (
	"net"
	"testing"

	"github.com/boynux/squid-exporter/types"
	"github.com/stretchr/testify/assert"
)

type mockConnectionHandler struct {
	server net.Conn

	buffer []byte
}

func (c *mockConnectionHandler) connect() (net.Conn, error) {
	var client net.Conn
	c.server, client = net.Pipe()

	return client, nil
}

func TestBuildBasicAuth(t *testing.T) {
	u := "test_username"
	p := "test_password"
	expectedAuthString := "dGVzdF91c2VybmFtZTp0ZXN0X3Bhc3N3b3Jk"

	ba := buildBasicAuthString(u, p)

	assert.Equal(t, expectedAuthString, ba, "Basic Auth format doesn't match")
}

func TestReadFromSquid(t *testing.T) {
	ch := &mockConnectionHandler{}

	go func() {
		b := make([]byte, 256)
		n, _ := ch.server.Read(b)
		ch.buffer = append(ch.buffer, b[:n]...)

		ch.server.Write(b[n:])
		ch.server.Close()
	}()

	coc := &CacheObjectClient{
		ch,
		"",
		[]string{},
	}
	expected := "GET cache_object://localhost/test HTTP/1.0\r\nHost: localhost\r\nUser-Agent: squidclient/3.5.12\r\nAccept: */*\r\n\r\n"
	coc.readFromSquid("test")

	assert.Equal(t, expected, string(ch.buffer))
}

func TestDecodeMetricStrings(t *testing.T) {
	tests := []struct {
		s string
		c types.Counter
		e string
		d func(string) (types.Counter, error)
	}{
		{"swap.files_cleaned=1", types.Counter{Key: "swap.files_cleaned", Value: 1}, "", decodeCounterStrings},
		{"client.http_requests=1", types.Counter{Key: "client.http_requests", Value: 1}, "", decodeCounterStrings},
		{"# test for invalid metric line", types.Counter{}, "counter - could not parse line: # test for invalid metric line", decodeCounterStrings},

		{"	HTTP Requests (All):  70%   10.00000  9.50000\n", types.Counter{Key: "HTTP_Requests_All_70", Value: 10}, "", decodeServiceTimeStrings},
		{"	Not-Modified Replies:  5%   12.00000  10.00000\n", types.Counter{Key: "Not-Modified_Replies_5", Value: 12}, "", decodeServiceTimeStrings},
		{"	ICP Queries:          85%   900.00000  1200.00000\n", types.Counter{Key: "ICP_Queries_85", Value: 900}, "", decodeServiceTimeStrings},
	}

	for _, tc := range tests {
		c, err := tc.d(tc.s)

		if tc.e != "" {
			assert.EqualError(t, err, tc.e)
		}
		assert.Equal(t, tc.c, c)
	}
}
