package main

import (
	"os"
	"strings"
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

func isPrintable(b byte) bool {
	return b >= 32 && b <= 126 // printable ASCII range
}

func sanitizeInput(input string) string {
	safe := ""
	for _, r := range input {
		if r >= 32 && r <= 126 && !strings.ContainsRune(`;&|><\`+"`"+`"'$*~(){}[]`, r) {
			safe += string(r)
		}
	}
	return safe
}
