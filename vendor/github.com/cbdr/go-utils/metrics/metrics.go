package metrics

import (
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

// Default is the default Client implementation
// we use a singleton to avoid injecting it everywhere
var Default Client = nullClient{}

// Client is a simplified datadog client decorated with a WithTags method
type Client interface {
	WithTags(tags []string) Client
	WithTag(tag string) Client
	Gauge(name string, value float64) error
	Incr(name string) error
	Histogram(name string, value float64) error
	Timing(name string, start time.Time) error
}

type client struct {
	datadog *statsd.Client
	tags    []string
}

// New creates a new datadog client
func New(addr, namespace string) (Client, error) {
	datadog, err := statsd.New(addr)
	if err != nil {
		return nil, err
	}
	datadog.Namespace = namespace + "."

	return &client{
		datadog: datadog,
		tags:    []string{},
	}, nil
}

// WithTags returns a new client with default tag values
func (c *client) WithTags(tags []string) Client {
	newClient := *c
	for _, tag := range tags {
		newClient.tags = append(newClient.tags, tag)
	}
	return &newClient
}

// WithTag returns a new client with a default tag value
func (c *client) WithTag(tag string) Client {
	return c.WithTags([]string{tag})
}

func (c *client) Gauge(name string, value float64) error {
	return c.datadog.Gauge(name, value, c.tags, 1.0)
}

func (c *client) Incr(name string) error {
	return c.datadog.Incr(name, c.tags, 1.0)
}

func (c *client) Histogram(name string, value float64) error {
	return c.datadog.Histogram(name, value, c.tags, 1.0)
}

func (c *client) Timing(name string, start time.Time) error {
	return c.datadog.Timing(name, time.Duration(time.Now().Sub(start)), c.tags, 1.0)
}

// Null implementation
type nullClient struct{}

func (c nullClient) WithTags(tags []string) Client {
	return c
}

func (c nullClient) WithTag(tag string) Client {
	return c
}

func (nullClient) Gauge(name string, value float64) error {
	return nil
}

func (nullClient) Incr(name string) error {
	return nil
}

func (nullClient) Histogram(name string, value float64) error {
	return nil
}

func (nullClient) Timing(name string, start time.Time) error {
	return nil
}

// WithTags calls WithTags on the default client
func WithTags(tags []string) Client {
	return Default.WithTags(tags)
}

// WithTag calls WithTag on the default client
func WithTag(tag string) Client {
	return Default.WithTag(tag)
}

// Gauge calls Gauge on the default client
func Gauge(name string, value float64) error {
	return Default.Gauge(name, value)
}

// Incr calls Incr on the default client
func Incr(name string) error {
	return Default.Incr(name)
}

// Histogram calls Histogram on the default client
func Histogram(name string, value float64) error {
	return Default.Histogram(name, value)
}

// Timing calls Timing on the default client
func Timing(name string, start time.Time) error {
	return Default.Timing(name, start)
}
