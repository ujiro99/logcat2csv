package main

import (
	"fmt"
	"os"

	"github.com/Maki-Daisuke/go-lines"
	"github.com/ujiro99/logcatf/logcat"
)

// MaxFailCount represents count for cancel conversion.
const MaxFailCount = 5

type logcat2csv struct{}

func (l *logcat2csv) execStream(params cmdParams) int {
	status := l.exec(params)
	if status == ExitCodeError {
		fmt.Fprintf(params.error, "Parse error. Conversion canceled.")
	}
	return status
}

func (l *logcat2csv) execFiles(params cmdParams) int {
	for _, path := range params.paths {
		r, e := os.Open(path)
		defer r.Close()
		if e != nil {
			fmt.Fprintf(params.error, "File open error: %s\n", path)
			continue
		}
		w, e := os.Create(path + ".csv")
		defer w.Close()
		if e != nil {
			fmt.Fprintf(params.error, "File create error: %s\n", path+".csv")
			continue
		}
		params.reader = r
		params.writer = w
		status := l.exec(params)
		if status == ExitCodeError {
			fmt.Fprintf(params.error, "Parse error. Conversion canceled: %s\n", path)
			os.Remove(path + ".csv")
		}
	}
	return ExitCodeOK
}

func (l *logcat2csv) exec(params cmdParams) int {
	csvWriter := NewWriter(params.writer, params.encode, params.osName)
	parser := logcat.NewParser()
	failCount := 0
	for line := range lines.Lines(params.reader) {
		entry, err := parser.Parse(line)
		if err == nil {
			csvWriter.Write(entry)
		}
		if err != nil || entry.Format() == "raw" {
			failCount = failCount + 1
		}
		if failCount > MaxFailCount {
			return ExitCodeError
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
