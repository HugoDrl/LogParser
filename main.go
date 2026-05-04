package main

import (
	"fmt"
	"errors"
	"flag"
	"os"
	"strings"
	"sync"
	"time"

	"git.hugoderlyn.com/Hugo/goLogParser.git/analyser"
	"git.hugoderlyn.com/Hugo/goLogParser.git/parser"
)

func ProcessFiles(settings *parser.ParseSettings, outChan chan<- *parser.Log, errsChan chan<- error, wg *sync.WaitGroup) {
	for _, filepath := range settings.Files {
		wg.Add(1)
		go func() {
			defer wg.Done()
			values, errors := parser.ParseFile(filepath, settings)
			for len(values) > 0 || len(errors) > 0 {
				var closableOutChan chan<- *parser.Log
				var closableErrsChan chan<- error
				var v *parser.Log
				var e error

				if len(values) > 0 {
					v = values[len(values)-1]
					closableOutChan = outChan
				}

				if len(errors) > 0 {
					e = errors[len(errors)-1]
					closableErrsChan = errsChan
				}

				select {
				case closableOutChan <- v:
					values = values[:len(values)-1]

				case closableErrsChan <- e:
					errors = errors[:len(errors)-1]
				}
			}
		}()
	}
}

func initSettings() (*parser.ParseSettings, *analyser.AnalyserSettings, error) {
	files := flag.String("files", "", "log files to analyse")
	startDate := flag.String("start", "", "logs date to start from")
	endDate := flag.String("end", "", "logs date to end to")
	service := flag.String("service", "", "filter logs by service")
	level := flag.String("level", "", "filter logs by level")
	slowestLogs := flag.Int("top", 0, "number of slowest logs to show")
	flag.Parse()

	if *files == "" {
		return nil, nil, errors.New("Please specify file(s) separated by a comma using --files flag")
	}

	var processedStartDate time.Time
	var processedEndDate time.Time
	var processErr error
	if *startDate != "" {
		processedStartDate, processErr = time.Parse(time.RFC3339, *startDate)
		if processErr != nil {
			return nil, nil, errors.New("Wrong format for starting date - excpected RFC3339")
		}
	}
	if *endDate != "" {
		processedEndDate, processErr = time.Parse(time.RFC3339, *endDate)
		if processErr != nil {
			return nil, nil, errors.New("Wrong format for starting date - excpected RFC3339")
		}
	}

	parsingSettings := parser.ParseSettings{
		Files: strings.Split(*files, ","),
		StartDate: processedStartDate,
		EndDate: processedEndDate,
		Level: parser.Level(*level),
		Service: *service,
	}
	analyserSettings := analyser.AnalyserSettings{
		SlowestLogsToRetrieve: *slowestLogs,
	}
	return &parsingSettings, &analyserSettings, nil
}

func main() {
	parsingSettings, analyserSettings, err := initSettings()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintln(err))
		os.Exit(1)
	}

	out := make(chan *parser.Log)
	errs := make(chan error)
	var wgFiles sync.WaitGroup
	ProcessFiles(parsingSettings, out, errs, &wgFiles)

	go func() {
		wgFiles.Wait()
		close(out)
		close(errs)
	}()

	metrics := analyser.AnalyseLogs(out, errs , analyserSettings)
	metrics.Display()
}
