package utils_test

import (
	"testing"
	"time"

	"github.com/HugoDrl/zebra/internal/analyser/utils"
	"github.com/HugoDrl/zebra/internal/parser"
	"github.com/google/go-cmp/cmp"
)

func feedThenCloseChan[T any](values []T) chan T {
	channel := make(chan T, len(values))
	defer close(channel)
	for _, value := range values {
		channel <- value
	}

	return channel
}

type ExtractLogAndErrToSlicesOutput struct {
	Logs []*parser.Log
	Errs []error
}

func TestExtractLogAndErrChanToSlices(t *testing.T) {
	tests := map[string]struct {
		inputLogs chan *parser.Log
		inputErrs chan error

		expected ExtractLogAndErrToSlicesOutput
	}{
		"two emtpy channels should return empty slices": {
			inputLogs: feedThenCloseChan([]*parser.Log{}),
			inputErrs: feedThenCloseChan([]error{}),
			expected: ExtractLogAndErrToSlicesOutput{
				Logs: []*parser.Log{},
				Errs: []error{},
			},
		},
		"channel with one log and one err should return both in slices": {
			inputLogs: feedThenCloseChan([]*parser.Log{{Time: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}}),
			inputErrs: feedThenCloseChan([]error{&parser.ParseError{Reason: "why not"}}),
			expected: ExtractLogAndErrToSlicesOutput{
				Logs: []*parser.Log{{Time: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}},
				Errs: []error{&parser.ParseError{Reason: "why not"}},
			},
		},
		"channel with only one log should return it in slice": {
			inputLogs: feedThenCloseChan([]*parser.Log{{Time: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}}),
			inputErrs: feedThenCloseChan([]error{}),
			expected: ExtractLogAndErrToSlicesOutput{
				Logs: []*parser.Log{{Time: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}},
				Errs: []error{},
			},
		},
		"channel with only one err should return it in slice": {
			inputLogs: feedThenCloseChan([]*parser.Log{}),
			inputErrs: feedThenCloseChan([]error{&parser.ParseError{Reason: "why not"}}),
			expected: ExtractLogAndErrToSlicesOutput{
				Logs: []*parser.Log{},
				Errs: []error{&parser.ParseError{Reason: "why not"}},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			logSlice, errSlice := utils.ExtractLogAndErrChanToSlices(test.inputLogs, test.inputErrs)
			output := ExtractLogAndErrToSlicesOutput{
				Logs: logSlice,
				Errs: errSlice,
			}

			if diff := cmp.Diff(test.expected, output); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
