package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

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
	fmt.Fprint(os.Stderr, cursor_show)
}
