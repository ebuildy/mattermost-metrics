package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/domain"
)

type InfoExporter struct {
	info             *prometheus.GaugeVec
	startTimeSeconds prometheus.Gauge
}

func newInfo(registry *prometheus.Registry) InfoExporter {
	return InfoExporter{
		info:             newSystemGaugeWithLabels(registry, MetricsSubsystemUsage, "info", "Mattermost server info", []string{"version", "edition", "sqldriver"}),
		startTimeSeconds: newSystemGauge(registry, MetricsSubsystemUsage, "start_time_seconds", "Start time of the process since unix epoch in seconds"),
	}
}

func (m InfoExporter) ExportMetrics(metrics *domain.MetricsData) error {
	infoMetrics := metrics.Info
	m.info.WithLabelValues(infoMetrics.MattermostVersion, infoMetrics.MattermostEdition, infoMetrics.SQLDriverName).Set(1.0)
	m.startTimeSeconds.Set(float64(infoMetrics.MattermostInstallationTime.Unix()))

	return nil
}
