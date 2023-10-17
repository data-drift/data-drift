package local_store

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func TableHandler(c *gin.Context) {
	store := c.Param("store")
	table := c.Param("table")
	tableColumns, err := getListOfColumnsFromTable(store, table)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	print(tableColumns)
	c.JSON(http.StatusOK, gin.H{
		"store":        store,
		"table":        table,
		"tableColumns": tableColumns,
	})
}

func getListOfColumnsFromTable(store string, table string) ([]string, error) {
	repoDir, err := getStoreDir(store)
	if err != nil {
		fmt.Println("Error getting store directory:", err)
		return nil, err
	}

	filename := table + ".csv"
	file, err := os.Open(filepath.Join(repoDir, filename))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	headers := records[0]
	return headers, nil
}
