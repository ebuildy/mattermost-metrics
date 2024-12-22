package mattermost_gateway

import (
	"database/sql"
	"fmt"
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/ports"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/pkg/errors"
)

// Client allow to communicate with Mattermost
type Client struct {
	Bot                *model.Bot
	APIUserAccessToken *model.UserAccessToken

	logger ports.Logger

	DB        *sql.DB
	API       *model.Client4
	PluginAPI *pluginapi.Client
}

func NewClient(pluginAPI *pluginapi.Client) (*Client, error) {

	logger := &pluginAPI.Log

	bot, err := getBot(pluginAPI)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to ensure bot")
	}

	dbClient, err := getDBClient(pluginAPI)

	if err != nil {
		return nil, err
	}

	apiClient, accessTokenResp, err := getAPIClient(pluginAPI, bot)

	if err != nil {
		return nil, err
	}

	return &Client{
		Bot:                bot,
		APIUserAccessToken: accessTokenResp,
		logger:             logger,
		DB:                 dbClient,
		API:                apiClient,
		PluginAPI:          pluginAPI,
	}, nil
}

// SQLValue executes a SQL query and retrieves a single integer value from the database, returning 0 on error or failure.
func (c *Client) SQLValue(query string) int64 {
	rows, err := c.DB.Query(query)

	if err != nil {
		c.logger.Error("Error while querying DB", "error", err)
		return 0
	}
	defer rows.Close()

	var count int64

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			c.logger.Error("Error while querying DB", "error", err)
			return 0
		}
	}

	return count
}

func getBot(pluginAPI *pluginapi.Client) (*model.Bot, error) {
	bot := &model.Bot{
		Username:    "minotor",
		DisplayName: "Minotor",
		Description: "Minotor bot to expose exporter.",
		OwnerId:     "minotor",
	}

	botID, err := pluginAPI.Bot.EnsureBot(bot,
		pluginapi.ProfileImagePath("assets/logo.png"),
	)

	if err != nil {
		return nil, err
	}

	bot.UserId = botID

	return bot, nil
}

func getDBClient(pluginAPI *pluginapi.Client) (*sql.DB, error) {
	pluginAPI.Store.GetMasterDB()
	mattermostDBClient, err := pluginAPI.Store.GetReplicaDB()

	if err != nil {
		return nil, err
	}

	if mattermostDBClient == nil {
		pluginAPI.Log.Error("dbDriver is nil")
		return nil, errors.New("dbDriver is nil")
	}

	return mattermostDBClient, nil
}

func getAPIClient(pluginAPI *pluginapi.Client, bot *model.Bot) (*model.Client4, *model.UserAccessToken, error) {
	botUser, err := pluginAPI.User.UpdateRoles(bot.UserId, model.SystemAdminRoleId)

	if err != nil {
		return nil, nil, err
	}

	pluginAPI.Log.Info(fmt.Sprintf("Bot user %s %s with roles %s", botUser.Username, botUser.Id, botUser.Roles))

	accessTokenResp, err := pluginAPI.User.CreateAccessToken(bot.UserId, "minotor api proxy")

	if err != nil {
		return nil, nil, errors.Wrapf(err, "Error creating access token")
	}

	mattermostAPIClient := model.NewAPIv4Client("http://localhost:8065")

	mattermostAPIClient.SetToken(accessTokenResp.Token)

	return mattermostAPIClient, accessTokenResp, nil
}
