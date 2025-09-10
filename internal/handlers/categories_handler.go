package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nfongster/ledger/internal/database"
	s "github.com/nfongster/ledger/internal/structs"
)

func GetCategoriesHandler(state *s.State) func(c *gin.Context) {
	return func(c *gin.Context) {
		categories, err := state.Database.GetAllCategories(c)
		if err != nil {
			c.String(http.StatusNotFound, "Failed to get categories from the database!")
		} else {
			c.IndentedJSON(http.StatusOK, categories)
		}
	}
}

func GetCurrentSpendingHandler(state *s.State) func(c *gin.Context) {
	return func(c *gin.Context) {
		idstr := c.Param("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			c.String(http.StatusNotFound, fmt.Sprintf("Failed to parse id %s!", idstr))
			return
		}

		start, end := c.Query("period_start"), c.Query("period_end")
		startTime, _ := time.Parse(time.DateOnly, start)
		endTime, _ := time.Parse(time.DateOnly, end)
		urlQuery := struct {
			hasStart bool
			start    time.Time
			hasEnd   bool
			end      time.Time
		}{
			hasStart: start != "",
			start:    startTime,
			hasEnd:   end != "",
			end:      endTime,
		}

		type Spending struct {
			CategoryID int     `json:"category_id"`
			TotalSpent float64 `json:"total_spent"`
		}

		if urlQuery.hasStart && urlQuery.hasEnd {
			amount, err := state.Database.GetSpendingBetweenStartAndEnd(c, database.GetSpendingBetweenStartAndEndParams{
				CategoryID: int32(id),
				StartDate:  urlQuery.start,
				EndDate:    urlQuery.end,
			})
			if err != nil {
				c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to query transaction count for id %d.", id))
			} else {
				spendingStartEnd := struct {
					Spending
					PeriodStart time.Time `json:"period_start"`
					PeriodEnd   time.Time `json:"period_end"`
				}{
					Spending: Spending{
						CategoryID: id,
						TotalSpent: amount,
					},
					PeriodStart: urlQuery.start,
					PeriodEnd:   urlQuery.end,
				}
				c.IndentedJSON(http.StatusOK, spendingStartEnd)
			}
		} else if urlQuery.hasStart {
			amount, err := state.Database.GetSpendingSinceStart(c, database.GetSpendingSinceStartParams{
				CategoryID: int32(id),
				StartDate:  urlQuery.start,
			})
			if err != nil {
				c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to query transaction count for id %d.", id))
			} else {
				spendingStart := struct {
					Spending
					PeriodStart time.Time `json:"period_start"`
				}{
					Spending: Spending{
						CategoryID: id,
						TotalSpent: amount,
					},
					PeriodStart: urlQuery.start,
				}
				c.IndentedJSON(http.StatusOK, spendingStart)
			}
		} else if urlQuery.hasEnd {
			amount, err := state.Database.GetSpendingUntilEnd(c, database.GetSpendingUntilEndParams{
				CategoryID: int32(id),
				EndDate:    urlQuery.end,
			})
			if err != nil {
				c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to query transaction count for id %d.", id))
			} else {
				spendingEnd := struct {
					Spending
					PeriodEnd time.Time `json:"period_start"`
				}{
					Spending: Spending{
						CategoryID: id,
						TotalSpent: amount,
					},
					PeriodEnd: urlQuery.end,
				}
				c.IndentedJSON(http.StatusOK, spendingEnd)
			}
		} else {
			amount, err := state.Database.GetSpendingAllTime(c, int32(id))
			if err != nil {
				c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to query transaction count for id %d.", id))
			} else {
				c.IndentedJSON(http.StatusOK, Spending{
					CategoryID: id,
					TotalSpent: amount,
				})
			}
		}
	}
}
