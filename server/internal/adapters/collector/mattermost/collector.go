package mattermost

import (
	"context"

	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/domain"
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/ports"
)

// Collector is responsible for collecting various metrics and aggregating data from different sources.
type Collector struct {
	logger  ports.Logger
	gateway *Client
	config  *domain.ConfigCollector
}

func NewCollector(logger ports.Logger, gateway *Client) *Collector {
	return &Collector{
		logger:  logger,
		gateway: gateway,
	}
}

func (c *Collector) Configure(config *domain.ConfigCollector) {
	c.config = config
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

	if c.config.ReactionEnabled {
		metrics.Reactions = c.collectReaction(c.config.ReactionCountByEmojiLimits)
	}

	return nil
}
