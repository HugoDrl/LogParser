package parser

import (
	"time"
)

type Level string

const (
	Debug   Level = "debug"
	Info    Level = "info"
	Warning Level = "warning"
	Error   Level = "error"
	Fatal   Level = "fatal"
)

type Log struct {
	Time     time.Time
	Level    Level
	Duration time.Duration
	Message  string
	Service  string
	extra    map[string]string
}

type ParseSettings struct {
	Files []string
	StartDate time.Time
	EndDate time.Time
	Level Level
	Service string
}