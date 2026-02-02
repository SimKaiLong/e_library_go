package models

import (
	"time"
)

type BookDetail struct {
	Title           string `json:"title" binding:"required"`
	AvailableCopies int    `json:"available_copies"`
}

type LoanDetail struct {
	NameOfBorrower string    `json:"name_of_borrower" binding:"required"`
	BookTitle      string    `json:"book_title" binding:"required"`
	LoanDate       time.Time `json:"loan_date"`
	ReturnDate     time.Time `json:"return_date"`
}
