package reports

import (
	"clubmembers/internal/members"
	"sort"
)

type FacultyCount struct {
	Faculty string
	Count   int
	Active  int
	Paused  int
}

func GroupByFaculty(records []members.MemberRecord) []FacultyCount {
	counts := make(map[string]FacultyCount)
	for _, record := range records {
		item := counts[record.Faculty]
		item.Faculty = record.Faculty
		item.Count++
		if record.Status == members.StatusActive {
			item.Active++
		} else if record.Status == members.StatusPaused {
			item.Paused++
		}
		counts[record.Faculty] = item
	}
	result := make([]FacultyCount, 0, len(counts))
	for _, item := range counts {
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Faculty < result[j].Faculty })
	return result
}

func ActiveRatio(records []members.MemberRecord) float64 {
	if len(records) == 0 {
		return 0
	}
	active := 0
	for _, record := range records {
		if record.Status == members.StatusActive {
			active++
		}
	}
	return float64(active) / float64(len(records))
}
