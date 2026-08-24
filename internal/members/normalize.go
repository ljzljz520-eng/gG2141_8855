package members

import (
	"regexp"
	"strings"
)

var phoneDigits = regexp.MustCompile(`[^0-9+]`)

func NormalizeMember(input MemberRecord) MemberRecord {
	input.ID = strings.TrimSpace(input.ID)
	input.Name = strings.TrimSpace(input.Name)
	input.StudentNumber = strings.TrimSpace(input.StudentNumber)
	input.Faculty = strings.TrimSpace(input.Faculty)
	input.Phone = normalizePhone(input.Phone)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.JoinedOn = strings.TrimSpace(input.JoinedOn)
	input.Status = Status(strings.ToLower(strings.TrimSpace(string(input.Status))))
	input.Tags = normalizeTags(input.Tags)
	input.Contacts = normalizeContacts(input.Contacts)
	input.Notes = normalizeNotes(input.Notes)
	return input
}

func normalizePhone(phone string) string {
	phone = strings.TrimSpace(phone)
	if strings.HasPrefix(phone, "00") {
		phone = "+" + strings.TrimPrefix(phone, "00")
	}
	return phoneDigits.ReplaceAllString(phone, "")
}

func normalizeTags(tags []string) []string {
	seen := make(map[string]bool, len(tags))
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.ToLower(strings.TrimSpace(tag))
		if tag == "" || seen[tag] {
			continue
		}
		seen[tag] = true
		result = append(result, tag)
	}
	return result
}

func normalizeContacts(contacts []ContactPoint) []ContactPoint {
	result := make([]ContactPoint, 0, len(contacts))
	preferred := false
	for _, contact := range contacts {
		contact.Label = strings.TrimSpace(contact.Label)
		contact.Name = strings.TrimSpace(contact.Name)
		contact.Phone = normalizePhone(contact.Phone)
		contact.Relationship = strings.TrimSpace(contact.Relationship)
		if contact.Preferred && preferred {
			contact.Preferred = false
		}
		if contact.Preferred {
			preferred = true
		}
		if contact.Name != "" && contact.Phone != "" {
			result = append(result, contact)
		}
	}
	return result
}

func normalizeNotes(notes []string) []string {
	result := make([]string, 0, len(notes))
	for _, note := range notes {
		note = strings.TrimSpace(note)
		if note != "" {
			result = append(result, note)
		}
	}
	return result
}
