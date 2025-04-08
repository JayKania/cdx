package main

func truncateToWidth(s string, width int) string {
	if len(s) <= width {
		return s
	}
	if width <= 1 {
		return ""
	}
	return s[:width-1] + "…"
}
