package members

import (
	"errors"
	"fmt"
	"strings"
)

type ChangeRequest struct {
	MemberID string
	Actor    string
	Phone    *string
	Email    *string
	Faculty  *string
	Status   *Status
	Reason   string
}

func ApplyChange(member MemberRecord, request ChangeRequest) (MemberRecord, []ChangeEvent, error) {
	if strings.TrimSpace(request.Actor) == "" {
		return member, nil, errors.New("actor is required")
	}
	if request.MemberID != member.ID {
		return member, nil, errors.New("member id mismatch")
	}
	updated := Snapshot(member)
	events := make([]ChangeEvent, 0, 4)
	if request.Phone != nil {
		phone := normalizePhone(*request.Phone)
		if phone != member.Phone {
			events = append(events, ChangeEvent{MemberID: member.ID, Actor: request.Actor, Field: "phone", Before: member.Phone, After: phone, Reason: request.Reason})
			updated.Phone = phone
		}
	}
	if request.Email != nil {
		email := strings.ToLower(strings.TrimSpace(*request.Email))
		if email != member.Email {
			events = append(events, ChangeEvent{MemberID: member.ID, Actor: request.Actor, Field: "email", Before: member.Email, After: email, Reason: request.Reason})
			updated.Email = email
		}
	}
	if request.Faculty != nil {
		faculty := strings.TrimSpace(*request.Faculty)
		if faculty != member.Faculty {
			events = append(events, ChangeEvent{MemberID: member.ID, Actor: request.Actor, Field: "faculty", Before: member.Faculty, After: faculty, Reason: request.Reason})
			updated.Faculty = faculty
		}
	}
	if request.Status != nil {
		if !validStatus(*request.Status) {
			return member, nil, fmt.Errorf("unsupported status %q", *request.Status)
		}
		if *request.Status != member.Status {
			events = append(events, ChangeEvent{MemberID: member.ID, Actor: request.Actor, Field: "status", Before: string(member.Status), After: string(*request.Status), Reason: request.Reason})
			updated.Status = *request.Status
		}
	}
	if len(events) == 0 {
		return member, nil, errors.New("change request has no differences")
	}
	updated.Revision++
	return updated, events, nil
}
