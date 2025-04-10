//go:build windows
// +build windows

package main

func handleResizing(visibleCount *int, selectedOption *int, startIndex *int, options *[]string, searchTerm *string) {
    // No-op for Windows as SIGWINCH is not supported
}