package members

import (
	"errors"
	"sort"
	"strings"
)

type Roster struct {
	records  []MemberRecord
	position map[string]int
}

func NewRoster(records []MemberRecord) (*Roster, error) {
	roster := &Roster{records: make([]MemberRecord, 0, len(records)), position: make(map[string]int)}
	for _, record := range records {
		if err := roster.Add(record); err != nil {
			return nil, err
		}
	}
	return roster, nil
}

func (r *Roster) Add(record MemberRecord) error {
	if r == nil {
		return errors.New("roster is nil")
	}
	record = NormalizeMember(record)
	if err := record.Validate(); err != nil {
		return err
	}
	if _, exists := r.position[record.ID]; exists {
		return errors.New("member already exists")
	}
	r.position[record.ID] = len(r.records)
	r.records = append(r.records, Snapshot(record))
	return nil
}

func (r *Roster) Replace(record MemberRecord) error {
	if r == nil {
		return errors.New("roster is nil")
	}
	index, exists := r.position[record.ID]
	if !exists {
		return errors.New("member does not exist")
	}
	record = NormalizeMember(record)
	if err := record.Validate(); err != nil {
		return err
	}
	r.records[index] = Snapshot(record)
	return nil
}

func (r *Roster) Find(id string) (MemberRecord, bool) {
	if r == nil {
		return MemberRecord{}, false
	}
	index, ok := r.position[id]
	if !ok || index < 0 || index >= len(r.records) {
		return MemberRecord{}, false
	}
	return Snapshot(r.records[index]), true
}

func (r *Roster) All() []MemberRecord {
	if r == nil {
		return nil
	}
	return SortByName(r.records)
}

func (r *Roster) ByFaculty(faculty string) []MemberRecord {
	if r == nil {
		return nil
	}
	return FilterByFaculty(r.records, strings.TrimSpace(faculty))
}

func (r *Roster) ByStatus(status Status) []MemberRecord {
	if r == nil {
		return nil
	}
	return FilterByStatus(r.records, status)
}

func (r *Roster) Search(query string) []MemberRecord {
	if r == nil {
		return nil
	}
	query = strings.ToLower(strings.TrimSpace(query))
	result := make([]MemberRecord, 0)
	for _, record := range r.records {
		if query == "" || strings.Contains(strings.ToLower(record.Name), query) || strings.Contains(strings.ToLower(record.StudentNumber), query) || strings.Contains(strings.ToLower(record.Phone), query) {
			result = append(result, Snapshot(record))
		}
	}
	return SortByName(result)
}

func (r *Roster) Tagged(tag string) []MemberRecord {
	if r == nil {
		return nil
	}
	result := make([]MemberRecord, 0)
	for _, record := range r.records {
		if record.HasTag(tag) {
			result = append(result, Snapshot(record))
		}
	}
	return SortByName(result)
}

func (r *Roster) UpdateNotes(id string, notes []string) error {
	record, ok := r.Find(id)
	if !ok {
		return errors.New("member does not exist")
	}
	record.Notes = normalizeNotes(notes)
	record.Revision++
	return r.Replace(record)
}

func (r *Roster) AddTag(id, tag string) error {
	record, ok := r.Find(id)
	if !ok {
		return errors.New("member does not exist")
	}
	if strings.TrimSpace(tag) == "" {
		return errors.New("tag is empty")
	}
	record.Tags = append(record.Tags, tag)
	record.Tags = normalizeTags(record.Tags)
	record.Revision++
	return r.Replace(record)
}

func (r *Roster) RemoveTag(id, tag string) error {
	record, ok := r.Find(id)
	if !ok {
		return errors.New("member does not exist")
	}
	tag = strings.ToLower(strings.TrimSpace(tag))
	filtered := make([]string, 0, len(record.Tags))
	for _, candidate := range record.Tags {
		if strings.ToLower(candidate) != tag {
			filtered = append(filtered, candidate)
		}
	}
	record.Tags = filtered
	record.Revision++
	return r.Replace(record)
}

func (r *Roster) Validate() []string {
	if r == nil {
		return []string{"roster is nil"}
	}
	issues := make([]string, 0)
	seenStudents := make(map[string]bool)
	for _, record := range r.records {
		if err := record.Validate(); err != nil {
			issues = append(issues, record.ID+": "+err.Error())
		}
		if seenStudents[record.StudentNumber] {
			issues = append(issues, record.ID+": duplicate student number")
		}
		seenStudents[record.StudentNumber] = true
	}
	return issues
}

func (r *Roster) FacultyCounts() map[string]int {
	counts := make(map[string]int)
	if r == nil {
		return counts
	}
	for _, record := range r.records {
		counts[record.Faculty]++
	}
	return counts
}

func (r *Roster) SortedIDs() []string {
	ids := make([]string, 0)
	if r != nil {
		for _, record := range r.records {
			ids = append(ids, record.ID)
		}
	}
	sort.Strings(ids)
	return ids
}
