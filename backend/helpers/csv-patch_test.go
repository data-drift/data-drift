package helpers

import (
	"strings"
	"testing"
)

func TestGenerateCsvPatch(t *testing.T) {
	// Define the test data
	currentCsv := [][]string{
		{"unique_key", "name", "date", "age"},
		{"2022-12-Alice", "Alice", "2022-12", "25"},
		{"2023-01-Bob", "Bob", "2023-01", "30"},
		{"2023-01-Charlie", "Charlie", "2023-01", "36"},
		{"2023-01-Charline", "Charline", "2023-01", "42"},
		{"2023-02-Antoine", "Antoine", "2023-02", "40"},
		{"2023-02-Didou", "Didou", "2023-02", "40"},
		{"2023-02-Philipe", "Philipe", "2023-02", "42"},
		{"2023-03-Cyril", "Cyril", "2023-03", "45"},
		{"2023-03-Victor", "Victor", "2023-03", "46"},
	}
	previousCsv := [][]string{
		{"unique_key", "name", "date", "age"},
		{"2022-12-Alice", "Alice", "2022-12", "25"},
		{"2023-01-Bob", "Bob", "2023-01", "30"},
		{"2023-01-Charlie", "Charlie", "2023-01", "36"},
		{"2023-02-Antoine", "Antoine", "2023-02", "40"},
		{"2023-02-Didier", "Didier", "2023-02", "40"},
		{"2023-02-Philipe", "Philipe", "2023-02", "42"},
		{"2023-03-Clement", "Clement", "2023-03", "45"},
		{"2023-03-Cyril", "Cyril", "2023-03", "45"},
		{"2023-03-Victor", "Victor", "2023-03", "46"},
	}

	// Call the GenerateCsvPatch function with the test data
	patch, err := GenerateCsvPatch(currentCsv, previousCsv)
	if err != nil {
		t.Errorf("GenerateCsvPatch returned an error: %v", err)
	}

	// Check that the patch string is correct
	expectedPatch := "@@ -2,9 +2,9 @@\n 2022-12-Alice,Alice,2022-12,25\n 2023-01-Bob,Bob,2023-01,30\n 2023-01-Charlie,Charlie,2023-01,36\n+2023-01-Charline,Charline,2023-01,42\n 2023-02-Antoine,Antoine,2023-02,40\n-2023-02-Didier,Didier,2023-02,40\n+2023-02-Didou,Didou,2023-02,40\n 2023-02-Philipe,Philipe,2023-02,42\n-2023-03-Clement,Clement,2023-03,45\n 2023-03-Cyril,Cyril,2023-03,45\n 2023-03-Victor,Victor,2023-03,46"

	for i, line := range strings.Split(patch, "\n") {
		if line != strings.Split(expectedPatch, "\n")[i] {
			t.Errorf("GenerateCsvPatch returned an incorrect patch string:\n%s\nExpected:\n%s", patch, expectedPatch)
			break
		}
	}
}
