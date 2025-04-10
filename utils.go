package main

import (
	"os"
	"os/exec"
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

func copyToClipboard(text string) error {
	cmd := exec.Command("pbcopy")
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	_, err = in.Write([]byte(text))
	if err != nil {
		return err
	}
	in.Close()
	return cmd.Wait()
}
