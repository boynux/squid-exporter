package types

/*Counter maps a squid conters */
type Counter struct {
	Key   string
	Value float64
}

/*Counters is a list of multiple squid counters */
type Counters []Counter
