package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nfongster/ledger/internal/database"
	util "github.com/nfongster/ledger/internal/util"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

func setupTestDB() *sql.DB {
	cfg, err := util.LoadConfig("../../config.json")
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", cfg.DbConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func TestGetTransactionsHandler(t *testing.T) {
	db := setupTestDB()
	defer db.Close()
	q := database.New(db)
	state := &util.State{Database: q}

	util.ResetDatabase(q)
	util.AddTransaction(q, "Coffee", "Groceries", time.Date(2025, time.October, 16, 0, 0, 0, 0, time.UTC), 9.99)

	router := gin.Default()
	router.GET("/api/transactions", GetTransactionsHandler(state))

	// Start a test server
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	// Make an HTTP request
	resp, err := http.Get(testServer.URL + "/api/transactions")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Assert the response
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected HTTP status code to be 200")

	var transactions []database.GetAllTransactionsRow
	err = json.NewDecoder(resp.Body).Decode(&transactions)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Assert the data returned from the database
	assert.Len(t, transactions, 1, "Expected one transaction to be returned")
	assert.Equal(t, "Coffee", transactions[0].Description)
	assert.Equal(t, "Groceries", transactions[0].Category)
	assert.Equal(t, time.Date(2025, time.October, 16, 0, 0, 0, 0, time.UTC), transactions[0].Date)
	assert.Equal(t, 9.99, transactions[0].Amount)
}

// TODO (budgets):
// GET /api/budgets/status
// POST /api/budgets
// DELETE /api/budgets/:id
// PUT /api/budgets/:id
