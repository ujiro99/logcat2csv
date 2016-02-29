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

func (l *logcat2csv) execFile(params cmdParams) int {
	r, e := os.Open(params.path)
	defer r.Close()
	if e != nil {
		fmt.Printf("File open error: %s\n", params.path)
		return ExitCodeError
	}
	w, e := os.Create(os.Args[1] + ".csv")
	defer w.Close()
	if e != nil {
		fmt.Printf("File create error: %s\n", params.path+".csv")
		return ExitCodeError
	}
	params.reader = r
	params.writer = w
	return l.exec(params)
}

func (l *logcat2csv) exec(params cmdParams) int {
	csvWriter := NewWriter(params.writer)
	parser := logcat.NewParser()
	for line := range lines.Lines(params.reader) {
		entry, _ := parser.Parse(line)
		csvWriter.Write(entry)
	}
	csvWriter.Flush()
	return ExitCodeOK
}

// Exec execute converting.
func (l *logcat2csv) Exec(params cmdParams) int {
	if params.reader != nil {
		return l.execStream(params)
	}
	return l.execFile(params)
}
