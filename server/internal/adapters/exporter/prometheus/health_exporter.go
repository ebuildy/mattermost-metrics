package prometheus

import (
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/domain"
	"github.com/prometheus/client_golang/prometheus"
)

type HealthExporter struct {
	system, systemDatabase, systemFilestore prometheus.Gauge
}

func newHealth(registry *prometheus.Registry) HealthExporter {
	return HealthExporter{
		system:          newSystemGauge(registry, MetricsSubsystemUsage, "status", "Global status"),
		systemDatabase:  newSystemGauge(registry, MetricsSubsystemSystem, "database_status", "Database component status"),
		systemFilestore: newSystemGauge(registry, MetricsSubsystemSystem, "filestore_status", "Filestore component status"),
	}
}

func (m HealthExporter) ExportMetrics(metrics *domain.MetricsData) error {
	healthMetrics := metrics.Health
	m.system.Set(boolToFloat64(healthMetrics.SystemHealth))
	m.systemDatabase.Set(boolToFloat64(healthMetrics.SystemHealthDatabase))
	m.systemFilestore.Set(boolToFloat64(healthMetrics.SystemHealthFilestore))

	return nil
}
