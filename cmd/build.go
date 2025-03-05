package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/cheggaaa/pb/v3"
	"github.com/willworth/trigram-code-atlas/internal/indexer"
	"github.com/willworth/trigram-code-atlas/internal/util"
)

func Build(fs *flag.FlagSet, args []string) {
	verbose := fs.Bool("verbose", false, "Print detailed logs")
	force := fs.Bool("force", false, "Overwrite existing atlas.json")
	fs.Parse(args)

	if len(args) != 1 {
		fmt.Println("Usage: tca build <dir> [--verbose] [--force]")
		os.Exit(1)
	}
	dir := args[0]

	atlasPath := "atlas.json"
	if !*force && util.FileExists(atlasPath) {
		fmt.Println("Error: atlas.json already exists. Use --force to overwrite.")
		os.Exit(1)
	}

	if *verbose {
		fmt.Printf("Building atlas for %s...\n", dir)
	}
	atlas, err := indexer.BuildAtlas(dir, *verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building atlas: %v\n", err)
		os.Exit(1)
	}

	bar := pb.StartNew(atlas.Metadata.FileCount)
	bar.SetWriter(os.Stderr)
	for range atlas.Files {
		bar.Increment()
	}
	bar.Finish()

	data, err := json.MarshalIndent(atlas, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(atlasPath, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing atlas.json: %v\n", err)
		os.Exit(1)
	}

	if *verbose {
		fmt.Printf("Atlas built with %d files\n", atlas.Metadata.FileCount)
	}
}
