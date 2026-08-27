package workflows

import (
	"clubmembers/internal/catalog"
	"clubmembers/internal/fixtures"
	"clubmembers/internal/store"
	"context"
	"testing"
)

func TestWorkflowSearch(t *testing.T) {
	db, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	service := NewService(db, catalog.DefaultDirectory())
	if _, errs := service.EnrollBatch(context.Background(), fixtures.SampleRecords()); len(errs) != 0 {
		t.Fatal(errs[0])
	}
	rows, err := service.List(context.Background(), MemberFilter{Faculty: "science", Tag: "摄影", Search: "林"})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Name != "林晓" {
		t.Fatalf("unexpected query result %#v", rows)
	}
}

func TestWorkflowReview(t *testing.T) {
	db, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	service := NewService(db, catalog.DefaultDirectory())
	if _, errs := service.EnrollBatch(context.Background(), fixtures.SampleRecords()); len(errs) != 0 {
		t.Fatal(errs[0])
	}
	review, err := service.Review(context.Background(), "m-linxiao")
	if err != nil {
		t.Fatal(err)
	}
	if !review.Ready || len(review.Issues) != 0 {
		t.Fatalf("unexpected review %#v", review)
	}
}
