package workflows

import (
	"clubmembers/internal/catalog"
	"clubmembers/internal/fixtures"
	"clubmembers/internal/store"
	"context"
	"testing"
)

func TestWorkflowPhoneChange(t *testing.T) {
	db, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	service := NewService(db, catalog.DefaultDirectory())
	if _, errs := service.EnrollBatch(context.Background(), fixtures.SampleRecords()); len(errs) != 0 {
		t.Fatal(errs[0])
	}
	updated, err := service.UpdatePhone(context.Background(), "admin", "m-linxiao", "139 1111 2222", "迎新通知更新", "2024-10-01")
	if err != nil {
		t.Fatal(err)
	}
	if updated.Phone != "13911112222" {
		t.Fatalf("unexpected phone %q", updated.Phone)
	}
	peer, err := service.Get(context.Background(), "m-chenyu")
	if err != nil {
		t.Fatal(err)
	}
	if peer.Phone != "13800138002" {
		t.Fatalf("peer changed unexpectedly to %q", peer.Phone)
	}
	history, err := service.History(context.Background(), "m-linxiao")
	if err != nil {
		t.Fatal(err)
	}
	if len(history) != 1 || history[0].Field != "phone" {
		t.Fatalf("unexpected history %#v", history)
	}
}
