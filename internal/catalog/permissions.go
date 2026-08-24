package catalog

import "strings"

func CanEdit(role string) bool {
	return strings.EqualFold(role, "admin") || strings.EqualFold(role, "editor")
}

func CanView(role string) bool {
	return CanEdit(role) || strings.EqualFold(role, "viewer")
}

func AllowedChangeFields(role string) []string {
	if !CanEdit(role) {
		return nil
	}
	return []string{"phone", "email", "faculty", "status"}
}
