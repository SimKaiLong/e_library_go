package repository

import (
	"database/sql"
	"e-library/internal/models"
	"errors"
	"strconv"
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
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("book not found")
	}
	return &b, err
}

func (p *PostgresRepo) BorrowBook(name, title string, days int) (*models.LoanDetail, error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var currentCopies int
	err = tx.QueryRow("SELECT available_copies FROM books WHERE title = $1 FOR UPDATE", title).Scan(&currentCopies)
	if err != nil {
		return nil, errors.New("book not found")
	}
	if currentCopies <= 0 {
		return nil, errors.New("no copies available")
	}

	_, _ = tx.Exec("UPDATE books SET available_copies = available_copies - 1 WHERE title = $1", title)

	loan := models.LoanDetail{
		NameOfBorrower: name,
		BookTitle:      title,
		LoanDate:       time.Now(),
		ReturnDate:     time.Now().AddDate(0, 0, days),
	}

	_, err = tx.Exec("INSERT INTO loans (borrower, title, loan_date, return_date) VALUES ($1, $2, $3, $4)",
		loan.NameOfBorrower, loan.BookTitle, loan.LoanDate, loan.ReturnDate)

	if err != nil {
		return nil, err
	}

	return &loan, tx.Commit()
}

func (p *PostgresRepo) ExtendLoan(name, title string, extrDays int) (*models.LoanDetail, error) {
	var l models.LoanDetail
	query := "UPDATE loans SET return_date = return_date + INTERVAL " + strconv.Itoa(extrDays) + "' days' WHERE borrower = $1 AND title = $2 RETURNING borrower, title, loan_date, return_date"
	err := p.DB.QueryRow(query, name, title).Scan(&l.NameOfBorrower, &l.BookTitle, &l.LoanDate, &l.ReturnDate)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("loan not found")
	}
	return &l, err
}

func (p *PostgresRepo) ReturnBook(name, title string) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, _ := tx.Exec("DELETE FROM loans WHERE borrower = $1 AND title = $2", name, title)
	count, _ := res.RowsAffected()
	if count == 0 {
		return errors.New("no active loan found")
	}

	_, err = tx.Exec("UPDATE books SET available_copies = available_copies + 1 WHERE title = $1", title)
	if err != nil {
		return err
	}
	return tx.Commit()
}
