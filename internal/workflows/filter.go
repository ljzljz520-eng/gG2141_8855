package workflows

import (
	"clubmembers/internal/members"
	"strings"
)

type MemberFilter struct {
	Faculty string
	Status  members.Status
	Tag     string
	Search  string
}

func (f MemberFilter) Matches(record members.MemberRecord) bool {
	if f.Faculty != "" && record.Faculty != f.Faculty {
		return false
	}
	if f.Status != "" && record.Status != f.Status {
		return false
	}
	if f.Tag != "" && !record.HasTag(f.Tag) {
		return false
	}
	if f.Search != "" {
		query := strings.ToLower(strings.TrimSpace(f.Search))
		if !strings.Contains(strings.ToLower(record.Name), query) && !strings.Contains(strings.ToLower(record.StudentNumber), query) && !strings.Contains(strings.ToLower(record.Phone), query) {
			return false
		}
	}
	return true
}
