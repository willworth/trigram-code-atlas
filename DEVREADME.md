# TrigramCodeAtlas: Developer Guide

This is the dev companion to `TrigramCodeAtlas` (`tca`), a Go CLI for indexing codebases with trigrams for AI use. Here's the current state, setup, bugs, and next steps as of March 5, 2025.

## Project Setup
- **Dir**: `~/trigram-code-atlas`
- **Go Version**: 1.23.5 (macOS ARM64)
- **Structure**:

  ```
  trigram-code-atlas/
  ├── cmd/
  │   ├── build.go    # tca build command
  │   ├── search.go   # Placeholder
  │   └── update.go   # Placeholder
  ├── internal/
  │   ├── indexer/
  │   │   ├── indexer.go # Core indexing logic
  │   │   └── types.go  # Atlas structs
  │   └── util/
  │       └── util.go   # Helpers (FileExists)
  ├── main.go         # CLI entry
  ├── go.mod          # Module: github.com/willworth/trigram-code-atlas
  └── go.sum
  ```

- **Dependencies**:
  - `github.com/cheggaaa/pb/v3` (progress bar).
  - Stdlib only otherwise (`os`, `filepath`, etc.).
- **Init**:
  ```bash
  go mod init github.com/willworth/trigram-code-atlas
  go get github.com/cheggaaa/pb/v3
  ```

## What's Implemented

- **Core Indexing**: internal/indexer/indexer.go
  - Scans dirs with filepath.WalkDir, uses 4 goroutines for concurrency.
  - Indexes files (.go, .js, etc.), skips .git, node_modules, binaries.
  - Supports .tcaignore (e.g., *.txt).
  - Outputs atlas.json with trigrams and metadata.
- **tca build**: Partially working.
  - `go run main.go build testrepo` works (no flags), creates atlas.json.
  - Sample output:
    ```json
    {
      "files": [
        {"path": "testrepo/main.js", "trigrams": ["fun", "unc", ...], "mtime": "..."}
      ],
      "metadata": {"version": "1.0", "created": "...", "file_count": 1}
    }
    ```
  - Progress bar shows post-indexing.

## Current Bugs

- ~~**Flag Parsing Fails**~~: Fixed in commit [args-bug]
  - Solution: Improved flag parsing in cmd/build.go by:
    1. Extracting directory argument first
    2. Parsing remaining flags separately
    3. Enhanced error messages and progress bar feedback
  - Status: All flags working (`--force`, `--verbose`)
  - Test cases passing:
    ```bash
    go run main.go build testrepo
    go run main.go build testrepo --force
    go run main.go build testrepo --verbose --force
    ```

## Next Steps

- ~~**Fix Flag Parsing**~~: Completed
- **Implement tca search**:
  - Read atlas.json, match trigrams, print results.
  - Flags: --verbose, --limit.
- **Implement tca update**:
  - Compare mtime in atlas.json with filesystem, update incrementally.

## Roadmap

- **MVP Completion**: Fix build flags, add search and update.
- **Enhancements**:
  - Live progress bar during indexing (not post).
  - Regex support in search.
  - MCP server (tca serve) for AI integration.
- **Distribution**: Build binaries, explore SaaS potential (e.g., paid pro features).

## Debugging Tips

- **Test Repo**: testrepo/main.js (content: function main() {}), .tcaignore (*.txt), ignoreme.txt.
- **Run**: `go run main.go build testrepo --force` to reproduce flag bug.
- **Editor**: If imports strip (e.g., VS Code), disable go.formatOnSave temporarily.

## Notes

