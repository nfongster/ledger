package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nfongster/ledger/internal/database"
	s "github.com/nfongster/ledger/internal/structs"
)

func GetTransactionsHandler(state *s.State) func(c *gin.Context) {
	return func(c *gin.Context) {
		var category = c.Query("category")
		fmt.Printf("Category: --%s--\n", category)

		switch {
		case category != "":
			transactions, err := state.Database.GetTransactionsByCategory(c, category)
			if err != nil {
				c.IndentedJSON(
					http.StatusNotFound,
					gin.H{"message": "Failed to get transactions from the database!"})
				return
			}
			c.IndentedJSON(http.StatusOK, transactions)
		default:
			transactions, err := state.Database.GetAllTransactions(c)
			if err != nil {
				c.IndentedJSON(
					http.StatusNotFound,
					gin.H{"message": "Failed to get transactions from the database!"})
				return
			}
			c.IndentedJSON(http.StatusOK, transactions)
		}
	}
}

func GetTransactionByIdHandler(state *s.State) func(c *gin.Context) {
	return func(c *gin.Context) {
		idstr := c.Param("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			c.IndentedJSON(
				http.StatusNotFound,
				gin.H{"message": fmt.Sprintf("Failed to parse id %s!", idstr)})
			return
		}

		t, err := state.Database.GetTransactionById(c, int32(id))
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

		id, err := state.Database.GetOrCreateCategory(c, transactionClient.Category)
		if err != nil {
			c.IndentedJSON(
				http.StatusInternalServerError,
				gin.H{"message": "Server encountered an issue creating your transaction."})
			return
		}

		t, err := state.Database.CreateTransaction(c, database.CreateTransactionParams{
			Date:        transactionClient.Date,
			Description: transactionClient.Description,
			Amount:      transactionClient.Amount,
			Notes: sql.NullString{
				String: transactionClient.Notes,
				Valid:  transactionClient.Notes != "",
			},
			CategoryID: id,
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

func DeleteTransactionHandler(state *s.State) func(c *gin.Context) {
	return func(c *gin.Context) {
		idstr := c.Param("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			c.IndentedJSON(
				http.StatusNotFound,
				gin.H{"message": fmt.Sprintf("Failed to parse id %s!", idstr)})
			return
		}

		if err = state.Database.DeleteTransaction(c, int32(id)); err != nil {
			c.IndentedJSON(
				http.StatusNotFound,
				gin.H{"message": fmt.Sprintf("Transaction id %d not found!", id)})
			return
		}

		c.String(http.StatusOK, fmt.Sprintf("Successfully deleted transaction id %d.", id))
	}
}
