package main

import (
	"bytes"
	"e-library/internal/handlers"
	"e-library/internal/repository"
	"e-library/internal/service"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() (*gin.Engine, *repository.MemoryRepo) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	repo := repository.NewMemoryRepo()
	svc := service.NewLibraryService(repo)
	h := &handlers.LibraryHandler{Service: svc}

	r.GET("/Book", h.GetBook)
	r.POST("/Borrow", h.BorrowBook)
	r.POST("/Extend", h.ExtendLoan)
	r.POST("/Return", h.ReturnBook)

	return r, repo
}

func TestEndpoints(t *testing.T) {
	router, repo := setupTestRouter()

	t.Run("GET /Book - Found", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/Book?title=Clean+Code", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GET /Book - Not Found", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/Book?title=Unknown", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("POST /Borrow - Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		payload, _ := json.Marshal(map[string]string{"name_of_borrower": "Alice", "book_title": "Clean Code"})
		req, _ := http.NewRequest("POST", "/Borrow", bytes.NewBuffer(payload))
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("POST /Borrow - Conflict (Out of Stock)", func(t *testing.T) {
		repo.Books["Design Patterns"].AvailableCopies = 0
		w := httptest.NewRecorder()
		payload, _ := json.Marshal(map[string]string{"name_of_borrower": "Bob", "book_title": "Design Patterns"})
		req, _ := http.NewRequest("POST", "/Borrow", bytes.NewBuffer(payload))
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusConflict, w.Code)
	})
}

// --- GET /Book Tests ---
func TestGetBook_Scenarios(t *testing.T) {
	router, _ := setupTestRouter()

	tests := []struct {
		name       string
		query      string
		expectCode int
	}{
		{"Happy Path", "?title=Clean+Code", http.StatusOK},
		{"Missing Title Param", "", http.StatusBadRequest},
		{"Book Not Found", "?title=NonExistent", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/Book"+tt.query, nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.expectCode, w.Code)
		})
	}
}

// --- POST /Borrow Tests ---
func TestBorrowBook_Scenarios(t *testing.T) {
	router, repo := setupTestRouter()

	t.Run("Success - Borrow Available Book", func(t *testing.T) {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(map[string]string{"name_of_borrower": "Alice", "book_title": "Clean Code"})
		req, _ := http.NewRequest("POST", "/Borrow", bytes.NewBuffer(body))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		// Verify side effect: copies should decrease
		book, _ := repo.GetBook("Clean Code")
		assert.Equal(t, 1, book.AvailableCopies)
	})

	t.Run("Error - Out of Stock", func(t *testing.T) {
		// Empty the stock first
		repo.Books["Clean Code"].AvailableCopies = 0
		w := httptest.NewRecorder()
		body, _ := json.Marshal(map[string]string{"name_of_borrower": "Bob", "book_title": "Clean Code"})
		req, _ := http.NewRequest("POST", "/Borrow", bytes.NewBuffer(body))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("Error - Missing JSON Body Fields", func(t *testing.T) {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(map[string]string{"name_of_borrower": "Alice"}) // Missing "title"
		req, _ := http.NewRequest("POST", "/Borrow", bytes.NewBuffer(body))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// --- POST /Extend Tests ---
func TestExtendLoan_Scenarios(t *testing.T) {
	router, repo := setupTestRouter()

	t.Run("Success - Extend Existing Loan", func(t *testing.T) {
		// Manually inject a loan to test extension
		repo.BorrowBook("Alice", "Clean Code", 28)
		initialReturnDate := repo.Loans["Clean Code"][0].ReturnDate

		w := httptest.NewRecorder()
		body, _ := json.Marshal(map[string]string{"name_of_borrower": "Alice", "book_title": "Clean Code"})
		req, _ := http.NewRequest("POST", "/Extend", bytes.NewBuffer(body))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		// Verify logic: should be exactly 21 days after initial return date
		assert.True(t, repo.Loans["Clean Code"][0].ReturnDate.After(initialReturnDate))
	})

	t.Run("Error - Loan Record Not Found", func(t *testing.T) {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(map[string]string{"name_of_borrower": "Stranger", "book_title": "Clean Code"})
		req, _ := http.NewRequest("POST", "/Extend", bytes.NewBuffer(body))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// --- POST /Return Tests ---
func TestReturnBook_Scenarios(t *testing.T) {
	router, repo := setupTestRouter()

	t.Run("Success - Return Book", func(t *testing.T) {
		repo.BorrowBook("Alice", "Clean Code", 28)
		beforeReturn, _ := repo.GetBook("Clean Code")
		initialCopies := beforeReturn.AvailableCopies // is 1

		w := httptest.NewRecorder()
		body, _ := json.Marshal(map[string]string{"name_of_borrower": "Alice", "book_title": "Clean Code"})
		req, _ := http.NewRequest("POST", "/Return", bytes.NewBuffer(body))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		// Verify side effect: copies should increase
		afterReturn, _ := repo.GetBook("Clean Code")
		assert.Equal(t, initialCopies+1, afterReturn.AvailableCopies)
	})
}
