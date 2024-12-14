package exporter

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

func TestMetrics_HTTPHandler(t *testing.T) {

	tests := []struct {
		name        string
		mustContain string
	}{
		{
			name:        "API returns 0 posts count",
			mustContain: "mattermost_usage_posts_total 0",
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

	o.UsageUsersCountSet(0).UsagePostsCountSet(10)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewFakeHTTPResponse("", 200)
			httpHandler.ServeHTTP(w, fakeReq)

			assert.Contains(t, w.String(), tt.mustContain)
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
