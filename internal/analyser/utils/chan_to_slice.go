package utils

import (
	"sync"

	"github.com/HugoDrl/zebra/internal/parser"
)

func ExtractLogAndErrChanToSlices(logChan <-chan *parser.Log, errChan <-chan error) ([]*parser.Log, []error) {
	logs := make([]*parser.Log, 0)
	errs := make([]error, 0)
	var wg sync.WaitGroup

	wg.Go(func() {
		for log := range logChan {
			logs = append(logs, log)
		}
	})
	wg.Go(func() {
		for err := range errChan {
			errs = append(errs, err)
		}
	})

	wg.Wait()
	return logs, errs
}
