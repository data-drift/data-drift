package helpers

import (
	"bytes"
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func GenerateCsvPatch(currentCsv [][]string, previousCsv [][]string) (string, error) {
	hash1 := md5.Sum([]byte(fmt.Sprintf("%v", previousCsv)))
	hashName1 := hex.EncodeToString(hash1[:])
	file1 := "dist/file-" + hashName1 + "-1.txt"
	hash2 := md5.Sum([]byte(fmt.Sprintf("%v", currentCsv)))
	hashName2 := hex.EncodeToString(hash2[:])
	file2 := "dist/file-" + hashName2 + "-2.txt"

	previousCsvString := csvToString(previousCsv)
	currentCsvString := csvToString(currentCsv)

	// Write the content to the files
	err := os.WriteFile(file1, []byte(previousCsvString), 0644)
	if err != nil {
		log.Fatalf("Failed to write to %v: %v", file1, err)
	}
	err = os.WriteFile(file2, []byte(currentCsvString), 0644)
	if err != nil {
		log.Fatalf("Failed to write to %v: %v", file2, err)
	}

	// Execute the diff command
	cmd := exec.Command("diff", "-u", file1, file2)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			log.Fatalf("Failed to execute command: %v", err)
		}
	}

	output := out.String()

	lines := strings.Split(output, "\n")

	if len(lines) > 2 {
		output = strings.TrimRight(strings.Join(lines[2:], "\n"), "\n")
	}
	return output, nil
}

func csvToString(csvData [][]string) string {
	var csvString strings.Builder
	writer := csv.NewWriter(&csvString)
	writer.WriteAll(csvData)
	writer.Flush()
	return csvString.String()
}
