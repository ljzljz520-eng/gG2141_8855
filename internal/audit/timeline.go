package audit

import (
	"sort"
	"strings"

	"clubmembers/internal/members"
)

func Chronological(events []members.ChangeEvent) []members.ChangeEvent {
	result := append([]members.ChangeEvent(nil), events...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].OccurredOn == result[j].OccurredOn {
			return result[i].ID < result[j].ID
		}
		return result[i].OccurredOn < result[j].OccurredOn
	})
	return result
}

func HumanLine(event members.ChangeEvent) string {
	reason := strings.TrimSpace(event.Reason)
	if reason == "" {
		reason = "未填写原因"
	}
	return event.OccurredOn + " " + event.Actor + " 将 " + event.Field + " 从 " + event.Before + " 改为 " + event.After + "（" + reason + "）"
}
