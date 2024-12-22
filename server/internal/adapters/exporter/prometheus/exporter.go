package prometheus

import (
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/domain"
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/ports"
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

	exporters []ports.MetricsExporter
}

type metrics struct {
	dbStats *DBStatsExporter
}

func NewExporter() *Exporter {
	registry := prometheus.NewRegistry()

	metrics := &metrics{
		dbStats: newDBStats(registry),
	}

	return &Exporter{
		registry,
		metrics,
		promhttp.HandlerFor(registry, promhttp.HandlerOpts{
			EnableOpenMetrics: false,
			Registry:          registry,
		}),
		[]ports.MetricsExporter{
			newHealth(registry),
			newInfo(registry),
			newKPI(registry),
			newJob(registry),
			newReaction(registry),
		},
	}
}

func (o *Exporter) ExportMetrics(metrics *domain.MetricsData) error {

	o.metrics.dbStats.bindDBStats(metrics.SQLStats)

	for _, exporter := range o.exporters {
		err := exporter.ExportMetrics(metrics)

		if err != nil {

		}
	}

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
