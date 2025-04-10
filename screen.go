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
	fmt.Print(clear_screen)
	fmt.Print(clear_scrollback)
	fmt.Print(cursor_top_left)
	fmt.Print(cursor_hide)
}

// renderList displays a portion of the list, highlighting the selected option.
func renderList(startIndex int, selected int, options []string, visibleCount int, searchTerm string) {
	clearScreen()

	endIndex := min(startIndex+visibleCount, len(options))
	cwd, _ := os.Getwd()
	fmt.Printf("  %scd %s %s \n", color_gray, cwd, color_reset)
	fmt.Print(carriage_return)
	fmt.Printf("  %sSearch: %s%s\n", color_gray, searchTerm, color_reset)

	for i := startIndex; i < endIndex; i++ {
		fmt.Print(clear_line)
		fmt.Print(carriage_return)

		width, _, _ := term.GetSize(int(os.Stdout.Fd()))
		visibleText := truncateToWidth(options[i], width-2) // -2 for "> " or padding
		if i == selected {
			fmt.Printf("> %s%s%s\n", color_green, visibleText, color_reset)
		} else {
			fmt.Printf("  %s\n", visibleText)
		}
	}
}
