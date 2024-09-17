package collector

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/boynux/squid-exporter/types"
)

func TestReadFromSquid(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.RequestURI, "/squid-internal-mgr/test")
	}))
	defer ts.Close()

	coc := &CacheObjectClient{
		ts.URL + "/squid-internal-mgr/",
		"",
		"",
		"",
	}

	_, err := coc.readFromSquid("test")
	if err != nil {
		t.Fatal(err)
	}
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
			require.EqualError(t, err, tc.e)
		}
		assert.Equal(t, tc.c, c)
	}
}
