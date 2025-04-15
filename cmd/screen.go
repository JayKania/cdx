package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// type UIState struct {
// 	VisibleCount   int
// 	SelectedOption int
// 	StartIndex     int
// 	Options        []string
// }

func clearScreen() {
	fmt.Fprint(os.Stderr, clear_screen)
	fmt.Fprint(os.Stderr, clear_scrollback)
	fmt.Fprint(os.Stderr, cursor_top_left)
	fmt.Fprint(os.Stderr, cursor_hide)
}

// renderList displays a portion of the list, highlighting the selected option.
func renderList(startIndex int, selected int, options []string, visibleCount int, searchTerm string) {
	clearScreen()

	endIndex := min(startIndex+visibleCount, len(options))
	cwd, _ := os.Getwd()
	fmt.Fprintf(os.Stderr, "  %scd %s %s \n", color_gray, cwd, color_reset)
	fmt.Fprint(os.Stderr, carriage_return)
	fmt.Fprintf(os.Stderr, "  %sSearch: %s%s\n", color_gray, searchTerm, color_reset)

	for i := startIndex; i < endIndex; i++ {
		fmt.Fprint(os.Stderr, clear_line)
		fmt.Fprint(os.Stderr, carriage_return)

		width, _, _ := term.GetSize(int(os.Stderr.Fd()))
		visibleText := truncateToWidth(options[i], width-2) // -2 for "> " or padding
		if i == selected {
			fmt.Fprintf(os.Stderr, "> %s%s%s\n", color_green, visibleText, color_reset)
		} else {
			fmt.Fprintf(os.Stderr, "  %s\n", visibleText)
		}
	}
}
