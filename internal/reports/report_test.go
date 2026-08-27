package reports

import (
	"bytes"
	"clubmembers/internal/fixtures"
	"strings"
	"testing"
)

func TestReportSummaryAndCSV(t *testing.T) {
	records := fixtures.SampleRecords()
	counts := GroupByFaculty(records)
	if len(counts) != 2 || counts[1].Faculty != "science" {
		t.Fatalf("unexpected counts %#v", counts)
	}
	if ActiveRatio(records) != 2.0/3.0 {
		t.Fatalf("unexpected active ratio %v", ActiveRatio(records))
	}
	var output bytes.Buffer
	if err := WriteCSV(&output, records); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "m-linxiao") || !strings.Contains(output.String(), "student_number") {
		t.Fatalf("unexpected csv %q", output.String())
	}
}
