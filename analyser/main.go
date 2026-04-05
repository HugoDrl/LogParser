package analyser

import (
	"sort"
	"time"

	"git.hugoderlyn.com/Hugo/goLogParser.git/parser"
)

func newMetrics() *CollectionMetric {
	metrics := CollectionMetric{}
	metrics.Lines = make(map[parser.Level]int)
	metrics.ServicePerformance = make(map[string]*ServiceMetric)
	return &metrics
}

func (m *CollectionMetric) handleService(log *parser.Log) {
	s := m.ServicePerformance[log.Service]
	if s == nil {
		s = &ServiceMetric{
			Name: log.Service,
		}
	}
	s.Lines++

	s.AverageDuration = s.AverageDuration * time.Duration(s.Lines-1)
	s.AverageDuration += log.Duration
	s.AverageDuration /= time.Duration(s.Lines)
	m.ServicePerformance[log.Service] = s
}

func (m *CollectionMetric) handleSlowestLogs(N int, log *parser.Log) {
	if len(m.SlowestInput) < N {
		m.SlowestInput = append(m.SlowestInput, log)
		return
	}
	sort.Slice(m.SlowestInput, func(i, j int) bool {
		return m.SlowestInput[i].Duration > m.SlowestInput[j].Duration
	})

	if log.Duration > m.SlowestInput[len(m.SlowestInput)-1].Duration {
		m.SlowestInput = append(m.SlowestInput[:len(m.SlowestInput)-1], log)
	}
}

func AnalyseLogs(logChan <-chan *parser.Log, slowestLogs int) *CollectionMetric {
	metrics := newMetrics()
	for log := range logChan {
		metrics.Lines[log.Level]++
		metrics.handleService(log)
		metrics.handleSlowestLogs(slowestLogs, log)
	}

	return metrics
}
