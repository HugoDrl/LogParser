package parser

import "time"

type Level int

const (
	Debug Level = iota
	Info
	Warning
	Error
	Fatal
)

type Log struct {
	Time     time.Time
	Level    Level
	Duration time.Duration
	Message  string
	Service  string
	extra    map[string]string
}

type CollectionMetric struct {
	Lines              int
	LogErrors          int
	ServicePerformance []ServiceMetric
	ParsingErrorCount  int
	SlowestInput       []Log
	Query              string
}

type ServiceMetric struct {
	Name          string
	Lines         int
	DuratedLines  int
	TotalDuration time.Duration
	LogErrors     int
}
