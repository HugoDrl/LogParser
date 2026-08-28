package parser_test

import (
	"testing"
	"time"

	"github.com/HugoDrl/zebra/internal/parser"
	"github.com/google/go-cmp/cmp"
)

func TestParseLineLog(t *testing.T) {
	type ParseLineOutput struct {
		Log parser.Log
		Err error
	}
	tests := map[string]struct {
		input    string
		expected ParseLineOutput
	}{
		"standard log": {
			input: "2026-01-01T00:00:00Z WARNING service=test message=hey duration=50ms",
			expected: ParseLineOutput{
				Log: parser.Log{
					Time:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
					Level:    parser.Warning,
					Service:  "test",
					Message:  "hey",
					Duration: parser.Duration(50 * time.Millisecond),
					Extra:    map[string]string{},
				},
			},
		},
		"invalid: invalid service": {
			input: "2026-01-01T00:00:00Z INVALID_SERVICE service=test message=hey duration=50ms",
			expected: ParseLineOutput{
				Err: &parser.ValueError{
					ExpectedValue: "",
					ErroredValue:  "level",
				},
			},
		},
		"invalid: not enough mandatory arguments": {
			input: "2026-01-01T00:00:00Z",
			expected: ParseLineOutput{
				Err: &parser.ParseError{Reason: "Not enough arguments - expected 2 - found 1"},
			},
		},
		"invalid: invalid date": {
			input: "INVALID_DATE WARNING service=test message=hey duration=50ms",
			expected: ParseLineOutput{
				Err: &parser.ValueError{
					ExpectedValue: "time format - RFC3339",
					ErroredValue:  "INVALID_DATE",
				},
			},
		},
		"invalid: invalid duration in seconds instead of ms": {
			// duration should (for now) only be in milliseconds
			input: "2026-01-01T00:00:00Z WARNING service=test message=hey duration=50s",
			expected: ParseLineOutput{
				Err: &parser.ValueError{
					ExpectedValue: "duration value in milliseconds on format xxxms",
					ErroredValue:  "50s",
				},
			},
		},
		"invalid: invalid duration is not a number": {
			input: "2026-01-01T00:00:00Z WARNING service=test message=hey duration=NaNms",
			expected: ParseLineOutput{
				Err: &parser.ValueError{
					ExpectedValue: "duration value in milliseconds on format xxxms",
					ErroredValue:  "NaNms",
				},
			},
		},
		"invalid: missing field": {
			input: "2026-01-01T00:00:00Z WARNING service=test duration=50ms",
			expected: ParseLineOutput{
				Err: &parser.ValueError{
					ExpectedValue: "message field, service field, and duration field, should not be empty.",
					ErroredValue:  "message: , service: test, duration: 50000000",
				},
			},
		},
		"wrong format key=value pairs": {
			input: "2026-01-01T00:00:00Z WARNING service=test=api duration=50ms",
			expected: ParseLineOutput{
				Err: &parser.ParseError{Reason: "Wrong format in key-values design (found [service test api])"},
			},
		},
		"log with sentence as message": {
			input: "2026-01-01T00:00:00Z WARNING service=test message=\"this log is very important\" duration=50ms",
			expected: ParseLineOutput{
				Log: parser.Log{
					Time:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
					Level:    parser.Warning,
					Service:  "test",
					Message:  "\"this log is very important\"",
					Duration: parser.Duration(50 * time.Millisecond),
					Extra:    map[string]string{},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			log, err := parser.ParseDefaultFormatLine(test.input)
			output := ParseLineOutput{
				Log: log,
				Err: err,
			}

			if diff := cmp.Diff(test.expected, output); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestParseJSONLine(t *testing.T) {
	type ParseLineOutput struct {
		Log parser.Log
		Err error
	}
	tests := map[string]struct {
		input    string
		expected ParseLineOutput
	}{
		"empty log should just return a json parsing error": {
			input: "",
			expected: ParseLineOutput{
				Err: &parser.ParseError{
					Reason: "unexpected end of JSON input",
				},
			},
		},
		"empty json log should return validation error": {
			input: "{}",
			expected: ParseLineOutput{
				Err: &parser.ValueError{
					ErroredValue: "level",
				},
			},
		},
		"valid json log with no extra fields should be parsed just fine": {
			input: `{"date": "2026-01-01T00:00:00Z", "level": "WARNING", "service": "database", "message": "drop db", "duration": "10ms"}`,
			expected: ParseLineOutput{
				Log: parser.Log{
					Time:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
					Level:    parser.Warning,
					Service:  "database",
					Message:  "drop db",
					Duration: parser.Duration(10 * time.Millisecond),
				},
			},
		},
		"valid json log with extra fields should put extra fields in 'extra' map": {
			input: `{"date": "2026-01-01T00:00:00Z", "level": "WARNING", "service": "database", "message": "drop db", "duration": "10ms", "this is an extra field": "hey hello"}`,
			expected: ParseLineOutput{
				Log: parser.Log{
					Time:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
					Level:    parser.Warning,
					Service:  "database",
					Message:  "drop db",
					Duration: parser.Duration(10 * time.Millisecond),
					Extra: map[string]string{
						"this is an extra field": "hey hello",
					},
				},
			},
		},
		"invalid level json log should return an error": {
			input: `{"date": "2026-01-01T00:00:00Z", "level": "HEY", "service": "database", "message": "drop db", "duration": "10ms"}`,
			expected: ParseLineOutput{
				Err: &parser.ValueError{
					ErroredValue: "level",
				},
			},
		},
		"invalid duration field for json log should raise an error": {
			input: `{"date": "2026-01-01T00:00:00Z", "level": "WARNING", "service": "database", "message": "drop db", "duration": "hey"}`,
			expected: ParseLineOutput{
				Err: &parser.ParseError{
					Reason: `time: invalid duration "hey"`,
				},
			},
		},
		"non-json format should return an error": {
			input: `{"date": "2026-01-01T00:00:00Z", "level": "WARNING", "service": "database", "message": "drop db", "duration": "10ms"`,
			expected: ParseLineOutput{
				Err: &parser.ParseError{
					Reason: "unexpected end of JSON input",
				},
			},
		},
		"invalid date field for json should return a validation error": {
			input: `{"date": "invalid date i guess", "level": "WARNING", "service": "database", "message": "drop db", "duration": "10ms"}`,
			expected: ParseLineOutput{
				Err: &parser.ParseError{
					Reason: `parsing time "invalid date i guess" as "2006-01-02T15:04:05Z07:00": cannot parse "invalid date i guess" as "2006"`,
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			log, err := parser.ParseJSONFormatLine(test.input)
			output := ParseLineOutput{
				Log: log,
				Err: err,
			}

			if diff := cmp.Diff(test.expected, output); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
