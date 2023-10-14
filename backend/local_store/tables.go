package local_store

import (
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func TablesHandler(c *gin.Context) {
	store := c.Param("store")
	tables := getListOfFilesFromStore(store)
	c.JSON(http.StatusOK, tables)
}

func getListOfFilesFromStore(store string) []string {
	repoDir, err := getStoreDir(store)
	if err != nil {
		fmt.Println("Error getting store directory:", err)
		return nil
	}

	fileNames := []string{}

	// Walk the directory tree and get all the files that end with ".csv"
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
