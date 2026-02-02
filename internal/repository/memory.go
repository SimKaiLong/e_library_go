package repository

import (
	"e-library/internal/models"
	"errors"
	"sync"
	"time"
)

type MemoryRepo struct {
	sync.RWMutex
	Books map[string]*models.BookDetail
	Loans map[string][]models.LoanDetail
}

func NewMemoryRepo() *MemoryRepo {
	repo := &MemoryRepo{
		Books: make(map[string]*models.BookDetail),
		Loans: make(map[string][]models.LoanDetail),
	}
	// Seed data
	repo.Books["The Go Programming Language"] = &models.BookDetail{Title: "The Go Programming Language", AvailableCopies: 5}
	repo.Books["Clean Code"] = &models.BookDetail{Title: "Clean Code", AvailableCopies: 2}
	repo.Books["Design Patterns"] = &models.BookDetail{Title: "Design Patterns", AvailableCopies: 1}
	return repo
}

func (m *MemoryRepo) GetBook(title string) (*models.BookDetail, error) {
	m.RLock()
	defer m.RUnlock()
	book, ok := m.Books[title]
	if !ok {
		return nil, errors.New("book not found")
	}
	return book, nil
}

func (m *MemoryRepo) BorrowBook(name, title string, days int) (*models.LoanDetail, error) {
	m.Lock()
	defer m.Unlock()

	book, ok := m.Books[title]
	if !ok {
		return nil, errors.New("book not found")
	}
	if book.AvailableCopies <= 0 {
		return nil, errors.New("no copies available")
	}

	book.AvailableCopies--
	loan := models.LoanDetail{
		NameOfBorrower: name,
		BookTitle:      title,
		LoanDate:       time.Now(),
		ReturnDate:     time.Now().AddDate(0, 0, days), // 4 weeks
	}
	m.Loans[title] = append(m.Loans[title], loan)
	return &loan, nil
}

func (m *MemoryRepo) ExtendLoan(name, title string, extrDays int) (*models.LoanDetail, error) {
	m.Lock()
	defer m.Unlock()

	loans, ok := m.Loans[title]
	if !ok {
		return nil, errors.New("no active loan found")
	}

	for i, l := range loans {
		if l.NameOfBorrower == name {
			m.Loans[title][i].ReturnDate = l.ReturnDate.AddDate(0, 0, extrDays)
			return &m.Loans[title][i], nil
		}
	}
	return nil, errors.New("loan record not found for user")
}

func (m *MemoryRepo) ReturnBook(name, title string) error {
	m.Lock()
	defer m.Unlock()

	loans, ok := m.Loans[title]
	if !ok {
		return errors.New("loan record not found")
	}

	for i, l := range loans {
		if l.NameOfBorrower == name {
			m.Loans[title] = append(loans[:i], loans[i+1:]...)
			m.Books[title].AvailableCopies++
			return nil
		}
	}
	return errors.New("loan record not found")
}
