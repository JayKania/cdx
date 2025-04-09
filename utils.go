package main

import (
	"os"
)

func truncateToWidth(s string, width int) string {
	if len(s) <= width {
		return s
	}
	if width <= 1 {
		return ""
	}
	return s[:width-1] + "…"
}

func hasSubdirs(path string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if entry.IsDir() {
			return true
		}
	}
	return false
}

func findIndexOfOption(options []string, option string) int {
	for i, opt := range options {
		if opt == option {
			return i
		}
	}
	return -1
}
