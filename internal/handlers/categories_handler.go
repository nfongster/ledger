package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	s "github.com/nfongster/ledger/internal/structs"
)

func GetCategoriesHandler(state *s.State) func(c *gin.Context) {
	return func(c *gin.Context) {
		categories, err := state.Database.GetAllCategories(c)
		if err != nil {
			c.IndentedJSON(
				http.StatusNotFound,
				gin.H{"message": "Failed to get categories from the database!"})
			return
		}
		c.IndentedJSON(http.StatusOK, categories)
	}
}
