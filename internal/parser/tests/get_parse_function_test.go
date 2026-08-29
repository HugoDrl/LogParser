package parser_test

import (
	"testing"

	"github.com/HugoDrl/zebra/internal/parser"
	"github.com/google/go-cmp/cmp"
)

func TestGetParseFunction(t *testing.T) {
	tests := map[string]struct {
		input    parser.ParseSettings
		expected parser.ParseFunction
	}{
		"default input should return default parsing function": {
			input:    parser.ParseSettings{},
			expected: parser.ParseDefaultFormatLine,
		},
		"json files with no specification on parsing should return default function": {
			input: parser.ParseSettings{
				Files: []string{"test.json", "parse.json"},
			},
			expected: parser.ParseDefaultFormatLine,
		},
		"explicit false on json field should return parsing function": {
			input: parser.ParseSettings{
				Json: false,
			},
			expected: parser.ParseDefaultFormatLine,
		},
		"json false with json files should return parsing function": {
			input: parser.ParseSettings{
				Json:  false,
				Files: []string{"test.json", "parse.json"},
			},
			expected: parser.ParseDefaultFormatLine,
		},
		"explicit true on json field should return json function": {
			input: parser.ParseSettings{
				Json: true,
			},
			expected: parser.ParseJSONFormatLine,
		},
		"json true with .log files should return json function": {
			input: parser.ParseSettings{
				Json:  true,
				Files: []string{"test.log", "parse.log"},
			},
			expected: parser.ParseJSONFormatLine,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			out := parser.GetParseFunction(test.input)
			if cmp.Equal(test.expected, out) {
				t.Fatalf("expected %v - got %v", test.expected, out)
			}
		})
	}
}
