package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
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

func TestGetBudgetsStatusHandler(t *testing.T) {
	db := setupTestDB()
	defer db.Close()
	q := database.New(db)
	state := &util.State{Database: q}

	util.ResetDatabase(q)
	util.AddTransaction(q,
		"Coffee",
		"Food",
		time.Date(2025, time.October, 16, 0, 0, 0, 0, time.UTC),
		9.99)
	util.AddBudget(q,
		"Food",
		500.00,
		database.PeriodMonthly,
		time.Date(2025, time.October, 16, 0, 0, 0, 0, time.UTC))

	router := gin.Default()
	router.GET("/api/budgets/status", GetBudgetsHandler(state))

	testServer := httptest.NewServer(router)
	defer testServer.Close()

	resp, err := http.Get(testServer.URL + "/api/budgets/status")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var budgetStatus []database.Budget
	err = json.NewDecoder(resp.Body).Decode(&budgetStatus)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	assert.Equal(t, 1, len(budgetStatus))
	assert.Equal(t, int32(1), budgetStatus[0].CategoryID)
	assert.Equal(t, database.PeriodMonthly, budgetStatus[0].TimePeriod)
	assert.Equal(t, time.Date(2025, time.October, 16, 0, 0, 0, 0, time.UTC), budgetStatus[0].StartDate)
}

func TestPostBudgetsHandler(t *testing.T) {
	db := setupTestDB()
	defer db.Close()
	q := database.New(db)
	state := &util.State{Database: q}

	util.ResetDatabase(q)
	util.AddTransaction(q,
		"Plane ticket",
		"Travel",
		time.Date(2025, time.November, 1, 0, 0, 0, 0, time.UTC),
		450.00)

	newBudget := util.BudgetClientParams{
		TargetAmount: 1200.00,
		TimePeriod:   "monthly",
		StartDate:    time.Date(2025, time.November, 1, 0, 0, 0, 0, time.UTC),
		Category:     "Travel",
	}

	jsonValue, err := json.Marshal(newBudget)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	router := gin.Default()
	router.POST("/api/budgets", PostBudgetHandler(state))

	testServer := httptest.NewServer(router)
	defer testServer.Close()

	// Make a POST request with the JSON payload
	resp, err := http.Post(testServer.URL+"/api/budgets", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected HTTP status code to be 201 Created")

	// Check the database directly to confirm the new budget was added
	budget, err := q.GetBudgetById(context.Background(), 1)
	if err != nil {
		t.Fatalf("Failed to retrieve budgets from DB: %v", err)
	}

	assert.Equal(t, int32(1), budget.CategoryID)
	assert.Equal(t, 1200.00, budget.TargetAmount)
	assert.Equal(t, database.PeriodMonthly, budget.TimePeriod)
	assert.Equal(t, time.Date(2025, time.November, 1, 0, 0, 0, 0, time.UTC), budget.StartDate.In(time.UTC))
}

func TestDeleteBudgetsHandler(t *testing.T) {
	db := setupTestDB()
	defer db.Close()
	q := database.New(db)
	state := &util.State{Database: q}

	util.ResetDatabase(q)
	// Add a budget to be deleted, assuming this helper returns the new ID
	util.AddTransaction(q,
		"Plane ticket",
		"Travel",
		time.Date(2025, time.November, 1, 0, 0, 0, 0, time.UTC),
		450.00)
	util.AddBudget(q,
		"Travel",
		1200.00,
		database.PeriodMonthly,
		time.Date(2025, time.November, 1, 0, 0, 0, 0, time.UTC))

	router := gin.Default()
	router.DELETE("/api/budgets/:id", DeleteBudgetHandler(state))

	testServer := httptest.NewServer(router)
	defer testServer.Close()

	req, err := http.NewRequest(http.MethodDelete, testServer.URL+"/api/budgets/1", nil)
	if err != nil {
		t.Fatalf("Failed to create DELETE request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make DELETE request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected HTTP status code to be 200 OK")

	// Check the database to confirm the budget was deleted
	_, err = q.GetBudgetById(context.Background(), 1)
	if err != nil && err != sql.ErrNoRows {
		t.Fatalf("Failed to get budget from DB: %v", err)
	}
	assert.ErrorIs(t, err, sql.ErrNoRows, "Expected the budget to be deleted")
}

func TestPutBudgetsHandler(t *testing.T) {
	db := setupTestDB()
	defer db.Close()
	q := database.New(db)
	state := &util.State{Database: q}

	util.ResetDatabase(q)
	// Add a budget to be updated
	util.AddTransaction(q,
		"Plane ticket",
		"Travel",
		time.Date(2025, time.November, 1, 0, 0, 0, 0, time.UTC),
		450.00)
	util.AddBudget(q,
		"Travel",
		1200.00,
		database.PeriodMonthly,
		time.Date(2025, time.November, 1, 0, 0, 0, 0, time.UTC))

	newBudget := util.BudgetClientParams{
		TargetAmount: 1100.00,
		TimePeriod:   "monthly",
		StartDate:    time.Date(2025, time.November, 1, 0, 0, 0, 0, time.UTC),
		Category:     "Travel",
	}
	jsonValue, err := json.Marshal(newBudget)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	router := gin.Default()
	router.PUT("/api/budgets/:id", PutBudgetHandler(state))

	testServer := httptest.NewServer(router)
	defer testServer.Close()

	req, err := http.NewRequest(http.MethodPut, testServer.URL+"/api/budgets/1", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Fatalf("Failed to create PUT request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make PUT request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected HTTP status code to be 200 OK")

	// Check the database to confirm the budget was updated
	budget, err := q.GetBudgetById(context.Background(), 1)
	if err != nil {
		t.Fatalf("Failed to retrieve budget from DB: %v", err)
	}
	assert.Equal(t, 1100.00, budget.TargetAmount, "Budget amount was not updated correctly")
}
