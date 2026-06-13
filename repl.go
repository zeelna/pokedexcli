package main

import "strings"

func cleanInput(text string) []string {
	trimmed := strings.TrimSpace(text)
	lowered := strings.ToLower(trimmed)
	if lowered == "" {
		return []string{}
	}
	collection := strings.Fields(lowered)
	return collection
}
