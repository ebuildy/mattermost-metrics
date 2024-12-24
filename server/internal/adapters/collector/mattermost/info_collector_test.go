package mattermost

import (
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/adapters/logger/fake"
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/adapters/services/mattermost_gateway"
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/domain"
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/utils"
	"github.com/jarcoal/httpmock"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDriver_collectInfo(t *testing.T) {
	httpmock.Activate()

	tests := []struct {
		name      string
		want      domain.MetricsDataInfo
		setupMock func(api *plugintest.API)
	}{
		{
			name: "all is good",
			want: domain.MetricsDataInfo{
				MattermostVersion:          "10.0",
				MattermostEdition:          "free",
				MattermostInstallationTime: time.Unix(1638224000, 0),
				SQLDriverName:              "mysql",
			},
			setupMock: func(api *plugintest.API) {
				api.On("GetConfig").Return(&model.Config{
					SqlSettings: model.SqlSettings{
						DriverName: utils.SafeRef("mysql"),
					},
				}, nil)
				api.On("GetSystemInstallDate").Return(int64(1638224000*1000), nil)
				api.On("GetServerVersion").Return("10.0", nil)
				api.On("GetLicense").Return(nil, nil)
			},
		},
		{
			name: "entreprise",
			want: domain.MetricsDataInfo{
				MattermostVersion:          "10.0",
				MattermostEdition:          "entreprise",
				MattermostInstallationTime: time.Unix(1638224000, 0),
				SQLDriverName:              "postgres",
			},
			setupMock: func(api *plugintest.API) {
				api.On("GetConfig").Return(&model.Config{
					SqlSettings: model.SqlSettings{
						DriverName: utils.SafeRef("postgres"),
					},
				}, nil)
				api.On("GetSystemInstallDate").Return(int64(1638224000*1000), nil)
				api.On("GetServerVersion").Return("10.0", nil)
				api.On("GetLicense").Return(&model.License{}, nil)
			},
		},
		{
			name: "API send a 500 error",
			want: domain.MetricsDataInfo{
				MattermostVersion:          "",
				MattermostEdition:          "free",
				MattermostInstallationTime: time.Unix(0, 0),
				SQLDriverName:              "",
			},
			setupMock: func(api *plugintest.API) {
				api.On("GetServerVersion").Return("", nil)
				api.On("GetSystemInstallDate").Return(int64(0), nil)
				api.On("GetConfig").Return(nil, nil)
				api.On("GetLicense").Return(nil, nil)
			},
		},
	}

	for _, tt := range tests {
		api := plugintest.NewAPI(t)

		mattermostGatewayClient := &mattermost_gateway.Client{
			API:       model.NewAPIv4Client(mattermostEndpointURL),
			PluginAPI: pluginapi.NewClient(api, &plugintest.Driver{}),
		}

		c := NewCollector(fake.NewFakeLogger(), mattermostGatewayClient)

		tt.setupMock(api)

		t.Run(tt.name, func(t *testing.T) {
			metrics := c.collectInfo()

			assert.Equal(t, tt.want.MattermostInstallationTime, metrics.MattermostInstallationTime)
			assert.Equal(t, tt.want.MattermostVersion, metrics.MattermostVersion)
			assert.Equal(t, tt.want.MattermostEdition, metrics.MattermostEdition)
			assert.Equal(t, tt.want.SQLDriverName, metrics.SQLDriverName)
		})
	}
}
