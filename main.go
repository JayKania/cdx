package main

import (
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

	oldState := enableRawMode()
	cleanup := func() {
		clearScreen()
		disableRawMode(oldState)
	}
	defer cleanup()

	selectedOption := 0
	startIndex := 0
	options := buildOptions()
	_, height, _ := term.GetSize(int(os.Stdout.Fd()))
	visibleCount := min(height-1, len(options))
	renderList(startIndex, selectedOption, options, visibleCount)

	// TODO: Figure out a better way to handle resizing, ranther than using this pointer spaghetti
	handleReszing(&visibleCount, &selectedOption, &startIndex, &options)

	buf := make([]byte, 32)

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
			switch buf[2] {
			case 'A': // Up
				if selectedOption > 0 {
					selectedOption--
					if selectedOption < startIndex {
						startIndex--
					}
				} else {
					selectedOption = len(options) - 1
					startIndex = max(len(options)-visibleCount, 0)
				}
			case 'B': // Down
				if selectedOption < len(options)-1 {
					selectedOption++
					if selectedOption >= startIndex+visibleCount {
						startIndex++
					}
				} else {
					selectedOption = 0
					startIndex = 0
				}
			case 'C': // Right
				cwd, _ := os.Getwd()
				if hasSubdirs(cwd + "/" + options[selectedOption]) {
					os.Chdir(cwd + "/" + options[selectedOption])
					selectedOption = 0
					startIndex = 0
					options = buildOptions()
					_, height, _ = term.GetSize(int(os.Stdout.Fd()))
					visibleCount = min(height-1, len(options))
					renderList(startIndex, selectedOption, options, visibleCount)
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
				visibleCount = min(height-1, len(options))
				startIndex = max(0, selectedOption-visibleCount+1)
				renderList(startIndex, selectedOption, options, visibleCount)
			}

			renderList(startIndex, selectedOption, options, visibleCount)
		}
	}
}
