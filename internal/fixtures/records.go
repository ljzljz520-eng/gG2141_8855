package fixtures

import "clubmembers/internal/members"

func ValidRecord(id, name, phone string) members.MemberRecord {
	return members.MemberRecord{ID: id, Name: name, StudentNumber: "S-" + id, Faculty: "science", Phone: phone, Email: id + "@example.org", JoinedOn: "2024-09-01", Status: members.StatusActive, Tags: []string{"摄影"}, Contacts: []members.ContactPoint{{Label: "家人", Name: "联系人", Phone: "13900000000", Relationship: "家人", Preferred: true}}, Notes: []string{"已完成入会登记"}}
}

func SampleRecords() []members.MemberRecord {
	return []members.MemberRecord{
		ValidRecord("m-linxiao", "林晓", "13800138001"),
		ValidRecord("m-chenyu", "陈宇", "13800138002"),
		{ID: "m-zhaonan", Name: "赵楠", StudentNumber: "S-m-zhaonan", Faculty: "engineering", Phone: "13800138003", Email: "m-zhaonan@example.org", JoinedOn: "2024-09-03", Status: members.StatusPaused, Tags: []string{"器材"}},
	}
}
