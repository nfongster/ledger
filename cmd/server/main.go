package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Transaction struct {
	ID          int       `json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Notes       string    `json:"notes"`
}

func main() {
	fmt.Println("Welcome to ledger!")

	// Setup Gin
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		html := `<h1>Hello There!</h1>`
		ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	})
	router.GET("/transactions", getTransactions)

	if err := router.Run("localhost:8080"); err != nil {
		log.Fatalf("Gin server failed to start: %v", err)
	}
}

func getTransactions(ctx *gin.Context) {
	transactions := []Transaction{
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
	}
	ctx.IndentedJSON(http.StatusOK, transactions)
}
