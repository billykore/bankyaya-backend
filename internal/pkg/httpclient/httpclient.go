package httpclient

import (
	"net"
	"net/http"
	"time"
)

const (
	timeout             = 10 * time.Second
	maxIdleConns        = 100
	maxIdleConnsPerHost = 10
	maxConnsPerHost     = 20
	idleConnTimeout     = 90 * time.Second
	dialTimeout         = 5 * time.Second
	keepAlive           = 30 * time.Second
)

// New initializes and returns a custom-configured http.Client.
func New() *http.Client {
	return &http.Client{
		Timeout: timeout, // Global timeout for HTTP requests
		Transport: &http.Transport{
			// Limit the maximum number of open connections
			MaxIdleConns:        maxIdleConns,        // Max idle connections across all hosts
			MaxIdleConnsPerHost: maxIdleConnsPerHost, // Max idle connections to the same host
			MaxConnsPerHost:     maxConnsPerHost,     // Max total connections to the same host

			// Reuse idle connections
			IdleConnTimeout: idleConnTimeout, // Close idle connections after timeout

			// Connection settings
			DialContext: (&net.Dialer{
				Timeout:   dialTimeout, // Connection timeout
				KeepAlive: keepAlive,
			}).DialContext,

			// Enable additional features like HTTP2
			ForceAttemptHTTP2: true, // Enforce use of HTTP/2 where applicable
		},
	}
}
