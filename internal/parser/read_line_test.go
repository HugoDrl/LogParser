package parser

import (
	"slices"
	"testing"
)

func TestSplitLine(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected []string
	}{
		"one word should return the word": {
			input:    "hello",
			expected: []string{"hello"},
		},
		"one sentence should split each words": {
			input:    "this is a sentence",
			expected: []string{"this", "is", "a", "sentence"},
		},
		"sentence with quotes should keep words in quotes": {
			input:    "this is a sentence \"with several words in quotes\"",
			expected: []string{"this", "is", "a", "sentence", "\"with several words in quotes\""},
		},
		"last word with quotes in it should not do anything": {
			input:    "this is a sen\"tence",
			expected: []string{"this", "is", "a", "sen\"tence"},
		},
		"standard log format should split by field": {
			input:    "2026-01-01T00:00:00Z WARNING service=test message=hey duration=50ms",
			expected: []string{"2026-01-01T00:00:00Z", "WARNING", "service=test", "message=hey", "duration=50ms"},
		},
		"standard log with quotes should split by field": {
			input:    "2026-01-01T00:00:00Z WARNING service=test message=\"this log is very important\" duration=50ms",
			expected: []string{"2026-01-01T00:00:00Z", "WARNING", "service=test", "message=\"this log is very important\"", "duration=50ms"},
		},
		"having a single quote should take all string from there as a single sentence": {
			input:    "2026-01-01T00:00:00Z WARNING service=test message=\"this log is very important duration=50ms",
			expected: []string{"2026-01-01T00:00:00Z", "WARNING", "service=test", "message=\"this log is very important duration=50ms"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			output := splitLine(test.input)
			if !slices.Equal(output, test.expected) {
				t.Fatalf("%s: expected %v, got %v", name, test.expected, output)
			}
		})
	}
}
