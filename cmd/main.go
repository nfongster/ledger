package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nfongster/ledger/internal/database"
	"github.com/nfongster/ledger/internal/handlers"
	s "github.com/nfongster/ledger/internal/structs"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Welcome to ledger!")

	// --- Initialize the database ---
	cfg, err := s.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", cfg.DbConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	q := database.New(db)

	// TODO: You may not want to reset the DB every time you open it!
	ResetDatabase(q)
	SeedDatabase(q)

	state := &s.State{
		Database: q,
	}
	router := gin.Default()

	// --- Setup webpages ---
	handlerIndex := func(c *gin.Context) {
		c.File("web/html/index.html")
	}
	router.GET("/", handlerIndex)
	router.GET("/index", handlerIndex)

	router.GET("/transactions", func(c *gin.Context) {
		c.File("web/html/transactions.html")
	})
	router.GET("/budgets", func(c *gin.Context) {
		c.File("web/html/budgets.html")
	})
	jsGroup := router.Group("/assets")
	{
		jsGroup.GET("/transactions.js", func(c *gin.Context) {
			c.File("web/js/transactions.js")
		})
	}

	// --- Setup API endpoints ---
	apiGroup := router.Group("/api")
	{
		apiGroup.GET("/transactions", handlers.GetTransactionsHandler(state))
		apiGroup.GET("/transactions/:id", handlers.GetTransactionByIdHandler(state))
		apiGroup.POST("/transactions", handlers.PostTransactionsHandler(state))
		apiGroup.PUT("/transactions/:id", handlers.PutTransactionHandler(state))
		apiGroup.DELETE("/transactions/:id", handlers.DeleteTransactionHandler(state))

		apiGroup.GET("/categories", handlers.GetCategoriesHandler(state))
		apiGroup.GET("/categories/:id/spending", handlers.GetCurrentSpendingHandler(state))

		apiGroup.GET("/budgets", handlers.GetBudgetsHandler(state))
		apiGroup.GET("/budgets/:id", handlers.GetBudgetByIdHandler(state))
		apiGroup.GET("/budgets/:id/status", handlers.GetBudgetStatusHandler(state))
		apiGroup.POST("/budgets", handlers.PostBudgetHandler(state))
		apiGroup.PUT("/budgets/:id", handlers.PutBudgetHandler(state))
		apiGroup.DELETE("/budgets/:id", handlers.DeleteBudgetHandler(state))
	}

	// --- Run server and accept connections from any IP address on host machine (WSL, Windows host, etc.) ---
	if err := router.Run("0.0.0.0:8080"); err != nil {
		log.Fatalf("Gin server failed to start: %v", err)
	}
}

func ResetDatabase(q *database.Queries) {
	// Truncate will permanently delete rows (i.e., the removal cannot be rolled back) and reset the primary key.
	if err := q.TruncateAllTables(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func SeedDatabase(q *database.Queries) {
	// --- Set Transactions ---
	groceries_id, err := q.GetOrCreateCategory(context.Background(), "Groceries")
	if err != nil {
		log.Fatal(err)
	}
	_, err = q.CreateTransaction(context.Background(), database.CreateTransactionParams{
		Date:        time.Date(2025, time.October, 16, 0, 0, 0, 0, time.UTC),
		Description: "Lembas",
		Amount:      9.99,
		CategoryID:  groceries_id,
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = q.GetOrCreateCategory(context.Background(), "Groceries")
	if err != nil {
		log.Fatal(err)
	}
	_, err = q.CreateTransaction(context.Background(), database.CreateTransactionParams{
		Date:        time.Date(2025, time.October, 16, 0, 0, 0, 0, time.UTC),
		Description: "Old Toby",
		Amount:      49.99,
		CategoryID:  groceries_id,
	})
	if err != nil {
		log.Fatal(err)
	}

	eatingout_id, err := q.GetOrCreateCategory(context.Background(), "Eating Out")
	if err != nil {
		log.Fatal(err)
	}
	_, err = q.CreateTransaction(context.Background(), database.CreateTransactionParams{
		Date:        time.Date(2025, time.October, 17, 0, 0, 0, 0, time.UTC),
		Description: "Shultzy's",
		Amount:      25.00,
		CategoryID:  eatingout_id,
		Notes: sql.NullString{
			String: "Bratz n beerz with the boyz",
			Valid:  true},
	})
	if err != nil {
		log.Fatal(err)
	}

	// --- Set Budgets ---
	_, err = q.CreateBudget(context.Background(), database.CreateBudgetParams{
		TargetAmount: 50.00,
		TimePeriod:   database.PeriodWeekly,
		StartDate:    time.Date(2025, time.October, 13, 0, 0, 0, 0, time.UTC),
		CategoryID:   groceries_id,
	})
	if err != nil {
		log.Fatal(err)
	}
}
