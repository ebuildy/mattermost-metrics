package ports

import (
	"net/http"

	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/domain"
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
