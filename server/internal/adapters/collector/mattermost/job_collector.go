package mattermost

import (
	"time"

	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/domain"
)

// collectJob retrieves job data from the database, including job counts grouped by type and status, and returns metric data.
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
		Last:               time.Unix(c.gateway.SQLValue("SELECT MAX(createat) FROM jobs"), 0),
	}
}
