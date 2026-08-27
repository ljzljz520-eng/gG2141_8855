package catalog

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type Term struct {
	Code      string
	Label     string
	StartDate string
	EndDate   string
	Open      bool
}

type TermCalendar struct {
	terms map[string]Term
}

func NewTermCalendar(terms []Term) (*TermCalendar, error) {
	calendar := &TermCalendar{terms: make(map[string]Term, len(terms))}
	for _, term := range terms {
		if err := calendar.Add(term); err != nil {
			return nil, err
		}
	}
	return calendar, nil
}

func (c *TermCalendar) Add(term Term) error {
	if c == nil {
		return errors.New("term calendar is nil")
	}
	term.Code = strings.TrimSpace(term.Code)
	term.Label = strings.TrimSpace(term.Label)
	term.StartDate = strings.TrimSpace(term.StartDate)
	term.EndDate = strings.TrimSpace(term.EndDate)
	if term.Code == "" || term.Label == "" {
		return errors.New("term code and label are required")
	}
	if term.StartDate == "" || term.EndDate == "" || term.StartDate > term.EndDate {
		return fmt.Errorf("term %q has invalid dates", term.Code)
	}
	if _, exists := c.terms[term.Code]; exists {
		return fmt.Errorf("term %q already exists", term.Code)
	}
	c.terms[term.Code] = term
	return nil
}

func (c *TermCalendar) Get(code string) (Term, bool) {
	if c == nil {
		return Term{}, false
	}
	term, ok := c.terms[strings.TrimSpace(code)]
	return term, ok
}

func (c *TermCalendar) SetOpen(code string, open bool) error {
	term, ok := c.Get(code)
	if !ok {
		return fmt.Errorf("term %q not found", code)
	}
	term.Open = open
	c.terms[term.Code] = term
	return nil
}

func (c *TermCalendar) IsOpen(code string) bool {
	term, ok := c.Get(code)
	return ok && term.Open
}

func (c *TermCalendar) ActiveTerms() []Term {
	result := make([]Term, 0)
	if c == nil {
		return result
	}
	for _, term := range c.terms {
		if term.Open {
			result = append(result, term)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].StartDate < result[j].StartDate })
	return result
}

func (c *TermCalendar) CloseAll() int {
	if c == nil {
		return 0
	}
	closed := 0
	for code, term := range c.terms {
		if term.Open {
			term.Open = false
			c.terms[code] = term
			closed++
		}
	}
	return closed
}

func (c *TermCalendar) Labels() []string {
	labels := make([]string, 0)
	if c != nil {
		for _, term := range c.terms {
			labels = append(labels, term.Label)
		}
	}
	sort.Strings(labels)
	return labels
}

func DefaultCalendar() *TermCalendar {
	calendar, err := NewTermCalendar([]Term{{Code: "2024-fall", Label: "2024 秋季", StartDate: "2024-09-01", EndDate: "2025-01-31", Open: true}, {Code: "2025-spring", Label: "2025 春季", StartDate: "2025-02-01", EndDate: "2025-07-31", Open: false}})
	if err != nil {
		return &TermCalendar{terms: make(map[string]Term)}
	}
	return calendar
}
