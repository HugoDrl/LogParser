package parser

import (
	"strings"
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

func toLevel(input string) (Level, bool) {
	switch inputLevel := Level(strings.ToLower(input)); inputLevel {
	case Debug:
		return Debug, true
	case Info:
		return Info, true
	case Warning:
		return Warning, true
	case Error:
		return Error, true
	case Fatal:
		return Fatal, true
	default:
		return "", false
	}
}

// This is necessary as current format for duration is not supported by JSON native unmarshaller.
type Duration time.Duration

func (d *Duration) UnmarshalJSON(data []byte) error {
	sdata := string(data)
	sdata = strings.Trim(sdata, `"`)
	durationParsed, err := time.ParseDuration(sdata)
	if err != nil {
		return err
	}
	*d = Duration(durationParsed)
	return nil
}

type Log struct {
	Time     time.Time         `json:"date"`
	Level    Level             `json:"level"`
	Duration Duration          `json:"duration"`
	Message  string            `json:"message"`
	Service  string            `json:"service"`
	Extra    map[string]string `json:"-"`
}

type ParseSettings struct {
	Files     []string
	Json      bool
	StartDate time.Time
	EndDate   time.Time
	Level     Level
	Service   string
}
