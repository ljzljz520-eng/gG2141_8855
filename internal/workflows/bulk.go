package workflows

import (
	"bytes"
	"context"
	"fmt"
	"sort"

	"clubmembers/internal/members"
	"clubmembers/internal/reports"
)

type BatchResult struct {
	Updated []members.MemberRecord
	Failed  map[string]error
}

type Page struct {
	Items      []members.MemberRecord
	Page       int
	PageSize   int
	Total      int
	TotalPages int
}

func (s *Service) ChangeStatuses(ctx context.Context, role string, ids []string, status members.Status, reason, occurredOn string) BatchResult {
	result := BatchResult{Updated: make([]members.MemberRecord, 0, len(ids)), Failed: make(map[string]error)}
	for _, id := range ids {
		changed, err := s.Change(ctx, role, members.ChangeRequest{MemberID: id, Actor: role, Status: &status, Reason: reason}, occurredOn)
		if err != nil {
			result.Failed[id] = err
			continue
		}
		result.Updated = append(result.Updated, changed)
	}
	return result
}

func (s *Service) SearchPage(ctx context.Context, filter MemberFilter, page, pageSize int) (Page, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	rows, err := s.List(ctx, filter)
	if err != nil {
		return Page{}, err
	}
	rows = members.SortByName(rows)
	start := (page - 1) * pageSize
	if start > len(rows) {
		start = len(rows)
	}
	end := start + pageSize
	if end > len(rows) {
		end = len(rows)
	}
	totalPages := 0
	if len(rows) > 0 {
		totalPages = (len(rows) + pageSize - 1) / pageSize
	}
	return Page{Items: rows[start:end], Page: page, PageSize: pageSize, Total: len(rows), TotalPages: totalPages}, nil
}

func (s *Service) Export(ctx context.Context, filter MemberFilter) (string, error) {
	rows, err := s.List(ctx, filter)
	if err != nil {
		return "", err
	}
	var output bytes.Buffer
	if err := reports.WriteCSV(&output, rows); err != nil {
		return "", err
	}
	return output.String(), nil
}

func (s *Service) CompareFaculties(ctx context.Context, first, second string) (int, error) {
	firstRows, err := s.List(ctx, MemberFilter{Faculty: first})
	if err != nil {
		return 0, err
	}
	secondRows, err := s.List(ctx, MemberFilter{Faculty: second})
	if err != nil {
		return 0, err
	}
	return len(firstRows) - len(secondRows), nil
}

func (s *Service) ValidateRoster(ctx context.Context) ([]string, error) {
	rows, err := s.List(ctx, MemberFilter{})
	if err != nil {
		return nil, err
	}
	roster, err := members.NewRoster(rows)
	if err != nil {
		return nil, err
	}
	return roster.Validate(), nil
}

func SortByRevision(records []members.MemberRecord) []members.MemberRecord {
	result := members.CloneRecords(records)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Revision == result[j].Revision {
			return result[i].ID < result[j].ID
		}
		return result[i].Revision > result[j].Revision
	})
	return result
}

func RequirePage(page Page) error {
	if page.Page < 1 {
		return fmt.Errorf("page must be positive")
	}
	if page.PageSize < 1 {
		return fmt.Errorf("page size must be positive")
	}
	if page.Page > page.TotalPages && page.TotalPages > 0 {
		return fmt.Errorf("page exceeds result pages")
	}
	return nil
}
