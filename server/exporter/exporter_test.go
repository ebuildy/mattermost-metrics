package exporter

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetrics_HTTPHandler(t *testing.T) {
	tests := []struct {
		name       string
		assertFunc func(t *testing.T, body string)
	}{
		{
			name: "No API no metrics",
			assertFunc: func(t *testing.T, body string) {
				assert.NotEmpty(t, body)
			},
		},
	}

	o := NewExporter()
	httpHandler := o.HTTPHandler

	u, _ := url.Parse("http://localhost/metrics")

	fakeReq := &http.Request{
		Method: "GET",
		URL:    u,
		Header: map[string][]string{
			"Accept": {"text/html,application/xhtml+xml"},
		},
	}

	// o.UsageUsersCountSet(0).UsagePostsCountSet(10)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewFakeHTTPResponse("", 200)
			httpHandler.ServeHTTP(w, fakeReq)

			tt.assertFunc(t, w.String())
		})
	}
}

type FakeHTTPResponse struct {
	body       string
	statusCode int
}

func NewFakeHTTPResponse(body string, statusCode int) *FakeHTTPResponse {
	return &FakeHTTPResponse{body, statusCode}
}

func (f *FakeHTTPResponse) Header() http.Header {
	return http.Header{}
}

func (f *FakeHTTPResponse) StatusCode() int {
	return f.statusCode
}

func (f *FakeHTTPResponse) Body() []byte {
	return []byte(f.body)
}

func (f *FakeHTTPResponse) String() string {
	return f.body
}

func (f *FakeHTTPResponse) WriteHeader(statusCode int) {
	f.statusCode = statusCode
}

func (f *FakeHTTPResponse) Write(body []byte) (int, error) {
	f.body = string(body)

	return len(body), nil
}
