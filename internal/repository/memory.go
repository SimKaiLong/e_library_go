package repository

import (
	"e-library-api/internal/errors"
	"e-library-api/internal/models"
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
		return nil, errors.ErrBookNotFound
	}
	return book, nil
}

func (m *MemoryRepo) GetLoan(name, title string) (*models.LoanDetail, error) {
	m.RLock()
	defer m.RUnlock()

	loans, ok := m.Loans[title]
	if !ok {
		return nil, errors.ErrLoanNotFound
	}

	for _, l := range loans {
		if l.NameOfBorrower == name {
			return &l, nil
		}
	}
	return nil, errors.ErrLoanNotFound
}

func (m *MemoryRepo) BorrowBook(loan *models.LoanDetail) (*models.LoanDetail, error) {
	m.Lock()
	defer m.Unlock()

	book, ok := m.Books[loan.BookTitle]
	if !ok {
		return nil, errors.ErrBookNotFound
	}
	if book.AvailableCopies <= 0 {
		return nil, errors.ErrNoCopies
	}

	for _, l := range m.Loans[loan.BookTitle] {
		if l.NameOfBorrower == loan.NameOfBorrower {
			return nil, errors.ErrDuplicateLoan
		}
	}

	book.AvailableCopies--
	m.Loans[loan.BookTitle] = append(m.Loans[loan.BookTitle], *loan)
	return loan, nil
}

func (m *MemoryRepo) ExtendLoan(name, title string, newReturnDate time.Time) (*models.LoanDetail, error) {
	m.Lock()
	defer m.Unlock()

	loans, ok := m.Loans[title]
	if !ok {
		return nil, errors.ErrLoanNotFound
	}

	for i, l := range loans {
		if l.NameOfBorrower == name {
			m.Loans[title][i].ReturnDate = newReturnDate
			return &m.Loans[title][i], nil
		}
	}
	return nil, errors.ErrLoanNotFound
}

func (m *MemoryRepo) ReturnBook(name, title string) error {
	m.Lock()
	defer m.Unlock()

	loans, ok := m.Loans[title]
	if !ok {
		return errors.ErrLoanNotFound
	}

	for i, l := range loans {
		if l.NameOfBorrower == name {
			m.Loans[title] = append(loans[:i], loans[i+1:]...)
			m.Books[title].AvailableCopies++
			return nil
		}
	}
	return errors.ErrLoanNotFound
}

func (m *MemoryRepo) Ping() error {
	return nil
}
