package local_store

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func TablesHandler(c *gin.Context) {
	store := c.Param("store")
	tables := getListOfFilesFromStore(store)

	acceptHeader := c.GetHeader("Accept")
	if strings.Contains(acceptHeader, "text/html") {
		html := "<ul>"
		for _, table := range tables {
			encodedTable := url.PathEscape(table)
			url := "./tables/" + encodedTable
			html += "<li><a href=\"" + url + "\">" + table + "</a></li>"
		}
		html += "</ul>"

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	} else {
		c.JSON(http.StatusOK, gin.H{
			"store":  store,
			"tables": tables,
		})
	}
}

func getListOfFilesFromStore(store string) []string {
	repoDir, err := getStoreDir(store)
	if err != nil {
		fmt.Println("Error getting store directory:", err)
		return nil
	}

	fileNames := []string{}

	err = filepath.Walk(repoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".csv") {
			relPath, err := filepath.Rel(repoDir, path)
			relPathWithoutExt := strings.TrimSuffix(relPath, ".csv")
			if err != nil {
				return err
			}
			fileNames = append(fileNames, relPathWithoutExt)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	return fileNames
}

func getStoreDir(store string) (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", nil
	}
	homeDir := currentUser.HomeDir

	repoDir := filepath.Join(homeDir, ".datadrift", store)
	return repoDir, nil
}
