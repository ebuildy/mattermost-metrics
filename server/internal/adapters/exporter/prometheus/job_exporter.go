package prometheus

import (
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/domain"
	"github.com/prometheus/client_golang/prometheus"
)

type JobExporter struct {
	countByStatusType *prometheus.GaugeVec
}

func newJob(registry *prometheus.Registry) JobExporter {
	return JobExporter{
		countByStatusType: newSystemGaugeWithLabels(registry, MetricsSubsystemUsage, "job", "Jobs count by status and type", []string{"type", "status"}),
	}
}

func (m JobExporter) ExportMetrics(metrics *domain.MetricsData) error {
	jobMetrics := metrics.Jobs

	m.countByStatusType.Reset()

	for _, item := range jobMetrics.CountByTypesStatus {
		m.countByStatusType.WithLabelValues(item.Type, item.Status).Set(float64(item.Count))
	}

	return nil
}
