package store

import "context"

func (s *Store) createTables(ctx context.Context) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS members (id TEXT PRIMARY KEY, name TEXT NOT NULL, student_number TEXT NOT NULL UNIQUE, faculty TEXT NOT NULL, phone TEXT NOT NULL, email TEXT NOT NULL, joined_on TEXT NOT NULL, status TEXT NOT NULL, revision INTEGER NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS member_tags (member_id TEXT NOT NULL REFERENCES members(id) ON DELETE CASCADE, tag TEXT NOT NULL, PRIMARY KEY(member_id, tag))`,
		`CREATE TABLE IF NOT EXISTS contact_points (id INTEGER PRIMARY KEY AUTOINCREMENT, member_id TEXT NOT NULL REFERENCES members(id) ON DELETE CASCADE, label TEXT NOT NULL, name TEXT NOT NULL, phone TEXT NOT NULL, relationship TEXT NOT NULL, preferred INTEGER NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS member_notes (id INTEGER PRIMARY KEY AUTOINCREMENT, member_id TEXT NOT NULL REFERENCES members(id) ON DELETE CASCADE, body TEXT NOT NULL, position INTEGER NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS change_events (id TEXT PRIMARY KEY, member_id TEXT NOT NULL REFERENCES members(id) ON DELETE CASCADE, actor TEXT NOT NULL, field TEXT NOT NULL, before_value TEXT NOT NULL, after_value TEXT NOT NULL, occurred_on TEXT NOT NULL, reason TEXT NOT NULL)`,
		`CREATE INDEX IF NOT EXISTS idx_members_faculty ON members(faculty)`,
		`CREATE INDEX IF NOT EXISTS idx_events_member ON change_events(member_id, occurred_on)`,
	}
	for _, statement := range statements {
		if _, err := s.db.ExecContext(ctx, statement); err != nil {
			return err
		}
	}
	return nil
}
