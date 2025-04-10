//go:build windows
// +build windows

package main

func handleResizing() {
    // No-op for Windows as SIGWINCH is not supported
}