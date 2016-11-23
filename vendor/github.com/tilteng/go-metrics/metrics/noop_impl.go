package metrics

import (
	"errors"
	"time"
)

type noopClient struct {
	inited    bool
	tags      map[string]string
	namespace string
}

func (self *noopClient) GetAddr() string {
	return ""
}

func (self *noopClient) GetNamespace() string {
	return self.namespace
}

func (self *noopClient) SetNamespace(namespace string) {
	self.namespace = namespace
}

func (self *noopClient) GetTags() map[string]string {
	return self.tags
}

func (self *noopClient) SetTags(tags map[string]string) {
	self.tags = tags
}

func (self *noopClient) Init() error {
	if self.inited {
		return errors.New("Client has already been initialized")
	}
	self.inited = true
	return nil
}

func (self *noopClient) Gauge(name string, value float64, rate float64, tags map[string]string) error {
	return nil
}

func (self *noopClient) Count(name string, value int64, rate float64, tags map[string]string) error {
	return nil
}

func (self *noopClient) Histogram(name string, value float64, rate float64, tags map[string]string) error {
	return nil
}

func (self *noopClient) Decr(name string, rate float64, tags map[string]string) error {
	return nil
}

func (self *noopClient) Incr(name string, rate float64, tags map[string]string) error {
	return nil
}

func (self *noopClient) Set(name string, value string, rate float64, tags map[string]string) error {
	return nil
}

func (self *noopClient) Timing(name string, value time.Duration, rate float64, tags map[string]string) error {
	return nil
}

func (self *noopClient) TimingMS(name string, value float64, rate float64, tags map[string]string) error {
	return nil
}

func NewNOOPClient() *noopClient {
	return &noopClient{
		tags: make(map[string]string),
	}
}
