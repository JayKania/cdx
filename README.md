# CDX

A minimal, fast terminal-based fuzzy directory navigator written in Go. Built with a simple goal: let you **navigate using arrow keys**, search with fuzzy matching, and **copy the `cd` command** of your chosen path to your clipboard.

## Features

- 🔍 Fuzzy search through directories
- ⌨️ Navigate with arrow keys:
  - ↑/↓ to move up and down
  - → to go into a directory
  - ← to move back to parent directory
- ⏎ Press Enter to copy the `cd` command to clipboard
- ✂️ Copies paths in quotes to handle spaces and special characters
- 📏 Handles terminal resizing gracefully
- 🚪 Exits on `Ctrl+C` and shows cursor again

## Installation

### Download the Binary

1. [Download the latest binary](https://github.com/JayKania/cdx/releases) for your OS and architecture.
2. Make it executable:
   ```bash
   chmod +x cdx
3. Move it to a location in your `$PATH` and run it from anywhere, for example:

   ```bash
   sudo mv cdx /usr/local/bin
