package ports

import (
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/domain"
	"net/http"
)

type MetricsCollector interface {
	Configure(config *domain.ConfigCollector)
	CollectMetrics(metrics *domain.MetricsData) error
}

type MetricsExporter interface {
	ExportMetrics(metrics *domain.MetricsData) error
}

type MetricsHandler interface {
	MetricsExporter
	ServeMetrics(w http.ResponseWriter, r *http.Request)
}
