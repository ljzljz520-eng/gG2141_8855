package members

import "testing"

func TestMemberValidationAndDisplay(t *testing.T) {
	record := MemberRecord{ID: "m1", Name: "林晓", StudentNumber: "S1", Faculty: "science", Phone: "13800138001", Status: StatusActive}
	if err := record.Validate(); err != nil {
		t.Fatalf("valid record rejected: %v", err)
	}
	if got := record.DisplayLabel(); got != "林晓 (S1)" {
		t.Fatalf("unexpected label %q", got)
	}
	if record.HasTag("摄影") {
		t.Fatal("unexpected tag")
	}
}

func TestContactSummaryPrefersMarkedContact(t *testing.T) {
	record := MemberRecord{Contacts: []ContactPoint{{Name: "A", Phone: "1"}, {Name: "B", Phone: "2", Preferred: true}}}
	if got := record.ContactSummary(); got != "B: 2" {
		t.Fatalf("unexpected contact summary %q", got)
	}
}
