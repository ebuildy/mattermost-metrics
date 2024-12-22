package prometheus

import (
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/domain"
	"github.com/prometheus/client_golang/prometheus"
)

type ReactionExporter struct {
	countByEmoji *prometheus.GaugeVec
	last         prometheus.Gauge
}

func newReaction(registry *prometheus.Registry) ReactionExporter {
	return ReactionExporter{
		countByEmoji: newSystemGaugeWithLabels(registry, MetricsSubsystemKPI, "reaction_total", "Count by emoji (top 5)", []string{"emoji"}),
		last:         newSystemGauge(registry, MetricsSubsystemKPI, "reaction_last_seconds", "Last reaction time - unix timestamp"),
	}
}

func (m ReactionExporter) ExportMetrics(metrics *domain.MetricsData) error {
	reactionMetrics := metrics.Reactions

	m.countByEmoji.Reset()

	if reactionMetrics != nil {
		for _, item := range reactionMetrics.CountByEmoji {
			m.countByEmoji.WithLabelValues(item.Emoji).Set(float64(item.Count))
		}

		m.last.Set(float64(reactionMetrics.Last.Unix()))
	}

	return nil
}
