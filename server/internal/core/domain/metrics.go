package domain

import (
	"database/sql"
	"time"
)

// MetricsData represents metrics data between collector and exporter
type MetricsData struct {
	UsagePostsCount, UsageUsersCount, UsageStorage int64

	SQLStats sql.DBStats

	Info   *MetricsDataInfo
	Health *MetricsDataHealth
	KPI    *MetricsDataKPI
	Jobs   *MetricsDataJobs
}

type MetricsDataInfo struct {
	MattermostVersion, MattermostEdition string
	MattermostInstallationTime           time.Time

	SQLDriverName string
}

type MetricsDataUsage struct{}

type MetricsDataHealth struct {
	SystemHealth, SystemHealthDatabase, SystemHealthFilestore bool
}

type MetricsDataKPI struct {
	KPILastPostDate, KPIChannelsLastCreationDate, KPIPostsCount, KPISessionsCount int64
	KPIPrivateChannelsCount, KPIPublicChannelsCount, KPIDirectMessagesCount       int64
}

type MetricsDataJobs struct {
	CountByTypesStatus []JobCountByStatusType
}

type JobCountByStatusType struct {
	Status, Type string
	Count        int64
}
