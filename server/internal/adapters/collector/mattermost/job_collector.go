package mattermost

import "github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/domain"

func (c *Collector) collectJob() *domain.MetricsDataJobs {
	var countByTypesStatus []domain.JobCountByStatusType

	rows, err := c.gateway.DB.Query("SELECT type, status, COUNT(*) FROM jobs GROUP BY type, status")

	if err != nil {
		c.logger.Error("failed to get jobs:", "error", err.Error())
	} else {
		for rows.Next() {
			var typeJob string
			var status string
			var count int64
			if err := rows.Scan(&typeJob, &status, &count); err != nil {
				c.logger.Error("failed to scan jobs:", "error", err.Error())
			}

			countByTypesStatus = append(countByTypesStatus, domain.JobCountByStatusType{
				Status: status,
				Type:   typeJob,
				Count:  count,
			})
		}
	}

	return &domain.MetricsDataJobs{
		CountByTypesStatus: countByTypesStatus,
	}
}
