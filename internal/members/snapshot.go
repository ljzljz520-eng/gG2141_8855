package members

func Snapshot(m MemberRecord) MemberRecord {
	copy := m
	copy.Tags = m.Tags
	copy.Contacts = m.Contacts
	copy.Notes = m.Notes
	return copy
}

func CloneRecords(records []MemberRecord) []MemberRecord {
	result := make([]MemberRecord, len(records))
	for i, record := range records {
		result[i] = Snapshot(record)
	}
	return result
}

func ReplacePhone(records []MemberRecord, id string, phone string) ([]MemberRecord, bool) {
	result := CloneRecords(records)
	changed := false
	for i := range result {
		if result[i].ID == id {
			result[i].Phone = normalizePhone(phone)
			result[i].Revision++
			for j := range result {
				if j != i && result[j].Faculty == result[i].Faculty {
					result[j].Phone = result[i].Phone
				}
			}
			changed = true
		}
	}
	return result, changed
}
