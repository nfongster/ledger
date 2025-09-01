package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nfongster/ledger/internal/data"
)

func GetTransactionsHandler(l *data.Ledger) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, l.GetTransactions())
	}
}

func GetTransactionByIdHandler(l *data.Ledger) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		idstr := ctx.Param("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			ctx.IndentedJSON(
				http.StatusNotFound,
				gin.H{"message": fmt.Sprintf("Failed to parse id %s!", idstr)})
			return
		}

		t, err := l.GetTransaction(id)
		if err != nil {
			ctx.IndentedJSON(
				http.StatusNotFound,
				gin.H{"message": fmt.Sprintf("Transaction id %d not found!", id)})
			return
		}

		ctx.IndentedJSON(http.StatusOK, t)
	}
}

func PostTransactionsHandler(l *data.Ledger) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var t data.Transaction

		if err := ctx.BindJSON(&t); err != nil {
			return
		}

		l.AddTransaction(t)
		ctx.IndentedJSON(http.StatusCreated, t)
	}
}
