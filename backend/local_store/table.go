package local_store

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func TableHandler(c *gin.Context) {
	store := c.Param("store")
	table := c.Param("table")
	c.JSON(http.StatusOK, gin.H{
		"store": store,
		"table": table,
	})
}
