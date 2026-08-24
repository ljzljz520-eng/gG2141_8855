package members

import "testing"

func TestMemberPhoneChangeIsIsolated(t *testing.T) {
	records := []MemberRecord{{ID: "lin", Name: "林晓", Faculty: "science", Phone: "13800138001"}, {ID: "chen", Name: "陈宇", Faculty: "science", Phone: "13800138002"}}
	updated, changed := ReplacePhone(records, "lin", "13900139001")
	if !changed {
		t.Fatal("expected a member change")
	}
	peers := FilterByFaculty(updated, "science")
	for _, peer := range peers {
		if peer.ID == "chen" && peer.Phone != "13800138002" {
			t.Fatalf("peer phone changed to %q", peer.Phone)
		}
	}
}

func TestNormalizeMemberRemovesDuplicateTags(t *testing.T) {
	record := NormalizeMember(MemberRecord{Tags: []string{"摄影", " 摄影 ", ""}})
	if len(record.Tags) != 1 || record.Tags[0] != "摄影" {
		t.Fatalf("unexpected tags %#v", record.Tags)
	}
}
