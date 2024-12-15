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

	pluginAPIClient *pluginapi.Client

	bot *model.Bot

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

	apiAccessTokenId *string
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}

// OnActivate is called when plugin is enabled
//
// -> Setup stuff (bot, accessToken ...)
// -> Run HTTP controller
func (p *Plugin) OnActivate() error {
	p.pluginAPIClient = pluginapi.NewClient(p.API, p.Driver)
	p.config = config.NewConfigService(p.pluginAPIClient, manifest)
	p.bot = &model.Bot{
		Username:    "minotor",
		DisplayName: "Minotor",
		Description: "Minotor bot to expose exporter.",
		OwnerId:     "minotor",
	}

	botID, err := p.pluginAPIClient.Bot.EnsureBot(p.bot,
		pluginapi.ProfileImagePath("assets/logo.png"),
	)
	if err != nil {
		return errors.Wrapf(err, "failed to ensure bot")
	}

	p.bot.UserId = botID

	botUser, err := p.pluginAPIClient.User.UpdateRoles(botID, model.SystemAdminRoleId)

	if err != nil {
		return err
	}

	p.pluginAPIClient.Log.Info(fmt.Sprintf("Bot user %s %s with roles %s", botUser.Username, botUser.Id, botUser.Roles))

	err = p.config.UpdateConfiguration(func(c *config.Configuration) {
		c.BotUserID = botID
	})
	if err != nil {
		return errors.Wrapf(err, "failed save bot to config")
	}

	accessTokenResp, err := p.pluginAPIClient.User.CreateAccessToken(botID, "minotor api proxy")

	if err != nil {
		return errors.Wrapf(err, "Error creating access token")
	}

	p.apiAccessTokenId = &accessTokenResp.Id
	p.router = mux.NewRouter()
	p.driver = mattermost.NewAuthenticatedDriver(&p.pluginAPIClient.Log, accessTokenResp.Token, "http://localhost:8065")
	p.exporter = exporter.NewExporter()
	p.collector = controller.NewCollector(&p.pluginAPIClient.Log, p.driver, p.exporter)

	p.router.Handle("/metrics", p.collector)

	return nil
}

// OnDeactivate is called when plugin is disabled
//
// -> revoke Mattermost API access token
func (p *Plugin) OnDeactivate() error {
	if p.apiAccessTokenId != nil {
		err := p.pluginAPIClient.User.RevokeAccessToken(*p.apiAccessTokenId)
		if err != nil {
			p.pluginAPIClient.Log.Warn("Failed to revoke access token", "error", err)
		}
	}
	return nil
}
