package parser

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetConsumedJSONFields(t *testing.T) {
	tests := map[string]struct {
		input    reflect.Type
		expected map[string]struct{}
	}{
		"type with all fields consumed by json should return nil": {
			input: reflect.TypeFor[struct {
				Test1 int `json:"test1"`
				Test2 int `json:"test2"`
			}](),
			expected: map[string]struct{}{
				"test1": {},
				"test2": {},
			},
		},
		"no consumed field should return an empty map": {
			input: reflect.TypeFor[struct {
				Test1 int
				Test2 int
			}](),
			expected: map[string]struct{}{},
		},
		"partially consumed struct return only consumed fields": {
			input: reflect.TypeFor[struct {
				Test1 int `json:"test1"`
				Test2 int
			}](),
			expected: map[string]struct{}{
				"test1": {},
			},
		},
		"field with different consumed name than field name should return the consumed name": {
			input: reflect.TypeFor[struct {
				Test int `json:"test1"`
			}](),
			expected: map[string]struct{}{
				"test1": {},
			},
		},
		"unexported field with consumed name should return it properly": {
			input: reflect.TypeFor[struct {
				test1 int `json:"test1"`
			}](),
			expected: map[string]struct{}{
				"test1": {},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			output := getConsumedJSONFields(test.input)
			if diff := cmp.Diff(test.expected, output); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestParseJSONExtraArgument(t *testing.T) {
	type ParseJSONInput struct {
		T    reflect.Type
		Json string
	}
	type ParseJSONOutput struct {
		Expected    map[string]string
		ExpectedErr error
	}
	tests := map[string]struct {
		input    ParseJSONInput
		expected ParseJSONOutput
	}{
		"standard log with no consumed field should return all fields": {
			input: ParseJSONInput{
				T: reflect.TypeFor[struct {
					Test int
				}](),
				Json: `{"test": "this is a test"}`,
			},
			expected: ParseJSONOutput{
				Expected: map[string]string{
					"test": "this is a test",
				},
				ExpectedErr: nil,
			},
		},
		"standard log with only consumed field should return no field": {
			input: ParseJSONInput{
				T: reflect.TypeFor[struct {
					Test int `json:"test"`
				}](),
				Json: `{"test": "this is a test"}`,
			},
			expected: ParseJSONOutput{
				Expected:    nil,
				ExpectedErr: nil,
			},
		},
		"invalid json log should return an error": {
			input: ParseJSONInput{
				T: reflect.TypeFor[struct {
					Test int `json:"test"`
				}](),
				Json: `{"test": "this is a test"`,
			},
			expected: ParseJSONOutput{
				Expected: nil,
				ExpectedErr: &ParseError{
					Reason: "unexpected end of JSON input",
				},
			},
		},
		"partially consumed fields should return all extras (not consumed)": {
			input: ParseJSONInput{
				T: reflect.TypeFor[struct {
					Test  int `json:"test"`
					Test2 int `json:"test2"`
					Test3 int
				}](),
				Json: `{"test": "this is a test", "test2": "new test", "test3": "again test"}`,
			},
			expected: ParseJSONOutput{
				Expected: map[string]string{
					"test3": "again test",
				},
				ExpectedErr: nil,
			},
		},
		"struct fields that are not in the json should not be a part of extra fields": {
			input: ParseJSONInput{
				T: reflect.TypeFor[struct {
					Test  int `json:"test"`
					Test2 int // if present, test2 is going to be an extra field
				}](),
				Json: `{"test": "hey hello"}`,
			},
			expected: ParseJSONOutput{
				Expected:    nil,
				ExpectedErr: nil,
			},
		},
		"fields in json that are not in struct should be an extra field": {
			input: ParseJSONInput{
				T: reflect.TypeFor[struct {
					Test  int `json:"test"`
					Test2 int `json:"test2"`
				}](),
				Json: `{"test": "hey hello", "hello": "how are you", "test2": "test again"}`,
			},
			expected: ParseJSONOutput{
				Expected: map[string]string{
					"hello": "how are you",
				},
				ExpectedErr: nil,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			args, err := parseJSONExtraArguments(test.input.T, test.input.Json)
			output := ParseJSONOutput{
				Expected:    args,
				ExpectedErr: err,
			}
			if diff := cmp.Diff(test.expected, output); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
