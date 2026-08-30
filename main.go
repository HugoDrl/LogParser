package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/HugoDrl/zebra/internal/analyser"
	"github.com/HugoDrl/zebra/internal/filter"
	"github.com/HugoDrl/zebra/internal/parser"
)

type extractLinesFromFileInput struct {
	reader        *bufio.Reader
	parseFunction parser.ParseFunction
	logsChan      chan<- *parser.Log
	errsChan      chan<- error
}

func extractLinesFromFile(input extractLinesFromFileInput) {
	scanner := bufio.NewScanner(input.reader)

	lineNo := 0
	for {
		if ok := scanner.Scan(); !ok {
			if err := scanner.Err(); err != nil {
				input.errsChan <- err
			}
			return
		}
		lineNo++
		contentLine := scanner.Text()

		log, err := input.parseFunction(contentLine)
		if err != nil {
			input.errsChan <- &parser.ParseError{
				Line: lineNo,
				Err:  err,
			}
		} else {
			input.logsChan <- &log
		}
	}
}

func ProcessFiles(
	settings *parser.ParseSettings,
) (<-chan *parser.Log, <-chan error) {
	logsChan := make(chan *parser.Log)
	errsChan := make(chan error)
	go func() {
		defer close(logsChan)
		defer close(errsChan)
		var wg sync.WaitGroup
		root, err := os.OpenRoot(".")
		if err != nil {
			errsChan <- err
			return
		}
		defer root.Close()
		defer wg.Wait()
		for _, filepath := range settings.Files {
			wg.Go(func() {
				reader, err := root.OpenFile(filepath, os.O_RDONLY, 0o000)
				if err != nil {
					errsChan <- &parser.FileError{
						File: filepath,
						Err:  err,
					}
					return
				}
				defer reader.Close()
				r := bufio.NewReader(reader)

				extractLinesFromFile(extractLinesFromFileInput{
					reader:        r,
					parseFunction: parser.GetParseFunction(*settings),
					logsChan:      logsChan,
					errsChan:      errsChan,
				})
			})
		}
	}()
	return logsChan, errsChan
}

func getLogFilesFromDir(dirName string) ([]string, error) {
	var logFiles []string
	err := filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if _, ok := strings.CutSuffix(path, ".log"); ok {
			logFiles = append(logFiles, path)
		}
		return nil
	})
	return logFiles, err
}

func initSettings() (*parser.ParseSettings, *filter.Filters, *analyser.AnalyserSettings, error) {
	files := flag.String("files", "", "log files to analyse")
	dirs := flag.String("dirs", "", "dirs containing log files")
	json := flag.Bool("json", false, "wether or not format to parse is json")
	startDate := flag.String("start", "", "logs date to start from")
	endDate := flag.String("end", "", "logs date to end to")
	service := flag.String("service", "", "filter logs by service")
	level := flag.String("level", "", "filter logs by level")
	slowestLogs := flag.Int("top", 0, "number of slowest logs to show")
	flag.Parse()

	if *files == "" && *dirs == "" {
		return nil, nil, nil, errors.New("Please specify file(s) separated by a comma using --files or --dirs flag")
	}

	var processedStartDate time.Time
	var processedEndDate time.Time
	var processErr error
	if *startDate != "" {
		processedStartDate, processErr = time.Parse(time.RFC3339, *startDate)
		if processErr != nil {
			return nil, nil, nil, errors.New("Wrong format for starting date - excpected RFC3339")
		}
	}
	if *endDate != "" {
		processedEndDate, processErr = time.Parse(time.RFC3339, *endDate)
		if processErr != nil {
			return nil, nil, nil, errors.New("Wrong format for starting date - excpected RFC3339")
		}
	}

	logFiles := strings.Split(*files, ",")
	if *dirs != "" {
		for dir := range strings.SplitSeq(*dirs, ",") {
			foundFiles, err := getLogFilesFromDir(dir)
			if err != nil {
				return nil, nil, nil, err
			}
			logFiles = append(logFiles, foundFiles...)
		}
	}

	parsingSettings := parser.ParseSettings{
		Files: logFiles,
		Json:  *json,
	}
	filters := filter.Filters{
		StartDate: processedStartDate,
		EndDate:   processedEndDate,
		Level:     parser.Level(*level),
		Service:   *service,
	}
	analyserSettings := analyser.AnalyserSettings{
		SlowestLogsToRetrieve: *slowestLogs,
	}
	return &parsingSettings, &filters, &analyserSettings, nil
}

func main() {
	parsingSettings, filters, analyserSettings, err := initSettings()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	logsChan, errsChan := ProcessFiles(parsingSettings)
	filteredLogs := filter.ProcessFilter(logsChan, filters)

	metrics := analyser.AnalyseLogs(filteredLogs, errsChan, analyserSettings)
	if payload, err := json.Marshal(metrics); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	} else {
		fmt.Println(string(payload))
	}
}
