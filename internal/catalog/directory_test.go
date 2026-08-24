package catalog

import "testing"

func TestDirectoryListsActiveFaculties(t *testing.T) {
	directory := DefaultDirectory()
	codes := directory.FacultyCodes()
	if len(codes) != 4 || codes[0] != "arts" || codes[3] != "science" {
		t.Fatalf("unexpected faculty codes %#v", codes)
	}
	if !directory.ValidateFaculty("science") || directory.ValidateFaculty("unknown") {
		t.Fatal("faculty validation mismatch")
	}
}

func TestPermissions(t *testing.T) {
	if !CanEdit("admin") || CanEdit("viewer") || !CanView("viewer") {
		t.Fatal("permission policy mismatch")
	}
	if len(AllowedChangeFields("viewer")) != 0 {
		t.Fatal("viewer received edit fields")
	}
}
