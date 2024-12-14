package mattermost

import (
	"fmt"
	"github.com/ebuildy/mattermost-plugin-minotor/server/controller"
	"github.com/ebuildy/mattermost-plugin-minotor/server/logger"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	mattermostEndpointURL = "http://test:8080"
)

func TestDriver_collectMetricsUsaget(t *testing.T) {
	httpmock.Activate()

	tests := []struct {
		name            string
		want            int64
		apiMockResponse httpmock.Responder
	}{
		{
			name:            "API returns 0 posts count",
			want:            0,
			apiMockResponse: httpmock.NewStringResponder(200, `{"count":0}`),
		},
		{
			name:            "API returns 89 posts count",
			want:            89,
			apiMockResponse: httpmock.NewStringResponder(200, `{"count":89}`),
		},
		{
			name:            "API send a 500 error",
			want:            0,
			apiMockResponse: httpmock.NewStringResponder(500, `{"error":"unexpected error"`),
		},
	}

	c := NewDriver(logger.NewFakeLogger(), mattermostEndpointURL, "")

	for _, tt := range tests {
		httpmock.RegisterResponder("GET", fmt.Sprintf("%s/api/v4/usage/posts", mattermostEndpointURL), tt.apiMockResponse)
		t.Run(tt.name, func(t *testing.T) {
			metrics := controller.Metrics{}
			c.collectMetricsUsage(&metrics)

			assert.Equal(t, tt.want, metrics.UsagePostsCount)
		})
	}
}

func TestDriver_collectMetricsSystem(t *testing.T) {
	httpmock.Activate()

	tests := []struct {
		name            string
		want            controller.Metrics
		apiMockResponse httpmock.Responder
	}{
		{
			name:            "all is good",
			want:            controller.Metrics{SystemHealth: true, SystemHealthFilestore: true, SystemHealthDatabase: true},
			apiMockResponse: httpmock.NewStringResponder(200, `{"ActiveSearchBackend": "database","status": "OK", "filestore_status": "OK", "database_status": "OK"}`),
		},
		{
			name:            "filestore is down",
			want:            controller.Metrics{SystemHealth: false, SystemHealthFilestore: false, SystemHealthDatabase: true},
			apiMockResponse: httpmock.NewStringResponder(200, `{"ActiveSearchBackend": "database","status": "err", "filestore_status": "err", "database_status": "OK"}`),
		},
		{
			name:            "API send a 500 error",
			want:            controller.Metrics{SystemHealth: false, SystemHealthFilestore: false, SystemHealthDatabase: false},
			apiMockResponse: httpmock.NewStringResponder(500, `{"error":"unexpected error"`),
		},
	}

	c := NewDriver(logger.NewFakeLogger(), mattermostEndpointURL, "")

	for _, tt := range tests {
		httpmock.RegisterResponder("GET", fmt.Sprintf("%s/api/v4/system/ping?get_server_status=true&use_rest_semantics=false", mattermostEndpointURL), tt.apiMockResponse)
		t.Run(tt.name, func(t *testing.T) {
			metrics := controller.Metrics{}
			c.collectMetricsSystem(&metrics)

			assert.Equal(t, tt.want.SystemHealth, metrics.SystemHealth)
			assert.Equal(t, tt.want.SystemHealthFilestore, metrics.SystemHealthFilestore)
			assert.Equal(t, tt.want.SystemHealthDatabase, metrics.SystemHealthDatabase)
		})
	}
}
