package service

import (
	"e-library/internal/models"
	"e-library/internal/repository"
)

// LibraryService handles business logic such as 4 weeks duration for books borrowed and 3 weeks extension
type LibraryService struct {
	Repo repository.LibraryRepository
}

func NewLibraryService(r repository.LibraryRepository) *LibraryService {
	return &LibraryService{Repo: r}
}

func (s *LibraryService) GetBook(title string) (*models.BookDetail, error) {
	return s.Repo.GetBook(title)
}

func (s *LibraryService) BorrowBook(name, title string) (*models.LoanDetail, error) {
	return s.Repo.BorrowBook(name, title, 28)
}

func (s *LibraryService) ExtendLoan(name, title string) (*models.LoanDetail, error) {
	return s.Repo.ExtendLoan(name, title, 21) // 3 weeks extension rule
}

func (s *LibraryService) ReturnBook(name, title string) error {
	return s.Repo.ReturnBook(name, title)
}
