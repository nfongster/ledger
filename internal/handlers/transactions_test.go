package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/nfongster/ledger/internal/database"
	s "github.com/nfongster/ledger/internal/structs"
	"github.com/stretchr/testify/assert"
)

func TestGetTransactionsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	queries := database.New(db)
	state := &s.State{Database: queries}

	router.GET("/api/transactions", GetTransactionsHandler(state))

	t.Run("Get all transactions", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "date", "description", "amount", "notes", "category"}).
			AddRow(1, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), "Transaction 1", 50.00, "Notes 1", "Groceries").
			AddRow(2, time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC), "Transaction 2", 25.50, "Notes 2", "Utilities")
		const sqlQuery = `SELECT t.id, t.date, t.description, t.amount, t.notes, c.name AS category
FROM transactions AS t
JOIN categories AS c ON t.category_id = c.id
`
		mock.ExpectQuery(sqlQuery).WillReturnRows(rows)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/transactions", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"Description": "Transaction 1"`)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Get transactions by category", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "date", "description", "amount", "notes", "category"}).
			AddRow(1, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), "Transaction 1", 50.00, "Notes 1", "Groceries")
		mock.ExpectQuery("SELECT t.id, t.date, t.description, t.amount, t.notes, c.name AS category FROM transactions AS t JOIN categories AS c ON t.category_id = c.id WHERE c.name = \\$1").WithArgs("Groceries").WillReturnRows(rows)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/transactions?category=Groceries", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"Category": "Groceries"`)
		assert.NotContains(t, w.Body.String(), `"Category": "Utilities"`)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TODO (transactions):
// POST /api/transactions
// DELETE /api/transactions/:id
// PUT /api/transactions/:id

// TODO (budgets):
// GET /api/budgets/status
// POST /api/budgets
// DELETE /api/budgets/:id
// PUT /api/budgets/:id
