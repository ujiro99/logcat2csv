package main

import (
	"fmt"
	"os"

	"github.com/Maki-Daisuke/go-lines"
	"github.com/ujiro99/logcatf/logcat"
)

type logcat2csv struct{}

func (l *logcat2csv) execStream(params cmdParams) int {
	return l.exec(params)
}

func (l *logcat2csv) execFiles(params cmdParams) int {
	for _, path := range params.paths {
		r, e := os.Open(path)
		defer r.Close()
		if e != nil {
			fmt.Printf("File open error: %s\n", path)
			return ExitCodeError
		}
		w, e := os.Create(path + ".csv")
		defer w.Close()
		if e != nil {
			fmt.Printf("File create error: %s\n", path+".csv")
			return ExitCodeError
		}
		params.reader = r
		params.writer = w
		l.exec(params)
	}
	return ExitCodeOK
}

func (l *logcat2csv) exec(params cmdParams) int {
	csvWriter := NewWriter(params.writer)
	parser := logcat.NewParser()
	for line := range lines.Lines(params.reader) {
		entry, err := parser.Parse(line)
		if err == nil {
			csvWriter.Write(entry)
		}
	}
	csvWriter.Flush()
	return ExitCodeOK
}

// Exec execute converting.
func (l *logcat2csv) Exec(params cmdParams) int {
	if params.reader != nil {
		return l.execStream(params)
	}
	return l.execFiles(params)
}
