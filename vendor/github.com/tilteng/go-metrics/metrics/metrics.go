package metrics

import "time"

type MetricsClient interface {
	GetNamespace() string
	SetNamespace(string)
	GetTags() map[string]string
	SetTags(map[string]string)
	Init() error
	Gauge(name string, value float64, rate float64, tags map[string]string) error
	Count(name string, value int64, rate float64, tags map[string]string) error
	Histogram(name string, value float64, rate float64, tags map[string]string) error
	Decr(name string, rate float64, tags map[string]string) error
	Incr(name string, rate float64, tags map[string]string) error
	Set(name string, value string, rate float64, tags map[string]string) error
	Timing(name string, value time.Duration, rate float64, tags map[string]string) error
	TimingMS(name string, value float64, rate float64, tags map[string]string) error
}

func NewMetricsClient(addr string) MetricsClient {
	if len(addr) == 0 {
		addr = "172.0.0.1:8125"
	}
	return NewDDClient(addr)
}
