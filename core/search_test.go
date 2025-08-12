package core

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func setupSearchTest(t *testing.T) string {
	tempDir := t.TempDir()

	// Create some test files
	err := os.WriteFile(filepath.Join(tempDir, "file1.txt"), []byte("hello world\nline two has another hello"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(tempDir, "file2.txt"), []byte("HELLO\nthis is a test"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(tempDir, "file4.txt"), []byte("this is helloworld"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Create a subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(subDir, "file3.txt"), []byte("another hello for the test"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Create a binary file
	err = os.WriteFile(filepath.Join(tempDir, "binary.bin"), []byte{0, 1, 2, 3, 'h', 'e', 'l', 'l', 'o'}, 0644)
	if err != nil {
		t.Fatal(err)
	}

	return tempDir
}

// assertResults checks if the search results match the expected results.
func assertResults(t *testing.T, results []SearchResult, expected []SearchResult) {
	t.Helper()

	if len(results) != len(expected) {
		t.Fatalf("Expected %d results, but got %d", len(expected), len(results))
	}

	// Sort results to have a predictable order for comparison
	sort.Slice(results, func(i, j int) bool {
		if results[i].FilePath != results[j].FilePath {
			return results[i].FilePath < results[j].FilePath
		}
		return results[i].LineNumber < results[j].LineNumber
	})
	sort.Slice(expected, func(i, j int) bool {
		if expected[i].FilePath != expected[j].FilePath {
			return expected[i].FilePath < expected[j].FilePath
		}
		return expected[i].LineNumber < expected[j].LineNumber
	})

	for i, res := range results {
		if res.FilePath != expected[i].FilePath || res.FileName != expected[i].FileName || res.LineNumber != expected[i].LineNumber || res.LineContent != expected[i].LineContent {
			t.Errorf("Result %d mismatch.\nGot:      %+v\nExpected: %+v", i, res, expected[i])
		}
	}
}

func TestSearchSimple(t *testing.T) {
	tempDir := setupSearchTest(t)
	options := SearchOptions{MatchCase: false, MatchWholeWord: false, UseRegex: false}
	results, err := SearchFiles(tempDir, "hello", options)
	if err != nil {
		t.Fatalf("SearchFiles returned an error: %v", err)
	}

	expected := []SearchResult{
		{FilePath: filepath.Join(tempDir, "file1.txt"), FileName: "file1.txt", LineNumber: 1, LineContent: "hello world"},
		{FilePath: filepath.Join(tempDir, "file1.txt"), FileName: "file1.txt", LineNumber: 2, LineContent: "line two has another hello"},
		{FilePath: filepath.Join(tempDir, "file2.txt"), FileName: "file2.txt", LineNumber: 1, LineContent: "HELLO"},
		{FilePath: filepath.Join(tempDir, "subdir", "file3.txt"), FileName: "file3.txt", LineNumber: 1, LineContent: "another hello for the test"},
		{FilePath: filepath.Join(tempDir, "file4.txt"), FileName: "file4.txt", LineNumber: 1, LineContent: "this is helloworld"},
	}
	assertResults(t, results, expected)
}

func TestSearchCaseSensitive(t *testing.T) {
	tempDir := setupSearchTest(t)
	options := SearchOptions{MatchCase: true, MatchWholeWord: false, UseRegex: false}
	results, err := SearchFiles(tempDir, "hello", options)
	if err != nil {
		t.Fatalf("SearchFiles returned an error: %v", err)
	}

	expected := []SearchResult{
		{FilePath: filepath.Join(tempDir, "file1.txt"), FileName: "file1.txt", LineNumber: 1, LineContent: "hello world"},
		{FilePath: filepath.Join(tempDir, "file1.txt"), FileName: "file1.txt", LineNumber: 2, LineContent: "line two has another hello"},
		{FilePath: filepath.Join(tempDir, "subdir", "file3.txt"), FileName: "file3.txt", LineNumber: 1, LineContent: "another hello for the test"},
		{FilePath: filepath.Join(tempDir, "file4.txt"), FileName: "file4.txt", LineNumber: 1, LineContent: "this is helloworld"},
	}
	assertResults(t, results, expected)
}

func TestSearchWholeWord(t *testing.T) {
	tempDir := setupSearchTest(t)
	options := SearchOptions{MatchCase: false, MatchWholeWord: true, UseRegex: false}
	results, err := SearchFiles(tempDir, "hello", options)
	if err != nil {
		t.Fatalf("SearchFiles returned an error: %v", err)
	}

	expected := []SearchResult{
		{FilePath: filepath.Join(tempDir, "file1.txt"), FileName: "file1.txt", LineNumber: 1, LineContent: "hello world"},
		{FilePath: filepath.Join(tempDir, "file1.txt"), FileName: "file1.txt", LineNumber: 2, LineContent: "line two has another hello"},
		{FilePath: filepath.Join(tempDir, "file2.txt"), FileName: "file2.txt", LineNumber: 1, LineContent: "HELLO"},
		{FilePath: filepath.Join(tempDir, "subdir", "file3.txt"), FileName: "file3.txt", LineNumber: 1, LineContent: "another hello for the test"},
	}
	assertResults(t, results, expected)
}

func TestSearchRegex(t *testing.T) {
	tempDir := setupSearchTest(t)
	options := SearchOptions{MatchCase: false, UseRegex: true}
	results, err := SearchFiles(tempDir, "h.llo", options)
	if err != nil {
		t.Fatalf("SearchFiles returned an error: %v", err)
	}

	expected := []SearchResult{
		{FilePath: filepath.Join(tempDir, "file1.txt"), FileName: "file1.txt", LineNumber: 1, LineContent: "hello world"},
		{FilePath: filepath.Join(tempDir, "file1.txt"), FileName: "file1.txt", LineNumber: 2, LineContent: "line two has another hello"},
		{FilePath: filepath.Join(tempDir, "file2.txt"), FileName: "file2.txt", LineNumber: 1, LineContent: "HELLO"},
		{FilePath: filepath.Join(tempDir, "subdir", "file3.txt"), FileName: "file3.txt", LineNumber: 1, LineContent: "another hello for the test"},
		{FilePath: filepath.Join(tempDir, "file4.txt"), FileName: "file4.txt", LineNumber: 1, LineContent: "this is helloworld"},
	}
	assertResults(t, results, expected)
}

func TestSearchNoResults(t *testing.T) {
	tempDir := setupSearchTest(t)
	options := SearchOptions{}
	results, err := SearchFiles(tempDir, "nonexistent", options)
	if err != nil {
		t.Fatalf("SearchFiles returned an error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

func TestSearchInvalidRegex(t *testing.T) {
	tempDir := setupSearchTest(t)
	options := SearchOptions{UseRegex: true}
	_, err := SearchFiles(tempDir, "[", options)
	if err == nil {
		t.Error("Expected an error for invalid regex, but got nil")
	}
}
