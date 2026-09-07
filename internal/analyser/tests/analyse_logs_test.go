package analyser_test

import (
	"testing"
	"time"

	"github.com/HugoDrl/zebra/internal/analyser"
	"github.com/HugoDrl/zebra/internal/analyser/utils"
	"github.com/HugoDrl/zebra/internal/parser"
	"github.com/google/go-cmp/cmp"
)

func feedThenCloseChan[T any](values ...T) chan T {
	channel := make(chan T, len(values))
	for _, value := range values {
		channel <- value
	}
	close(channel)
	return channel
}

func TestLogAnalyze(t *testing.T) {
	type LogAnalyzeInput struct {
		logChan  chan *parser.Log
		errChan  chan error
		settings *analyser.AnalyserSettings
	}

	tests := map[string]struct {
		input    LogAnalyzeInput
		expected *analyser.CollectionMetric
	}{
		"empty everything should result in empty metrics": {
			input: LogAnalyzeInput{
				logChan: feedThenCloseChan([]*parser.Log{}...),
				errChan: feedThenCloseChan([]error{}...),
			},
			expected: &analyser.CollectionMetric{
				Lines:              map[parser.Level]int{},
				ServicePerformance: map[string]analyser.ServiceMetric{},
			},
		},
		"one file error should append to error slice": {
			input: LogAnalyzeInput{
				logChan: feedThenCloseChan([]*parser.Log{}...),
				errChan: feedThenCloseChan([]error{
					&parser.FileError{
						File: "test.txt",
						Err:  nil,
					},
				}...),
			},
			expected: &analyser.CollectionMetric{
				Lines:              map[parser.Level]int{},
				ServicePerformance: map[string]analyser.ServiceMetric{},
				FileErrors: []*parser.FileError{
					{File: "test.txt", Err: nil},
				},
			},
		},
		"any other error than a file error should increase log parsing error": {
			input: LogAnalyzeInput{
				logChan: feedThenCloseChan([]*parser.Log{}...),
				errChan: feedThenCloseChan([]error{
					&parser.ValueError{},
					&parser.ParseError{},
				}...),
			},
			expected: &analyser.CollectionMetric{
				Lines:              map[parser.Level]int{},
				ServicePerformance: map[string]analyser.ServiceMetric{},
				ParsingErrorCount:  2,
			},
		},
		"adding a log should increase metrics stats": {
			input: LogAnalyzeInput{
				logChan: feedThenCloseChan([]*parser.Log{
					{
						Time:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
						Service:  "database",
						Level:    parser.Warning,
						Duration: parser.Duration(10 * time.Millisecond),
					},
				}...),
				errChan:  feedThenCloseChan([]error{}...),
				settings: &analyser.AnalyserSettings{},
			},
			expected: &analyser.CollectionMetric{
				Lines: map[parser.Level]int{
					parser.Warning: 1,
				},
				ServicePerformance: map[string]analyser.ServiceMetric{
					"database": {
						Name:            "database",
						Lines:           1,
						AverageDuration: 10 * time.Millisecond,
					},
				},
			},
		},
		"adding a log with slowest logs gestion should increase metrics stats and slowest logs": {
			input: LogAnalyzeInput{
				logChan: feedThenCloseChan([]*parser.Log{
					{
						Time:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
						Service:  "database",
						Level:    parser.Warning,
						Duration: parser.Duration(10 * time.Millisecond),
					},
				}...),
				errChan: feedThenCloseChan([]error{}...),
				settings: &analyser.AnalyserSettings{
					SlowestLogsToRetrieve: 1,
				},
			},
			expected: &analyser.CollectionMetric{
				Lines: map[parser.Level]int{
					parser.Warning: 1,
				},
				ServicePerformance: map[string]analyser.ServiceMetric{
					"database": {
						Name:            "database",
						Lines:           1,
						AverageDuration: 10 * time.Millisecond,
					},
				},
				SlowestInput: []*parser.Log{
					{
						Time:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
						Service:  "database",
						Level:    parser.Warning,
						Duration: parser.Duration(10 * time.Millisecond),
					},
				},
			},
		},
		// TODO: make sure this behavior is ok ?
		// Should this function handle validation for log ?
		// Else, correct it
		"adding an invalid log should not handle it correctly i guess": {
			input: LogAnalyzeInput{
				logChan: feedThenCloseChan([]*parser.Log{
					{
						Service:  "database",
						Level:    "INVALID",
						Duration: parser.Duration(10 * time.Millisecond),
					},
				}...),
				errChan:  feedThenCloseChan([]error{}...),
				settings: &analyser.AnalyserSettings{},
			},
			expected: &analyser.CollectionMetric{
				Lines: map[parser.Level]int{
					"INVALID": 1,
				},
				ServicePerformance: map[string]analyser.ServiceMetric{
					"database": {
						Name:            "database",
						Lines:           1,
						AverageDuration: 10 * time.Millisecond,
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			logs, errs := utils.ExtractLogAndErrChanToSlices(test.input.logChan, test.input.errChan)
			output := analyser.AnalyseLogs(logs, errs, test.input.settings)
			if diff := cmp.Diff(
				test.expected,
				output,
			); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
