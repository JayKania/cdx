package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"

	"golang.org/x/term"
)

// buildOptions returns a sorted list of directories and files in the current directory.
func buildOptions() []string {
	children, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	dirOptions := []string{}

	for _, child := range children {
		if child.IsDir() {
			dirOptions = append(dirOptions, child.Name()+"/")
		}
	}

	slices.Sort(dirOptions)

	return dirOptions
}

func main() {

	fuzzyFlagPtr := flag.Bool("fuzzy-search", false, "Enable fuzzy search")
	flag.Parse()

	oldState := enableRawMode()
	cleanup := func() {
		clearScreen()
		disableRawMode(oldState)
	}
	defer cleanup()

	selectedOption := 0
	startIndex := 0
	options := buildOptions()
	matches := []string{}
	_, height, _ := term.GetSize(int(os.Stdout.Fd()))
	visibleCount := min(height-3, len(options))
	renderList(startIndex, selectedOption, options, visibleCount, "")

	buf := make([]byte, 32)
	searchTerm := ""

	// TODO: Figure out a better way to handle resizing, ranther than using this pointer spaghetti
	handleResizing(&visibleCount, &selectedOption, &startIndex, &options, &searchTerm)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			panic(err)
		}

		// Handle Ctrl+C
		if n == 1 && buf[0] == 0x03 {
			return
		}

		// Handle navigation
		if n == 3 && buf[0] == 0x1b && buf[1] == '[' {
			selectedList := options
			if len(matches) > 0 {
				selectedList = matches
			}

			switch buf[2] {
			case 'A': // Up
				if selectedOption > 0 {
					selectedOption--
					if selectedOption < startIndex {
						startIndex--
					}
				} else {
					selectedOption = len(selectedList) - 1
					startIndex = max(len(selectedList)-visibleCount, 0)
				}
				renderList(startIndex, selectedOption, selectedList, visibleCount, searchTerm)
			case 'B': // Down
				if selectedOption < len(selectedList)-1 {
					selectedOption++
					if selectedOption >= startIndex+visibleCount {
						startIndex++
					}
				} else {
					selectedOption = 0
					startIndex = 0
				}
				renderList(startIndex, selectedOption, selectedList, visibleCount, searchTerm)
			case 'C': // Right
				cwd, _ := os.Getwd()
				if hasSubdirs(cwd + "/" + selectedList[selectedOption]) {
					os.Chdir(cwd + "/" + selectedList[selectedOption])
					selectedOption = 0
					startIndex = 0
					options = buildOptions()
					_, height, _ = term.GetSize(int(os.Stdout.Fd()))
					visibleCount = min(height-3, len(options))
					searchTerm = ""
					matches = matches[:0]
					renderList(startIndex, selectedOption, options, visibleCount, searchTerm)
				}
			case 'D': // Left
				currentPath, _ := filepath.Abs(".")
				parentName := filepath.Base(currentPath) + "/"
				os.Chdir("..")
				options = buildOptions()
				selectedOption = findIndexOfOption(options, parentName)
				if selectedOption == -1 {
					selectedOption = 0
				}
				_, height, _ := term.GetSize(int(os.Stdout.Fd()))
				visibleCount = min(height-3, len(options))
				startIndex = max(0, selectedOption-visibleCount+1)
				matches = matches[:0]
				searchTerm = ""
				renderList(startIndex, selectedOption, options, visibleCount, searchTerm)
			}
		}

		// Handle input
		if n == 1 && isPrintable(buf[0]) {
			searchTerm += sanitizeInput(string(buf[0]))
			matches = search(searchTerm, options, fuzzyFlagPtr)
			_, height, _ := term.GetSize(int(os.Stdout.Fd()))
			if len(matches) > 0 {
				visibleCount = min(height-3, len(matches))
				startIndex = 0
				selectedOption = 0
				renderList(startIndex, selectedOption, matches, visibleCount, searchTerm)
			} else {
				visibleCount = min(height-3, len(options))
				startIndex = max(0, selectedOption-visibleCount+1)
				renderList(startIndex, selectedOption, options, visibleCount, searchTerm)
			}
		}

		// Handle backspace
		if n == 1 && buf[0] == 0x7f {
			if len(searchTerm) > 0 {
				searchTerm = searchTerm[:len(searchTerm)-1]
				matches = search(searchTerm, options, fuzzyFlagPtr)
				_, height, _ := term.GetSize(int(os.Stdout.Fd()))
				if len(matches) > 0 && len(searchTerm) > 0 {
					visibleCount = min(height-3, len(matches))
					startIndex = 0
					selectedOption = 0
					renderList(startIndex, selectedOption, matches, visibleCount, searchTerm)
				} else {
					matches = matches[:0]
					visibleCount = min(height-3, len(options))
					startIndex = max(0, selectedOption-visibleCount+1)
					renderList(startIndex, selectedOption, options, visibleCount, searchTerm)
				}
			} else {
				matches = matches[:0]
				visibleCount = min(height-3, len(options))
				startIndex = max(0, selectedOption-visibleCount+1)
				renderList(startIndex, selectedOption, options, visibleCount, searchTerm)
			}
		}

		// Handle Enter
		if n == 1 && buf[0] == 0x0D {
			cwd, _ := os.Getwd()
			disableRawMode(oldState)
			clearScreen()
			fmt.Print(cursor_show)

			var path string
			if len(matches) > 0 {
				path = fmt.Sprintf("cd \"%s/%s\"", cwd, matches[selectedOption])
			} else {
				path = fmt.Sprintf("cd \"%s/%s\"", cwd, options[selectedOption])
			}

			copyToClipboard(path)
			fmt.Println("✅ Path copied to clipboard:", path)
			os.Exit(0)
		}
	}
}
