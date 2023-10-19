package local_store

import (
	"log"
	"net/http"

	"github.com/data-drift/data-drift/common"
	"github.com/gin-gonic/gin"
)

type MeasurementRequest struct {
	Metric    string           `json:"metric"`
	TimeGrain common.TimeGrain `json:"timegrain"`
}

func MeasurementHandler(c *gin.Context) {
	store := c.Param("store")
	measurementId := c.Param("measurementId")

	var req MeasurementRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println(store, measurementId, req.TimeGrain)
	// get the record, compute every period in the timegrain
}
