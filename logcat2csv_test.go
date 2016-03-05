package main

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestLogcat2csv_Exec_Stdio(t *testing.T) {
	expect := "01-01 00:00:00.000,930,931,I,tag_value,message_value\n"
	out := new(bytes.Buffer)
	params := cmdParams{
		reader: strings.NewReader("01-01 00:00:00.000   930   931 I tag_value  : message_value"),
		writer: out,
	}

	logcat2csv := logcat2csv{}
	logcat2csv.Exec(params)
	if out.String() != expect {
		t.Errorf("\n  result: %q\n  expect: %q", out.String(), expect)
	}
}

func TestLogcat2csv_Exec_File(t *testing.T) {
	expect := []string{
		"01-01 00:00:00.000,930,931,I,tag_value,message_value_1",
		"01-01 00:00:01.000,930,931,I,tag_value,message_value_2",
	}
	paths := []string{"./test/logcat.txt"}
	params := cmdParams{
		paths: paths,
	}

	logcat2csv := logcat2csv{}
	logcat2csv.Exec(params)

	checkFile(paths[0], expect, t)
}

func TestLogcat2csv_Exec_Multiple_File(t *testing.T) {
	expect0 := []string{
		"01-01 00:00:00.000,930,931,I,tag_value,message_value_1",
		"01-01 00:00:01.000,930,931,I,tag_value,message_value_2",
	}
	expect1 := []string{
		"01-01 00:00:00.000,930,931,I,tag_value,message_value_3",
		"01-01 00:00:01.000,930,931,I,tag_value,message_value_4",
	}
	paths := []string{"./test/logcat.txt", "./test/logcat2.txt"}
	params := cmdParams{
		paths: paths,
	}

	logcat2csv := logcat2csv{}
	logcat2csv.Exec(params)

	checkFile(paths[0], expect0, t)
	checkFile(paths[1], expect1, t)
}

func checkFile(file string, expect []string, t *testing.T) {
	var out string
	fp, err := os.Open(file + ".csv")
	if err != nil {
		t.Errorf("os.Open: %v", err)
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)

	i := 0
	for scanner.Scan() {
		out = scanner.Text()
		if out != expect[i] {
			t.Errorf("\n  result: %q\n  expect: %q", out, expect[i])
		}
		i = i + 1
	}
	if err := scanner.Err(); err != nil {
		t.Errorf("scanner.Err: %v", err)
	} else {
		os.Remove(file + ".csv")
	}
}
