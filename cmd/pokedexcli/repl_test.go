package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCleanInput(t *testing.T) {
	// Step#1: Create a slice of test case structs
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		// add more cases here
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "ChariZard ",
			expected: []string{"charizard"},
		},
		{
			input:    "  ",
			expected: []string{},
		},
		{
			input:    "",
			expected: []string{},
		},
	}

	// Step#2: Loop over cases and run tests
	for i, c := range cases {
		t.Run(fmt.Sprintf("Testcase:%d", i+1), func(t *testing.T) {
			actual := cleanInput(c.input)
			// Check the length of the actual slice against the expected slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			if !reflect.DeepEqual(len(actual), len(c.expected)) {
				t.Fatalf("FAILED TEST: Length mismatch. \n Expected: %d words \n Actual: %d words \n", len(c.expected), len(actual))
				return
			}

			for i := range actual {
				word := actual[i]
				expectedWord := c.expected[i]
				if len(actual) != len(c.expected) {
					// Check each word in the slice
					// if they don't match, use t.Errorf to print an error message
					// and fail the test
					t.Fatalf("FAILED TEST: Word does not match. \n Expected: %s\n Actual: %s\n", expectedWord, word)
					return
				}
			}
		}) // end of t.Run
	}
}
