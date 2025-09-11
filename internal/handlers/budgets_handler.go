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
		idstr := c.Param("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("Error parsing budget id %s!", idstr))
		}

		budget, err := state.Database.GetBudgetById(c, int32(id))
		if err != nil {
			c.String(http.StatusNotFound, fmt.Sprintf("Failed to get budget id %s from the database!", idstr))
		} else {
			c.IndentedJSON(http.StatusOK, budget)
		}
	}
}

func PostBudgetHandler(state *s.State) func(c *gin.Context) {
	return func(c *gin.Context) {
		var budgetClient s.BudgetClientParams
		if err := c.BindJSON(&budgetClient); err != nil {
			c.String(http.StatusBadRequest, "Failed to parse your JSON!")
			return
		}

		b, err := state.Database.CreateBudget(c, database.CreateBudgetParams{
			TargetAmount: budgetClient.TargetAmount,
			TimePeriod:   database.Period(budgetClient.TimePeriod),
			StartDate:    budgetClient.StartDate,
			Notes: sql.NullString{
				String: budgetClient.Notes,
				Valid:  budgetClient.Notes != "",
			},
			CategoryID: int32(budgetClient.CategoryId),
		})

		if err != nil {
			c.String(http.StatusInternalServerError, "Server encountered an issue creating your budget.")
			return
		}

		c.IndentedJSON(http.StatusCreated, b)
	}
}

func PutBudgetHandler(state *s.State) func(c *gin.Context) {
	return func(c *gin.Context) {
		idstr := c.Param("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			c.String(http.StatusNotFound, fmt.Sprintf("Failed to parse id %s!", idstr))
			return
		}

		var budgetClient s.BudgetClientParams
		if err := c.BindJSON(&budgetClient); err != nil {
			c.String(http.StatusBadRequest, "Failed to parse your JSON!")
			return
		}

		b, err := state.Database.UpdateBudget(c, database.UpdateBudgetParams{
			ID:           int32(id),
			TargetAmount: budgetClient.TargetAmount,
			TimePeriod:   database.Period(budgetClient.TimePeriod),
			StartDate:    budgetClient.StartDate,
			Notes: sql.NullString{
				String: budgetClient.Notes,
				Valid:  budgetClient.Notes != "",
			},
			CategoryID: int32(budgetClient.CategoryId),
		})

		if err != nil {
			c.String(http.StatusInternalServerError, "Server encountered an issue updating your budget.")
			return
		}

		c.IndentedJSON(http.StatusCreated, b)
	}
}

func DeleteBudgetHandler(state *s.State) func(c *gin.Context) {
	return func(c *gin.Context) {
		idstr := c.Param("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			c.String(http.StatusNotFound, fmt.Sprintf("Failed to parse id %s!", idstr))
			return
		}

		if err = state.Database.DeleteBudget(c, int32(id)); err != nil {
			c.String(http.StatusNotFound, fmt.Sprintf("Transaction id %d not found!", id))
		} else {
			c.String(http.StatusOK, fmt.Sprintf("Successfully deleted budget id %d.", id))
		}
	}
}

func GetBudgetStatusHandler(state *s.State) func(c *gin.Context) {
	return func(c *gin.Context) {
		idstr := c.Param("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			c.String(http.StatusNotFound, fmt.Sprintf("Failed to parse id %s!", idstr))
			return
		}

		budget_status, err := state.Database.GetBudgetStatus(c, int32(id))
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to get budget status for id %d", id))
		} else {
			c.IndentedJSON(http.StatusOK, budget_status)
		}
	}
}
