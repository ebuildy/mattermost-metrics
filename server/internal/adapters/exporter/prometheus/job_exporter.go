package prometheus

import (
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/domain"
	"github.com/prometheus/client_golang/prometheus"
)

type JobExporter struct {
	countByStatusType *prometheus.GaugeVec
	last              prometheus.Gauge
}

func newJob(registry *prometheus.Registry) JobExporter {
	return JobExporter{
		countByStatusType: newSystemGaugeWithLabels(registry, MetricsSubsystemUsage, "job_total", "Jobs count by status and type", []string{"type", "status"}),
		last:              newSystemGauge(registry, MetricsSubsystemUsage, "job_last_seconds", "Last job execution time - unix timestamp"),
	}
}

func (m JobExporter) ExportMetrics(metrics *domain.MetricsData) error {
	jobMetrics := metrics.Jobs

	m.countByStatusType.Reset()

	if jobMetrics != nil {
		for _, item := range jobMetrics.CountByTypesStatus {
			m.countByStatusType.WithLabelValues(item.Type, item.Status).Set(float64(item.Count))
		}

		m.last.Set(float64(jobMetrics.Last.Unix()))
	}

	return nil
}
