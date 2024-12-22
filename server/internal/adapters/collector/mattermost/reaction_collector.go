package mattermost

import (
	"fmt"
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/domain"
	"time"
)

func (c *Collector) collectReaction(countByEmojiLimits int) *domain.MetricsDataReactions {
	var countByEmoji []domain.ReactionCountByEmoji

	query := fmt.Sprintf("SELECT emojiname, COUNT(*) FROM reactions GROUP BY emojiname ORDER BY COUNT(*) DESC LIMIT %d", countByEmojiLimits)
	rows, err := c.gateway.DB.Query(query)

	if err != nil {
		c.logger.Error("failed to get jobs:", "error", err.Error())
	} else {
		for rows.Next() {
			var emoji string
			var count int64
			if err := rows.Scan(&emoji, &count); err != nil {
				c.logger.Error("failed to scan jobs:", "error", err.Error())
			}

			countByEmoji = append(countByEmoji, domain.ReactionCountByEmoji{
				Emoji: emoji,
				Count: count,
			})
		}
	}

	countByEmoji = append(countByEmoji, domain.ReactionCountByEmoji{
		Emoji: "",
		Count: c.gateway.SQLValue("SELECT COUNT(*) FROM reactions"),
	})

	return &domain.MetricsDataReactions{
		CountByEmoji: countByEmoji,
		Last:         time.Unix(c.gateway.SQLValue("SELECT MAX(createat) FROM reactions"), 0),
	}
}
