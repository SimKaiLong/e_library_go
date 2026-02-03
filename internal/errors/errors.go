package errors

import "errors"

var (
	ErrBookNotFound  = errors.New("book not found")
	ErrNoCopies      = errors.New("no copies available")
	ErrLoanNotFound  = errors.New("loan not found")
	ErrDuplicateLoan = errors.New("borrower already has an active loan for this book")
)
