package mattermost

import "github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/domain"

func (c *Collector) collectKPI() *domain.MetricsDataKPI {
	return &domain.MetricsDataKPI{
		KPILastPostDate:             c.gateway.SQLValue("SELECT MAX(createat) FROM posts"),
		KPIChannelsLastCreationDate: c.gateway.SQLValue("SELECT MAX(createat) FROM channels"),
		KPIPostsCount:               c.gateway.SQLValue("SELECT COUNT(*) FROM posts"),
		KPISessionsCount:            c.gateway.SQLValue("SELECT COUNT(*) FROM sessions"),
		KPIPrivateChannelsCount:     c.gateway.SQLValue("SELECT COUNT(*) FROM channels WHERE type = 'P'"),
		KPIPublicChannelsCount:      c.gateway.SQLValue("SELECT COUNT(*) FROM channels WHERE type = 'O'"),
		KPIDirectMessagesCount:      c.gateway.SQLValue("SELECT COUNT(*) FROM channels WHERE type = 'D'"),
	}
}
