package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/willworth/trigram-code-atlas/cmd"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tca <build|search|update> [args]")
		fmt.Println("Commands:")
		fmt.Println("  build <dir>  - Build an atlas from a directory")
		fmt.Println("  search <query> - Search the atlas")
		fmt.Println("  update <dir> - Update an existing atlas")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "build":
		buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
		buildCmd.Parse(os.Args[2:])          // Parse flags like --verbose, --force
		cmd.Build(buildCmd, buildCmd.Args()) // Pass only positional args (e.g., "testrepo")
	case "search":
		searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
		searchCmd.Parse(os.Args[2:])
		cmd.Search(searchCmd, searchCmd.Args())
	case "update":
		updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
		updateCmd.Parse(os.Args[2:])
		cmd.Update(updateCmd, updateCmd.Args())
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
