package analyser

import (
	"errors"
	"sync"

	"github.com/HugoDrl/zebra/internal/parser"
)

func AnalyseLogs(
	logs []*parser.Log,
	errs []error,
	settings *AnalyserSettings,
) *CollectionMetric {
	metrics := newMetrics()
	var wg sync.WaitGroup

	for _, log := range logs {
		metrics.handleService(log)
		metrics.handleSlowestLogs(settings.SlowestLogsToRetrieve, log)
	}

	for _, err := range errs {
		var fileErr *parser.FileError
		if errors.As(err, &fileErr) {
			metrics.FileErrors = append(metrics.FileErrors, fileErr)
		} else {
			metrics.ParsingErrorCount++
		}
	}

	wg.Wait()

	return metrics
}
