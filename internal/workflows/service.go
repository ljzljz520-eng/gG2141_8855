package workflows

import (
	"context"
	"fmt"

	"clubmembers/internal/audit"
	"clubmembers/internal/catalog"
	"clubmembers/internal/members"
	"clubmembers/internal/store"
	"clubmembers/internal/validation"
)

type Service struct {
	store     *store.Store
	directory *catalog.Directory
	ledger    *audit.Ledger
}

func NewService(s *store.Store, directory *catalog.Directory) *Service {
	return &Service{store: s, directory: directory, ledger: audit.NewLedger(s)}
}

func (s *Service) Register(ctx context.Context, record members.MemberRecord) (members.MemberRecord, error) {
	normalized, err := validation.NormalizeAndValidate(record)
	if err != nil {
		return normalized, err
	}
	if !s.directory.ValidateFaculty(normalized.Faculty) {
		return normalized, fmt.Errorf("faculty %q is not active", normalized.Faculty)
	}
	if normalized.Revision == 0 {
		normalized.Revision = 1
	}
	if err := s.store.SaveMember(ctx, normalized); err != nil {
		return normalized, err
	}
	return normalized, nil
}

func (s *Service) Get(ctx context.Context, id string) (members.MemberRecord, error) {
	return s.store.GetMember(ctx, id)
}

func (s *Service) List(ctx context.Context, filter MemberFilter) ([]members.MemberRecord, error) {
	rows, err := s.store.ListMembers(ctx, filter.Faculty, filter.Status)
	if err != nil {
		return nil, err
	}
	if filter.Tag != "" {
		filtered := make([]members.MemberRecord, 0, len(rows))
		for _, row := range rows {
			if row.HasTag(filter.Tag) {
				filtered = append(filtered, row)
			}
		}
		rows = filtered
	}
	if filter.Search != "" {
		filtered := make([]members.MemberRecord, 0, len(rows))
		for _, row := range rows {
			if filter.Matches(row) {
				filtered = append(filtered, row)
			}
		}
		rows = filtered
	}
	return rows, nil
}

func (s *Service) Change(ctx context.Context, role string, request members.ChangeRequest, occurredOn string) (members.MemberRecord, error) {
	if err := validation.ValidateChange(role, request, s.directory); err != nil {
		return members.MemberRecord{}, err
	}
	current, err := s.store.GetMember(ctx, request.MemberID)
	if err != nil {
		return current, err
	}
	updated, events, err := members.ApplyChange(current, request)
	if err != nil {
		return current, err
	}
	if err := s.store.SaveMember(ctx, updated); err != nil {
		return current, err
	}
	if err := s.ledger.Record(ctx, events, occurredOn); err != nil {
		return current, err
	}
	return updated, nil
}

func (s *Service) History(ctx context.Context, memberID string) ([]members.ChangeEvent, error) {
	return s.ledger.History(ctx, memberID)
}
