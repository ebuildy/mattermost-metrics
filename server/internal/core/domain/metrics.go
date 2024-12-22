package domain

import (
	"database/sql"
	"time"
)

// MetricsData represents metrics data between collector and exporter
type MetricsData struct {
	UsagePostsCount, UsageUsersCount, UsageStorage int64

	SQLStats sql.DBStats

	Info      *MetricsDataInfo
	Health    *MetricsDataHealth
	KPI       *MetricsDataKPI
	Jobs      *MetricsDataJobs
	Reactions *MetricsDataReactions
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
	Last               time.Time
}

type JobCountByStatusType struct {
	Type, Status string
	Count        int64
}

type MetricsDataReactions struct {
	CountByEmoji []ReactionCountByEmoji
	Last         time.Time
}

type ReactionCountByEmoji struct {
	Emoji string
	Count int64
}
