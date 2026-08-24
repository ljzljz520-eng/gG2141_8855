package audit

import (
	"sort"
	"strings"

	"clubmembers/internal/members"
)

type RetentionView struct {
	Events      []members.ChangeEvent
	FieldCounts []FieldCount
	Actors      []string
}

func BuildRetentionView(events []members.ChangeEvent) RetentionView {
	ordered := Chronological(events)
	actors := make(map[string]bool)
	for _, event := range ordered {
		actor := strings.TrimSpace(event.Actor)
		if actor != "" {
			actors[actor] = true
		}
	}
	actorList := make([]string, 0, len(actors))
	for actor := range actors {
		actorList = append(actorList, actor)
	}
	sort.Strings(actorList)
	return RetentionView{Events: ordered, FieldCounts: CountFields(ordered), Actors: actorList}
}

func KeepLatest(events []members.ChangeEvent, limit int) []members.ChangeEvent {
	if limit <= 0 {
		return nil
	}
	ordered := Chronological(events)
	if len(ordered) <= limit {
		return ordered
	}
	return append([]members.ChangeEvent(nil), ordered[len(ordered)-limit:]...)
}

func FieldsChanged(events []members.ChangeEvent) []string {
	latest := LatestPerField(events)
	fields := make([]string, 0, len(latest))
	for field := range latest {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	return fields
}
