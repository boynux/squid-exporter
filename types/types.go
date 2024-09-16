package types

// VarLabel maps key value Prometheus labels.
type VarLabel struct {
	Key   string
	Value string
}

// Counter maps a Squid counter.
type Counter struct {
	Key       string
	Value     float64
	VarLabels []VarLabel
}

// Counters is a list of multiple Squid counters.
type Counters []Counter
