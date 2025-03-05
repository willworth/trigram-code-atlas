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
		cmd.Build(buildCmd, os.Args[2:])
	case "search":
		searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
		cmd.Search(searchCmd, os.Args[2:])
	case "update":
		updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
		cmd.Update(updateCmd, os.Args[2:])
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
