package parser_test

import (
	"errors"
	"os"
	"testing"

	"github.com/HugoDrl/zebra/internal/parser"
	"github.com/google/go-cmp/cmp"
)

func TestFileError(t *testing.T) {
	tests := map[string]struct {
		input    parser.FileError
		expected string
	}{
		"nil error should not work i guess": {
			input:    parser.FileError{},
			expected: "Error encountered on file  - %!s(<nil>)",
		},
		"valid error with complete informations should return the error": {
			input: parser.FileError{
				File: "test.txt",
				Err:  errors.New("this is an error"),
			},
			expected: "Error encountered on file test.txt - this is an error",
		},
		"valid error with os Path error": {
			input: parser.FileError{
				File: "test.txt",
				Err: &os.PathError{
					Op:   "READ",
					Path: "./test.txt",
					Err:  os.ErrNotExist,
				},
			},
			expected: "Error encountered on file test.txt - READ ./test.txt: file does not exist",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			output := test.input.Error()

			if diff := cmp.Diff(test.expected, output); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestParseError(t *testing.T) {
	tests := map[string]struct {
		input    parser.ParseError
		expected string
	}{
		"empty error should return a weird string": {
			input:    parser.ParseError{},
			expected: "Error encountered while parsing file { <nil>} on line 0 - ",
		},
		// TODO: correct this behavior as this is not normal
		// Let's by the way rethink about this parseError as it is not a struct I like
		"normal parse error should format the string properly": {
			input: parser.ParseError{
				File: parser.FileError{
					File: "./test.txt",
					Err:  os.ErrNotExist,
				},
				Line:   10,
				Reason: os.ErrNotExist.Error(),
				Err:    os.ErrNotExist,
			},
			expected: "Error encountered while parsing file {./test.txt file does not exist} on line 10 - file does not exist",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			output := test.input.Error()
			if diff := cmp.Diff(test.expected, output); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestValueError(t *testing.T) {
	tests := map[string]struct {
		input    parser.ValueError
		expected string
	}{
		"empty value error should return a weird string": {
			input:    parser.ValueError{},
			expected: "{ <nil>} : Wrong value on line 0: expected value :  (got )",
		},
		"normal value error should return a formatted string": {
			input: parser.ValueError{
				File: parser.FileError{
					File: "./test.txt",
					Err:  os.ErrNotExist,
				},
				Line:          10,
				ErroredValue:  "time",
				ExpectedValue: "hey",
				Err:           os.ErrNotExist,
			},
			expected: "{./test.txt file does not exist} : Wrong value on line 10: expected value : hey (got time)",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if diff := cmp.Diff(test.expected, test.input.Error()); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
