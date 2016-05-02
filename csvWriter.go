package main

import (
	"encoding/csv"
	"io"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"bytes"

	"github.com/ujiro99/logcatf/logcat"
	"errors"
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
	// Message represents message key of Logcat item's entry
	Message = "message"
)

// CsvWriter is wrapper of csv.Writer to writing logcat.Entry.
type CsvWriter struct {
	encodedWriter *csv.Writer
	writer        *csv.Writer
	encoder       io.Writer
	buff          *bytes.Buffer
}

// NewWriter creates new csvWriter.
func NewWriter(w io.Writer, encode string, osName string) *CsvWriter {
	if osName == Windows && encode == "" {
		encode = ShiftJIS
	}
	res := &CsvWriter{
		encodedWriter: csv.NewWriter(generateEncoder(w, encode)),
	}

	// for fail-safe of encoding.
	if encode != UTF8 {
		buff := new(bytes.Buffer)
		res.writer = csv.NewWriter(w)
		res.encoder = generateEncoder(buff, encode)
		res.buff = buff
	}

	if osName == Windows {
		res.encodedWriter.UseCRLF = true
	}
	return res
}

// Write write logcat.Entry.
func (f *CsvWriter) Write(item logcat.Entry) (err error) {
	if item == nil {
		return nil
	}

	err = f.canEncode(item[Message])
	if err == nil {
		f.encodedWriter.Write(item.Values())
	} else {
		// If the message can't be encoded, output with UTF8.
		f.encodedWriter.Flush()
		f.writer.Write(item.Values())
		f.writer.Flush()
	}

	return err
}

// Flush flushes buffer to file.
func (f *CsvWriter) Flush() {
	f.encodedWriter.Flush()
}

// Check the `str` can be encoded.
func (f *CsvWriter) canEncode(str string) (err error) {
	if f.buff == nil {
		return nil
	}
	_, err = f.encoder.Write([]byte(str))
	f.buff.Reset()
	if err != nil {
		return errors.New("Not supported encoding.")
	}
	return nil
}

func generateEncoder(w io.Writer, encode string) io.Writer {
	var enc io.Writer

	switch encode {
	case ShiftJIS:
		enc = transform.NewWriter(w, japanese.ShiftJIS.NewEncoder())
	case EUCJP:
		enc = transform.NewWriter(w, japanese.EUCJP.NewEncoder())
	case ISO2022JP:
		enc = transform.NewWriter(w, japanese.ISO2022JP.NewEncoder())
	default:
		enc = w
	}

	return enc
}
