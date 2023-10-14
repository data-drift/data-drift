package local_store

import (
	"fmt"
	"net/http"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/src-d/go-git.v4"
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

	repo, err := git.PlainOpen(repoDir)
	if err != nil {
		fmt.Println("Error opening repo:", err)
		return nil
	}

	worktree, err := repo.Worktree()
	if err != nil {
		fmt.Println("Error getting worktree:", err)
		return nil
	}
	files, err := worktree.Filesystem.ReadDir(".")
	if err != nil {
		fmt.Println("Error getting files:", err)
		return nil
	}
	fileNames := []string{}

	for _, file := range files {
		fileName := file.Name()
		if strings.HasSuffix(fileName, ".csv") {
			fileNameWithoutExt := strings.TrimSuffix(fileName, ".csv")
			fileNames = append(fileNames, fileNameWithoutExt)
		}
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
