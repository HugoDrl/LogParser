package parser

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func getConsumedJSONFields(t reflect.Type) map[string]struct{} {
	foundFields := make(map[string]struct{}, t.NumField())
	for field := range t.Fields() {
		value, ok := field.Tag.Lookup("json")
		if !ok {
			continue
		}
		foundFields[value] = struct{}{}
	}
	return foundFields
}

func parseJSONExtraArguments(t reflect.Type, line string) (map[string]string, error) {
	var extras map[string]string
	if err := json.Unmarshal([]byte(line), &extras); err != nil {
		return nil, &ParseError{
			Reason: err.Error(),
		}
	}
	logFields := getConsumedJSONFields(t)
	for key := range extras {
		if _, ok := logFields[key]; ok {
			delete(extras, key)
		}
	}
	if len(extras) == 0 {
		extras = nil
	}
	return extras, nil
}

func ParseJSONFormatLine(line string) (Log, error) {
	var log Log
	if err := json.Unmarshal([]byte(line), &log); err != nil {
		return Log{}, &ParseError{
			Reason: err.Error(),
		}
	}
	if level, ok := toLevel(string(log.Level)); !ok {
		return Log{}, &ValueError{ErroredValue: "level"}
	} else {
		log.Level = level
	}

	if extras, err := parseJSONExtraArguments(reflect.TypeFor[Log](), line); err != nil {
		return Log{}, err
	} else {
		log.Extra = extras
	}

	return log, nil
}

func ParseDefaultFormatLine(line string) (Log, error) {
	words := splitLine(line)
	if len(words) < 2 {
		return Log{}, &ParseError{Reason: fmt.Sprintf("Not enough arguments - expected 2 - found %d", len(words))}
	}

	date, err := time.Parse(time.RFC3339, words[0])
	if err != nil {
		return Log{}, &ValueError{ExpectedValue: "time format - RFC3339", ErroredValue: words[0]}
	}

	l := strings.TrimFunc(words[1], func(l rune) bool { return l == '[' || l == ']' })
	level, ok := toLevel(strings.ToLower(l))
	if !ok {
		return Log{}, &ValueError{
			ErroredValue: "level",
		}
	}
	fields := make(map[string]string)
	var message string
	var service string
	var duration time.Duration
	for i, word := range words {
		if i < 2 {
			continue
		}
		f := strings.Split(word, "=")
		if len(f) != 2 {
			return Log{}, &ParseError{Reason: fmt.Sprintf("Wrong format in key-values design (found %v)", f)}
		}
		title := strings.ToLower(f[0])

		switch title {
		case "message":
			message = f[1]
		case "service":
			service = f[1]
		case "duration":
			if f[1][len(f[1])-2:] != "ms" {
				return Log{}, &ValueError{ExpectedValue: "duration value in milliseconds on format xxxms", ErroredValue: f[1]}
			}
			value, err := strconv.Atoi(f[1][:len(f[1])-2])
			if err != nil {
				return Log{}, &ValueError{ExpectedValue: "duration value in milliseconds on format xxxms", ErroredValue: f[1]}
			}

			duration = time.Duration(value * int(time.Millisecond))
		default:
			fields[title] = f[1]
		}
	}
	if message == "" || service == "" || duration == 0 {
		return Log{}, &ValueError{ExpectedValue: "message field, service field, and duration field, should not be empty.", ErroredValue: fmt.Sprintf("message: %s, service: %s, duration: %d", message, service, duration)}
	}
	return Log{
		Time:     date,
		Level:    level,
		Message:  message,
		Service:  service,
		Extra:    fields,
		Duration: Duration(duration),
	}, nil
}
