package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/ebuildy/mattermost-plugin-minotor/server/config"
	"github.com/ebuildy/mattermost-plugin-minotor/server/controller"
	"github.com/ebuildy/mattermost-plugin-minotor/server/exporter"
	"github.com/ebuildy/mattermost-plugin-minotor/server/mattermost"

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

	exporter  *exporter.Exporter
	config    *config.ServiceImpl
	driver    *mattermost.Driver
	collector *controller.Controller

	router *mux.Router
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}

func (p *Plugin) OnActivate() error {
	pluginAPIClient := pluginapi.NewClient(p.API, p.Driver)

	p.config = config.NewConfigService(pluginAPIClient, manifest)

	botID, err := pluginAPIClient.Bot.EnsureBot(&model.Bot{
		Username:    "minotor",
		DisplayName: "Minotor",
		Description: "Minotor bot to expose exporter.",
		OwnerId:     "minotor",
	},
		pluginapi.ProfileImagePath("assets/logo.png"),
	)
	if err != nil {
		return errors.Wrapf(err, "failed to ensure bot")
	}

	botUser, err := pluginAPIClient.User.UpdateRoles(botID, model.SystemAdminRoleId)

	if err != nil {
		return err
	}

	pluginAPIClient.Log.Info(fmt.Sprintf("Bot user %s %s with roles %s", botUser.Username, botUser.Id, botUser.Roles))

	err = p.config.UpdateConfiguration(func(c *config.Configuration) {
		c.BotUserID = botID
	})
	if err != nil {
		return errors.Wrapf(err, "failed save bot to config")
	}

	p.router = mux.NewRouter()
	p.driver = mattermost.NewDriver(&pluginAPIClient.Log, "http://localhost:8065", botID)
	p.exporter = exporter.NewExporter()
	p.collector = controller.NewCollector(&pluginAPIClient.Log, p.driver, p.exporter)

	p.router.Handle("/metrics", p.collector)

	return nil
}
