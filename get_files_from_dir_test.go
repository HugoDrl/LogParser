package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testDir struct {
	dirName   string
	filenames []string
}

func prepareTestDirsSetup(dir testDir) error {
	root, err := os.OpenRoot(".")
	if err != nil {
		return err
	}

	// This is just for the case of a empty files dir, as it would have been created in the loop
	if err := root.Mkdir(dir.dirName, 0o777); err != nil {
		return err
	}

	for _, file := range dir.filenames {
		completeName := dir.dirName + "/" + file
		if err := root.MkdirAll(filepath.Dir(completeName), 0o777); err != nil {
			return err
		}
		if _, err := root.Create(completeName); err != nil {
			return err
		}
	}
	return nil
}

func TestFetchLogFilesFromDir(t *testing.T) {
	// I chose to not test whether or not it was possible to get log files from current directory, just because this is too much impredictable.
	// I can not be sure user will not have created some logs to try the program, and get a test fail.
	// Maybe I will add it later, because I know CI will not have tests and garbage log files, so it would pass, but for dev comfort I am not sure.
	type testFetchFilesOutput struct {
		Files []string
		Err   error
	}
	tests := map[string]struct {
		dir      testDir
		expected testFetchFilesOutput
	}{
		"emtpy dir should return nothing": {
			dir: testDir{
				dirName: "temp-empty-dir",
			},
			expected: testFetchFilesOutput{
				Files: nil,
				Err:   nil,
			},
		},
		"dir with one log file should return it": {
			dir: testDir{
				dirName:   "temp-one-log",
				filenames: []string{"temp.log"},
			},
			expected: testFetchFilesOutput{
				Files: []string{"temp-one-log/temp.log"},
				Err:   nil,
			},
		},
		"dir with a single file that is not log should return nothing": {
			dir: testDir{
				dirName:   "temp-no-log",
				filenames: []string{"temp.txt"},
			},
			expected: testFetchFilesOutput{
				Files: nil,
				Err:   nil,
			},
		},
		"dir with several log files should return them": {
			dir: testDir{
				dirName:   "temp-multiple-log",
				filenames: []string{"temp.log", "temp2.log"},
			},
			expected: testFetchFilesOutput{
				Files: []string{"temp-multiple-log/temp.log", "temp-multiple-log/temp2.log"},
				Err:   nil,
			},
		},
		"dir with one log and several others should return only log": {
			dir: testDir{
				dirName:   "temp-one-log-and-others",
				filenames: []string{"temp.log", "temp.txt", "temp2.txt"},
			},
			expected: testFetchFilesOutput{
				Files: []string{"temp-one-log-and-others/temp.log"},
				Err:   nil,
			},
		},
		"dir with log in subdir should return it": {
			dir: testDir{
				dirName:   "temp-subdir-log",
				filenames: []string{"temp-dir/temp.log"},
			},
			expected: testFetchFilesOutput{
				Files: []string{"temp-subdir-log/temp-dir/temp.log"},
				Err:   nil,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := prepareTestDirsSetup(test.dir)
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(test.dir.dirName)

			files, err := getLogFilesFromDir(test.dir.dirName)
			output := testFetchFilesOutput{
				Files: files,
				Err:   err,
			}
			if diff := cmp.Diff(test.expected, output); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
