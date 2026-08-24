package store

import (
	"clubmembers/internal/members"
	"context"
)

func (s *Store) AppendEvents(ctx context.Context, events []members.ChangeEvent, occurredOn string) error {
	if len(events) == 0 {
		return nil
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for index, event := range events {
		if event.ID == "" {
			event.ID = event.MemberID + "-" + event.Field + "-" + occurredOn + "-" + string(rune('a'+index))
		}
		if event.OccurredOn == "" {
			event.OccurredOn = occurredOn
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO change_events(id,member_id,actor,field,before_value,after_value,occurred_on,reason) VALUES(?,?,?,?,?,?,?,?)`, event.ID, event.MemberID, event.Actor, event.Field, event.Before, event.After, event.OccurredOn, event.Reason); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) ListEvents(ctx context.Context, memberID string) ([]members.ChangeEvent, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id,member_id,actor,field,before_value,after_value,occurred_on,reason FROM change_events WHERE member_id=? ORDER BY occurred_on,id`, memberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]members.ChangeEvent, 0)
	for rows.Next() {
		var event members.ChangeEvent
		if err := rows.Scan(&event.ID, &event.MemberID, &event.Actor, &event.Field, &event.Before, &event.After, &event.OccurredOn, &event.Reason); err != nil {
			return nil, err
		}
		result = append(result, event)
	}
	return result, rows.Err()
}
