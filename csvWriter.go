package main

import (
	"encoding/csv"
	"io"
	"runtime"

	"github.com/ujiro99/logcatf/logcat"
)

const (
	// Windows represents windows os.
	Windows string = "windows"
)

// CsvWriter is wrapper of csv.Writer to writing logcat.Entry.
type CsvWriter struct {
	w *csv.Writer
}

// NewWriter creates new csvWriter.
func NewWriter(w io.Writer) *CsvWriter {
	res := &CsvWriter{csv.NewWriter(w)}
	if runtime.GOOS == Windows {
		res.w.UseCRLF = true
	}
	return res
}

// Write write logcat.Entry.
func (f *CsvWriter) Write(item logcat.Entry) {
	f.w.Write(item.Values())
}

// Flush flushes buffer to file.
func (f *CsvWriter) Flush() {
	f.w.Flush()
}
