package validation

import (
	"sort"
	"strings"

	"clubmembers/internal/members"
)

type FieldIssue struct {
	Field    string
	Severity string
	Message  string
}

type Completeness struct {
	MemberID string
	Score    int
	Issues   []FieldIssue
}

func CheckCompleteness(record members.MemberRecord) Completeness {
	issues := make([]FieldIssue, 0)
	score := 0
	if strings.TrimSpace(record.Name) == "" {
		issues = append(issues, FieldIssue{Field: "name", Severity: "error", Message: "姓名不能为空"})
	} else {
		score++
	}
	if strings.TrimSpace(record.StudentNumber) == "" {
		issues = append(issues, FieldIssue{Field: "student_number", Severity: "error", Message: "学号不能为空"})
	} else {
		score++
	}
	if strings.TrimSpace(record.Faculty) == "" {
		issues = append(issues, FieldIssue{Field: "faculty", Severity: "error", Message: "学院不能为空"})
	} else {
		score++
	}
	if strings.TrimSpace(record.Phone) == "" {
		issues = append(issues, FieldIssue{Field: "phone", Severity: "error", Message: "手机号不能为空"})
	} else {
		score++
	}
	if strings.TrimSpace(record.Email) == "" {
		issues = append(issues, FieldIssue{Field: "email", Severity: "warning", Message: "邮箱尚未填写"})
	} else {
		score++
	}
	if len(record.Contacts) == 0 {
		issues = append(issues, FieldIssue{Field: "contacts", Severity: "warning", Message: "建议填写紧急联系人"})
	} else {
		score++
	}
	return Completeness{MemberID: record.ID, Score: score, Issues: issues}
}

func (c Completeness) Complete() bool {
	for _, issue := range c.Issues {
		if issue.Severity == "error" {
			return false
		}
	}
	return true
}

func (c Completeness) Warnings() []FieldIssue {
	result := make([]FieldIssue, 0)
	for _, issue := range c.Issues {
		if issue.Severity == "warning" {
			result = append(result, issue)
		}
	}
	return result
}

func SortIssues(issues []FieldIssue) []FieldIssue {
	result := append([]FieldIssue(nil), issues...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Severity == result[j].Severity {
			return result[i].Field < result[j].Field
		}
		return result[i].Severity == "error"
	})
	return result
}

func MissingFields(record members.MemberRecord) []string {
	check := CheckCompleteness(record)
	missing := make([]string, 0, len(check.Issues))
	for _, issue := range check.Issues {
		missing = append(missing, issue.Field)
	}
	sort.Strings(missing)
	return missing
}

func CompareCompleteness(before, after members.MemberRecord) []string {
	oldMissing := make(map[string]bool)
	for _, field := range MissingFields(before) {
		oldMissing[field] = true
	}
	newMissing := make(map[string]bool)
	for _, field := range MissingFields(after) {
		newMissing[field] = true
	}
	resolved := make([]string, 0)
	for field := range oldMissing {
		if !newMissing[field] {
			resolved = append(resolved, field)
		}
	}
	sort.Strings(resolved)
	return resolved
}
