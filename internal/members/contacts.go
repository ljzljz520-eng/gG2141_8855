package members

import (
	"errors"
	"strings"
)

type ContactBook struct {
	contacts []ContactPoint
}

func NewContactBook(contacts []ContactPoint) *ContactBook {
	book := &ContactBook{contacts: make([]ContactPoint, 0, len(contacts))}
	for _, contact := range contacts {
		book.Add(contact)
	}
	return book
}

func (b *ContactBook) Add(contact ContactPoint) error {
	if b == nil {
		return errors.New("contact book is nil")
	}
	contact.Label = strings.TrimSpace(contact.Label)
	contact.Name = strings.TrimSpace(contact.Name)
	contact.Phone = normalizePhone(contact.Phone)
	contact.Relationship = strings.TrimSpace(contact.Relationship)
	if contact.Name == "" || contact.Phone == "" {
		return errors.New("contact name and phone are required")
	}
	if contact.Preferred {
		for i := range b.contacts {
			b.contacts[i].Preferred = false
		}
	}
	b.contacts = append(b.contacts, contact)
	return nil
}

func (b *ContactBook) Replace(index int, contact ContactPoint) error {
	if b == nil {
		return errors.New("contact book is nil")
	}
	if index < 0 || index >= len(b.contacts) {
		return errors.New("contact index out of range")
	}
	contact.Label = strings.TrimSpace(contact.Label)
	contact.Name = strings.TrimSpace(contact.Name)
	contact.Phone = normalizePhone(contact.Phone)
	contact.Relationship = strings.TrimSpace(contact.Relationship)
	if contact.Name == "" || contact.Phone == "" {
		return errors.New("contact name and phone are required")
	}
	if contact.Preferred {
		for i := range b.contacts {
			b.contacts[i].Preferred = i == index
		}
	}
	b.contacts[index] = contact
	return nil
}

func (b *ContactBook) Remove(index int) error {
	if b == nil {
		return errors.New("contact book is nil")
	}
	if index < 0 || index >= len(b.contacts) {
		return errors.New("contact index out of range")
	}
	b.contacts = append(b.contacts[:index], b.contacts[index+1:]...)
	if len(b.contacts) > 0 && !b.HasPreferred() {
		b.contacts[0].Preferred = true
	}
	return nil
}

func (b *ContactBook) HasPreferred() bool {
	if b == nil {
		return false
	}
	for _, contact := range b.contacts {
		if contact.Preferred {
			return true
		}
	}
	return false
}

func (b *ContactBook) Preferred() (ContactPoint, bool) {
	if b == nil {
		return ContactPoint{}, false
	}
	for _, contact := range b.contacts {
		if contact.Preferred {
			return contact, true
		}
	}
	return ContactPoint{}, false
}

func (b *ContactBook) All() []ContactPoint {
	if b == nil {
		return nil
	}
	result := make([]ContactPoint, len(b.contacts))
	copy(result, b.contacts)
	return result
}

func (b *ContactBook) FindByLabel(label string) []ContactPoint {
	result := make([]ContactPoint, 0)
	if b == nil {
		return result
	}
	label = strings.ToLower(strings.TrimSpace(label))
	for _, contact := range b.contacts {
		if strings.ToLower(contact.Label) == label {
			result = append(result, contact)
		}
	}
	return result
}
