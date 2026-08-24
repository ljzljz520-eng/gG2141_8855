package workflows

import (
	"context"
	"fmt"

	"clubmembers/internal/members"
)

func (s *Service) AddEmergencyContact(ctx context.Context, role, memberID string, contact members.ContactPoint, reason, occurredOn string) (members.MemberRecord, error) {
	if role != "admin" && role != "editor" {
		return members.MemberRecord{}, fmt.Errorf("role %q cannot edit contacts", role)
	}
	record, err := s.Get(ctx, memberID)
	if err != nil {
		return record, err
	}
	book := members.NewContactBook(record.Contacts)
	if err := book.Add(contact); err != nil {
		return record, err
	}
	record.Contacts = book.All()
	record.Revision++
	if err := s.store.SaveMember(ctx, record); err != nil {
		return record, err
	}
	return record, nil
}

func (s *Service) ReplaceEmergencyContact(ctx context.Context, role, memberID string, index int, contact members.ContactPoint, reason, occurredOn string) (members.MemberRecord, error) {
	if role != "admin" && role != "editor" {
		return members.MemberRecord{}, fmt.Errorf("role %q cannot edit contacts", role)
	}
	record, err := s.Get(ctx, memberID)
	if err != nil {
		return record, err
	}
	book := members.NewContactBook(record.Contacts)
	if err := book.Replace(index, contact); err != nil {
		return record, err
	}
	record.Contacts = book.All()
	record.Revision++
	if err := s.store.SaveMember(ctx, record); err != nil {
		return record, err
	}
	return record, nil
}

func (s *Service) RemoveEmergencyContact(ctx context.Context, role, memberID string, index int, reason, occurredOn string) (members.MemberRecord, error) {
	if role != "admin" && role != "editor" {
		return members.MemberRecord{}, fmt.Errorf("role %q cannot edit contacts", role)
	}
	record, err := s.Get(ctx, memberID)
	if err != nil {
		return record, err
	}
	book := members.NewContactBook(record.Contacts)
	if err := book.Remove(index); err != nil {
		return record, err
	}
	record.Contacts = book.All()
	record.Revision++
	if err := s.store.SaveMember(ctx, record); err != nil {
		return record, err
	}
	return record, nil
}
