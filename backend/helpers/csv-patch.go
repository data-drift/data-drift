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
	"sort"
	"strings"
)

func GenerateCsvPatch(currentCsv [][]string, previousCsv [][]string) (string, error) {
	hash1 := md5.Sum([]byte(fmt.Sprintf("%v", previousCsv)))
	hashName1 := hex.EncodeToString(hash1[:])
	file1 := "dist/file-" + hashName1 + "-1.txt"
	hash2 := md5.Sum([]byte(fmt.Sprintf("%v", currentCsv)))
	hashName2 := hex.EncodeToString(hash2[:])
	file2 := "dist/file-" + hashName2 + "-2.txt"

	// Add a space to the last column of the previous csv to make sure it will be present in the diff
	previousCsv[0][len(previousCsv[0])-1] = previousCsv[0][len(previousCsv[0])-1] + " "
	sortedPreviousCsv, err := sortCsvData(previousCsv)
	if err != nil {
		return "", err
	}

	previousCsvString := csvToString(sortedPreviousCsv)
	sortedCurrentCsv, err := sortCsvData(currentCsv)
	if err != nil {
		return "", err
	}
	currentCsvString := csvToString(sortedCurrentCsv)

	// Write the content to the files
	os.Mkdir("dist", 0755)

	err = os.WriteFile(file1, []byte(previousCsvString), 0644)
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

type CsvRecord struct {
	Date      string
	UniqueKey string
	OtherData []string
}

func sortCsvData(csvData [][]string) ([][]string, error) {
	if len(csvData) < 2 {
		return csvData, nil
	}

	dateIndex, uniqueKeyIndex := -1, -1
	for i, columnName := range csvData[0] {
		if columnName == "date" {
			dateIndex = i
		} else if columnName == "unique_key" {
			uniqueKeyIndex = i
		}
	}
	if dateIndex == -1 || uniqueKeyIndex == -1 {
		return csvData, nil
	}

	var records []CsvRecord
	for _, row := range csvData[1:] {
		record := CsvRecord{
			Date:      row[dateIndex],
			UniqueKey: row[uniqueKeyIndex],
			OtherData: append([]string{}, row...),
		}
		records = append(records, record)
	}

	sort.Slice(records, func(i, j int) bool {
		if records[i].Date == records[j].Date {
			return records[i].UniqueKey < records[j].UniqueKey
		}
		return records[i].Date < records[j].Date
	})

	sortedCsvData := [][]string{csvData[0]}
	for _, record := range records {
		sortedCsvData = append(sortedCsvData, record.OtherData)
	}

	return sortedCsvData, nil
}

func csvToString(csvData [][]string) string {
	var csvString strings.Builder
	writer := csv.NewWriter(&csvString)
	writer.WriteAll(csvData)
	writer.Flush()
	return csvString.String()
}
