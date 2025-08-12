package core

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

// SearchResult represents a single search result.
type SearchResult struct {
	FilePath    string `json:"file_path"`
	FileName    string `json:"file_name"`
	LineNumber  int    `json:"line_number"`
	LineContent string `json:"line_content"`
}

// SearchOptions defines the options for a search operation.
type SearchOptions struct {
	MatchCase      bool `json:"match_case"`
	MatchWholeWord bool `json:"match_whole_word"`
	UseRegex       bool `json:"use_regex"`
}

// worker is a goroutine that processes files from the files channel and sends results to the results channel.
func worker(wg *sync.WaitGroup, files <-chan string, results chan<- SearchResult, matcher func(string) bool) {
	defer wg.Done()
	for file := range files {
		if IsBinary(file) {
			continue
		}

		f, err := os.Open(file)
		if err != nil {
			// skip files we can't open
			continue
		}

		scanner := bufio.NewScanner(f)
		lineNumber := 0
		for scanner.Scan() {
			lineNumber++
			line := scanner.Text()
			if matcher(line) {
				results <- SearchResult{
					FilePath:    file,
					FileName:    filepath.Base(file),
					LineNumber:  lineNumber,
					LineContent: line,
				}
			}
		}
		f.Close()
	}
}

func search(rootPath, query string, options SearchOptions) ([]SearchResult, error) {
	var wg sync.WaitGroup
	results := make(chan SearchResult)
	files := make(chan string)

	var matcher func(string) bool

	if options.UseRegex {
		regexStr := query
		if !options.MatchCase {
			regexStr = "(?i)" + regexStr
		}
		re, err := regexp.Compile(regexStr)
		if err != nil {
			return nil, fmt.Errorf("invalid regular expression: %w", err)
		}
		matcher = re.MatchString
	} else if options.MatchWholeWord {
		regexStr := `\b` + regexp.QuoteMeta(query) + `\b`
		if !options.MatchCase {
			regexStr = "(?i)" + regexStr
		}
		re, err := regexp.Compile(regexStr)
		if err != nil {
			return nil, fmt.Errorf("internal error compiling regex for whole word search: %w", err)
		}
		matcher = re.MatchString
	} else {
		if options.MatchCase {
			matcher = func(s string) bool {
				return strings.Contains(s, query)
			}
		} else {
			lowerQuery := strings.ToLower(query)
			matcher = func(s string) bool {
				return strings.Contains(strings.ToLower(s), lowerQuery)
			}
		}
	}

	// Start a pool of workers.
	numWorkers := runtime.NumCPU()
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(&wg, files, results, matcher)
	}

	// Walk the directory tree and send file paths to the files channel.
	go func() {
		defer close(files)
		filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if !d.IsDir() {
				files <- path
			}
			return nil
		})
	}()

	// Start a goroutine to wait for all workers to finish, then close the results channel.
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect all results from the results channel.
	var searchResults []SearchResult
	for result := range results {
		searchResults = append(searchResults, result)
	}

	return searchResults, nil
}
