package main

import (
	"fmt"
	"log"
	"os"
	"slices"

	"golang.org/x/term"
)

func clearScreen() {
	// Sequences for clearing both screen and scrollback buffer
	fmt.Print("\033[2J")   // Clear entire screen
	fmt.Print("\033[3J")   // Clear scrollback buffer
	fmt.Print("\033[H")    // Move cursor to top-left corner
}

func renderList(selected int, options []string) {
	clearScreen()
	for i, opt := range options {
		fmt.Print("\033[2K\r") // clear the line before printing
		if i == selected {
			fmt.Printf("> \033[32m%s\033[0m\n", opt) // Green and highlighted
		} else {
			fmt.Printf("  %s\n", opt)
		}
	}
	fmt.Print("\033[2K\r") // clear the line before printing
	fmt.Printf("\033[%d;1H", selected+1)
}

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	defer term.Restore(int(os.Stdin.Fd()), oldState)

	buf := make([]byte, 3)
	selectedOption := 0
	options := buildOptions()
	renderList(selectedOption, options)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			panic(err)
		}

		if buf[0] == 'q' {
			break
		}

		if n == 3 && buf[0] == 0x1b && buf[1] == '[' {
			switch buf[2] {
			case 'A':
				{
					if selectedOption == 0 {
						selectedOption = len(options) - 1
					} else {
						selectedOption--
					}
					renderList(selectedOption, options)
				}
			case 'B':
				{
					if selectedOption == len(options) - 1 {
						selectedOption = 0
					} else {
						selectedOption++
					}
					renderList(selectedOption, options)
				}
			}
		}
	}
}

func buildOptions() []string {
	children, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	// separating directories and files as os.ReadDir(".") does not maintain the order
	dirOptions := []string{}
	fileOptions := []string{}

	for i := 0; i < len(children); i++ {
		if children[i].IsDir() {
			dirOptions = append(dirOptions, "/"+children[i].Name())
		} else {
			fileOptions = append(fileOptions, children[i].Name())
		}
	}
	slices.Sort(dirOptions)
	slices.Sort(fileOptions)

	options := append(dirOptions, fileOptions...)
	return options
}