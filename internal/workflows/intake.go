package workflows

import (
	"clubmembers/internal/members"
	"context"
	"fmt"
)

func (s *Service) EnrollBatch(ctx context.Context, records []members.MemberRecord) ([]members.MemberRecord, []error) {
	accepted := make([]members.MemberRecord, 0, len(records))
	errorsFound := make([]error, 0)
	for _, record := range records {
		created, err := s.Register(ctx, record)
		if err != nil {
			errorsFound = append(errorsFound, fmt.Errorf("%s: %w", record.ID, err))
			continue
		}
		accepted = append(accepted, created)
	}
	return accepted, errorsFound
}

func (s *Service) SeedDemo(ctx context.Context) error {
	_, errorsFound := s.EnrollBatch(ctx, []members.MemberRecord{
		{ID: "m-linxiao", Name: "林晓", StudentNumber: "S202401", Faculty: "science", Phone: "13800138001", Email: "lin.xiao@example.org", JoinedOn: "2024-09-01", Status: members.StatusActive, Tags: []string{"摄影", "迎新"}, Contacts: []members.ContactPoint{{Label: "家人", Name: "林妈妈", Phone: "13900139001", Relationship: "母亲", Preferred: true}}, Notes: []string{"负责秋季招新签到"}},
		{ID: "m-chenyu", Name: "陈宇", StudentNumber: "S202402", Faculty: "science", Phone: "13800138002", Email: "chen.yu@example.org", JoinedOn: "2024-09-02", Status: members.StatusActive, Tags: []string{"摄影"}},
		{ID: "m-zhaonan", Name: "赵楠", StudentNumber: "S202403", Faculty: "engineering", Phone: "13800138003", Email: "zhao.nan@example.org", JoinedOn: "2024-09-03", Status: members.StatusPaused, Tags: []string{"器材"}},
	})
	if len(errorsFound) > 0 {
		return errorsFound[0]
	}
	return nil
}
