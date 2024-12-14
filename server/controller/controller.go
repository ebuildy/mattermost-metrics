package controller

import (
	"github.com/ebuildy/mattermost-plugin-minotor/server/logger"
	"net/http"
)

type Controller struct {
	logger    logger.Logger
	collector MetricsCollector
	exporter  MetricsExporter
}

type Metrics struct {
	UsagePostsCount, UsageUsersCount, UsageStorage int64

	SystemHealth, SystemHealthDatabase, SystemHealthFilestore bool

	KPILastPostDate, KPIChannelsLastCreationDate, KPIChannelsCount, KPIPostsCount int64
}

type MetricsCollector interface {
	CollectMetrics(metrics *Metrics) error
}

type MetricsExporter interface {
	ExportMetrics(metrics *Metrics) error
	ServeMetrics(w http.ResponseWriter, r *http.Request)
}

func NewCollector(logger logger.Logger, collector MetricsCollector, exporter MetricsExporter) *Controller {
	return &Controller{
		logger:    logger,
		collector: collector,
		exporter:  exporter,
	}
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	metricsData := &Metrics{}

	err := c.collector.CollectMetrics(metricsData)
	if err != nil {
		return
	}

	err = c.exporter.ExportMetrics(metricsData)
	if err != nil {
		return
	}

	c.exporter.ServeMetrics(w, r)
}
