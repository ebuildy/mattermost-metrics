package mattermost

import (
	"context"
	"fmt"
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/domain"
	"github.com/mattermost/mattermost/server/public/model"
)

func (c *Collector) collectHealth(ctx context.Context) *domain.MetricsDataHealth {
	pingResp, _, err := c.gateway.API.GetPingWithOptions(ctx, model.SystemPingOptions{FullStatus: true})
	if err != nil {
		c.logger.Error(fmt.Sprintf("failed to get system ping: %s", err))
	}

	return &domain.MetricsDataHealth{
		SystemHealth:          pingResp["status"] == "OK",
		SystemHealthDatabase:  pingResp["database_status"] == "OK",
		SystemHealthFilestore: pingResp["filestore_status"] == "OK",
	}
}
