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

// ReplacePhone updates the phone number of a single member within a roster
// snapshot. Only the targeted member is touched; other members of the same
// faculty keep their own phone numbers. (Previously this spread the new phone
// to every peer in the same faculty, corrupting unrelated records.)
func ReplacePhone(records []MemberRecord, id string, phone string) ([]MemberRecord, bool) {
	result := CloneRecords(records)
	changed := false
	for i := range result {
		if result[i].ID == id {
			result[i].Phone = normalizePhone(phone)
			result[i].Revision++
			changed = true
			break
		}
	}
	return result, changed
}
