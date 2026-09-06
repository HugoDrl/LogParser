package parser

import (
	"testing"
)

func TestLevelParse(t *testing.T) {
	tests := map[string]struct {
		input      string
		expected   Level
		expectedOk bool
	}{
		"WARNING level should exist": {
			input:      "WARNING",
			expected:   Warning,
			expectedOk: true,
		},
		"warning level should exist": {
			input:      "warning",
			expected:   Warning,
			expectedOk: true,
		},
		"waRNiNg level minimal should exist": {
			input:      "waRNiNg",
			expected:   Warning,
			expectedOk: true,
		},
		"error level should exist": {
			input:      "error",
			expected:   Error,
			expectedOk: true,
		},
		"info level should exist": {
			input:      "info",
			expected:   Info,
			expectedOk: true,
		},
		"debug level should exist": {
			input:      "DEBUG",
			expected:   Debug,
			expectedOk: true,
		},
		"fatal level should exist": {
			input:      "Fatal",
			expected:   Fatal,
			expectedOk: true,
		},
		"incorrect level should not exist": {
			input:      "incorrect",
			expected:   "",
			expectedOk: false,
		},
		"emtpy level should not exist": {
			input:      "",
			expected:   "",
			expectedOk: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			level, ok := toLevel(test.input)

			if ok != test.expectedOk {
				t.Fatalf("%s: error in OK behavior: expected %t - got %t", name, test.expectedOk, ok)
			}
			if level != test.expected {
				t.Fatalf("%s: expected %s - got %s", name, test.expected, level)
			}
		})
	}
}
