package workflows

import (
	"clubmembers/internal/members"
	"clubmembers/internal/validation"
	"context"
	"fmt"
)

type Review struct {
	Member members.MemberRecord
	Ready  bool
	Issues []string
}

func (s *Service) Review(ctx context.Context, memberID string) (Review, error) {
	record, err := s.Get(ctx, memberID)
	if err != nil {
		return Review{}, err
	}
	issues := validation.ValidateRecord(record)
	return Review{Member: record, Ready: validation.IsReadyForActivation(record), Issues: issues}, nil
}

func (s *Service) Activate(ctx context.Context, role, memberID, reason, occurredOn string) (members.MemberRecord, error) {
	review, err := s.Review(ctx, memberID)
	if err != nil {
		return members.MemberRecord{}, err
	}
	if !review.Ready {
		return review.Member, fmt.Errorf("member is not ready: %v", review.Issues)
	}
	status := members.StatusActive
	return s.Change(ctx, role, members.ChangeRequest{MemberID: memberID, Actor: role, Status: &status, Reason: reason}, occurredOn)
}
