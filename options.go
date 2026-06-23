package siyuan

import (
	"fmt"
	"net/http"
	"strings"
)

type clientOptions struct {
	endpoint   string
	token      string
	httpClient *http.Client
	userAgent  string
	headers    http.Header
}

// Option configures a Client.
type Option func(*clientOptions) error

// WithEndpoint configures the SiYuan API endpoint.
func WithEndpoint(endpoint string) Option {
	return func(opts *clientOptions) error {
		if strings.TrimSpace(endpoint) == "" {
			return fmt.Errorf("endpoint is empty")
		}
		opts.endpoint = endpoint
		return nil
	}
}

// WithToken configures the SiYuan API token.
func WithToken(token string) Option {
	return func(opts *clientOptions) error {
		opts.token = token
		return nil
	}
}

// WithHTTPClient configures the HTTP client used by the SDK.
func WithHTTPClient(client *http.Client) Option {
	return func(opts *clientOptions) error {
		if client == nil {
			return fmt.Errorf("http client is nil")
		}
		opts.httpClient = client
		return nil
	}
}

// WithUserAgent configures the User-Agent header.
func WithUserAgent(userAgent string) Option {
	return func(opts *clientOptions) error {
		opts.userAgent = userAgent
		return nil
	}
}

// WithHeader configures an additional header sent on every request.
func WithHeader(name string, value string) Option {
	return func(opts *clientOptions) error {
		name = strings.TrimSpace(name)
		if name == "" {
			return fmt.Errorf("header name is empty")
		}
		opts.headers.Add(name, value)
		return nil
	}
}
