package repository

import (
	"e-library/internal/models"
)

type LibraryRepository interface {
	GetBook(title string) (*models.BookDetail, error)
	BorrowBook(name, title string, days int) (*models.LoanDetail, error)
	ExtendLoan(name, title string, extraDays int) (*models.LoanDetail, error)
	ReturnBook(name, title string) error
}
