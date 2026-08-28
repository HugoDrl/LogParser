package parser

import (
	"testing"
	"time"
)

func TestLogFilter(t *testing.T) {
	testLog := &Log{
		Time:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Level:    Error,
		Duration: Duration(5 * time.Second),
		Message:  "\"response time out \"",
		Service:  "api",
		Extra: map[string]string{
			"path": "/",
		},
	}
	tests := map[string]struct {
		input    ParseSettings
		expected bool
	}{
		"filter by date that includes the log": {
			input: ParseSettings{
				StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: true,
		},
		"filter by end date that does not include the log": {
			input: ParseSettings{
				StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: false,
		},
		"filter by start date that does not include the log": {
			input: ParseSettings{
				StartDate: time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: false,
		},
		"filter by service that includes the log": {
			input: ParseSettings{
				Service: "api",
			},
			expected: true,
		},
		"filter by date that includes, but service that does not include log": {
			input: ParseSettings{
				StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				Service:   "db",
			},
			expected: false,
		},
		"filter by service that includes, but level that does not include log": {
			input: ParseSettings{
				Service: "api",
				Level:   Warning,
			},
			expected: false,
		},
		"filter by level that includes log": {
			input: ParseSettings{
				Level: Error,
			},
			expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			output := filterLog(testLog, &test.input)
			if output != test.expected {
				t.Fatalf("%s: expected %t - got %t", name, test.expected, output)
			}
		})
	}
}

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
