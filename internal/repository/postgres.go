package repository

import (
	"database/sql"
	"e-library-api/internal/errors"
	"e-library-api/internal/models"
	stdErrors "errors"
	"time"
)

type PostgresRepo struct {
	DB *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{DB: db}
}

func (p *PostgresRepo) GetBook(title string) (*models.BookDetail, error) {
	var b models.BookDetail
	err := p.DB.QueryRow("SELECT title, available_copies FROM books WHERE title = $1", title).Scan(&b.Title, &b.AvailableCopies)
	if err != nil {
		if stdErrors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrBookNotFound
		}
		return nil, err
	}
	return &b, nil
}

func (p *PostgresRepo) GetLoan(name, title string) (*models.LoanDetail, error) {
	var l models.LoanDetail
	err := p.DB.QueryRow("SELECT borrower, title, loan_date, return_date FROM loans WHERE borrower = $1 AND title = $2", name, title).Scan(&l.NameOfBorrower, &l.BookTitle, &l.LoanDate, &l.ReturnDate)
	if err != nil {
		if stdErrors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrLoanNotFound
		}
		return nil, err
	}
	return &l, nil
}

func (p *PostgresRepo) BorrowBook(loan *models.LoanDetail) (*models.LoanDetail, error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var currentCopies int
	err = tx.QueryRow("SELECT available_copies FROM books WHERE title = $1 FOR UPDATE", loan.BookTitle).Scan(&currentCopies)
	if err != nil {
		if stdErrors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrBookNotFound
		}
		return nil, err
	}
	if currentCopies <= 0 {
		return nil, errors.ErrNoCopies
	}

	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM loans WHERE borrower = $1 AND title = $2)", loan.NameOfBorrower, loan.BookTitle).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrDuplicateLoan
	}

	if _, err = tx.Exec("UPDATE books SET available_copies = available_copies - 1 WHERE title = $1", loan.BookTitle); err != nil {
		return nil, err
	}

	_, err = tx.Exec("INSERT INTO loans (borrower, title, loan_date, return_date) VALUES ($1, $2, $3, $4)",
		loan.NameOfBorrower, loan.BookTitle, loan.LoanDate, loan.ReturnDate)

	if err != nil {
		return nil, err
	}

	return loan, tx.Commit()
}

func (p *PostgresRepo) ExtendLoan(name, title string, newReturnDate time.Time) (*models.LoanDetail, error) {
	var l models.LoanDetail
	query := "UPDATE loans SET return_date = $1 WHERE borrower = $2 AND title = $3 RETURNING borrower, title, loan_date, return_date"
	err := p.DB.QueryRow(query, newReturnDate, name, title).Scan(&l.NameOfBorrower, &l.BookTitle, &l.LoanDate, &l.ReturnDate)
	if err != nil {
		if stdErrors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrLoanNotFound
		}
		return nil, err
	}
	return &l, nil
}

func (p *PostgresRepo) ReturnBook(name, title string) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec("DELETE FROM loans WHERE borrower = $1 AND title = $2", name, title)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.ErrLoanNotFound
	}

	_, err = tx.Exec("UPDATE books SET available_copies = available_copies + 1 WHERE title = $1", title)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (p *PostgresRepo) Ping() error {
	return p.DB.Ping()
}
