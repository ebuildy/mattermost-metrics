package collector

import (
	"context"
	"fmt"
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/domain"
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/ports"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/pluginapi"
)

type Driver struct {
	logger          ports.Logger
	client          *model.Client4
	pluginAPIClient *pluginapi.Client
}

func NewAuthenticatedDriver(logger ports.Logger, pluginAPIClient *pluginapi.Client, accessToken string, endpointURL string) *Driver {
	client := model.NewAPIv4Client(endpointURL)

	client.SetToken(accessToken)

	return &Driver{
		logger:          logger,
		client:          client,
		pluginAPIClient: pluginAPIClient,
	}
}

func NewDriver(logger ports.Logger, endpointURL string) *Driver {
	client := model.NewAPIv4Client(endpointURL)

	return &Driver{
		logger: logger,
		client: client,
	}
}

// CollectMetrics is the public API to run metrics harvest
//
// This call all sub metrics collectors
func (c *Driver) CollectMetrics(metrics *domain.MetricsData) error {
	ctx := context.Background()

	metrics.MattermostVersion = c.pluginAPIClient.System.GetServerVersion()
	metrics.MattermostInstallationTime, _ = c.pluginAPIClient.System.GetSystemInstallDate()
	metrics.SQLDriverName = *c.pluginAPIClient.Configuration.GetConfig().SqlSettings.DriverName

	if c.pluginAPIClient.System.GetLicense() == nil {
		metrics.MattermostEdition = "free"
	} else {
		metrics.MattermostEdition = "entreprise"
	}

	c.collectMetricsUsage(ctx, metrics)
	c.collectMetricsSystem(ctx, metrics)
	c.collectKPIMetrics(ctx, metrics)

	return nil
}

func (c *Driver) collectKPIMetrics(ctx context.Context, metrics *domain.MetricsData) {
	var channels []*model.ChannelWithTeamData

	channelsResp, channelsCountResp, _, err := c.client.GetAllChannelsWithCount(ctx, 0, 100, "1")

	if err != nil {
		c.logger.Error(fmt.Sprintf("failed to get all channels: %s", err))
	}

	channels = append(channels, channelsResp...)

	postsCount := int64(0)
	lastPostDate := int64(0)
	lastChannelDate := int64(0)

	for _, channel := range channels {
		if channel != nil {
			postsCount += channel.TotalMsgCount

			if channel.LastPostAt > lastPostDate {
				lastPostDate = channel.LastPostAt
			}

			if channel.CreateAt > lastChannelDate {
				lastChannelDate = channel.CreateAt
			}
		}
	}

	metrics.KPIPostsCount = postsCount
	metrics.KPILastPostDate = lastPostDate
	metrics.KPIChannelsLastCreationDate = lastChannelDate
	metrics.KPIChannelsCount = channelsCountResp
}

func (c *Driver) collectMetricsSystem(ctx context.Context, metrics *domain.MetricsData) {
	pingResp, _, err := c.client.GetPingWithOptions(ctx, model.SystemPingOptions{FullStatus: true})
	if err != nil {
		c.logger.Error(fmt.Sprintf("failed to get system ping: %s", err))
	}

	metrics.SystemHealth = pingResp["status"] == "OK"
	metrics.SystemHealthDatabase = pingResp["database_status"] == "OK"
	metrics.SystemHealthFilestore = pingResp["filestore_status"] == "OK"
}

func (c *Driver) collectMetricsUsage(ctx context.Context, metrics *domain.MetricsData) {
	if postUsage, _, err := c.client.GetPostsUsage(ctx); err != nil {
		c.logger.Error(fmt.Sprintf("failed to get posts usage count: %s", err))
		metrics.UsagePostsCount = 0
	} else {
		metrics.UsagePostsCount = postUsage.Count
	}

	if postUsage, _, err := c.client.GetTotalUsersStats(ctx, "1"); err != nil {
		c.logger.Error(fmt.Sprintf("failed to get posts usage count: %s", err))
		metrics.UsageUsersCount = 0
	} else {
		metrics.UsageUsersCount = postUsage.TotalUsersCount
	}

	if storageUsage, _, err := c.client.GetStorageUsage(ctx); err != nil {
		c.logger.Error(fmt.Sprintf("failed to get posts usage count: %s", err))
		metrics.UsageStorage = 0
	} else {
		metrics.UsageStorage = storageUsage.Bytes
	}
}
