package gojpcal

import (
	"fmt"
	"net/http"

	"github.com/syumai/workers/cloudflare/fetch"
)

func NewWorkerTransport() *WorkerTransport {
	return &WorkerTransport{}
}

// WorkerTransport is a http.RoundTripper that uses github.com/syumai/workers/cloudflare/fetch to make HTTP requests
type WorkerTransport struct{}

var _ http.RoundTripper = (*WorkerTransport)(nil)

func (t *WorkerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Create a new cloudflare/fetch request from the http.Request
	fetchReq, err := fetch.NewRequest(req.Context(), req.Method, req.URL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	for k, v := range req.Header {
		for _, vv := range v {
			fetchReq.Header.Add(k, vv)
		}
	}

	fetchClient := fetch.NewClient()
	return fetchClient.Do(fetchReq, nil)
}
