package handlers

import (
	"bytes"
	"context"
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

func TestPostTransactionsHandler(t *testing.T) {
	db := setupTestDB()
	defer db.Close()
	q := database.New(db)
	state := &util.State{Database: q}

	util.ResetDatabase(q)

	newTransaction := util.TransactionClientParams{
		Date:        time.Date(2025, time.October, 16, 0, 0, 0, 0, time.UTC),
		Description: "Milk",
		Amount:      5.50,
		Category:    "Groceries",
		Notes:       "Whole milk",
	}

	// Marshal the struct into a JSON byte slice
	jsonValue, err := json.Marshal(newTransaction)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	router := gin.Default()
	router.POST("/api/transactions", PostTransactionsHandler(state))

	testServer := httptest.NewServer(router)
	defer testServer.Close()

	// Make a POST request with the JSON payload
	resp, err := http.Post(testServer.URL+"/api/transactions", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	// Assertions
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected HTTP status code to be 201")

	// Optional: Check the database to confirm the new transaction was added
	transactions, err := q.GetAllTransactions(context.Background())
	if err != nil {
		t.Fatalf("Failed to retrieve transactions from DB: %v", err)
	}
	assert.Len(t, transactions, 1, "Expected one transaction in the database")
	assert.Equal(t, "Milk", transactions[0].Description)
	assert.Equal(t, "Groceries", transactions[0].Category)
	assert.Equal(t, time.Date(2025, time.October, 16, 0, 0, 0, 0, time.UTC), transactions[0].Date.In(time.UTC))
	assert.Equal(t, 5.50, transactions[0].Amount)
	assert.Equal(t, sql.NullString(sql.NullString{String: "Whole milk", Valid: true}), transactions[0].Notes)
}

func TestDeleteTransactionsHandler(t *testing.T) {
	db := setupTestDB()
	defer db.Close()
	q := database.New(db)
	state := &util.State{Database: q}

	util.ResetDatabase(q)
	util.AddTransaction(q, "Coffee", "Groceries", time.Now(), 9.99)

	router := gin.Default()
	router.DELETE("/api/transactions/:id", DeleteTransactionHandler(state))

	testServer := httptest.NewServer(router)
	defer testServer.Close()

	// Use the transaction's ID in the URL path
	transactionID := "1" // Assuming AddTransaction returns an ID or you get it another way
	req, err := http.NewRequest(http.MethodDelete, testServer.URL+"/api/transactions/"+transactionID, nil)
	if err != nil {
		t.Fatalf("Failed to create DELETE request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make DELETE request: %v", err)
	}
	defer resp.Body.Close()

	// Assertions
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected HTTP status code to be 200") // Or http.StatusNoContent (204)

	// Check the database to confirm the transaction was deleted
	transactions, err := q.GetAllTransactions(context.Background())
	if err != nil {
		t.Fatalf("Failed to retrieve transactions from DB: %v", err)
	}
	assert.Len(t, transactions, 0, "Expected no transactions in the database after deletion")
}

func TestPutTransactionsHandler(t *testing.T) {
	db := setupTestDB()
	defer db.Close()
	q := database.New(db)
	state := &util.State{Database: q}

	util.ResetDatabase(q)
	// Add a transaction to be updated
	util.AddTransaction(q, "Coffee", "Groceries", time.Now(), 9.99)

	// Define the updated data
	updatedTransaction := util.TransactionClientParams{
		Date:        time.Date(2025, time.October, 17, 0, 0, 0, 0, time.UTC),
		Description: "Fancy Coffee",
		Amount:      12.50,
		Category:    "Groceries",
		Notes:       "Updated",
	}
	jsonValue, err := json.Marshal(updatedTransaction)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	router := gin.Default()
	router.PUT("/api/transactions/:id", PutTransactionHandler(state))

	testServer := httptest.NewServer(router)
	defer testServer.Close()

	transactionID := "1" // The ID of the transaction to update
	req, err := http.NewRequest(http.MethodPut, testServer.URL+"/api/transactions/"+transactionID, bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Fatalf("Failed to create PUT request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make PUT request: %v", err)
	}
	defer resp.Body.Close()

	// Assertions
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected HTTP status code to be 200")

	// Check the database to confirm the transaction was updated
	transactions, err := q.GetAllTransactions(context.Background())
	if err != nil {
		t.Fatalf("Failed to retrieve transactions from DB: %v", err)
	}
	assert.Equal(t, "Fancy Coffee", transactions[0].Description, "Transaction description was not updated correctly")
}
