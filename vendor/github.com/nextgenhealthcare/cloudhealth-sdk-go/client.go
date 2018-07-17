// Package cloudhealth is a wrapper for the CloudHealth API.
package cloudhealth

import (
	"errors"
	"net/url"
)

// Client communicates with the CloudHealth API.
type Client struct {
	ApiKey      string
	EndpointURL *url.URL
}

// ErrClientAuthenticationError is returned for authentication errors with the API.
var ErrClientAuthenticationError = errors.New("Authentication Error with CloudHealth")

// NewClient returns a new cloudhealth.Client for accessing the CloudHealth API.
func NewClient(apiKey string, defaultEndpointURL string) (*Client, error) {
	s := &Client{
		ApiKey: apiKey,
	}
	endpointURL, err := url.Parse(defaultEndpointURL)
	if err != nil {
		return nil, err
	}
	s.EndpointURL = endpointURL
	return s, nil
}
