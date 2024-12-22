package mattermost

import (
	"context"
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/adapters/services/mattermost_gateway"
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/domain"
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/ports"
)

type Collector struct {
	logger  ports.Logger
	gateway *mattermost_gateway.Client
}

func NewCollector(logger ports.Logger, gateway *mattermost_gateway.Client) *Collector {
	return &Collector{
		logger:  logger,
		gateway: gateway,
	}
}

// CollectMetrics is the public API to run metrics harvest
//
// This call all sub metrics collectors
func (c *Collector) CollectMetrics(metrics *domain.MetricsData) error {
	ctx := context.Background()

	metrics.Info = c.collectInfo()
	metrics.KPI = c.collectKPI()
	metrics.Health = c.collectHealth(ctx)
	metrics.Jobs = c.collectJob()

	return nil
}
