package validation

import (
	"fmt"
	"regexp"
	"strings"

	"clubmembers/internal/members"
)

var emailPattern = regexp.MustCompile(`^[^@[:space:]]+@[^@[:space:]]+\.[^@[:space:]]+$`)

func ValidateRecord(record members.MemberRecord) []string {
	issues := make([]string, 0)
	if err := record.Validate(); err != nil {
		issues = append(issues, err.Error())
	}
	if record.Email != "" && !emailPattern.MatchString(record.Email) {
		issues = append(issues, "email format is invalid")
	}
	if len(record.Phone) < 6 {
		issues = append(issues, "phone is too short")
	}
	for i, contact := range record.Contacts {
		if strings.TrimSpace(contact.Name) == "" || strings.TrimSpace(contact.Phone) == "" {
			issues = append(issues, fmt.Sprintf("contact %d is incomplete", i+1))
		}
	}
	return issues
}

func IsReadyForActivation(record members.MemberRecord) bool {
	return record.Name != "" && record.Phone != "" && record.Email != "" && len(ValidateRecord(record)) == 0
}

func RequireReason(reason string) error {
	if len(strings.TrimSpace(reason)) < 6 {
		return fmt.Errorf("change reason must be at least 6 characters")
	}
	return nil
}
