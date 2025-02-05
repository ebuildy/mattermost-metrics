package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/domain"
)

type CallExporter struct {
	totalCount, totalDuration, lastCreatedDate prometheus.Gauge
}

func newCall(registry *prometheus.Registry) CallExporter {
	return CallExporter{
		totalCount:      newSystemGauge(registry, "calls", "count_total", "Total number of calls"),
		totalDuration:   newSystemGauge(registry, "calls", "duration_seconds", "Total duration of calls in seconds"),
		lastCreatedDate: newSystemGauge(registry, "calls", "last_date_seconds", "Last calls creation date"),
	}
}

func (m CallExporter) ExportMetrics(metrics *domain.MetricsData) error {
	if metrics.Calls != nil {
		m.totalCount.Set(float64(metrics.Calls.Count))
		m.lastCreatedDate.Set(float64(metrics.Calls.Last))
		m.totalDuration.Set(float64(metrics.Calls.Duration))
	}

	return nil
}
