package metrics

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

type ddClient struct {
	client              *statsd.Client
	inited              bool
	tags                map[string]string
	addr                string
	namespace           string
	numBufferedCommands int
}

func (self *ddClient) GetNamespace() string {
	return self.namespace
}

func (self *ddClient) SetNamespace(namespace string) {
	if self.client != nil {
		self.client.Namespace = namespace
	}
	self.namespace = namespace
}

func (self *ddClient) GetTags() map[string]string {
	return self.tags
}

func (self *ddClient) SetTags(tags map[string]string) {
	if self.client != nil {
		self.client.Tags = self.tagsToSlice(tags)
	}
	self.tags = tags
}

func (self *ddClient) Init() error {
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

func (self *ddClient) tagsToSlice(tags map[string]string) []string {
	sl := make([]string, 0, len(tags))
	for k, v := range tags {
		sl = append(sl, k+":"+v)
	}
	return sl
}

func (self *ddClient) Gauge(name string, value float64, rate float64, tags map[string]string) error {
	return self.client.Gauge(name, value, self.tagsToSlice(tags), rate)
}

func (self *ddClient) Count(name string, value int64, rate float64, tags map[string]string) error {
	return self.client.Count(name, value, self.tagsToSlice(tags), rate)
}

func (self *ddClient) Histogram(name string, value float64, rate float64, tags map[string]string) error {
	return self.client.Histogram(name, value, self.tagsToSlice(tags), rate)
}

func (self *ddClient) Decr(name string, rate float64, tags map[string]string) error {
	return self.client.Decr(name, self.tagsToSlice(tags), rate)
}

func (self *ddClient) Incr(name string, rate float64, tags map[string]string) error {
	return self.client.Incr(name, self.tagsToSlice(tags), rate)
}

func (self *ddClient) Set(name string, value string, rate float64, tags map[string]string) error {
	return self.client.Set(name, value, self.tagsToSlice(tags), rate)
}

func (self *ddClient) Timing(name string, value time.Duration, rate float64, tags map[string]string) error {
	return self.client.Timing(name, value, self.tagsToSlice(tags), rate)
}

func (self *ddClient) TimingMS(name string, value float64, rate float64, tags map[string]string) error {
	return self.client.TimeInMilliseconds(name, value, self.tagsToSlice(tags), rate)
}

func NewDDClient(addr string) (*ddClient, error) {
	url, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("Couldn't parse client address: %s", err)
	}

	if len(url.Scheme) > 0 && url.Scheme != "udp" {
		return nil, errors.New(
			"Unsupported scheme '" + url.Scheme + "'. Datadog requires udp.",
		)
	}

	addr = url.Host

	if len(addr) == 0 {
		addr = "172.0.0.1:8125"
	} else if !strings.Contains(addr, ":") {
		addr += ":8125"
	}

	return &ddClient{
		addr: addr,
		tags: make(map[string]string),
	}, nil
}
