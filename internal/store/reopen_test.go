package store

import (
	"clubmembers/internal/fixtures"
	"context"
	"path/filepath"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "members.db")
	record := fixtures.ValidRecord("m-reopen", "重开成员", "13800138009")
	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.SaveMember(context.Background(), record); err != nil {
		db.Close()
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	got, err := db.GetMember(context.Background(), record.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Phone != record.Phone || got.Name != record.Name {
		t.Fatalf("reopened record mismatch %#v", got)
	}
}
