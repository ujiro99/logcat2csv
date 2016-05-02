package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/Maki-Daisuke/go-lines"
	"github.com/ujiro99/logcatf/logcat"
)

// MaxFailCount represents count for cancel conversion.
const MaxFailCount = 5

type logcat2csv struct{}

func (l *logcat2csv) execStream(params cmdParams) int {
	err := l.exec(params)
	if err != nil {
		fmt.Fprintf(params.error, err.Error())
		return ExitCodeError
	}
	return ExitCodeOK
}

func (l *logcat2csv) execFiles(params cmdParams) int {
	success := false
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
		err := l.exec(params)
		if err != nil {
			fmt.Fprintf(params.error, "%s: %s\n", err, path)
			os.Remove(path + ".csv")
		} else {
			success = true
		}
	}
	if success {
		return ExitCodeOK
	} else {
		return ExitCodeError
	}
}

func (l *logcat2csv) exec(params cmdParams) error {
	csvWriter := NewWriter(params.writer, params.encode, params.osName)
	parser := logcat.NewParser()
	fail := 0
	success := 0
	for line := range lines.Lines(params.reader) {
		entry, err := parser.Parse(line)
		if err == nil {
			csvWriter.Write(entry) // TODO: handle err
			if entry.Format() == "raw" {
				fail++ // Not a logcat format.
			} else {
				success++
			}
		} else {
			fail++ // Can't parse to logcat.
		}
		if fail > MaxFailCount {
			return errors.New("Parse error. Conversion canceled")
		}
	}
	if success <= 0 {
		return errors.New("Format error. Conversion canceled")
	}
	csvWriter.Flush()
	return nil
}

// Exec execute converting.
func (l *logcat2csv) Exec(params cmdParams) int {
	if params.reader != nil {
		return l.execStream(params)
	}
	return l.execFiles(params)
}
