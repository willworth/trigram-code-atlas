package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/willworth/trigram-code-atlas/internal/indexer"
	"github.com/willworth/trigram-code-atlas/internal/util"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

func Build(fs *flag.FlagSet, args []string) {
	// Define flags
	verbose := fs.Bool("verbose", false, "Print detailed logs")
	force := fs.Bool("force", false, "Overwrite existing output file")
	output := fs.String("output", "", "Specify output file (default: <dir>-atlas-YYYY-MM-DD.json)")

	// Extract directory argument first
	var dir string
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		dir = args[0]
		args = args[1:]
	} else {
		dir = "." // Default to current directory
	}

	// Parse remaining flags
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	// Interactive mode if no args/flags provided
	if len(os.Args) == 2 && os.Args[1] == "build" {
		dir, *verbose, *force, *output = runInteractiveMode()
	}

	// Default output name if not specified
	if *output == "" {
		*output = fmt.Sprintf("%s-atlas-%s.json", filepath.Base(dir), time.Now().Format("2006-01-02"))
	}
	atlasPath := filepath.Join(dir, *output)

	// Debug output
	if *verbose {
		fmt.Fprintf(os.Stderr, "Build using directory: %s, output: %s\n", dir, atlasPath)
	}

	// Check output file
	if !*force && util.FileExists(atlasPath) {
		fmt.Fprintf(os.Stderr, "%sWarning: '%s' already exists.%s\n", colorYellow, atlasPath, colorReset)
		fmt.Print("Would you like to overwrite it? [y/N]: ")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			fmt.Fprintf(os.Stderr, "%sOperation cancelled for '%s'.%s\n", colorRed, atlasPath, colorReset)
			os.Exit(1)
		}
		fmt.Printf("%sContinuing with overwrite of '%s'...%s\n", colorGreen, atlasPath, colorReset)
	}

	// Build with timing
	if *verbose {
		fmt.Printf("%sBuilding atlas for '%s'...%s\n", colorBlue, dir, colorReset)
	}
	start := time.Now()
	atlas, err := indexer.BuildAtlas(dir, *verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to build atlas: %v\n", err)
		os.Exit(1)
	}

	bar := pb.StartNew(atlas.Metadata.FileCount)
	bar.SetWriter(os.Stderr)
	bar.Set("prefix", "Indexing files: ")
	bar.Set("unit", "files")
	for range atlas.Files {
		bar.Increment()
	}
	bar.Finish()

	data, err := json.MarshalIndent(atlas, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to marshal JSON: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(atlasPath, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to write '%s': %v\n", atlasPath, err)
		os.Exit(1)
	}

	elapsed := time.Since(start)
	if *verbose {
		fmt.Printf("%sSuccess: Atlas built with %d files in %v to '%s'%s\n", colorGreen, atlas.Metadata.FileCount, elapsed, atlasPath, colorReset)
	} else {
		fmt.Printf("%sSuccess: Atlas built in %v to '%s'%s\n", colorGreen, elapsed, atlasPath, colorReset)
	}
}

// runInteractiveMode prompts the user for build options with colors
func runInteractiveMode() (dir string, verbose, force bool, output string) {
	defaultOutput := fmt.Sprintf("%s-atlas-%s.json", filepath.Base("."), time.Now().Format("2006-01-02"))

	fmt.Printf("%sWhich directory to index? (default: .): %s", colorBlue, colorReset)
	fmt.Scanln(&dir)
	if dir == "" {
		dir = "."
	}

	fmt.Printf("%sOutput file? (default: %s): %s", colorBlue, defaultOutput, colorReset)
	fmt.Scanln(&output)
	if output == "" || len(output) < 3 { // Basic validation
		output = fmt.Sprintf("%s-atlas-%s.json", filepath.Base(dir), time.Now().Format("2006-01-02"))
	}

	fmt.Printf("%sForce overwrite if file exists? [y/N]: %s", colorYellow, colorReset)
	var forceResp string
	fmt.Scanln(&forceResp)
	force = strings.ToLower(forceResp) == "y"

	fmt.Printf("%sVerbose output? [y/N]: %s", colorGreen, colorReset)
	var verboseResp string
	fmt.Scanln(&verboseResp)
	verbose = strings.ToLower(verboseResp) == "y"

	return dir, verbose, force, output
}
