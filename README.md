# TrigramCodeAtlas (`tca`)

**TrigramCodeAtlas** (`tca`) is a command-line tool written in Go that indexes codebases using trigrams, making them AI-ready for tools like LLMs (e.g., Grok, Claude). It's built for developers who want to automate workflows across multi-language projects, inspired by Yacine's vision of AI-driven productivity.

## Features
- Indexes source files (`.go`, `.js`, `.ts`, etc.) into a JSON atlas (`atlas.json`) with trigrams.
- Ignores binary files, hidden dirs (`.git`), and dependency dirs (`node_modules`).
- Supports custom exclusions via `.tcaignore` (like `.gitignore`).
- Fast, concurrent file scanning with Go goroutines.

## Installation
1. **Prerequisites**: Go 1.23+ (e.g., `brew install go` on macOS).
2. **Clone & Build**:
   ```bash
   git clone https://github.com/willworth/trigram-code-atlas.git
   cd trigram-code-atlas
   go build
   ```
3. **Run**: Use `./tca` or `go run main.go`.

## Usage
### Commands
* `tca build <dir>`: Indexes the directory into atlas.json.
   * `--verbose`: Show detailed logs.
   * `--force`: Overwrite existing atlas.json.
   * Example: `tca build ./myrepo --verbose --force`
* `tca search <query>`: (Planned) Search the atlas for trigrams matching `<query>`.
* `tca update <dir>`: (Planned) Update an existing atlas with changes.

### Example
```bash
# Index a repo
tca build ./testrepo --force
# Output: atlas.json with trigram-indexed files
```

### .tcaignore
Add patterns to `.tcaignore` in your repo root to exclude files:

```
*.txt
dist/
```

## Roadmap
* **MVP Completion**: Fix flag parsing for build, implement search and update.
* **Enhancements**: Live progress bar during indexing, regex search, MCP server integration.


## Contributing
Fork, tweak, PR—let's make it better! See DEV-README.md for dev details.

## License
MIT.