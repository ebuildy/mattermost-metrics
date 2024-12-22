package ports

import (
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/domain"
	"net/http"
)

type MetricsCollector interface {
	CollectMetrics(metrics *domain.MetricsData) error
}

type MetricsExporter interface {
	ExportMetrics(metrics *domain.MetricsData) error
	ServeMetrics(w http.ResponseWriter, r *http.Request)
}
