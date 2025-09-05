package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	if err := q.DeleteAllTransactions(context.Background()); err != nil {
		log.Fatal(err)
	}

	t, err := q.CreateTransaction(context.Background(), database.CreateTransactionParams{
		ID:          uuid.New(),
		Date:        time.Date(2025, time.October, 16, 0, 0, 0, 0, time.UTC),
		Description: "Lembas",
		Amount:      9.99,
		Category:    "Groceries",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Transaction: %v", t)

	t, err = q.CreateTransaction(context.Background(), database.CreateTransactionParams{
		ID:          uuid.New(),
		Date:        time.Date(2025, time.October, 16, 0, 0, 0, 0, time.UTC),
		Description: "Old Toby",
		Amount:      49.99,
		Category:    "Groceries",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Transaction: %v", t)

	t, err = q.CreateTransaction(context.Background(), database.CreateTransactionParams{
		ID:          uuid.New(),
		Date:        time.Date(2025, time.October, 17, 0, 0, 0, 0, time.UTC),
		Description: "Shultzy's",
		Amount:      25.00,
		Category:    "Fun",
		Notes: sql.NullString{
			String: "Bratz n beerz with the boyz",
			Valid:  true},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Transaction: %v", t)

	state := &s.State{
		Database: q,
	}
	router := gin.Default()

	// --- Setup webpages ---
	router.GET("/", func(ctx *gin.Context) {
		html := `<h1>Hello There!</h1>`
		ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	})

	// --- Setup API endpoints ---
	router.GET("/api/transactions", handlers.GetTransactionsHandler(state))
	router.GET("/api/transactions/:id", handlers.GetTransactionByIdHandler(state))
	router.POST("/api/transactions", handlers.PostTransactionsHandler(state))

	if err := router.Run("localhost:8080"); err != nil {
		log.Fatalf("Gin server failed to start: %v", err)
	}
}
