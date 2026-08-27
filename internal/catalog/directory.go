package catalog

import "sort"

type Faculty struct {
	Code        string
	Name        string
	Coordinator string
	Active      bool
}

type Directory struct {
	faculties map[string]Faculty
	roles     map[string]string
}

func DefaultDirectory() *Directory {
	return &Directory{faculties: map[string]Faculty{
		"arts":        {Code: "arts", Name: "人文学院", Coordinator: "顾老师", Active: true},
		"science":     {Code: "science", Name: "理学院", Coordinator: "何老师", Active: true},
		"engineering": {Code: "engineering", Name: "工学院", Coordinator: "陆老师", Active: true},
		"business":    {Code: "business", Name: "商学院", Coordinator: "章老师", Active: true},
	}, roles: map[string]string{"admin": "社团管理员", "editor": "档案编辑员", "viewer": "只读查看员"}}
}

func (d *Directory) Faculty(code string) (Faculty, bool) {
	faculty, ok := d.faculties[code]
	return faculty, ok && faculty.Active
}

func (d *Directory) FacultyCodes() []string {
	codes := make([]string, 0, len(d.faculties))
	for code, faculty := range d.faculties {
		if faculty.Active {
			codes = append(codes, code)
		}
	}
	sort.Strings(codes)
	return codes
}

func (d *Directory) RoleName(role string) (string, bool) {
	name, ok := d.roles[role]
	return name, ok
}

func (d *Directory) ValidateFaculty(code string) bool {
	_, ok := d.Faculty(code)
	return ok
}
