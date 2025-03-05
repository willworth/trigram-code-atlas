package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/cheggaaa/pb/v3"
	"github.com/willworth/trigram-code-atlas/internal/indexer"
	"github.com/willworth/trigram-code-atlas/internal/util"
)

func Build(fs *flag.FlagSet, args []string) {
	// Define flags
	verbose := fs.Bool("verbose", false, "Print detailed logs")
	force := fs.Bool("force", false, "Overwrite existing atlas.json")

	// Extract directory argument first
	var dir string
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		dir = args[0]
		args = args[1:] // Remove the directory from args
	}

	// Parse remaining flags
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	// Check if directory was provided
	if dir == "" {
		fmt.Fprintf(os.Stderr, "Error: Directory argument is required\n")
		fmt.Println("Usage: tca build <dir> [--verbose] [--force]")
		fmt.Println("Example: tca build ./myproject --verbose")
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Build using directory: %s\n", dir) // Debug output

	atlasPath := "atlas.json"
	if !*force && util.FileExists(atlasPath) {
		fmt.Fprintf(os.Stderr, "Error: '%s' already exists. Use --force to overwrite.\n", atlasPath)
		os.Exit(1)
	}

	if *verbose {
		fmt.Printf("Building atlas for '%s'...\n", dir)
	}
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

	if *verbose {
		fmt.Printf("Success: Atlas built with %d files\n", atlas.Metadata.FileCount)
	} else {
		fmt.Println("Success: Atlas built")
	}
}
