package filter

import "github.com/HugoDrl/zebra/internal/parser"

func ProcessFilter(inputChan <-chan *parser.Log, filters *Filters) <-chan *parser.Log {
	filteredChan := make(chan *parser.Log)
	go func() {
		defer close(filteredChan)
		for log := range inputChan {
			if filterLog(log, filters) {
				filteredChan <- log
			}
		}
	}()
	return filteredChan
}
