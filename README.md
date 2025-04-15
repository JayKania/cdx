# CDX

A minimal, fast terminal-based directory navigator written in Go. Built with a simple goal: let you **navigate using arrow keys**, **search directories**, and **change directories easily**.

## Features

- 🔍 Regular search by default
- Fuzzy search
- ⌨️ Navigate with arrow keys:
  - ↑/↓ to move up and down
  - → to go into a directory
  - ← to move back to parent directory
- ⏎ Press Enter to change the directory and exit
- 📏 Handles terminal resizing gracefully
- 🚪 Exits on `Ctrl+C` and shows cursor again

## Installation

### Download the Binary

1. [Download the latest binary](https://github.com/JayKania/cdx/releases) for your OS and architecture.
2. Extract the archive — you'll find two binaries: `cdx` and `setup`.
3. Make both the files executable:
   ```bash
   chmod +x cdx setup
   ```
4. Move the cdx file to a location in your `$PATH`, for example:
   ```bash
   sudo mv cdx /usr/local/bin
   ```
5. Run the setup file to install and integrate CDX:
   ```bash
   setup
6. Restart your terminal and run your cdx program from anywhere