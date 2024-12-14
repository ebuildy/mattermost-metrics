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
	botID  string
}

func NewDriver(logger logger.Logger, endpointURL string, botID string) *Driver {
	return &Driver{
		logger: logger,
		botID:  botID,
		client: model.NewAPIv4Client(endpointURL),
	}
}

func (c *Driver) CollectMetrics(metrics *controller.Metrics) error {
	c.collectMetricsUsage(metrics)
	c.collectMetricsSystem(metrics)
	c.collectKPIMetrics(metrics)

	return nil
}

func (c *Driver) collectKPIMetrics(metrics *controller.Metrics) error {
	ctx := context.Background()

	var channels []*model.ChannelWithTeamData

	channelsResp, _, err := c.client.GetAllChannels(ctx, 0, 1000, "1")

	if err != nil {
		c.logger.Error(fmt.Sprintf("failed to get all channels: %s", err))
	}

	channels = append(channels, channelsResp...)

	postsCount := int64(0)
	lastPostDate := int64(0)

	for _, channel := range channels {
		postsCount += channel.TotalMsgCount

		if channel.LastPostAt > lastPostDate {
			lastPostDate = channel.LastPostAt
		}
	}

	metrics.KPIPostsCount = postsCount
	metrics.KPILastPostDate = lastPostDate
	metrics.KPIChannelsCount = int64(len(channels))

	return nil
}

func (c *Driver) collectMetricsSystem(metrics *controller.Metrics) error {
	ctx := context.Background()

	pingResp, _, err := c.client.GetPingWithOptions(ctx, model.SystemPingOptions{FullStatus: true})
	if err != nil {
		return err
	}

	metrics.SystemHealth = pingResp["status"] == "OK"
	metrics.SystemHealthDatabase = pingResp["database_status"] == "OK"
	metrics.SystemHealthFilestore = pingResp["filestore_status"] == "OK"

	return nil
}

func (c *Driver) collectMetricsUsage(metrics *controller.Metrics) error {
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

	return nil
}
