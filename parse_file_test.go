package main

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/HugoDrl/zebra/internal/parser"
	"github.com/google/go-cmp/cmp"
)

func emptyBothChannels(logChan <-chan *parser.Log, errsChan <-chan error) ([]*parser.Log, []error) {
	logs := make([]*parser.Log, 0)
	errs := make([]error, 0)
	c := 2
	for {
		select {
		case log, ok := <-logChan:
			if !ok {
				c--
				if c == 0 {
					return logs, errs
				}
				continue
			}
			logs = append(logs, log)
		case err, ok := <-errsChan:
			if !ok {
				c--
				if c == 0 {
					return logs, errs
				}
				continue
			}
			errs = append(errs, err)
		}
	}
}

func prepareFilesForTests(fileName string, content []byte) error {
	root, err := os.OpenRoot(".")
	if err != nil {
		return err
	}
	defer root.Close()

	// Allow subdirs in file name
	root.MkdirAll(filepath.Dir(fileName), 0o777)
	if err := os.WriteFile(fileName, content, 0o777); err != nil {
		return err
	}

	return nil
}

func TestParseLogsFromFile(t *testing.T) {
	type expectedStruct struct {
		Logs []*parser.Log
		Errs []error
	}
	tests := map[string]struct {
		inputFileContent   []string
		inputParseSettings parser.ParseSettings
		expected           expectedStruct
	}{
		"empty file should return nothing": {
			inputFileContent: []string{""},
			inputParseSettings: parser.ParseSettings{
				Files: []string{"temp_parse_logs_from_file.txt"},
			},
			expected: expectedStruct{
				Logs: []*parser.Log{},
				Errs: []error{},
			},
		},
		"file with a valid log should parse it": {
			inputFileContent: []string{
				"2026-01-01T00:00:00Z WARNING service=database message=\"hello from database\" duration=10ms",
			},
			inputParseSettings: parser.ParseSettings{
				Files: []string{"temp_parse_logs_from_file.txt"},
			},
			expected: expectedStruct{
				Logs: []*parser.Log{
					{
						Time:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
						Level:    parser.Warning,
						Service:  "database",
						Message:  `"hello from database"`,
						Duration: parser.Duration(10 * time.Millisecond),
						Extra:    map[string]string{},
					},
				},
				Errs: []error{},
			},
		},
		"several files with several valid logs should parse them all": {
			inputFileContent: []string{
				"2026-01-01T00:00:00Z WARNING service=database message=\"hello from database\" duration=10ms",
				"2026-01-01T00:00:10Z INFO service=database message=\"hello again from database\" duration=12ms",
			},
			inputParseSettings: parser.ParseSettings{
				Files: []string{
					"temp_parse_logs_from_file1.txt",
					"temp_parse_logs_from_file2.txt",
				},
			},
			expected: expectedStruct{
				Logs: []*parser.Log{
					{
						Time:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
						Level:    parser.Warning,
						Service:  "database",
						Message:  `"hello from database"`,
						Duration: parser.Duration(10 * time.Millisecond),
						Extra:    map[string]string{},
					},
					{
						Time:     time.Date(2026, 1, 1, 0, 0, 10, 0, time.UTC),
						Level:    parser.Info,
						Service:  "database",
						Message:  `"hello again from database"`,
						Duration: parser.Duration(12 * time.Millisecond),
						Extra:    map[string]string{},
					},
				},
				Errs: []error{},
			},
		},
		"several files, but one has invalid content should return both logs and errors": {
			inputFileContent: []string{
				"2026-01-01T00:00:00Z WARNING service=database message=\"hello from database\" duration=10ms",
				"invalid content for logs",
			},
			inputParseSettings: parser.ParseSettings{
				Files: []string{
					"temp_parse_logs_from_file1.txt",
					"temp_parse_logs_from_file2.txt",
				},
			},
			expected: expectedStruct{
				Logs: []*parser.Log{
					{
						Time:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
						Level:    parser.Warning,
						Service:  "database",
						Message:  `"hello from database"`,
						Duration: parser.Duration(10 * time.Millisecond),
						Extra:    map[string]string{},
					},
				},
				Errs: []error{&parser.ParseError{
					Line: 1,
					Err: &parser.ValueError{
						ErroredValue:  "invalid",
						ExpectedValue: "time format - RFC3339",
					},
				}},
			},
		},
		"one file with several logs should parse them all": {
			inputFileContent: []string{
				`2026-01-01T00:00:00Z WARNING service=database message="hello from database" duration=10ms
2026-01-01T00:00:10Z INFO service=database message="hello again from database" duration=12ms`,
			},
			inputParseSettings: parser.ParseSettings{
				Files: []string{
					"temp_parse_logs_from_file_parse_all.txt",
				},
			},
			expected: expectedStruct{
				Logs: []*parser.Log{
					{
						Time:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
						Level:    parser.Warning,
						Service:  "database",
						Message:  `"hello from database"`,
						Duration: parser.Duration(10 * time.Millisecond),
						Extra:    map[string]string{},
					},
					{
						Time:     time.Date(2026, 1, 1, 0, 0, 10, 0, time.UTC),
						Level:    parser.Info,
						Service:  "database",
						Message:  `"hello again from database"`,
						Duration: parser.Duration(12 * time.Millisecond),
						Extra:    map[string]string{},
					},
				},
				Errs: []error{},
			},
		},
		"file in dir should be read correctly": {
			inputFileContent: []string{
				"2026-01-01T00:00:00Z WARNING service=database message=\"hello from database\" duration=10ms",
			},
			inputParseSettings: parser.ParseSettings{
				Files: []string{
					"temp-nested-dir/temp.txt",
				},
			},
			expected: expectedStruct{
				Logs: []*parser.Log{
					{
						Time:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
						Level:    parser.Warning,
						Service:  "database",
						Message:  `"hello from database"`,
						Duration: parser.Duration(10 * time.Millisecond),
						Extra:    map[string]string{},
					},
				},
				Errs: []error{},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			for i, fileName := range test.inputParseSettings.Files {
				fileContent := []byte(test.inputFileContent[i])
				if err := prepareFilesForTests(fileName, fileContent); err != nil {
					t.Fatal(err)
				}
				if dir := filepath.Dir(fileName); dir != "" {
					defer os.RemoveAll(dir)
				} else {
					defer os.Remove(fileName)
				}
			}

			logChan, errsChan := processFiles(&test.inputParseSettings)

			logs, errs := emptyBothChannels(logChan, errsChan)
			// Clean files
			for _, file := range test.inputParseSettings.Files {
				os.Remove(file)
			}

			// Sort logs by time to be predictible for comparison
			sort.Slice(logs, func(i, j int) bool {
				return logs[i].Time.Compare(logs[j].Time) < 0
			})

			output := expectedStruct{
				Logs: logs,
				Errs: errs,
			}

			if diff := cmp.Diff(test.expected, output); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
