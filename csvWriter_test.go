package main

import (
	"bytes"
	"testing"

	"github.com/ujiro99/logcatf/logcat"
)

func TestCsvWriter_Write_Full(t *testing.T) {

	entry := logcat.Entry{
		"time":     "12-28 18:54:07.180",
		"pid":      "930",
		"tid":      "931",
		"priority": "I",
		"tag":      "auditd",
		"message":  "  test Message",
	}
	expected := "12-28 18:54:07.180,930,931,I,auditd,\"  test Message\"\n"

	writer := new(bytes.Buffer)
	csvWriter := NewWriter(writer, "", "")
	csvWriter.Write(entry)
	csvWriter.Flush()

	if !(writer.String() == expected) {
		t.Errorf("expected %q to eq %q", writer.String(), expected)
	}
}

func TestCsvWriter_Write(t *testing.T) {

	entry := logcat.Entry{
		"time":    "12-28 18:54:07.180",
		"message": "  test Message",
		"tag":     "auditd",
	}
	expected := "12-28 18:54:07.180,auditd,\"  test Message\"\n"

	writer := new(bytes.Buffer)
	csvWriter := NewWriter(writer, "", "")
	csvWriter.Write(entry)
	csvWriter.Flush()

	if !(writer.String() == expected) {
		t.Errorf("expected %q to eq %q", writer.String(), expected)
	}
}

func TestCsvWriter_Write_Windows(t *testing.T) {

	entry := logcat.Entry{
		"time":    "12-28 18:54:07.180",
		"message": "test Message:漢字",
		"tag":     "auditd",
	}
	expected := convertTo("12-28 18:54:07.180,auditd,test Message:漢字\r\n", ShiftJIS)

	writer := new(bytes.Buffer)
	csvWriter := NewWriter(writer, "", "windows")
	csvWriter.Write(entry)
	csvWriter.Flush()

	if !(writer.String() == expected) {
		t.Errorf("expected %q to eq %q", writer.String(), expected)
	}
}

func TestCsvWriter_Write_Encode(t *testing.T) {

	entry := logcat.Entry{
		"time":    "12-28 18:54:07.180",
		"message": "test Message:あ亜Ａア￥凜熙♪堯",
		"tag":     "auditd",
	}
	expected := ("12-28 18:54:07.180,auditd,test Message:あ亜Ａア￥凜熙♪堯\n")

	for _, encode := range []string{UTF8, ShiftJIS, EUCJP, ISO2022JP} {
		writer := new(bytes.Buffer)
		csvWriter := NewWriter(writer, encode, "")
		csvWriter.Write(entry)
		csvWriter.Flush()

		if !(writer.String() == convertTo(expected, encode)) {
			t.Errorf("expected %q to eq %q", writer.String(), convertTo(expected, encode))
		}
	}
}

func TestCsvWriter_Empty(t *testing.T) {

	entry := logcat.Entry{}
	expected := "\n"

	writer := new(bytes.Buffer)
	csvWriter := NewWriter(writer, "", "")
	csvWriter.Write(entry)
	csvWriter.Flush()

	if !(writer.String() == expected) {
		t.Errorf("expected %q to eq %q", writer.String(), expected)
	}
}

func TestCsvWriter_Nil(t *testing.T) {

	expected := ""

	writer := new(bytes.Buffer)
	csvWriter := NewWriter(writer, "", "")
	csvWriter.Write(nil)
	csvWriter.Flush()

	if !(writer.String() == expected) {
		t.Errorf("expected %q to eq %q", writer.String(), expected)
	}
}
