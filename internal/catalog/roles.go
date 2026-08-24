package catalog

import (
	"fmt"
	"sort"
	"strings"
)

type RolePolicy struct {
	Role        string
	CanView     bool
	CanEdit     bool
	CanArchive  bool
	CanExport   bool
	Description string
}

func DefaultRolePolicies() []RolePolicy {
	return []RolePolicy{
		{Role: "admin", CanView: true, CanEdit: true, CanArchive: true, CanExport: true, Description: "管理全部档案和历史"},
		{Role: "editor", CanView: true, CanEdit: true, CanArchive: false, CanExport: true, Description: "编辑成员资料"},
		{Role: "viewer", CanView: true, CanEdit: false, CanArchive: false, CanExport: false, Description: "查看经过授权的档案"},
	}
}

func ResolvePolicy(role string) (RolePolicy, error) {
	role = strings.ToLower(strings.TrimSpace(role))
	for _, policy := range DefaultRolePolicies() {
		if policy.Role == role {
			return policy, nil
		}
	}
	return RolePolicy{}, fmt.Errorf("unknown role %q", role)
}

func CanArchive(role string) bool {
	policy, err := ResolvePolicy(role)
	return err == nil && policy.CanArchive
}

func CanExport(role string) bool {
	policy, err := ResolvePolicy(role)
	return err == nil && policy.CanExport
}

func RoleNames() []string {
	policies := DefaultRolePolicies()
	names := make([]string, 0, len(policies))
	for _, policy := range policies {
		names = append(names, policy.Role)
	}
	sort.Strings(names)
	return names
}

func FilterPolicies(canEdit bool) []RolePolicy {
	result := make([]RolePolicy, 0)
	for _, policy := range DefaultRolePolicies() {
		if !canEdit || policy.CanEdit {
			result = append(result, policy)
		}
	}
	return result
}
