package main

import (
	"net/http"
	"sync"

	"github.com/ebuildy/mattermost-plugin-minotor/server/config"
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/adapters/collector/mattermost"
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/adapters/exporter/prometheus"
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/adapters/handler"
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/domain"
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/ports"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/pluginapi"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	pluginAPIClient *pluginapi.Client

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	exporter            *prometheus.Exporter
	config              *config.ServiceImpl
	mattermostCollector *mattermost.Collector
	metricsHandler      *handler.MetricsHandler

	mattermostGatewayClient *mattermost.Client
}

func (p *Plugin) ServeMetrics(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.metricsHandler.ServeMetrics(w, r)
}

// OnActivate is called when plugin is enabled
//
// -> Setup stuff (bot, accessToken ...)
// -> Run HTTP controller
func (p *Plugin) OnActivate() error {
	p.pluginAPIClient = pluginapi.NewClient(p.API, p.Driver)
	p.config = config.NewConfigService(p.pluginAPIClient, manifest)

	logger := &p.pluginAPIClient.Log

	mattermostGatewayClient, err := mattermost.NewClient(p.pluginAPIClient)

	if err != nil {
		return errors.Wrapf(err, "failed to create mattermost gateway client")
	}

	err = p.config.UpdateConfiguration(func(c *config.Configuration) {
		c.BotUserID = mattermostGatewayClient.Bot.UserId
	})
	if err != nil {
		return errors.Wrapf(err, "failed save bot to config")
	}

	p.mattermostCollector = mattermost.NewCollector(logger, mattermostGatewayClient)

	p.mattermostCollector.Configure(&domain.ConfigCollector{
		ReactionEnabled:            true,
		ReactionCountByEmojiLimits: 5,
	})

	p.exporter = prometheus.NewExporter(logger)
	p.metricsHandler = handler.NewMetricsHandler(logger, []ports.MetricsCollector{p.mattermostCollector}, p.exporter)

	return nil
}

// OnDeactivate is called when plugin is disabled
//
// -> revoke Mattermost API access token
func (p *Plugin) OnDeactivate() error {
	if p.mattermostGatewayClient != nil && p.mattermostGatewayClient.APIUserAccessToken != nil {
		err := p.pluginAPIClient.User.RevokeAccessToken(p.mattermostGatewayClient.APIUserAccessToken.Id)
		if err != nil {
			p.pluginAPIClient.Log.Warn("Failed to revoke access token", "error", err)
		}
	}
	return nil
}
