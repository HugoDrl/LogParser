package main

import (
	"errors"
	"flag"
	"fmt"
	"sync"

	"git.hugoderlyn.com/Hugo/goLogParser.git/analyser"
	"git.hugoderlyn.com/Hugo/goLogParser.git/parser"
)

func ProcessFiles(filepaths []string, outChan chan<- *parser.Log, errsChan chan<- error, wg *sync.WaitGroup) {
	for _, filepath := range filepaths {
		wg.Add(1)
		go func() {
			defer wg.Done()
			values, errors := parser.ParseFile(filepath)
			for len(values) > 0 || len(errors) > 0 {
				var closableOutChan chan<- *parser.Log
				var closableErrsChan chan<- error
				var v *parser.Log
				var e error

				if len(values) > 0 {
					v = values[len(values)-1]
					closableOutChan = outChan
				}

				if len(errors) > 0 {
					e = errors[len(errors)-1]
					closableErrsChan = errsChan
				}

				select {
				case closableOutChan <- v:
					values = values[:len(values)-1]

				case closableErrsChan <- e:
					errors = errors[:len(errors)-1]
				}
			}
		}()
	}
}

func main() {
	files := flag.String("files", "", "log files to analyse")
	files := flag.String("files", "", "log files to analyse")
	files := flag.String("files", "", "log files to analyse")
	files := flag.String("files", "", "log files to analyse")
	files := flag.String("files", "", "log files to analyse")

	out := make(chan *parser.Log, 10)
	errs := make(chan error, 10)
	var wgFiles sync.WaitGroup
	var wgProcess sync.WaitGroup
	ProcessFiles([]string{"test1.log", "test2.log", "test3.log"}, out, errs, &wgFiles)

	go func() {
		wgFiles.Wait()
		close(out)
		close(errs)
	}()

	wgProcess.Add(1)
	go func() {
		defer wgProcess.Done()
		metrics := analyser.AnalyseLogs(out, 3)
		fmt.Println(metrics)
		for _, s := range metrics.ServicePerformance {
			fmt.Println(s)
		}
		for _, l := range metrics.SlowestInput {
			fmt.Println(l)
		}
	}()

	wgProcess.Add(1)
	go func() {
		linesErrCounter := 0
		defer wgProcess.Done()
		for err := range errs {
			var fileErr *parser.FileError
			if errors.As(err, &fileErr) {
				fmt.Println(err.Error())
			} else {
				linesErrCounter++
			}
		}
		fmt.Printf("%d errors on lines encountered.\n", linesErrCounter)
	}()

	wgProcess.Wait()
}
