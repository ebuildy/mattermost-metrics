package domain

import "time"

// MetricsData represents metrics data between collector and exporter
type MetricsData struct {
	UsagePostsCount, UsageUsersCount, UsageStorage int64

	SystemHealth, SystemHealthDatabase, SystemHealthFilestore bool

	KPILastPostDate, KPIChannelsLastCreationDate, KPIChannelsCount, KPIPostsCount int64

	MattermostVersion, MattermostEdition string
	MattermostInstallationTime           time.Time

	SQLDriverName string
}
