package mattermost

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ebuildy/mattermost-plugin-minotor/server/controller"
	"github.com/ebuildy/mattermost-plugin-minotor/server/logger"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
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

	c := NewDriver(logger.NewFakeLogger(), mattermostEndpointURL)

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

	c := NewDriver(logger.NewFakeLogger(), mattermostEndpointURL)

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

func TestDriver_collectKPIMetrics(t *testing.T) {
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	tests := []struct {
		name            string
		want            controller.Metrics
		apiMockResponse httpmock.Responder
	}{
		{
			name: "No channel",
			want: controller.Metrics{KPILastPostDate: 0, KPIPostsCount: 0, KPIChannelsCount: 0, KPIChannelsLastCreationDate: 0},
			apiMockResponse: func(req *http.Request) (*http.Response, error) {
				resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
					"channels":    []any{},
					"total_count": 0,
				})
				if err != nil {
					return httpmock.NewStringResponse(500, ""), nil
				}
				return resp, nil
			},
		},
		{
			name: "8 channels, no pagination",
			want: controller.Metrics{KPILastPostDate: 1734186913838, KPIPostsCount: 88, KPIChannelsCount: 8, KPIChannelsLastCreationDate: 1733941707169},
			apiMockResponse: func(req *http.Request) (*http.Response, error) {
				var channels = make([]any, 8)

				for i := 0; i < 8; i++ {
					channel := map[string]any{
						"id":                   fmt.Sprintf("%daqj7ttnji8hjjgrcgfrfndyyo", i),
						"create_at":            1733941707162 + i,
						"update_at":            1733941707162 + i,
						"delete_at":            0,
						"team_id":              "91fsqtw98388tbt8ka5d1fgzhy",
						"type":                 "O",
						"display_name":         "Off-Topic",
						"name":                 "off-topic",
						"header":               "",
						"purpose":              "",
						"last_post_at":         1734186913831 + i,
						"total_msg_count":      11,
						"extra_update_at":      0,
						"creator_id":           "",
						"total_msg_count_root": 5,
						"last_root_post_at":    1734186913831 + i,
						"team_display_name":    "test",
						"team_name":            "test",
						"team_update_at":       1733941707151 + i,
					}

					channels[i] = channel
				}

				resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
					"channels":    channels,
					"total_count": 8,
				})
				if err != nil {
					return httpmock.NewStringResponse(500, ""), nil
				}
				return resp, nil
			},
		},
		// {
		//	name:            "200 channels, must paginate",
		//	want:            controller.Metrics{KPILastPostDate: 1734186913838, KPIPostsCount: 1, KPIChannelsCount: 1, KPIChannelsLastCreationDate: 1733941707169},
		//	apiMockResponse: httpmock.NewStringResponder(200, `{"ActiveSearchBackend": "database","status": "err", "filestore_status": "err", "database_status": "OK"}`),
		// },
		// {
		//	name:            "API send a 500 error",
		//	want:            controller.Metrics{KPILastPostDate: 1734186913838, KPIPostsCount: 1, KPIChannelsCount: 1, KPIChannelsLastCreationDate: 1733941707169},
		//	apiMockResponse: httpmock.NewStringResponder(500, `{"error":"unexpected error"`),
		// },
	}

	c := NewDriver(logger.NewFakeLogger(), mattermostEndpointURL)

	for _, tt := range tests {
		httpmock.RegisterResponder("GET", fmt.Sprintf("%s/api/v4/channels?page=0&per_page=100&include_total_count=true", mattermostEndpointURL), tt.apiMockResponse)
		t.Run(tt.name, func(t *testing.T) {
			metrics := controller.Metrics{}
			c.collectKPIMetrics(&metrics)

			assert.Equal(t, tt.want.KPILastPostDate, metrics.KPILastPostDate, "KPILastPostDate")
			assert.Equal(t, tt.want.KPIPostsCount, metrics.KPIPostsCount, "KPIPostsCount")
			assert.Equal(t, tt.want.KPIChannelsCount, metrics.KPIChannelsCount, "KPIChannelsCount")
			assert.Equal(t, tt.want.KPIChannelsLastCreationDate, metrics.KPIChannelsLastCreationDate, "KPIChannelsLastCreationDate")
		})
	}
}
