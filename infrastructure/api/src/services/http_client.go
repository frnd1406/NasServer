package services

import (
	"net/http"
	"time"
)

// NewSecureHTTPClient creates an HTTP client that automatically injects
// the X-Internal-Secret header into all requests.
func NewSecureHTTPClient(secret string, timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &secureTransport{
			secret: secret,
			base:   http.DefaultTransport,
		},
	}
}

type secureTransport struct {
	secret string
	base   http.RoundTripper
}

func (t *secureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.secret != "" {
		req.Header.Set("X-Internal-Secret", t.secret)
	}
	if t.base == nil {
		return http.DefaultTransport.RoundTrip(req)
	}
	return t.base.RoundTrip(req)
}
