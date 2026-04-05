package analyser

import (
	"time"

	"git.hugoderlyn.com/Hugo/goLogParser.git/parser"
)

type CollectionMetric struct {
	Lines              map[parser.Level]int
	ServicePerformance map[string]*ServiceMetric
	FileErrors         []*parser.FileError
	ParsingErrorCount  int
	SlowestInput       []*parser.Log
	Query              string
}

type ServiceMetric struct {
	Name            string
	Lines           int
	AverageDuration time.Duration
}
