//go:build darwin || linux
// +build darwin linux

package main

import (
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/term"
)

func handleResizing(visibleCount *int, selectedOption *int, startIndex *int, options *[]string, searchTerm *string) {
	// resize signal handling
	resizeCh := make(chan os.Signal, 1)
	signal.Notify(resizeCh, syscall.SIGWINCH)

	go func() {
		for range resizeCh {
			_, height, _ := term.GetSize(int(os.Stderr.Fd()))
			*visibleCount = min(height-3, len(*options))
			if *selectedOption < *startIndex {
				*startIndex = *selectedOption
			} else if *selectedOption >= *startIndex+*visibleCount {
				*startIndex = *selectedOption - *visibleCount + 1
			}
			if *startIndex < 0 {
				*startIndex = 0
			}
			renderList(*startIndex, *selectedOption, *options, *visibleCount, *searchTerm)
		}
	}()
}

// Add your resize handling logic here
