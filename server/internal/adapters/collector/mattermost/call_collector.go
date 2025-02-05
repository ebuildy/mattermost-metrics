package mattermost

import "github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/domain"

func (c *Collector) collectCall() *domain.MetricsDataCalls {
	return &domain.MetricsDataCalls{
		Count:    c.gateway.SQLValue("SELECT COUNT(*) FROM calls"),
		Duration: c.gateway.SQLValue("SELECT FLOOR(SUM(endat - startat)/1000) FROM calls WHERE startat IS NOT NULL AND endat IS NOT NULL"),
		Last:     c.gateway.SQLValue("SELECT MAX(startat) FROM calls"),
	}
}
