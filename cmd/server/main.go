package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nfongster/ledger/internal/data"
	"github.com/nfongster/ledger/internal/handlers"
)

func main() {
	fmt.Println("Welcome to ledger!")

	// Setup data in memory
	ledger := data.NewLedger()
	ledger.AddTransactions([]data.Transaction{
		{
			ID:          0,
			Date:        time.Date(2025, time.October, 16, 0, 0, 0, 0, time.UTC),
			Description: "Trader Joe's",
			Amount:      40.00,
			Category:    "Groceries",
		},
		{
			ID:          1,
			Date:        time.Date(2025, time.October, 16, 0, 0, 0, 0, time.UTC),
			Description: "Grocery Outlet",
			Amount:      10.55,
			Category:    "Groceries",
		},
		{
			ID:          2,
			Date:        time.Date(2025, time.October, 17, 0, 0, 0, 0, time.UTC),
			Description: "Shultzy's",
			Amount:      25.00,
			Category:    "Fun",
			Notes:       "Bratz n beerz with the boyz",
		},
	})

	router := gin.Default()

	// Setup webpages
	router.GET("/", func(ctx *gin.Context) {
		html := `<h1>Hello There!</h1>`
		ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	})

	// Setup API endpoints
	router.GET("/transactions", handlers.GetTransactionsHandler(ledger))
	router.GET("/transactions/:id", handlers.GetTransactionByIdHandler(ledger))
	router.POST("/transactions", handlers.PostTransactionsHandler(ledger))

	if err := router.Run("localhost:8080"); err != nil {
		log.Fatalf("Gin server failed to start: %v", err)
	}
}
