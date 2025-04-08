package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"slices"
	"syscall"

	"golang.org/x/term"
)

func clearScreen() {
	fmt.Print(clear_screen)
	fmt.Print(clear_scrollback)
	fmt.Print(cursor_top_left)
	fmt.Print(cursor_hide)
}

// renderList displays a portion of the list, highlighting the selected option.
func renderList(startIndex int, selected int, options []string, visibleCount int) {
	clearScreen()

	endIndex := startIndex + visibleCount
	if endIndex > len(options) {
		endIndex = len(options)
	}

	for i := startIndex; i < endIndex; i++ {
		fmt.Print(clear_line)
		fmt.Print(carriage_return)
		if i == selected {
			fmt.Printf("> %s%s%s\n", color_green, options[i], color_reset)
		} else {
			fmt.Printf("  %s\n", options[i])
		}
	}

	fmt.Print(cursor_top_left)
	fmt.Print(cursor_hide)
}

// enableRawMode switches the terminal to raw mode.
func enableRawMode() *term.State {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	return oldState
}

// disableRawMode restores the terminal to its original state.
func disableRawMode(oldState *term.State) {
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	fmt.Print(cursor_show) 
}

// buildOptions returns a sorted list of directories and files in the current directory.
func buildOptions() []string {
	children, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	dirOptions := []string{}
	fileOptions := []string{}

	for _, child := range children {
		if child.IsDir() {
			dirOptions = append(dirOptions, "/"+child.Name())
		} else {
			fileOptions = append(fileOptions, child.Name())
		}
	}

	slices.Sort(dirOptions)
	slices.Sort(fileOptions)

	options := append(dirOptions, fileOptions...)
	options = append([]string{".."}, options...)
	return options
}

func main() {

	oldState := enableRawMode()
	cleanup := func() {
		clearScreen()
		disableRawMode(oldState)
	}

	defer cleanup()

	buf := make([]byte, 32)
	selectedOption := 0
	options := buildOptions()
	startIndex := 0

	_, height, _ := term.GetSize(int(os.Stdout.Fd()))
	visibleCount := height - 1 // subtract one for padding at the top

	if visibleCount > len(options) {
		visibleCount = len(options)
	}

	renderList(startIndex, selectedOption, options, visibleCount)

	// resize signal handling
	resizeCh := make(chan os.Signal, 1)
	signal.Notify(resizeCh, syscall.SIGWINCH)

	go func() {
		for range resizeCh {
			_, height, _ := term.GetSize(int(os.Stdout.Fd()))
			visibleCount = height - 1 // subtract one for padding at the top
			if selectedOption < startIndex {
				startIndex = selectedOption
			} else if selectedOption >= startIndex+visibleCount {
				startIndex = selectedOption - visibleCount + 1
			}
			if startIndex < 0 {
				startIndex = 0
			}
			renderList(startIndex, selectedOption, options, visibleCount)
		}
	}()

	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			panic(err)
		}

		if n == 1 && buf[0] == 3 {
			return
		}

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
					startIndex = len(options) - visibleCount
					if startIndex < 0 {
						startIndex = 0
					}
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
			}
			renderList(startIndex, selectedOption, options, visibleCount)
		}
	}
}
