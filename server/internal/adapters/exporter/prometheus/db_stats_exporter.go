package prometheus

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

type DBStatsExporter struct {
	MaxOpenConnections prometheus.Gauge

	// Pool Status
	OpenConnections prometheus.Gauge
	InUse           prometheus.Gauge
	Idle            prometheus.Gauge

	// Counters
	WaitCount         prometheus.Gauge
	WaitDuration      prometheus.Gauge
	MaxIdleClosed     prometheus.Gauge
	MaxIdleTimeClosed prometheus.Gauge
	MaxLifetimeClosed prometheus.Gauge
}

func newDBStats(registry *prometheus.Registry) *DBStatsExporter {
	return &DBStatsExporter{
		MaxOpenConnections: newSystemGauge(registry, "db", "max_open_connections", "Maximum number of open connections to the database."),
		OpenConnections:    newSystemGauge(registry, "db", "open_connections", "The number of established connections both in use and idle."),
		InUse:              newSystemGauge(registry, "db", "connections_in_use", "The number of connections currently in use."),
		Idle:               newSystemGauge(registry, "db", "connections_idle", "The number of idle connections."),
		WaitCount:          newSystemGauge(registry, "db", "wait_count", "The total number of connections waited for."),
		WaitDuration:       newSystemGauge(registry, "db", "wait_duration_seconds", "The total time blocked waiting for a new connection (seconds)."),
		MaxIdleClosed:      newSystemGauge(registry, "db", "idle_connections_closed_total", "The total number of connections closed due to SetMaxIdleConns."),
		MaxIdleTimeClosed:  newSystemGauge(registry, "db", "idle_time_connections_closed_total", "The total number of connections closed due to SetConnMaxIdleTime."),
		MaxLifetimeClosed:  newSystemGauge(registry, "db", "lifetime_connections_closed_total", "The total number of connections closed due to SetConnMaxLifetime."),
	}
}

func (e *DBStatsExporter) bindDBStats(dbStats sql.DBStats) {
	e.MaxOpenConnections.Set(float64(dbStats.MaxOpenConnections))
	e.OpenConnections.Set(float64(dbStats.OpenConnections))
	e.InUse.Set(float64(dbStats.InUse))
	e.Idle.Set(float64(dbStats.Idle))
	e.WaitCount.Set(float64(dbStats.WaitCount))
	e.WaitDuration.Set(float64(dbStats.WaitDuration.Seconds()))
	e.MaxIdleClosed.Set(float64(dbStats.MaxIdleClosed))
	e.MaxIdleTimeClosed.Set(float64(dbStats.MaxIdleTimeClosed))
	e.MaxLifetimeClosed.Set(float64(dbStats.MaxLifetimeClosed))
}
