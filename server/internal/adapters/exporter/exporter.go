package exporter

import (
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/domain"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	MetricsNamespace       = "mattermost"
	MetricsSubsystemUsage  = "usage"
	MetricsSubsystemSystem = "system"
	MetricsSubsystemKPI    = "kpi"
)

type Exporter struct {
	Registry    *prometheus.Registry
	metrics     *metrics
	HTTPHandler http.Handler
}

type metrics struct {
	info                                                      *prometheus.GaugeVec
	startTimeSeconds                                          prometheus.Gauge
	usagePostCount, usageUsersCount, usageStorage             prometheus.Gauge
	systemHealth, systemDatabaseHealth, systemFilestoreHealth prometheus.Gauge
	kpiPostsCount, kpiChannelsCount, kpiLastPostDate          prometheus.Gauge
}

func newSystemGauge(reg *prometheus.Registry, subsystem, name, help string) prometheus.Gauge {
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   MetricsNamespace,
		Subsystem:   subsystem,
		Name:        name,
		Help:        help,
		ConstLabels: nil,
	})

	reg.MustRegister(g)

	return g
}

func newSystemGaugeWithLabels(reg *prometheus.Registry, subsystem, name, help string, labels []string) *prometheus.GaugeVec {
	g := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   MetricsNamespace,
		Subsystem:   subsystem,
		Name:        name,
		Help:        help,
		ConstLabels: nil,
	}, labels)

	reg.MustRegister(g)

	return g
}

func NewExporter() *Exporter {
	registry := prometheus.NewRegistry()

	metrics := &metrics{
		info:                  newSystemGaugeWithLabels(registry, MetricsSubsystemUsage, "info", "Mattermost server info", []string{"version", "edition", "sqldriver"}),
		startTimeSeconds:      newSystemGauge(registry, MetricsSubsystemUsage, "start_time_seconds", "Start time of the process since unix epoch in seconds"),
		usagePostCount:        newSystemGauge(registry, MetricsSubsystemUsage, "posts_total", "Total number of posts"),
		usageUsersCount:       newSystemGauge(registry, MetricsSubsystemUsage, "users_total", "Total number of users"),
		usageStorage:          newSystemGauge(registry, MetricsSubsystemUsage, "storage_bytes", "Storage usage in bytes"),
		systemHealth:          newSystemGauge(registry, MetricsSubsystemUsage, "status", "Global status"),
		systemDatabaseHealth:  newSystemGauge(registry, MetricsSubsystemSystem, "database_status", "Database component status"),
		systemFilestoreHealth: newSystemGauge(registry, MetricsSubsystemSystem, "filestore_status", "Filestore component status"),
		kpiPostsCount:         newSystemGauge(registry, MetricsSubsystemKPI, "posts_total", "Total number of posts"),
		kpiChannelsCount:      newSystemGauge(registry, MetricsSubsystemKPI, "channels_total", "Total number of channels"),
		kpiLastPostDate:       newSystemGauge(registry, MetricsSubsystemKPI, "last_post_date", "Timestamp of last post date"),
	}

	return &Exporter{
		registry,
		metrics,
		promhttp.HandlerFor(registry, promhttp.HandlerOpts{
			EnableOpenMetrics: false,
			Registry:          registry,
		}),
	}
}

func (o *Exporter) ExportMetrics(metrics *domain.MetricsData) error {
	o.metrics.info.WithLabelValues(metrics.MattermostVersion, metrics.MattermostEdition, metrics.SQLDriverName).Set(1.0)

	o.metrics.startTimeSeconds.Set(float64(metrics.MattermostInstallationTime.Unix()))

	o.metrics.usagePostCount.Set(float64(metrics.UsagePostsCount))
	o.metrics.usageUsersCount.Set(float64(metrics.UsageUsersCount))
	o.metrics.systemHealth.Set(boolToFloat64(metrics.SystemHealth))

	o.metrics.systemDatabaseHealth.Set(boolToFloat64(metrics.SystemHealthDatabase))
	o.metrics.systemFilestoreHealth.Set(boolToFloat64(metrics.SystemHealthFilestore))
	o.metrics.usageStorage.Set(float64(metrics.UsageStorage))

	o.metrics.kpiPostsCount.Set(float64(metrics.KPIPostsCount))
	o.metrics.kpiChannelsCount.Set(float64(metrics.KPIChannelsCount))
	o.metrics.kpiLastPostDate.Set(float64(metrics.KPILastPostDate))

	return nil
}

func (o *Exporter) ServeMetrics(w http.ResponseWriter, r *http.Request) {
	o.HTTPHandler.ServeHTTP(w, r)
}

func boolToFloat64(b bool) float64 {
	if b {
		return 1.0
	}

	return 0.0
}
