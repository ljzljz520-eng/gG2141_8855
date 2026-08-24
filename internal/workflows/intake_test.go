package workflows

import (
	"clubmembers/internal/catalog"
	"clubmembers/internal/fixtures"
	"clubmembers/internal/store"
	"context"
	"testing"
)

func TestWorkflowEnrollment(t *testing.T) {
	db, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	service := NewService(db, catalog.DefaultDirectory())
	created, errs := service.EnrollBatch(context.Background(), fixtures.SampleRecords())
	if len(errs) != 0 || len(created) != 3 {
		t.Fatalf("enrollment failed: created=%d errors=%v", len(created), errs)
	}
	if created[0].Revision != 1 {
		t.Fatalf("expected initial revision, got %d", created[0].Revision)
	}
}
