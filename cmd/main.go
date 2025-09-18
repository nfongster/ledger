package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nfongster/ledger/internal/database"
	"github.com/nfongster/ledger/internal/handlers"
	util "github.com/nfongster/ledger/internal/util"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Welcome to ledger!")

	// --- Initialize the database ---
	connStr := util.GetDbConnectionString()
	fmt.Printf("Opening database with connStr %s\n", connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	q := database.New(db)

	// TODO: You may not want to reset the DB every time you open it!
	util.ResetDatabase(q)
	util.SeedDatabase(q)

	state := &util.State{
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
		jsGroup.GET("/budgets.js", func(c *gin.Context) {
			c.File("web/js/budgets.js")
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
		apiGroup.GET("/budgets/status", handlers.GetAllBudgetStatusHandler(state))
		apiGroup.POST("/budgets", handlers.PostBudgetHandler(state))
		apiGroup.PUT("/budgets/:id", handlers.PutBudgetHandler(state))
		apiGroup.DELETE("/budgets/:id", handlers.DeleteBudgetHandler(state))
	}

	// --- Run server and accept connections from any IP address on host machine (WSL, Windows host, etc.) ---
	if err := router.Run("0.0.0.0:8080"); err != nil {
		log.Fatalf("Gin server failed to start: %v", err)
	}
}
