package reports

import (
	"fmt"
	"sort"
	"strings"

	"clubmembers/internal/members"
)

type DirectoryLine struct {
	Faculty string
	Name    string
	Phone   string
	Status  members.Status
}

func SortForExport(records []members.MemberRecord) []members.MemberRecord {
	result := members.CloneRecords(records)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Faculty == result[j].Faculty {
			if result[i].Status == result[j].Status {
				return result[i].Name < result[j].Name
			}
			return result[i].Status < result[j].Status
		}
		return result[i].Faculty < result[j].Faculty
	})
	return result
}

func ContactDirectory(records []members.MemberRecord) []DirectoryLine {
	result := make([]DirectoryLine, 0, len(records))
	for _, record := range SortForExport(records) {
		result = append(result, DirectoryLine{Faculty: record.Faculty, Name: record.Name, Phone: record.Phone, Status: record.Status})
	}
	return result
}

func MarkdownSummary(records []members.MemberRecord) string {
	counts := GroupByFaculty(records)
	var builder strings.Builder
	builder.WriteString("# 成员档案摘要\n\n")
	builder.WriteString(fmt.Sprintf("成员总数：%d\n\n", len(records)))
	builder.WriteString("| 学院 | 总数 | 活跃 | 暂停 |\n| --- | ---: | ---: | ---: |\n")
	for _, count := range counts {
		builder.WriteString(fmt.Sprintf("| %s | %d | %d | %d |\n", count.Faculty, count.Count, count.Active, count.Paused))
	}
	return builder.String()
}

func SearchText(records []members.MemberRecord, query string) []string {
	query = strings.ToLower(strings.TrimSpace(query))
	labels := make([]string, 0)
	for _, record := range SortForExport(records) {
		if query == "" || strings.Contains(strings.ToLower(record.Name), query) || strings.Contains(strings.ToLower(record.Email), query) {
			labels = append(labels, record.DisplayLabel())
		}
	}
	return labels
}

func StatusLabels(records []members.MemberRecord) map[members.Status][]string {
	result := make(map[members.Status][]string)
	for _, record := range records {
		result[record.Status] = append(result[record.Status], record.Name)
	}
	for status := range result {
		sort.Strings(result[status])
	}
	return result
}
