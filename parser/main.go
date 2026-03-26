package parser

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func ReadFile(filename string) ([]byte, error) {
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
		return nil, nil
	}

	date, err := time.Parse(time.RFC3339, words[0])
	if err != nil {
		return nil, err
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
			return nil, nil
		}
		title := strings.ToLower(f[0])

		switch title {
		case "message":
			message = f[1]
		case "service":
			service = f[1]
		case "duration":
			if f[1][len(f[1])-2:] != "ms" {
				return nil, nil
			}
			value, err := strconv.Atoi(f[1][:len(f[1])-2])
			if err != nil {
				return nil, nil
			}

			duration = time.Duration(value * int(time.Millisecond))
		default:
			fields[title] = f[1]
		}
	}
	if message == "" || service == "" || duration == 0 {
		return nil, nil
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

func ParseLog(content string) ([]*Log, []error) {
	lines := strings.Split(content, "\r\n")
	var logs []*Log
	var logsErrors []error

	for _, line := range lines {
		if log, err := parseLine(line); err != nil || log == nil {
			logsErrors = append(logsErrors, err)
		} else {
			logs = append(logs, log)
		}
	}
	for _, log := range logs {
		fmt.Printf("%v+\n", log)
	}
	return logs, logsErrors
}
