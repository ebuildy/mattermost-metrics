package main

import (
	"net/http"
	"sync"

	"github.com/pkg/errors"
	"github.com/gorilla/mux"
	// "github.com/sirupsen/logrus"

	"github.com/ebuildy/mattermost-plugin-minotor/server/api"
	"github.com/ebuildy/mattermost-plugin-minotor/server/config"
	"github.com/ebuildy/mattermost-plugin-minotor/server/registry"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/pluginapi"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	handler              *api.Handler
	config               *config.ServiceImpl

	router *mux.Router
}

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	// p.handler.ServeHTTP(w, r)
	p.router.ServeHTTP(w, r)
}


// OnActivate Called when this plugin is activated.
func (p *Plugin) OnActivate() error {
	pluginAPIClient := pluginapi.NewClient(p.API, p.Driver)

	p.config = config.NewConfigService(pluginAPIClient, manifest)
	p.handler = api.NewHandler(pluginAPIClient, p.config)

	// logger := logrus.StandardLogger()

	botID, err := pluginAPIClient.Bot.EnsureBot(&model.Bot{
		Username:    "minotor",
		DisplayName: "Minotor",
		Description: "Minotor bot to expose metrics.",
		OwnerId:     "minotor",
	},
		pluginapi.ProfileImagePath("assets/logo.png"),
	)
	if err != nil {
		return errors.Wrapf(err, "failed to ensure bot")
	}

	err = p.config.UpdateConfiguration(func(c *config.Configuration) {
		c.BotUserID = botID
	})
	if err != nil {
		return errors.Wrapf(err, "failed save bot to config")
	}

	p.router = mux.NewRouter()

	registry := registry.NewRegistry()

	p.router.Handle("/metrics", registry.HandleHTTP())

	return nil
}
