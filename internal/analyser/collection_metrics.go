package analyser

import (
	"sort"
	"time"

	"github.com/HugoDrl/zebra/internal/parser"
)

func validateLog(log parser.Log) bool {
	if log.Service == "" {
		return false
	}
	if log.Duration < 0 {
		return false
	}
	return true
}

func (m *CollectionMetric) handleService(log *parser.Log) {
	if log == nil {
		return
	}
	if !validateLog(*log) {
		return
	}

	s, ok := m.ServicePerformance[log.Service]
	if !ok {
		s = ServiceMetric{
			Name: log.Service,
		}
	}
	s.Lines++

	s.AverageDuration = s.AverageDuration * time.Duration(s.Lines-1)
	s.AverageDuration += time.Duration(log.Duration)
	s.AverageDuration /= time.Duration(s.Lines)
	m.ServicePerformance[log.Service] = s
	m.Lines[log.Level]++
}

func (m *CollectionMetric) handleSlowestLogs(
	slowestLogsToRetrieve int,
	log *parser.Log,
) {
	if slowestLogsToRetrieve == 0 {
		return
	}
	if len(m.SlowestInput) < slowestLogsToRetrieve {
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
