package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

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
			return
		}

		category, err := state.Database.GetCategory(c, budget_status.CategoryID)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to get category id %d", budget_status.CategoryID))
			return
		}

		budget, err := state.Database.GetBudgetById(c, budget_status.BudgetID)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to get budget for id %d", id))
		} else {
			c.IndentedJSON(http.StatusOK, struct {
				Category        string
				TimePeriod      database.Period
				StartDate       time.Time
				EndDate         time.Time
				TargetAmount    float64
				CurrentSpent    float64
				RemainingAmount float64
			}{
				Category:        category.Name,
				TimePeriod:      budget_status.TimePeriod,
				StartDate:       budget.StartDate,
				EndDate:         budget_status.EndDate,
				TargetAmount:    budget_status.TargetAmount,
				CurrentSpent:    budget_status.CurrentSpent,
				RemainingAmount: budget_status.TargetAmount - budget_status.CurrentSpent,
			})
		}
	}
}

func GetAllBudgetStatusHandler(state *s.State) func(c *gin.Context) {
	return func(c *gin.Context) {
		all_statuses, err := state.Database.GetAllBudgetStatuses(c)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to get all budget statuses!")
			return
		}

		type BudgetStatus struct {
			Category        string
			TimePeriod      database.Period
			StartDate       time.Time
			EndDate         time.Time
			TargetAmount    float64
			CurrentSpent    float64
			RemainingAmount float64
		}

		statusPayload := []BudgetStatus{}

		for _, status := range all_statuses {
			category, err := state.Database.GetCategory(c, status.CategoryID)
			if err != nil {
				c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to get category id %d", status.CategoryID))
				return
			}

			budget, err := state.Database.GetBudgetById(c, status.BudgetID)
			if err != nil {
				c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to get budget id %d", status.BudgetID))
				return
			}
			statusPayload = append(statusPayload, BudgetStatus{
				Category:        category.Name,
				TimePeriod:      status.TimePeriod,
				StartDate:       budget.StartDate,
				EndDate:         status.EndDate,
				TargetAmount:    status.TargetAmount,
				CurrentSpent:    status.CurrentSpent,
				RemainingAmount: status.TargetAmount - status.CurrentSpent,
			})
		}

		c.IndentedJSON(http.StatusOK, statusPayload)
	}
}
