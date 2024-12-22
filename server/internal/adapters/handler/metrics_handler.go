package handler

import (
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/domain"
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/ports"
	"net/http"
)

type MetricsHandler struct {
	logger    ports.Logger
	collector []ports.MetricsCollector
	exporter  ports.MetricsHandler
}

func NewMetricsHandler(logger ports.Logger, collector []ports.MetricsCollector, exporter ports.MetricsHandler) *MetricsHandler {
	return &MetricsHandler{
		logger:    logger,
		collector: collector,
		exporter:  exporter,
	}
}

// ServeMetrics is called when /metrics is called
//
// -> Collect metrics data
// -> Feed exporter
// -> Render exporter
func (c *MetricsHandler) ServeMetrics(w http.ResponseWriter, r *http.Request) {
	metricsData := &domain.MetricsData{}

	for _, collector := range c.collector {
		err := collector.CollectMetrics(metricsData)

		if err != nil {
			c.logger.Error("Error while collecting metrics", "error", err)
		}
	}

	err := c.exporter.ExportMetrics(metricsData)
	if err != nil {
		return
	}

	c.exporter.ServeMetrics(w, r)
}
