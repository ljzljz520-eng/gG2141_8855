package members

import (
	"errors"
	"fmt"
	"strings"
)

type Status string

const (
	StatusActive  Status = "active"
	StatusPaused  Status = "paused"
	StatusAlumni  Status = "alumni"
	StatusInvited Status = "invited"
)

type MemberRecord struct {
	ID            string
	Name          string
	StudentNumber string
	Faculty       string
	Phone         string
	Email         string
	JoinedOn      string
	Status        Status
	Tags          []string
	Contacts      []ContactPoint
	Notes         []string
	Revision      int
}

type ContactPoint struct {
	Label        string
	Name         string
	Phone        string
	Relationship string
	Preferred    bool
}

type ClubRole struct {
	Name      string
	Scope     string
	GrantedOn string
	Active    bool
}

type ChangeEvent struct {
	ID         string
	MemberID   string
	Actor      string
	Field      string
	Before     string
	After      string
	OccurredOn string
	Reason     string
}

func (m MemberRecord) Validate() error {
	if strings.TrimSpace(m.ID) == "" {
		return errors.New("member id is required")
	}
	if strings.TrimSpace(m.Name) == "" {
		return errors.New("member name is required")
	}
	if strings.TrimSpace(m.StudentNumber) == "" {
		return errors.New("student number is required")
	}
	if strings.TrimSpace(m.Faculty) == "" {
		return errors.New("faculty is required")
	}
	if !validStatus(m.Status) {
		return fmt.Errorf("unsupported status %q", m.Status)
	}
	if len(m.Tags) > 12 {
		return errors.New("a member can have at most 12 tags")
	}
	return nil
}

func validStatus(status Status) bool {
	switch status {
	case StatusActive, StatusPaused, StatusAlumni, StatusInvited:
		return true
	default:
		return false
	}
}

func (m MemberRecord) HasTag(tag string) bool {
	tag = strings.ToLower(strings.TrimSpace(tag))
	for _, candidate := range m.Tags {
		if strings.ToLower(candidate) == tag {
			return true
		}
	}
	return false
}

func (m MemberRecord) DisplayLabel() string {
	if m.StudentNumber == "" {
		return m.Name
	}
	return fmt.Sprintf("%s (%s)", m.Name, m.StudentNumber)
}

func (m MemberRecord) ContactSummary() string {
	if len(m.Contacts) == 0 {
		return "no emergency contact"
	}
	preferred := m.Contacts[0]
	for _, contact := range m.Contacts {
		if contact.Preferred {
			preferred = contact
			break
		}
	}
	return fmt.Sprintf("%s: %s", preferred.Name, preferred.Phone)
}
