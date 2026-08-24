package store

import (
	"clubmembers/internal/fixtures"
	"context"
	"testing"
)

func TestStoreSavesRelatedProfileData(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	record := fixtures.ValidRecord("m1", "林晓", "13800138001")
	if err := db.SaveMember(context.Background(), record); err != nil {
		t.Fatal(err)
	}
	got, err := db.GetMember(context.Background(), record.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != record.Name || len(got.Tags) != 1 || len(got.Contacts) != 1 || len(got.Notes) != 1 {
		t.Fatalf("related data did not round trip: %#v", got)
	}
}
