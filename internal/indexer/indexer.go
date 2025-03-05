package indexer

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/willworth/trigram-code-atlas/internal/util"
)

// Default file extensions to include
var includedExts = map[string]bool{
	".go": true, ".js": true, ".ts": true, ".jsx": true, ".tsx": true,
	".py": true, ".java": true, ".c": true, ".cpp": true, ".h": true,
	".html": true, ".css": true, ".json": true, ".yaml": true, ".md": true,
}

// Default directories to exclude
var excludedDirs = map[string]bool{
	".git": true, "node_modules": true, "vendor": true,
}

// BuildAtlas creates an atlas from a directory
func BuildAtlas(dir string, verbose bool) (*Atlas, error) {
	// Load .tcaignore patterns
	ignorePatterns, err := loadIgnorePatterns(dir)
	if err != nil && verbose {
		fmt.Fprintf(os.Stderr, "Warning: failed to load .tcaignore: %v\n", err)
	}

	// Channels for concurrency
	filesChan := make(chan string, 100)
	resultsChan := make(chan FileEntry, 100)
	var wg sync.WaitGroup

	// Count files for progress bar (first pass)
	fileCount := 0
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip errors
		}
		if !d.IsDir() && shouldIndex(path, ignorePatterns) {
			fileCount++
		}
		return nil
	})

	// Start 4 workers
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range filesChan {
				entry, err := indexFile(path)
				if err != nil {
					if verbose {
						fmt.Fprintf(os.Stderr, "Skipping %s: %v\n", path, err)
					}
					continue
				}
				resultsChan <- entry
			}
		}()
	}

	// Feed files to workers
	go func() {
		filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if !d.IsDir() && shouldIndex(path, ignorePatterns) {
				filesChan <- path
			}
			return nil
		})
		close(filesChan)
	}()

	// Collect results
	var files []FileEntry
	go func() {
		wg.Wait()
		close(resultsChan)
	}()
	for entry := range resultsChan {
		files = append(files, entry)
	}

	// Build atlas
	atlas := &Atlas{Files: files}
	atlas.Metadata.Version = "1.0"
	atlas.Metadata.Created = time.Now().Format(time.RFC3339)
	atlas.Metadata.FileCount = len(files)
	return atlas, nil
}

// shouldIndex decides if a file should be indexed
func shouldIndex(path string, ignorePatterns []string) bool {
	base := filepath.Base(path)
	ext := filepath.Ext(path)

	// Skip excluded dirs
	for dir := range excludedDirs {
		if strings.Contains(path, string(filepath.Separator)+dir+string(filepath.Separator)) {
			return false
		}
	}

	// Check .tcaignore patterns
	for _, pattern := range ignorePatterns {
		if matched, _ := filepath.Match(pattern, base); matched {
			return false
		}
	}

	// Include only specified extensions
	return includedExts[ext]
}

// indexFile generates trigrams for a single file
func indexFile(path string) (FileEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return FileEntry{}, err
	}
	if len(data) > 10<<20 { // 10MB limit
		return FileEntry{}, fmt.Errorf("file too large (>10MB)")
	}

	trigrams := make(map[string]bool)
	content := string(data)
	for i := 0; i < len(content)-2; i++ {
		trigram := content[i : i+3]
		trigrams[trigram] = true
	}
	trigramList := make([]string, 0, len(trigrams))
	for trigram := range trigrams {
		trigramList = append(trigramList, trigram)
	}

	info, err := os.Stat(path)
	if err != nil {
		return FileEntry{}, err
	}

	return FileEntry{
		Path:     path,
		Trigrams: trigramList,
		Mtime:    info.ModTime().Format(time.RFC3339),
	}, nil
}

// loadIgnorePatterns reads .tcaignore
func loadIgnorePatterns(dir string) ([]string, error) {
	path := filepath.Join(dir, ".tcaignore")
	if !util.FileExists(path) {
		return nil, nil // No ignore file, return empty patterns
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	var patterns []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			patterns = append(patterns, line)
		}
	}
	return patterns, nil
}
