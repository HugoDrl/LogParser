package filter

import (
	"testing"
	"time"

	"github.com/HugoDrl/zebra/internal/parser"
	"github.com/google/go-cmp/cmp"
)

func TestLogFilter(t *testing.T) {
	testLog := &parser.Log{
		Time:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Level:    parser.Error,
		Duration: parser.Duration(5 * time.Second),
		Message:  "\"response time out \"",
		Service:  "api",
		Extra: map[string]string{
			"path": "/",
		},
	}
	tests := map[string]struct {
		input    Filters
		expected bool
	}{
		"filter by date that includes the log": {
			input: Filters{
				StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: true,
		},
		"filter by end date that does not include the log": {
			input: Filters{
				StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: false,
		},
		"filter by start date that does not include the log": {
			input: Filters{
				StartDate: time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: false,
		},
		"filter by service that includes the log": {
			input: Filters{
				Service: "api",
			},
			expected: true,
		},
		"filter by date that includes, but service that does not include log": {
			input: Filters{
				StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				Service:   "db",
			},
			expected: false,
		},
		"filter by service that includes, but level that does not include log": {
			input: Filters{
				Service: "api",
				Level:   parser.Warning,
			},
			expected: false,
		},
		"filter by level that includes log": {
			input: Filters{
				Level: parser.Error,
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

type testProcessFilterInput struct {
	logs    []*parser.Log
	filters *Filters
}

func TestProcessFilter(t *testing.T) {
	tests := map[string]struct {
		input    testProcessFilterInput
		expected []*parser.Log
	}{
		"no logs and no filters should return nothing": {
			input: testProcessFilterInput{
				logs:    []*parser.Log{},
				filters: &Filters{},
			},
			expected: []*parser.Log{},
		},
		"some logs with no filters should return it": {
			input: testProcessFilterInput{
				logs: []*parser.Log{
					{
						Time:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
						Level:    parser.Warning,
						Duration: parser.Duration(20 * time.Millisecond),
					},
				},
				filters: &Filters{},
			},
			expected: []*parser.Log{
				{
					Time:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
					Level:    parser.Warning,
					Duration: parser.Duration(20 * time.Millisecond),
				},
			},
		},
		"logs under start date filter should not be retrieved": {
			input: testProcessFilterInput{
				logs: []*parser.Log{
					{
						Time: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						Time: time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
				filters: &Filters{
					StartDate: time.Date(2026, 2, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			expected: []*parser.Log{
				{
					Time: time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		"log with exactly same date should be retrieved": {
			input: testProcessFilterInput{
				logs: []*parser.Log{
					{
						Time: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						Time: time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
				filters: &Filters{
					StartDate: time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			expected: []*parser.Log{
				{
					Time: time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		"filtering by service should get concerned log": {
			input: testProcessFilterInput{
				logs: []*parser.Log{
					{
						Service: "WARNING",
					},
					{
						Service: "ERROR",
					},
				},
				filters: &Filters{
					Service: "ERROR",
				},
			},
			expected: []*parser.Log{
				{
					Service: "ERROR",
				},
			},
		},
		"filtering by service and date should get only logs concerned": {
			input: testProcessFilterInput{
				logs: []*parser.Log{
					{
						Time:    time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC),
						Service: "WARNING",
					},
					{
						Time:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
						Service: "ERROR",
					},
				},
				filters: &Filters{
					Service: "WARNING",
					EndDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			expected: []*parser.Log{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			inputChannel := make(chan *parser.Log, len(test.input.logs))
			for _, log := range test.input.logs {
				inputChannel <- log
			}
			close(inputChannel)

			outputChan := ProcessFilter(inputChannel, test.input.filters)

			outputLogs := make([]*parser.Log, 0)
			for log := range outputChan {
				outputLogs = append(outputLogs, log)
			}

			if diff := cmp.Diff(test.expected, outputLogs); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
