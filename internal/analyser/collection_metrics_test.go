package analyser

import (
	"testing"
	"time"

	"github.com/HugoDrl/zebra/internal/parser"
	"github.com/google/go-cmp/cmp"
)

func getDefaultMetric() CollectionMetric {
	return CollectionMetric{
		Lines: map[parser.Level]int{
			parser.Warning: 1,
			parser.Info:    10,
		},
		ServicePerformance: map[string]ServiceMetric{
			"api": {
				Name:            "api",
				Lines:           11,
				AverageDuration: 44 * time.Millisecond,
			},
		},
	}
}

func TestHandleService(t *testing.T) {
	tests := map[string]struct {
		input    *parser.Log
		expected CollectionMetric
	}{
		"empty log should not change anything": {
			input:    &parser.Log{},
			expected: getDefaultMetric(),
		},
		"nil pointer should not change anything": {
			input:    nil,
			expected: getDefaultMetric(),
		},
		"api log should increase the same ServiceMetric": {
			input: &parser.Log{
				Duration: parser.Duration(44 * time.Millisecond),
				Service:  "api",
				Level:    parser.Info,
			},
			expected: func() CollectionMetric {
				metric := getDefaultMetric()
				metric.Lines[parser.Info]++

				// Average duration should not move for this one
				perf := metric.ServicePerformance["api"]
				perf.Lines++
				metric.ServicePerformance["api"] = perf

				return metric
			}(),
		},
		"other service log should create a new ServiceMetric": {
			input: &parser.Log{
				Duration: parser.Duration(100 * time.Millisecond),
				Service:  "database",
				Level:    parser.Info,
			},
			expected: func() CollectionMetric {
				metric := getDefaultMetric()
				metric.Lines[parser.Info]++
				metric.ServicePerformance["database"] = ServiceMetric{
					Name:            "database",
					Lines:           1,
					AverageDuration: 100 * time.Millisecond,
				}
				return metric
			}(),
		},
		"happending a log with a different duration should modify average properly": {
			input: &parser.Log{
				Duration: parser.Duration(100 * time.Millisecond),
				Service:  "api",
				Level:    parser.Error,
			},
			expected: func() CollectionMetric {
				metric := getDefaultMetric()
				metric.Lines[parser.Error]++

				// Average duration should not move for this one
				perf := metric.ServicePerformance["api"]
				perf.Lines++
				perf.AverageDuration = (perf.AverageDuration*time.Duration(perf.Lines-1) + 100*time.Millisecond) / time.Duration(perf.Lines)
				metric.ServicePerformance["api"] = perf

				return metric
			}(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			metric := getDefaultMetric()
			metric.handleService(test.input)

			if diff := cmp.Diff(test.expected, metric); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestSlowestLogsHandler(t *testing.T) {
	type SlowestLogsInput struct {
		logs            []parser.Log
		slowestLogsSize int
	}
	tests := map[string]struct {
		input    SlowestLogsInput
		expected CollectionMetric
	}{
		"adding one log to metrics with size 1 should add it": {
			input: SlowestLogsInput{
				logs: []parser.Log{
					{Duration: parser.Duration(12 * time.Millisecond)},
				},
				slowestLogsSize: 1,
			},
			expected: func() CollectionMetric {
				metric := getDefaultMetric()
				metric.SlowestInput = []*parser.Log{
					{Duration: parser.Duration(12 * time.Millisecond)},
				}
				return metric
			}(),
		},
		"adding one entry to metrics with size 0 should not add it": {
			input: SlowestLogsInput{
				logs: []parser.Log{
					{Duration: parser.Duration(12 * time.Millisecond)},
				},
				slowestLogsSize: 0,
			},
			expected: getDefaultMetric(),
		},
		"adding more entries to metrics with size 1 should add the slowest": {
			input: SlowestLogsInput{
				logs: []parser.Log{
					{Duration: parser.Duration(20 * time.Millisecond)},
					{Duration: parser.Duration(10 * time.Millisecond)},
					{Duration: parser.Duration(12 * time.Millisecond)},
				},
				slowestLogsSize: 1,
			},
			expected: func() CollectionMetric {
				metric := getDefaultMetric()
				metric.SlowestInput = []*parser.Log{
					{Duration: parser.Duration(20 * time.Millisecond)},
				}
				return metric
			}(),
		},
		"adding two entries to metrics with size 5 should add both": {
			input: SlowestLogsInput{
				logs: []parser.Log{
					{Duration: parser.Duration(12 * time.Millisecond)},
					{Duration: parser.Duration(20 * time.Millisecond)},
				},
				slowestLogsSize: 5,
			},
			expected: func() CollectionMetric {
				metric := getDefaultMetric()
				metric.SlowestInput = []*parser.Log{
					{Duration: parser.Duration(12 * time.Millisecond)},
					{Duration: parser.Duration(20 * time.Millisecond)},
				}
				return metric
			}(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			metric := getDefaultMetric()
			for _, log := range test.input.logs {
				metric.handleSlowestLogs(
					test.input.slowestLogsSize,
					&log,
				)
			}

			if diff := cmp.Diff(test.expected, metric); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestValidateLog(t *testing.T) {
	tests := map[string]struct {
		input    parser.Log
		expected bool
	}{
		"valid log should just be valid": {
			input: parser.Log{
				Time:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
				Service:  "api",
				Duration: parser.Duration(50 * time.Millisecond),
			},
			expected: true,
		},
		"log with emtpy service should not be valid": {
			input: parser.Log{
				Time:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
				Service:  "",
				Duration: parser.Duration(50 * time.Millisecond),
			},
			expected: false,
		},
		"log with negative duration should not be valid": {
			input: parser.Log{
				Time:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
				Service:  "api",
				Duration: parser.Duration(-50 * time.Millisecond),
			},
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if diff := cmp.Diff(test.expected, validateLog(test.input)); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
