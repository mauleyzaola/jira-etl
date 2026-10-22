package proxy

import (
	"strings"

	"github.com/mauleyzaola/jira-etl/domains"
)

func ExtractStatuses(log domains.JiraLog) domains.Histories {
	return domains.Histories(log.Histories).Filter(func(h domains.History) bool {
		statuses := domains.HistoryItems(h.Items).Filter(func(v domains.HistoryItem) bool {
			return strings.EqualFold(v.Field, "status")
		})
		return statuses.Len() != 0
	})
}
