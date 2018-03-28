package oembed

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/dyatlov/go-oembed/oembed"
	log "github.com/sirupsen/logrus"
)

// LinkType repensents a type of embedded link
type LinkType string

// Known link types
// We retrict this list to providers actually used in this project
// The code should be able to gracefuly handle other values though
const (
	LinkTypePhoto LinkType = "photo"
	LinkTypeVideo LinkType = "video"
)

// Provider represents an oEmbed provider (Flickr, Vimeo, etc)
type Provider string

// Known provider list
// We retrict this list to providers actually used in this project
// The code should be able to gracefuly handle other values though
const (
	ProviderFlickr Provider = "Flickr"
	ProviderVimeo  Provider = "Vimeo"
)

// Link reprensents the result of an oEmbed query
// We do not implement all possible properties here but only the ones used in this project
type Link struct {
	URL        string    `json:"url"`
	Type       LinkType  `json:"type"`
	Provider   Provider  `json:"provider"`
	Title      string    `json:"title"`
	AuthorName string    `json:"author_name"`
	Width      StringInt `json:"width"`
	Height     StringInt `json:"height"`
	Duration   int       `json:"duration"`
}

// Fetcher uses the oEmbed protocol to fetch properties of a link
type Fetcher interface {
	Fetch(rawURL string) (*Link, error)
}

// ProvidersURL is the url where the providers list is located
const ProvidersURL = "https://oembed.com/providers.json"

// Default fetcher implementation
// based on the https://github.com/dyatlov/go-oembed library
// We confine this dependency here. It should not be referenced outside this package
type fetcher struct {
	oe     *oembed.Oembed
	logger log.FieldLogger
}

// NewFetcher returns a default fetch implementation
// This function loads the provider list from an URL
// For this reason it might fail
func NewFetcher(logger log.FieldLogger) (Fetcher, error) {
	providers, err := getProviders()
	if err != nil {
		return nil, err
	}

	oe := oembed.NewOembed()
	oe.ParseProviders(providers)

	return &fetcher{
		oe:     oe,
		logger: logger,
	}, nil
}

func getProviders() (io.Reader, error) {
	req, err := http.NewRequest("GET", ProvidersURL, nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res.Body, err
}

// Fetch implements the Fetcher interface
func (f *fetcher) Fetch(rawURL string) (*Link, error) {
	item := f.oe.FindItem(rawURL)
	if item == nil {
		return nil, &NotFoundError{err: errors.New("URL not found")}
	}

	// We're interrested in getting the duration field which is non-standard
	// it is not managed by the library so we have to handle the rest of the process manually
	// (but still using the library to parse url shemes)
	fullURL := fmt.Sprintf("%s?format=json&url=%s", item.EndpointURL, escapeURL(rawURL))
	f.logger.WithField("url", fullURL).Info("fetching URL...")

	body, err := apiCall(fullURL)
	if err != nil {
		return nil, err
	}

	f.logger.WithField("body", string(body)).Debug("provider's response")

	var l Link
	if err := json.Unmarshal(body, &l); err != nil {
		return nil, err
	}

	return &l, nil
}

func escapeURL(rawURL string) string {
	// Escape the URL
	u, err := url.Parse(rawURL)
	if err != nil {
		// go is very picky about compliance of URL formats
		// it might work even if the URL is not escaped. It worth trying
		return rawURL
	}
	return u.String()
}

func apiCall(fullURL string) ([]byte, error) {
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	switch {
	case res.StatusCode == 404:
		return nil, &NotFoundError{err: errors.New("URL not found")}
	case res.StatusCode >= 300:
		// TODO: error management is super basic here. Should be improved
		return nil, fmt.Errorf("provider returned a %d status code", res.StatusCode)
	}

	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

// NotFoundError is retured when no info was found for this URL
type NotFoundError struct {
	err error
}

// Error implements the Error interface
func (err *NotFoundError) Error() string {
	return err.err.Error()
}
