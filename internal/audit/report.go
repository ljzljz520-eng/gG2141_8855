package audit

import (
	"sort"
	"strings"

	"clubmembers/internal/members"
)

type FieldCount struct {
	Field string
	Count int
}

type MemberActivity struct {
	MemberID string
	Changes  int
	Fields   []string
}

func CountFields(events []members.ChangeEvent) []FieldCount {
	counts := make(map[string]int)
	for _, event := range events {
		counts[event.Field]++
	}
	result := make([]FieldCount, 0, len(counts))
	for field, count := range counts {
		result = append(result, FieldCount{Field: field, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count == result[j].Count {
			return result[i].Field < result[j].Field
		}
		return result[i].Count > result[j].Count
	})
	return result
}

func ActivityByMember(events []members.ChangeEvent) []MemberActivity {
	byMember := make(map[string]*MemberActivity)
	for _, event := range events {
		activity := byMember[event.MemberID]
		if activity == nil {
			activity = &MemberActivity{MemberID: event.MemberID, Fields: make([]string, 0)}
			byMember[event.MemberID] = activity
		}
		activity.Changes++
		known := false
		for _, field := range activity.Fields {
			if field == event.Field {
				known = true
				break
			}
		}
		if !known {
			activity.Fields = append(activity.Fields, event.Field)
		}
	}
	result := make([]MemberActivity, 0, len(byMember))
	for _, activity := range byMember {
		sort.Strings(activity.Fields)
		result = append(result, *activity)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].MemberID < result[j].MemberID })
	return result
}

func FilterActor(events []members.ChangeEvent, actor string) []members.ChangeEvent {
	result := make([]members.ChangeEvent, 0)
	actor = strings.TrimSpace(actor)
	for _, event := range events {
		if actor == "" || event.Actor == actor {
			result = append(result, event)
		}
	}
	return Chronological(result)
}

func LatestPerField(events []members.ChangeEvent) map[string]members.ChangeEvent {
	latest := make(map[string]members.ChangeEvent)
	for _, event := range Chronological(events) {
		latest[event.Field] = event
	}
	return latest
}

func HasSensitiveChange(events []members.ChangeEvent) bool {
	for _, event := range events {
		if event.Field == "phone" || event.Field == "email" {
			return true
		}
	}
	return false
}
