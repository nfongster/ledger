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

type Ledger struct {
	transactions []Transaction
}

func NewLedger() *Ledger {
	return &Ledger{
		transactions: make([]Transaction, 0),
	}
}

func (l *Ledger) AddTransaction(t Transaction) {
	l.transactions = append(l.transactions, t)
}

func (l *Ledger) AddTransactions(tSlice []Transaction) {
	l.transactions = append(l.transactions, tSlice...)
}

func (l *Ledger) GetTransactions() []Transaction {
	return l.transactions
}

func main() {
	fmt.Println("Welcome to ledger!")

	// Setup data in memory
	ledger := NewLedger()
	ledger.AddTransactions([]Transaction{
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

	// Setup Gin
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		html := `<h1>Hello There!</h1>`
		ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	})
	router.GET("/transactions", getTransactionsHandler(ledger))
	router.POST("/transactions", postTransactionsHandler(ledger))

	if err := router.Run("localhost:8080"); err != nil {
		log.Fatalf("Gin server failed to start: %v", err)
	}
}

func getTransactionsHandler(l *Ledger) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, l.GetTransactions())
	}
}

func postTransactionsHandler(l *Ledger) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var t Transaction

		if err := ctx.BindJSON(&t); err != nil {
			return
		}

		l.AddTransaction(t)
		ctx.IndentedJSON(http.StatusCreated, t)
	}
}
