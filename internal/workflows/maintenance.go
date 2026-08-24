package workflows

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"clubmembers/internal/members"
)

type MaintenanceItem struct {
	MemberID string
	Name     string
	Reason   string
	Priority int
}

func (s *Service) FindDuplicates(ctx context.Context) (map[string][]string, error) {
	rows, err := s.List(ctx, MemberFilter{})
	if err != nil {
		return nil, err
	}
	byStudent := make(map[string][]string)
	byEmail := make(map[string][]string)
	for _, record := range rows {
		student := strings.ToLower(strings.TrimSpace(record.StudentNumber))
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if student != "" {
			byStudent[student] = append(byStudent[student], record.ID)
		}
		if email != "" {
			byEmail[email] = append(byEmail[email], record.ID)
		}
	}
	duplicates := make(map[string][]string)
	for key, ids := range byStudent {
		if len(ids) > 1 {
			duplicates["student:"+key] = append([]string(nil), ids...)
		}
	}
	for key, ids := range byEmail {
		if len(ids) > 1 {
			duplicates["email:"+key] = append([]string(nil), ids...)
		}
	}
	return duplicates, nil
}

func (s *Service) PendingReviews(ctx context.Context) ([]MaintenanceItem, error) {
	rows, err := s.List(ctx, MemberFilter{})
	if err != nil {
		return nil, err
	}
	items := make([]MaintenanceItem, 0)
	for _, record := range rows {
		check := record.Validate()
		if check != nil {
			items = append(items, MaintenanceItem{MemberID: record.ID, Name: record.Name, Reason: check.Error(), Priority: 3})
			continue
		}
		if record.Status == members.StatusInvited {
			items = append(items, MaintenanceItem{MemberID: record.ID, Name: record.Name, Reason: "邀请状态待确认", Priority: 2})
			continue
		}
		if len(record.Contacts) == 0 {
			items = append(items, MaintenanceItem{MemberID: record.ID, Name: record.Name, Reason: "缺少紧急联系人", Priority: 1})
		}
	}
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Priority == items[j].Priority {
			return items[i].Name < items[j].Name
		}
		return items[i].Priority > items[j].Priority
	})
	return items, nil
}

func (s *Service) Archive(ctx context.Context, role, memberID, reason, occurredOn string) (members.MemberRecord, error) {
	if role != "admin" {
		return members.MemberRecord{}, fmt.Errorf("only admin can archive records")
	}
	status := members.StatusAlumni
	return s.Change(ctx, role, members.ChangeRequest{MemberID: memberID, Actor: role, Status: &status, Reason: reason}, occurredOn)
}

func (s *Service) Reactivate(ctx context.Context, role, memberID, reason, occurredOn string) (members.MemberRecord, error) {
	if role != "admin" && role != "editor" {
		return members.MemberRecord{}, fmt.Errorf("role %q cannot reactivate records", role)
	}
	status := members.StatusActive
	return s.Change(ctx, role, members.ChangeRequest{MemberID: memberID, Actor: role, Status: &status, Reason: reason}, occurredOn)
}

func (s *Service) ChangeTags(ctx context.Context, role, memberID string, add, remove []string) (members.MemberRecord, error) {
	if role != "admin" && role != "editor" {
		return members.MemberRecord{}, fmt.Errorf("role %q cannot change tags", role)
	}
	record, err := s.Get(ctx, memberID)
	if err != nil {
		return record, err
	}
	roster, err := members.NewRoster([]members.MemberRecord{record})
	if err != nil {
		return record, err
	}
	for _, tag := range add {
		if err := roster.AddTag(memberID, tag); err != nil {
			return record, err
		}
	}
	for _, tag := range remove {
		if err := roster.RemoveTag(memberID, tag); err != nil {
			return record, err
		}
	}
	updated, ok := roster.Find(memberID)
	if !ok {
		return record, fmt.Errorf("member disappeared from roster")
	}
	if err := s.store.SaveMember(ctx, updated); err != nil {
		return record, err
	}
	return updated, nil
}
