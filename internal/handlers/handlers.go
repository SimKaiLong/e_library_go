package handlers

import (
	"e-library-api/internal/errors"
	"e-library-api/internal/models"
	"e-library-api/internal/service"
	stdErrors "errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LibraryHandler struct {
	Service service.LibraryServiceInterface
}

// Create a private helper to handle the repetitive request binding logic
func (h *LibraryHandler) bindRequest(c *gin.Context) (*models.LoanDetail, bool) {
	var input models.LoanDetail
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, false
	}
	return &input, true
}

func (h *LibraryHandler) GetBook(c *gin.Context) {
	title := c.Query("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title parameter is required"})
		return
	}
	book, err := h.Service.GetBook(title)
	if err != nil {
		if stdErrors.Is(err, errors.ErrBookNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, book)
}

func (h *LibraryHandler) BorrowBook(c *gin.Context) {
	input, ok := h.bindRequest(c)
	if !ok {
		return
	}

	loan, err := h.Service.BorrowBook(input.NameOfBorrower, input.BookTitle)
	if err != nil {
		if stdErrors.Is(err, errors.ErrBookNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if stdErrors.Is(err, errors.ErrNoCopies) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, loan)
}

func (h *LibraryHandler) ExtendLoan(c *gin.Context) {
	input, ok := h.bindRequest(c)
	if !ok {
		return
	}

	loan, err := h.Service.ExtendLoan(input.NameOfBorrower, input.BookTitle)
	if err != nil {
		if stdErrors.Is(err, errors.ErrLoanNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, loan)
}

func (h *LibraryHandler) ReturnBook(c *gin.Context) {
	input, ok := h.bindRequest(c)
	if !ok {
		return
	}

	err := h.Service.ReturnBook(input.NameOfBorrower, input.BookTitle)
	if err != nil {
		if stdErrors.Is(err, errors.ErrLoanNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "book returned successfully"})
}
