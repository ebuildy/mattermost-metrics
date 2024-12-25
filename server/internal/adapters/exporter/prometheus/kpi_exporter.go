package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/domain"
)

type KPIExporter struct {
	usagePostCount, usageUsersCount, usageStorage prometheus.Gauge

	kpiPostsCount, kpiLastPostDate, sessionsCount prometheus.Gauge

	channelsCountByTypes *prometheus.GaugeVec
}

func newKPI(registry *prometheus.Registry) KPIExporter {
	return KPIExporter{
		usagePostCount:       newSystemGauge(registry, MetricsSubsystemUsage, "posts_total", "Total number of posts"),
		usageUsersCount:      newSystemGauge(registry, MetricsSubsystemUsage, "users_total", "Total number of users"),
		usageStorage:         newSystemGauge(registry, MetricsSubsystemUsage, "storage_bytes", "Storage usage in bytes"),
		kpiPostsCount:        newSystemGauge(registry, MetricsSubsystemKPI, "posts_total", "Total number of posts"),
		kpiLastPostDate:      newSystemGauge(registry, MetricsSubsystemKPI, "last_post_date", "Timestamp of last post date"),
		sessionsCount:        newSystemGauge(registry, MetricsSubsystemKPI, "sessions_total", "Total number of sessions"),
		channelsCountByTypes: newSystemGaugeWithLabels(registry, MetricsSubsystemKPI, "channels_total", "Number of channels by type", []string{"type"}),
	}
}

func (m KPIExporter) ExportMetrics(metrics *domain.MetricsData) error {
	kpiMetrics := metrics.KPI

	m.usagePostCount.Set(float64(metrics.UsagePostsCount))
	m.usageUsersCount.Set(float64(metrics.UsageUsersCount))
	m.usageStorage.Set(float64(metrics.UsageStorage))

	m.kpiPostsCount.Set(float64(kpiMetrics.KPIPostsCount))
	m.kpiLastPostDate.Set(float64(kpiMetrics.KPILastPostDate))

	m.sessionsCount.Set(float64(kpiMetrics.KPISessionsCount))

	m.channelsCountByTypes.WithLabelValues("public").Set(float64(kpiMetrics.KPIPublicChannelsCount))
	m.channelsCountByTypes.WithLabelValues("private").Set(float64(kpiMetrics.KPIPrivateChannelsCount))
	m.channelsCountByTypes.WithLabelValues("direct").Set(float64(kpiMetrics.KPIDirectMessagesCount))

	return nil
}
