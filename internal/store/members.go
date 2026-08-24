package store

import (
	"context"
	"database/sql"
	"fmt"

	"clubmembers/internal/members"
)

func (s *Store) SaveMember(ctx context.Context, record members.MemberRecord) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `INSERT INTO members(id,name,student_number,faculty,phone,email,joined_on,status,revision) VALUES(?,?,?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET name=excluded.name,student_number=excluded.student_number,faculty=excluded.faculty,phone=excluded.phone,email=excluded.email,joined_on=excluded.joined_on,status=excluded.status,revision=excluded.revision`, record.ID, record.Name, record.StudentNumber, record.Faculty, record.Phone, record.Email, record.JoinedOn, record.Status, record.Revision)
	if err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM member_tags WHERE member_id=?`, record.ID); err != nil {
		return err
	}
	for _, tag := range record.Tags {
		if _, err = tx.ExecContext(ctx, `INSERT INTO member_tags(member_id,tag) VALUES(?,?)`, record.ID, tag); err != nil {
			return err
		}
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM contact_points WHERE member_id=?`, record.ID); err != nil {
		return err
	}
	for _, contact := range record.Contacts {
		if _, err = tx.ExecContext(ctx, `INSERT INTO contact_points(member_id,label,name,phone,relationship,preferred) VALUES(?,?,?,?,?,?)`, record.ID, contact.Label, contact.Name, contact.Phone, contact.Relationship, boolInt(contact.Preferred)); err != nil {
			return err
		}
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM member_notes WHERE member_id=?`, record.ID); err != nil {
		return err
	}
	for i, note := range record.Notes {
		if _, err = tx.ExecContext(ctx, `INSERT INTO member_notes(member_id,body,position) VALUES(?,?,?)`, record.ID, note, i); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func (s *Store) GetMember(ctx context.Context, id string) (members.MemberRecord, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id,name,student_number,faculty,phone,email,joined_on,status,revision FROM members WHERE id=?`, id)
	var record members.MemberRecord
	var status string
	if err := row.Scan(&record.ID, &record.Name, &record.StudentNumber, &record.Faculty, &record.Phone, &record.Email, &record.JoinedOn, &status, &record.Revision); err != nil {
		if err == sql.ErrNoRows {
			return record, fmt.Errorf("member %q not found", id)
		}
		return record, err
	}
	record.Status = members.Status(status)
	if err := s.loadDetails(ctx, &record); err != nil {
		return record, err
	}
	return record, nil
}

func (s *Store) ListMembers(ctx context.Context, faculty string, status members.Status) ([]members.MemberRecord, error) {
	query := `SELECT id,name,student_number,faculty,phone,email,joined_on,status,revision FROM members`
	args := []any{}
	where := ""
	if faculty != "" {
		where = " WHERE faculty=?"
		args = append(args, faculty)
	}
	if status != "" {
		if where == "" {
			where = " WHERE status=?"
		} else {
			where += " AND status=?"
		}
		args = append(args, status)
	}
	query += where + " ORDER BY name, id"
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]members.MemberRecord, 0)
	for rows.Next() {
		var record members.MemberRecord
		var state string
		if err := rows.Scan(&record.ID, &record.Name, &record.StudentNumber, &record.Faculty, &record.Phone, &record.Email, &record.JoinedOn, &state, &record.Revision); err != nil {
			return nil, err
		}
		record.Status = members.Status(state)
		result = append(result, record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	for i := range result {
		if err := s.loadDetails(ctx, &result[i]); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (s *Store) loadDetails(ctx context.Context, record *members.MemberRecord) error {
	tagRows, err := s.db.QueryContext(ctx, `SELECT tag FROM member_tags WHERE member_id=? ORDER BY tag`, record.ID)
	if err != nil {
		return err
	}
	for tagRows.Next() {
		var tag string
		if err := tagRows.Scan(&tag); err != nil {
			tagRows.Close()
			return err
		}
		record.Tags = append(record.Tags, tag)
	}
	if err := tagRows.Close(); err != nil {
		return err
	}
	contactRows, err := s.db.QueryContext(ctx, `SELECT label,name,phone,relationship,preferred FROM contact_points WHERE member_id=? ORDER BY id`, record.ID)
	if err != nil {
		return err
	}
	for contactRows.Next() {
		var contact members.ContactPoint
		var preferred int
		if err := contactRows.Scan(&contact.Label, &contact.Name, &contact.Phone, &contact.Relationship, &preferred); err != nil {
			contactRows.Close()
			return err
		}
		contact.Preferred = preferred == 1
		record.Contacts = append(record.Contacts, contact)
	}
	if err := contactRows.Close(); err != nil {
		return err
	}
	noteRows, err := s.db.QueryContext(ctx, `SELECT body FROM member_notes WHERE member_id=? ORDER BY position`, record.ID)
	if err != nil {
		return err
	}
	for noteRows.Next() {
		var note string
		if err := noteRows.Scan(&note); err != nil {
			noteRows.Close()
			return err
		}
		record.Notes = append(record.Notes, note)
	}
	return noteRows.Close()
}
