package reports

import (
	"sort"
	"strings"

	"clubmembers/internal/members"
)

type FilterSummary struct {
	Label string
	Value string
	Count int
}

func SummarizeTags(records []members.MemberRecord) []FilterSummary {
	counts := make(map[string]int)
	for _, record := range records {
		seen := make(map[string]bool)
		for _, tag := range record.Tags {
			tag = strings.TrimSpace(tag)
			if tag != "" && !seen[tag] {
				counts[tag]++
				seen[tag] = true
			}
		}
	}
	result := make([]FilterSummary, 0, len(counts))
	for tag, count := range counts {
		result = append(result, FilterSummary{Label: "tag", Value: tag, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count == result[j].Count {
			return result[i].Value < result[j].Value
		}
		return result[i].Count > result[j].Count
	})
	return result
}

func SummarizeStatuses(records []members.MemberRecord) []FilterSummary {
	counts := make(map[members.Status]int)
	for _, record := range records {
		counts[record.Status]++
	}
	result := make([]FilterSummary, 0, len(counts))
	for status, count := range counts {
		result = append(result, FilterSummary{Label: "status", Value: string(status), Count: count})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Value < result[j].Value })
	return result
}

func RestrictToActive(records []members.MemberRecord) []members.MemberRecord {
	result := make([]members.MemberRecord, 0)
	for _, record := range records {
		if record.Status == members.StatusActive {
			result = append(result, members.Snapshot(record))
		}
	}
	return result
}
