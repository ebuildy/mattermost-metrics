package exporter

import (
	"net/http"

	"github.com/ebuildy/mattermost-plugin-minotor/server/controller"
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
	usagePostCount, usageUsersCount, usageStorage             prometheus.Gauge
	systemHealth, systemDatabaseHealth, systemFilestoreHealth prometheus.Gauge
	kpiPostsCount, kpiChannelsCount, kpiLastPostDate          prometheus.Gauge
}

func NewExporter() *Exporter {
	reg := prometheus.NewRegistry()

	// reg.MustRegister(
	//	controller.NewGoCollector(),
	//	controller.NewProcessCollector(controller.ProcessCollectorOpts{}),
	//)

	metrics := &metrics{
		usagePostCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   MetricsNamespace,
			Subsystem:   MetricsSubsystemUsage,
			Name:        "posts_total",
			Help:        "Total number of posts",
			ConstLabels: map[string]string{"group": "mattermost"},
		}),
		usageUsersCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   MetricsNamespace,
			Subsystem:   MetricsSubsystemUsage,
			Name:        "users_total",
			Help:        "Total number of users",
			ConstLabels: nil,
		}),
		usageStorage: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   MetricsNamespace,
			Subsystem:   MetricsSubsystemUsage,
			Name:        "storage_bytes",
			Help:        "Storage usage in bytes",
			ConstLabels: nil,
		}),

		systemHealth: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   MetricsNamespace,
			Subsystem:   MetricsSubsystemSystem,
			Name:        "status",
			Help:        "Global status",
			ConstLabels: nil,
		}),
		systemDatabaseHealth: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   MetricsNamespace,
			Subsystem:   MetricsSubsystemSystem,
			Name:        "database_status",
			Help:        "Database component status",
			ConstLabels: nil,
		}),
		systemFilestoreHealth: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   MetricsNamespace,
			Subsystem:   MetricsSubsystemSystem,
			Name:        "filestore_status",
			Help:        "Filestore component status",
			ConstLabels: nil,
		}),

		kpiPostsCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   MetricsNamespace,
			Subsystem:   MetricsSubsystemKPI,
			Name:        "posts_total",
			Help:        "Total number of posts",
			ConstLabels: nil,
		}),
		kpiChannelsCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   MetricsNamespace,
			Subsystem:   MetricsSubsystemKPI,
			Name:        "channels_total",
			Help:        "Total number of channels",
			ConstLabels: nil,
		}),
		kpiLastPostDate: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   MetricsNamespace,
			Subsystem:   MetricsSubsystemKPI,
			Name:        "last_post_date",
			Help:        "Timestamp of last post date",
			ConstLabels: nil,
		}),
	}

	reg.MustRegister(metrics.usagePostCount)
	reg.MustRegister(metrics.usageUsersCount)
	reg.MustRegister(metrics.usageStorage)

	reg.MustRegister(metrics.systemHealth)
	reg.MustRegister(metrics.systemDatabaseHealth)
	reg.MustRegister(metrics.systemFilestoreHealth)

	reg.MustRegister(metrics.kpiLastPostDate)
	reg.MustRegister(metrics.kpiPostsCount)
	reg.MustRegister(metrics.kpiChannelsCount)

	return &Exporter{
		reg,
		metrics,
		promhttp.HandlerFor(reg, promhttp.HandlerOpts{
			EnableOpenMetrics: false,
			Registry:          reg,
		}),
	}
}

func (o *Exporter) ExportMetrics(metrics *controller.Metrics) error {
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
