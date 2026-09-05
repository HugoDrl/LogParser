package flags

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/HugoDrl/zebra/internal/analyser"
	"github.com/HugoDrl/zebra/internal/filter"
	"github.com/HugoDrl/zebra/internal/parser"
)

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

func InitSettings() (*parser.ParseSettings, *filter.Filters, *analyser.AnalyserSettings, error) {
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
