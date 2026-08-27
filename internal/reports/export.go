package reports

import (
	"clubmembers/internal/members"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

func WriteCSV(w io.Writer, records []members.MemberRecord) error {
	writer := csv.NewWriter(w)
	if err := writer.Write([]string{"id", "name", "student_number", "faculty", "phone", "email", "status", "revision"}); err != nil {
		return err
	}
	for _, record := range records {
		if err := writer.Write([]string{record.ID, record.Name, record.StudentNumber, record.Faculty, record.Phone, record.Email, string(record.Status), strconv.Itoa(record.Revision)}); err != nil {
			return err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}
	return nil
}

func ProfileCard(record members.MemberRecord) string {
	return fmt.Sprintf("%s\n学院: %s\n电话: %s\n邮箱: %s\n状态: %s\n%s", record.DisplayLabel(), record.Faculty, record.Phone, record.Email, record.Status, record.ContactSummary())
}
