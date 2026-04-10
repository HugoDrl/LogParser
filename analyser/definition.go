package analyser

import (
	"fmt"
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

func (m *CollectionMetric)Display() {
	fmt.Printf("Number of lines : \n")
	for key, value := range m.Lines {
		fmt.Printf("Level %s : %d\n", key, value)
	}
	for _, value := range m.ServicePerformance {
		value.Display()
	}
	if len(m.SlowestInput) > 0 {
		fmt.Printf("Slowest inputs are the following :\n")
		for _, s := range m.SlowestInput{
			fmt.Printf("%s %v %v : %s (%s)\n", s.Service, s.Duration, s.Time, s.Message, s.Level)
		}
	}

	fmt.Printf("%d errors encountered while parsing files", m.ParsingErrorCount)
	if len(m.FileErrors) > 0 {
		fmt.Println("Following errors encountered while opening files :")
		for _, e := range m.FileErrors {
			fmt.Println(e.Error())
		}
	}
}

type ServiceMetric struct {
	Name            string
	Lines           int
	AverageDuration time.Duration
}

func (s *ServiceMetric)Display() {
	fmt.Printf("-- SERVICE %s --\n", s.Name)
	fmt.Printf("%d lines in service\n", s.Lines)
	fmt.Printf("Average duration : %dms\n", s.AverageDuration.Milliseconds())
}
type AnalyserSettings struct {
	SlowestLogsToRetrieve int
}