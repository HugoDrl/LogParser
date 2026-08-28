package analyser

import (
	"time"

	"github.com/HugoDrl/zebra/internal/parser"
)

type CollectionMetric struct {
	Lines              map[parser.Level]int     `json:"number_of_lines"`
	ServicePerformance map[string]ServiceMetric `json:"service_performance"`
	FileErrors         []*parser.FileError      `json:"file_errors"`
	ParsingErrorCount  int                      `json:"parse_errors_count"`
	SlowestInput       []*parser.Log            `json:"slowest_logs"`
}

func newMetrics() *CollectionMetric {
	metrics := CollectionMetric{}
	metrics.Lines = make(map[parser.Level]int)
	metrics.ServicePerformance = make(map[string]ServiceMetric)
	return &metrics
}

type ServiceMetric struct {
	Name            string        `json:"name"`
	Lines           int           `json:"number_of_lines"`
	AverageDuration time.Duration `json:"average_duration"`
}

type AnalyserSettings struct {
	SlowestLogsToRetrieve int
}
