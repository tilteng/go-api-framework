package metrics

import (
	"errors"
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

type DDClient struct {
	client              *statsd.Client
	inited              bool
	tags                map[string]string
	addr                string
	namespace           string
	numBufferedCommands int
}

func (self *DDClient) GetNamespace() string {
	return self.namespace
}

func (self *DDClient) SetNamespace(namespace string) {
	if self.client != nil {
		self.client.Namespace = namespace
	}
	self.namespace = namespace
}

func (self *DDClient) GetTags() map[string]string {
	return self.tags
}

func (self *DDClient) SetTags(tags map[string]string) {
	if self.client != nil {
		self.client.Tags = self.tagsToSlice(tags)
	}
	self.tags = tags
}

func (self *DDClient) Init() error {
	if self.inited {
		return errors.New("Client has already been initialized")
	}

	var cli *statsd.Client
	var err error

	if self.numBufferedCommands > 0 {
		cli, err = statsd.NewBuffered(self.addr, self.numBufferedCommands)
	} else {
		cli, err = statsd.New(self.addr)
	}
	if err != nil {
		return err
	}
	cli.Namespace = self.namespace
	cli.Tags = self.tagsToSlice(self.tags)
	self.client = cli
	self.inited = true
	return nil
}

func (self *DDClient) tagsToSlice(tags map[string]string) []string {
	sl := make([]string, 0, len(tags))
	for k, v := range tags {
		sl = append(sl, k+":"+v)
	}
	return sl
}

func (self *DDClient) Gauge(name string, value float64, rate float64, tags map[string]string) error {
	return self.client.Gauge(name, value, self.tagsToSlice(tags), rate)
}

func (self *DDClient) Count(name string, value int64, rate float64, tags map[string]string) error {
	return self.client.Count(name, value, self.tagsToSlice(tags), rate)
}

func (self *DDClient) Histogram(name string, value float64, rate float64, tags map[string]string) error {
	return self.client.Histogram(name, value, self.tagsToSlice(tags), rate)
}

func (self *DDClient) Decr(name string, rate float64, tags map[string]string) error {
	return self.client.Decr(name, self.tagsToSlice(tags), rate)
}

func (self *DDClient) Incr(name string, rate float64, tags map[string]string) error {
	return self.client.Incr(name, self.tagsToSlice(tags), rate)
}

func (self *DDClient) Set(name string, value string, rate float64, tags map[string]string) error {
	return self.client.Set(name, value, self.tagsToSlice(tags), rate)
}

func (self *DDClient) Timing(name string, value time.Duration, rate float64, tags map[string]string) error {
	return self.client.Timing(name, value, self.tagsToSlice(tags), rate)
}

func (self *DDClient) TimingMS(name string, value float64, rate float64, tags map[string]string) error {
	return self.client.TimeInMilliseconds(name, value, self.tagsToSlice(tags), rate)
}

func NewDDClient(addr string) *DDClient {
	return &DDClient{
		addr: addr,
		tags: make(map[string]string),
	}
}
