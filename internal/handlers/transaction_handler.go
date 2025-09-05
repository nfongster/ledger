package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nfongster/ledger/internal/database"
	s "github.com/nfongster/ledger/internal/structs"
)

func GetTransactionsHandler(state *s.State) func(c *gin.Context) {
	return func(c *gin.Context) {
		var category = c.Query("category")
		fmt.Printf("Category: --%s--\n", category)

		var transactions []database.Transaction
		var err error
		switch {
		case category != "":
			transactions, err = state.Database.GetTransactionsByCategory(c, category)
		default:
			transactions, err = state.Database.GetAllTransactions(c)
		}

		if err != nil {
			c.IndentedJSON(
				http.StatusNotFound,
				gin.H{"message": "Failed to get transactions from the database!"})
			return
		}
		c.IndentedJSON(http.StatusOK, transactions)
	}
}

func GetTransactionByIdHandler(state *s.State) func(c *gin.Context) {
	return func(c *gin.Context) {
		idstr := c.Param("id")
		id, err := uuid.Parse(idstr)
		if err != nil {
			c.IndentedJSON(
				http.StatusNotFound,
				gin.H{"message": fmt.Sprintf("Failed to parse id %s!", idstr)})
			return
		}

		t, err := state.Database.GetTransactionById(c, uuid.UUID(id))
		if err != nil {
			c.IndentedJSON(
				http.StatusNotFound,
				gin.H{"message": fmt.Sprintf("Transaction id %d not found!", id)})
			return
		}

		c.IndentedJSON(http.StatusOK, t)
	}
}

func PostTransactionsHandler(state *s.State) func(c *gin.Context) {
	return func(c *gin.Context) {
		var transactionClient s.TransactionClientParams

		if err := c.BindJSON(&transactionClient); err != nil {
			c.IndentedJSON(
				http.StatusBadRequest,
				gin.H{"message": "Failed to parse your JSON!"})
			return
		}

		t, err := state.Database.CreateTransaction(c, database.CreateTransactionParams{
			ID:          uuid.New(),
			Date:        transactionClient.Date,
			Description: transactionClient.Description,
			Amount:      transactionClient.Amount,
			Category:    transactionClient.Category,
			Notes: sql.NullString{
				String: transactionClient.Notes,
				Valid:  transactionClient.Notes != "",
			},
		})

		if err != nil {
			c.IndentedJSON(
				http.StatusInternalServerError,
				gin.H{"message": "Server encountered an issue creating your transaction."})
			return
		}

		c.IndentedJSON(http.StatusCreated, t)
	}
}
