package main

import (
	"fmt"
	"log"
	"os"
	"slices"

	"golang.org/x/term"
)

// clearScreen clears the terminal screen and scrollback buffer, moving the cursor to the top-left.
func clearScreen() {
	fmt.Print("\033[2J") // Clear the entire screen
	fmt.Print("\033[3J") // Clear the scrollback buffer
	fmt.Print("\033[H")  // Move cursor to top-left
}

// renderList displays a list of options, highlighting the selected one.
// It clears the terminal before rendering.
func renderList(selected int, options []string, windowSize int) {
	clearScreen() // Clear the terminal

	start := selected - (windowSize / 2) // Calculate the start of the visible window
    if start < 0 {
        start = 0
    }
    end := start + windowSize
    if end > len(options) {
        end = len(options)
        start = end - windowSize
        if start < 0 {
            start = 0
        }
    }

	// Render only the visible portion of the list
    for i := start; i < end; i++ {
        fmt.Print("\033[2K\r") // Clear the current line
        if i == selected {
            // Highlight the selected option in green
            fmt.Printf("> \033[32m%s\033[0m\n", options[i])
        } else {
            fmt.Printf("  %s\n", options[i]) // Print non-selected options
        }
    }

    // Move the cursor back to the top of the terminal
    fmt.Print("\033[H")
	// fmt.Print("\033[2K\r")              // Clear the last line
	// fmt.Printf("\033[%d;1H", selected+1) // Move cursor to the selected option's row
}

func main() {
	// Switch terminal to raw mode for direct keypress capture.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err) // Exit if raw mode cannot be enabled.
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState) // Restore terminal mode on exit.

	buf := make([]byte, 3) // Buffer for reading keypresses (supports escape sequences).
	selectedOption := 0    // Index of the selected option.
	options := buildOptions() // Get options from current directory.
	windowSize := 10          // Number of visible lines in the terminal

	renderList(selectedOption, options, windowSize) // Initial rendering of the list.

	for {
		n, err := os.Stdin.Read(buf) // Read keypresses.
		if err != nil {
			panic(err) // Exit on read error.
		}

		if buf[0] == 'q' {
			clearScreen() // Clear the screen before exiting.
			break // Exit program on 'q' press.
		}

		// Handle arrow key input (escape sequence of 3 bytes: Esc, '[', 'A' or 'B').
		if n == 3 && buf[0] == 0x1b && buf[1] == '[' {
			switch buf[2] {
			case 'A': // Up arrow key.
				if selectedOption == 0 {
					selectedOption = len(options) - 1 // Wrap to last option.
				} else {
					selectedOption-- // Move selection up.
				}
				renderList(selectedOption, options, windowSize) // Re-render list.
			case 'B': // Down arrow key.
				if selectedOption == len(options)-1 {
					selectedOption = 0 // Wrap to first option.
				} else {
					selectedOption++ // Move selection down.
				}
				renderList(selectedOption, options, windowSize) // Re-render list.
			}
		}
	}
}

// buildOptions reads the current directory, sorts entries, and prefixes directories with '/'.
func buildOptions() []string {
	children, err := os.ReadDir(".") // Read current directory contents.
	if err != nil {
		log.Fatal(err) // Exit if directory cannot be read.
	}

	dirOptions := []string{}   // List for directories.
	fileOptions := []string{}  // List for files.

	for _, child := range children {
		if child.IsDir() {
			dirOptions = append(dirOptions, "/"+child.Name()) // Prefix directories.
		} else {
			fileOptions = append(fileOptions, child.Name())
		}
	}

	slices.Sort(dirOptions)
	slices.Sort(fileOptions)

	options := append(dirOptions, fileOptions...) // Combine directories and files.
	return options
}