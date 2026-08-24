package audit

import (
	"context"
	"fmt"

	"clubmembers/internal/members"
	"clubmembers/internal/store"
)

type Ledger struct{ store *store.Store }

func NewLedger(s *store.Store) *Ledger { return &Ledger{store: s} }

func (l *Ledger) Record(ctx context.Context, events []members.ChangeEvent, occurredOn string) error {
	if l == nil || l.store == nil {
		return fmt.Errorf("audit ledger is unavailable")
	}
	return l.store.AppendEvents(ctx, events, occurredOn)
}

func (l *Ledger) History(ctx context.Context, memberID string) ([]members.ChangeEvent, error) {
	if l == nil || l.store == nil {
		return nil, fmt.Errorf("audit ledger is unavailable")
	}
	return l.store.ListEvents(ctx, memberID)
}

func Summarize(events []members.ChangeEvent) map[string]int {
	counts := make(map[string]int)
	for _, event := range events {
		counts[event.Field]++
	}
	return counts
}
