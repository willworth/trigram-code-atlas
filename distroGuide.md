# Distribution Guides
1. Install on Your Local Machine
To run tca from anywhere on your Mac:

    Build the Binary:
    bash

    cd ~/trigram-code-atlas
    go build -o tca

    Move to PATH:
    bash

    sudo mv tca /usr/local/bin/

        /usr/local/bin is in your PATH by default on macOS.
    Test:
    bash

    cd ~/some/other/repo
    tca build . --verbose

        Indexes the current dir, writes atlas.json there.

2. Create a GitHub Release
Assuming your repo is github.com/willworth/trigram-code-atlas:

    Commit Changes:
    bash

    git add .
    git commit -m "Add output file option, timing, and current dir default"
    git push origin main

    Tag a Version:
    bash

    git tag v0.1.0
    git push origin v0.1.0

    Build Binaries:
    bash

    GOOS=darwin GOARCH=arm64 go build -o tca-darwin-arm64
    GOOS=linux GOARCH=amd64 go build -o tca-linux-amd64
    GOOS=windows GOARCH=amd64 go build -o tca-windows-amd64.exe

    Create Release:
        Go to github.com/willworth/trigram-code-atlas/releases.
        Click “Draft a new release”.
        Enter tag: v0.1.0, title: “v0.1.0 - Initial Release”.
        Upload the binaries (tca-darwin-arm64, etc.).
        Write release notes (e.g., “Initial build command with custom output and timing”).
        Publish!


