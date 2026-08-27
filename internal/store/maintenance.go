package store

import (
	"context"
	"fmt"
	"strings"

	"clubmembers/internal/members"
)

func (s *Store) AddTag(ctx context.Context, memberID, tag string) error {
	tag = strings.ToLower(strings.TrimSpace(tag))
	if tag == "" {
		return fmt.Errorf("tag is empty")
	}
	if exists, err := s.MemberExists(ctx, memberID); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf("member %q not found", memberID)
	}
	_, err := s.db.ExecContext(ctx, `INSERT OR IGNORE INTO member_tags(member_id,tag) VALUES(?,?)`, memberID, tag)
	return err
}

func (s *Store) RemoveTag(ctx context.Context, memberID, tag string) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM member_tags WHERE member_id=? AND tag=?`, memberID, strings.ToLower(strings.TrimSpace(tag)))
	if err != nil {
		return err
	}
	if count, err := result.RowsAffected(); err != nil {
		return err
	} else if count == 0 {
		return fmt.Errorf("tag not found")
	}
	return nil
}

func (s *Store) ReplaceContacts(ctx context.Context, memberID string, contacts []members.ContactPoint) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `DELETE FROM contact_points WHERE member_id=?`, memberID); err != nil {
		return err
	}
	for _, contact := range contacts {
		if contact.Name == "" || contact.Phone == "" {
			return fmt.Errorf("contact is incomplete")
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO contact_points(member_id,label,name,phone,relationship,preferred) VALUES(?,?,?,?,?,?)`, memberID, contact.Label, contact.Name, contact.Phone, contact.Relationship, boolInt(contact.Preferred)); err != nil {
			return err
		}
	}
	return tx.Commit()
}
