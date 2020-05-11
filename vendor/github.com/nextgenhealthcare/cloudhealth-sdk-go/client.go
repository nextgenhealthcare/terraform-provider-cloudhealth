// Package cloudhealth is a wrapper for the CloudHealth API.
package cloudhealth

import (
	"errors"
	"net/url"
)

var defaultTimeout int = 15

// Client communicates with the CloudHealth API.
type Client struct {
	ApiKey      string
	EndpointURL *url.URL
	Timeout     int
}

// ErrClientAuthenticationError is returned for authentication errors with the API.
var ErrClientAuthenticationError = errors.New("Authentication Error with CloudHealth")

// NewClient returns a new cloudhealth.Client for accessing the CloudHealth API.
func NewClient(apiKey string, defaultEndpointURL string, timeout ...int) (*Client, error) {
	s := &Client{
		ApiKey: apiKey,
	}
	endpointURL, err := url.Parse(defaultEndpointURL)
	if err != nil {
		return nil, err
	}
	s.EndpointURL = endpointURL
	s.Timeout = defaultTimeout
	if len(timeout) > 0 {
		s.Timeout = timeout[0]
	}
	return s, nil
}
