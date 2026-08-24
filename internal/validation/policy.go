package validation

import (
	"fmt"
	"strings"

	"clubmembers/internal/catalog"
	"clubmembers/internal/members"
)

func ValidateChange(actorRole string, request members.ChangeRequest, directory *catalog.Directory) error {
	if !catalog.CanEdit(actorRole) {
		return fmt.Errorf("role %q cannot edit member records", actorRole)
	}
	if directory == nil {
		return fmt.Errorf("faculty directory is unavailable")
	}
	if request.Faculty != nil && !directory.ValidateFaculty(strings.TrimSpace(*request.Faculty)) {
		return fmt.Errorf("faculty %q is not active", *request.Faculty)
	}
	return RequireReason(request.Reason)
}

func NormalizeAndValidate(record members.MemberRecord) (members.MemberRecord, error) {
	record = members.NormalizeMember(record)
	if issues := ValidateRecord(record); len(issues) > 0 {
		return record, fmt.Errorf("record invalid: %s", strings.Join(issues, "; "))
	}
	return record, nil
}
