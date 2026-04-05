package parser

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func readFile(filename string) ([]byte, error) {
	c, err := os.ReadFile(filename)
	if err != nil {
		return nil, &FileError{File: filename, Err: err}
	}
	return c, nil
}

func splitLine(line string) []string {
	words := make([]string, 0)
	currWord := ""
	for _, letter := range line {
		if letter == ' ' && strings.Count(currWord, "\"")%2 == 0 {
			words = append(words, currWord)
			currWord = ""
			continue
		}
		currWord += string(letter)

	}
	if currWord != "" {
		words = append(words, currWord)
	}
	return words
}

func parseLine(line string) (*Log, error) {
	words := splitLine(line)
	if len(words) < 2 {
		return nil, &ParseError{Reason: fmt.Sprintf("Not enough arguments - expected 2 - found %d", len(words))}
	}

	date, err := time.Parse(time.RFC3339, words[0])
	if err != nil {
		return nil, &ValueError{ExpectedValue: "time format - RFC3339", ErroredValue: words[0]}
	}

	l := strings.TrimFunc(words[1], func(l rune) bool { return l == '[' || l == ']' })
	level := Level(strings.ToLower(l))
	fields := make(map[string]string, 0)
	var message string
	var service string
	var duration time.Duration
	for i, word := range words {
		if i < 2 {
			continue
		}
		f := strings.Split(word, "=")
		if len(f) != 2 {
			return nil, &ParseError{Reason: fmt.Sprintf("Wrong format in key-values design (found %v)", f)}
		}
		title := strings.ToLower(f[0])

		switch title {
		case "message":
			message = f[1]
		case "service":
			service = f[1]
		case "duration":
			if f[1][len(f[1])-2:] != "ms" {
				return nil, &ValueError{ExpectedValue: "duration value in milliseconds on format xxxms", ErroredValue: f[1]}
			}
			value, err := strconv.Atoi(f[1][:len(f[1])-2])
			if err != nil {
				return nil, &ValueError{ExpectedValue: "duration value in milliseconds on format xxxms", ErroredValue: f[1]}
			}

			duration = time.Duration(value * int(time.Millisecond))
		default:
			fields[title] = f[1]
		}
	}
	if message == "" || service == "" || duration == 0 {
		return nil, &ValueError{ExpectedValue: "message field, service field, and duration field, should not be empty.", ErroredValue: fmt.Sprintf("message: %s, service: %s, duration: %d", message, service, duration)}
	}
	return &Log{
		Time:     date,
		Level:    level,
		Message:  message,
		Service:  service,
		extra:    fields,
		Duration: duration,
	}, nil
}

func parseLog(content string) ([]*Log, []error) {
	lines := strings.Split(content, "\r\n")
	var logs []*Log
	var logsErrors []error

	for line_no, line := range lines {
		if log, err := parseLine(line); err != nil {
			var valueErr *ValueError
			if errors.As(err, &valueErr) {
				valueErr.Line = line_no
			}
			logsErrors = append(logsErrors, err)
		} else {
			logs = append(logs, log)
		}
	}

	return logs, logsErrors
}

func ParseFile(file string) ([]*Log, []error) {
	logsParsed := make([]*Log, 0)
	logsErr := make([]error, 0)
	content, err := readFile(file)
	if err != nil {
		return nil, []error{err}
	}

	logs, errs := parseLog(string(content))
	if errs != nil {
		var fileErr *FileError
		for _, err := range errs {
			if errors.As(err, &fileErr) {
				fileErr.File = file
			}
			logsErr = append(logsErr, err)
		}
	}
	logsParsed = append(logsParsed, logs...)

	return logsParsed, logsErr
}
