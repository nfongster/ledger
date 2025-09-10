package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	s "github.com/nfongster/ledger/internal/structs"
)

func GetBudgetsHandler(state *s.State) func(c *gin.Context) {
	return func(c *gin.Context) {
		budgets, err := state.Database.GetBudgets(c)
		if err != nil {
			c.String(http.StatusNotFound, "Failed to get budgets from the database!")
		} else {
			c.IndentedJSON(http.StatusOK, budgets)
		}
	}
}

func GetBudgetByIdHandler(state *s.State) func(c *gin.Context) {
	return func(c *gin.Context) {
		// TODO
	}
}

func PostBudgetHandler(state *s.State) func(c *gin.Context) {
	return func(c *gin.Context) {
		// TODO
	}
}

func PutBudgetHandler(state *s.State) func(c *gin.Context) {
	return func(c *gin.Context) {
		// TODO
	}
}

func DeleteBudgetHandler(state *s.State) func(c *gin.Context) {
	return func(c *gin.Context) {
		// TODO
	}
}
