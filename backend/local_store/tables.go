package local_store

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func TablesHandler(c *gin.Context) {
	store := c.Param("store")
	c.JSON(http.StatusOK, store)
}
