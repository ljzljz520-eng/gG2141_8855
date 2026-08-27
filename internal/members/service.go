package members

import "sort"

func SortByName(records []MemberRecord) []MemberRecord {
	result := CloneRecords(records)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Name == result[j].Name {
			return result[i].ID < result[j].ID
		}
		return result[i].Name < result[j].Name
	})
	return result
}

func FilterByFaculty(records []MemberRecord, faculty string) []MemberRecord {
	result := make([]MemberRecord, 0)
	for _, record := range records {
		if record.Faculty == faculty {
			result = append(result, Snapshot(record))
		}
	}
	return result
}

func FilterByStatus(records []MemberRecord, status Status) []MemberRecord {
	result := make([]MemberRecord, 0)
	for _, record := range records {
		if record.Status == status {
			result = append(result, Snapshot(record))
		}
	}
	return result
}
