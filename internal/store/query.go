package store

import (
	"context"
	"fmt"
	"sort"

	"clubmembers/internal/members"
)

type StatusCount struct {
	Status members.Status
	Count  int
}

type FacultySummary struct {
	Faculty string
	Total   int
	Active  int
	Paused  int
	Alumni  int
}

func (s *Store) CountByStatus(ctx context.Context) ([]StatusCount, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT status, COUNT(*) FROM members GROUP BY status`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]StatusCount, 0)
	for rows.Next() {
		var state string
		var count int
		if err := rows.Scan(&state, &count); err != nil {
			return nil, err
		}
		result = append(result, StatusCount{Status: members.Status(state), Count: count})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Status < result[j].Status })
	return result, nil
}

func (s *Store) SummarizeFaculties(ctx context.Context) ([]FacultySummary, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT faculty,status,COUNT(*) FROM members GROUP BY faculty,status ORDER BY faculty,status`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	byFaculty := make(map[string]FacultySummary)
	for rows.Next() {
		var faculty, state string
		var count int
		if err := rows.Scan(&faculty, &state, &count); err != nil {
			return nil, err
		}
		summary := byFaculty[faculty]
		summary.Faculty = faculty
		summary.Total += count
		switch members.Status(state) {
		case members.StatusActive:
			summary.Active += count
		case members.StatusPaused:
			summary.Paused += count
		case members.StatusAlumni:
			summary.Alumni += count
		}
		byFaculty[faculty] = summary
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	result := make([]FacultySummary, 0, len(byFaculty))
	for _, summary := range byFaculty {
		result = append(result, summary)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Faculty < result[j].Faculty })
	return result, nil
}

func (s *Store) DeleteMember(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM members WHERE id=?`, id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("member %q not found", id)
	}
	return nil
}

func (s *Store) MemberExists(ctx context.Context, id string) (bool, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM members WHERE id=?`, id).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *Store) FindByStudentNumber(ctx context.Context, studentNumber string) (members.MemberRecord, error) {
	var id string
	if err := s.db.QueryRowContext(ctx, `SELECT id FROM members WHERE student_number=?`, studentNumber).Scan(&id); err != nil {
		return members.MemberRecord{}, err
	}
	return s.GetMember(ctx, id)
}

func (s *Store) FindByEmail(ctx context.Context, email string) (members.MemberRecord, error) {
	var id string
	if err := s.db.QueryRowContext(ctx, `SELECT id FROM members WHERE email=?`, email).Scan(&id); err != nil {
		return members.MemberRecord{}, err
	}
	return s.GetMember(ctx, id)
}
