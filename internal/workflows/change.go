package workflows

import (
	"clubmembers/internal/members"
	"context"
	"fmt"
)

func (s *Service) UpdatePhone(ctx context.Context, role, memberID, phone, reason, occurredOn string) (members.MemberRecord, error) {
	if phone == "" {
		return members.MemberRecord{}, fmt.Errorf("phone is required")
	}
	return s.Change(ctx, role, members.ChangeRequest{MemberID: memberID, Actor: role, Phone: &phone, Reason: reason}, occurredOn)
}

func (s *Service) UpdateEmail(ctx context.Context, role, memberID, email, reason, occurredOn string) (members.MemberRecord, error) {
	if email == "" {
		return members.MemberRecord{}, fmt.Errorf("email is required")
	}
	return s.Change(ctx, role, members.ChangeRequest{MemberID: memberID, Actor: role, Email: &email, Reason: reason}, occurredOn)
}

func (s *Service) TransferFaculty(ctx context.Context, role, memberID, faculty, reason, occurredOn string) (members.MemberRecord, error) {
	return s.Change(ctx, role, members.ChangeRequest{MemberID: memberID, Actor: role, Faculty: &faculty, Reason: reason}, occurredOn)
}
