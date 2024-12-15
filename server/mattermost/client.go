package mattermost

import (
	"context"
	"fmt"
	"github.com/ebuildy/mattermost-plugin-minotor/server/controller"
	"github.com/ebuildy/mattermost-plugin-minotor/server/logger"

	"github.com/mattermost/mattermost/server/public/model"
)

type Driver struct {
	logger logger.Logger
	client *model.Client4
}

func NewAuthenticatedDriver(logger logger.Logger, accessToken string, endpointURL string) *Driver {
	client := model.NewAPIv4Client(endpointURL)

	client.SetToken(accessToken)

	return &Driver{
		logger: logger,
		client: client,
	}
}

func NewDriver(logger logger.Logger, endpointURL string) *Driver {
	client := model.NewAPIv4Client(endpointURL)

	return &Driver{
		logger: logger,
		client: client,
	}
}

// CollectMetrics is the public API to run metrics harvest
//
// This call all sub metrics collectors
func (c *Driver) CollectMetrics(metrics *controller.Metrics) error {
	c.collectMetricsUsage(metrics)
	c.collectMetricsSystem(metrics)
	c.collectKPIMetrics(metrics)

	return nil
}

func (c *Driver) collectKPIMetrics(metrics *controller.Metrics) {
	ctx := context.Background()

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

func (c *Driver) collectMetricsSystem(metrics *controller.Metrics) {
	ctx := context.Background()

	pingResp, _, err := c.client.GetPingWithOptions(ctx, model.SystemPingOptions{FullStatus: true})
	if err != nil {
		c.logger.Error(fmt.Sprintf("failed to get system ping: %s", err))
	}

	metrics.SystemHealth = pingResp["status"] == "OK"
	metrics.SystemHealthDatabase = pingResp["database_status"] == "OK"
	metrics.SystemHealthFilestore = pingResp["filestore_status"] == "OK"
}

func (c *Driver) collectMetricsUsage(metrics *controller.Metrics) {
	ctx := context.Background()

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
