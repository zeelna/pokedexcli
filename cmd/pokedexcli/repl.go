package main

import "strings"

// Verify the input is trimmed of whitespaced, lowercased and split into slice by each whitespace between the words.
func cleanInput(text string) []string {
	trimmed := strings.TrimSpace(text)
	lowered := strings.ToLower(trimmed)
	if lowered == "" {
		return []string{}
	}
	collection := strings.Fields(lowered)
	return collection
}
