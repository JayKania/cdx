package main

import (
	"sort"
	"strings"
)

// LevenshteinDistance calculates the edit distance between two strings
func LevenshteinDistance(a, b string) int {
	a = strings.ToLower(a)
	b = strings.ToLower(b)

	// Create a matrix with len(a)+1 rows and len(b)+1 columns
	d := make([][]int, len(a)+1)
	for i := range d {
		d[i] = make([]int, len(b)+1)
	}

	// Initialize the first row and column
	for i := 0; i <= len(a); i++ {
		d[i][0] = i
	}
	for j := 0; j <= len(b); j++ {
		d[0][j] = j
	}

	// Fill the matrix
	for j := 1; j <= len(b); j++ {
		for i := 1; i <= len(a); i++ {
			if a[i-1] == b[j-1] {
				d[i][j] = d[i-1][j-1] // No operation
			} else {
				d[i][j] = min(
					d[i-1][j]+1,   // Deletion
					d[i][j-1]+1,   // Insertion
					d[i-1][j-1]+1, // Substitution
				)
			}
		}
	}

	return d[len(a)][len(b)]
}

// FuzzySearch performs fuzzy matching using Levenshtein distance
func FuzzySearch(query string, items []string, maxDistance int) []MatchResult {
	var results []MatchResult

	for _, item := range items {
		distance := LevenshteinDistance(query, item)

		// Only include items within the maximum distance
		if distance <= maxDistance {
			score := 100 - (distance * 10) // Higher score for smaller distance
			results = append(results, MatchResult{
				Text:     item,
				Score:    score,
				Distance: distance,
			})
		}
	}

	// Sort by score (higher is better)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

// MatchResult represents a fuzzy match result with score and distance
type MatchResult struct {
	Text     string
	Score    int
	Distance int
}

func min(args ...int) int {
	min := args[0]
	for _, arg := range args[1:] {
		if arg < min {
			min = arg
		}
	}
	return min
}

// regularSearch performs a regular search for the search term in the options
func search(searchTerm string, options []string) []string {
	matches := []MatchResult{}
	remainingOptionsAfterExactMatch := []string{}
	remainingOptionsAfterPartialMatch := []string{}
	for _, option := range options {
		if strings.HasPrefix(strings.ToLower(option), strings.ToLower(searchTerm)) {
			matches = append(matches, MatchResult{
				Text:     option,
				Score:    100,
				Distance: 0,
			})
		} else {
			remainingOptionsAfterExactMatch = append(remainingOptionsAfterExactMatch, option)
		}
	}

	for _, option := range remainingOptionsAfterExactMatch {
		if strings.Contains(strings.ToLower(option), strings.ToLower(searchTerm)) {
			matches = append(matches, MatchResult{
				Text:     option,
				Score:    100,
				Distance: 0,
			})
		} else {
			remainingOptionsAfterPartialMatch = append(remainingOptionsAfterPartialMatch, option)
		}
	}

	maxDistance := 4
	matches = append(matches, FuzzySearch(searchTerm, remainingOptionsAfterPartialMatch, maxDistance)...)

	results := []string{}
	for _, matches := range matches {
		results = append(results, matches.Text)
	}
	
	return results;
}
