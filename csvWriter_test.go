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
	csvWriter := NewWriter(writer)
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
	csvWriter := NewWriter(writer)
	csvWriter.Write(entry)
	csvWriter.Flush()

	if !(writer.String() == expected) {
		t.Errorf("expected %q to eq %q", writer.String(), expected)
	}
}

func TestCsvWriter_Empty(t *testing.T) {

	entry := logcat.Entry{}
	expected := "\n"

	writer := new(bytes.Buffer)
	csvWriter := NewWriter(writer)
	csvWriter.Write(entry)
	csvWriter.Flush()

	if !(writer.String() == expected) {
		t.Errorf("expected %q to eq %q", writer.String(), expected)
	}
}

func TestCsvWriter_Nil(t *testing.T) {

	expected := ""

	writer := new(bytes.Buffer)
	csvWriter := NewWriter(writer)
	csvWriter.Write(nil)
	csvWriter.Flush()

	if !(writer.String() == expected) {
		t.Errorf("expected %q to eq %q", writer.String(), expected)
	}
}
