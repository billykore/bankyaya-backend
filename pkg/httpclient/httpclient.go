package httpclient

import (
	"net"
	"net/http"
	"time"
)

// New initializes and returns a custom-configured http.Client.
func New() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second, // Global timeout for HTTP requests
		Transport: &http.Transport{
			// Limit the maximum number of open connections
			MaxIdleConns:        100, // Max idle connections across all hosts
			MaxIdleConnsPerHost: 10,  // Max idle connections to the same host
			MaxConnsPerHost:     20,  // Max total connections to the same host

			// Reuse idle connections
			IdleConnTimeout: 90 * time.Second, // Close idle connections after timeout

			// Connection settings
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second, // Connection timeout
				KeepAlive: 30 * time.Second,
			}).DialContext,

			// Enable additional features like HTTP2
			ForceAttemptHTTP2: true, // Enforce use of HTTP/2 where applicable
		},
	}
}
