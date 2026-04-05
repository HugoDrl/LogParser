package parser

import "fmt"

type FileError struct {
	File string
	Err  error
}

func (e *FileError) Error() string {
	return fmt.Sprintf("Error encountered on file %s - %s", e.File, e.Err)
}

type ParseError struct {
	File   FileError
	Line   int
	Reason string
	Err    error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("Error encoutered while parsing file %s on line %d - %s", e.File, e.Line, e.Reason)
}

type ValueError struct {
	File          FileError
	Line          int
	ErroredValue  string
	ExpectedValue string
	Err           error
}

func (e *ValueError) Error() string {
	return fmt.Sprintf("%s : Wrong value on line %d: expected value : %s (got %s)", e.File, e.Line, e.ExpectedValue, e.ErroredValue)
}
