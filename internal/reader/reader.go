package reader

import (
	"bufio"

	"github.com/HugoDrl/zebra/internal/parser"
)

type ExtractLinesFromFileInput struct {
	Reader        *bufio.Reader
	ParseFunction parser.ParseFunction
	LogsChan      chan<- *parser.Log
	ErrsChan      chan<- error
}

func ExtractLinesFromFile(input ExtractLinesFromFileInput) {
	scanner := bufio.NewScanner(input.Reader)

	lineNo := 0
	for {
		if ok := scanner.Scan(); !ok {
			if err := scanner.Err(); err != nil {
				input.ErrsChan <- err
			}
			return
		}
		lineNo++
		contentLine := scanner.Text()

		log, err := input.ParseFunction(contentLine)
		if err != nil {
			input.ErrsChan <- &parser.ParseError{
				Line: lineNo,
				Err:  err,
			}
		} else {
			input.LogsChan <- &log
		}
	}
}
