package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/HugoDrl/zebra/internal/analyser"
	"github.com/HugoDrl/zebra/internal/filter"
	"github.com/HugoDrl/zebra/internal/flags"
	"github.com/HugoDrl/zebra/internal/parser"
	"github.com/HugoDrl/zebra/internal/reader"
)

func processFiles(
	settings *parser.ParseSettings,
) (<-chan *parser.Log, <-chan error) {
	logsChan := make(chan *parser.Log)
	errsChan := make(chan error)
	go func() {
		defer close(logsChan)
		defer close(errsChan)
		var wg sync.WaitGroup
		root, err := os.OpenRoot(".")
		if err != nil {
			errsChan <- err
			return
		}
		defer root.Close()
		defer wg.Wait()
		for _, filepath := range settings.Files {
			wg.Go(func() {
				filesReader, err := root.OpenFile(filepath, os.O_RDONLY, 0o000)
				if err != nil {
					errsChan <- &parser.FileError{
						File: filepath,
						Err:  err,
					}
					return
				}
				defer filesReader.Close()
				r := bufio.NewReader(filesReader)

				reader.ExtractLinesFromFile(reader.ExtractLinesFromFileInput{
					Reader:        r,
					ParseFunction: parser.GetParseFunction(*settings),
					LogsChan:      logsChan,
					ErrsChan:      errsChan,
				})
			})
		}
	}()
	return logsChan, errsChan
}

func main() {
	parsingSettings, filters, analyserSettings, err := flags.InitSettings()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	logsChan, errsChan := processFiles(parsingSettings)
	filteredLogs := filter.ProcessFilter(logsChan, filters)

	metrics := analyser.AnalyseLogs(filteredLogs, errsChan, analyserSettings)
	if payload, err := json.Marshal(metrics); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	} else {
		fmt.Println(string(payload))
	}
}
