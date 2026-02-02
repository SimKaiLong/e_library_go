package handlers

import (
	"e-library/models"
	"e-library/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LibraryHandler struct {
	Repo repository.LibraryRepository
}

func (h *LibraryHandler) GetBook(c *gin.Context) {
	title := c.Query("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title parameter is required"})
		return
	}
	book, err := h.Repo.GetBook(title)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, book)
}

func (h *LibraryHandler) BorrowBook(c *gin.Context) {
	var input models.LoanDetail
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	loan, err := h.Repo.BorrowBook(input.NameOfBorrower, input.BookTitle)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, loan)
}

func (h *LibraryHandler) ExtendLoan(c *gin.Context) {
	var input models.LoanDetail
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	loan, err := h.Repo.ExtendLoan(input.NameOfBorrower, input.BookTitle)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, loan)
}

func (h *LibraryHandler) ReturnBook(c *gin.Context) {
	var input models.LoanDetail
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.Repo.ReturnBook(input.NameOfBorrower, input.BookTitle)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "book returned successfully"})
}
