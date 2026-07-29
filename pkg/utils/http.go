package utils

import (
	"crypto/tls"
	"net/http"
	"time"
)

// NewHTTPClient creates a new HTTP client with custom settings
func NewHTTPClient(timeout time.Duration, skipVerify bool) *http.Client {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	if skipVerify {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	return &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}
}

// AddBasicAuth adds Basic Authentication header to request
func AddBasicAuth(req *http.Request, username, password string) {
	req.SetBasicAuth(username, password)
}
