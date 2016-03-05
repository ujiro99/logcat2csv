package main

import (
	"encoding/csv"
	"io"
	"runtime"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/ujiro99/logcatf/logcat"
)

const (
	// Windows represents windows os.
	Windows string = "windows"
	// UTF8 represents encode `utf-8`
	UTF8 = "utf-8"
	// ShiftJIS represents encode `shift-jis`
	ShiftJIS = "shift-jis"
	// EUCJP represents encode `euc-jp`
	EUCJP = "euc-jp"
	// ISO2022JP represents encode `iso-2022-jp`
	ISO2022JP = "iso-2022-jp"
)

// CsvWriter is wrapper of csv.Writer to writing logcat.Entry.
type CsvWriter struct {
	w *csv.Writer
}

// NewWriter creates new csvWriter.
func NewWriter(w io.Writer, encode string) *CsvWriter {
	if runtime.GOOS == Windows && encode == "" {
		encode = ShiftJIS
	}
	switch encode {
	case ShiftJIS:
		w = transform.NewWriter(w, japanese.ShiftJIS.NewEncoder())
	case EUCJP:
		w = transform.NewWriter(w, japanese.EUCJP.NewEncoder())
	case ISO2022JP:
		w = transform.NewWriter(w, japanese.ISO2022JP.NewEncoder())
	}

	res := &CsvWriter{csv.NewWriter(w)}
	if runtime.GOOS == Windows {
		res.w.UseCRLF = true
	}
	return res
}

// Write write logcat.Entry.
func (f *CsvWriter) Write(item logcat.Entry) {
	if item == nil {
		return
	}
	f.w.Write(item.Values())
}

// Flush flushes buffer to file.
func (f *CsvWriter) Flush() {
	f.w.Flush()
}
