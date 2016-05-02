package main

import (
	"bufio"
	"bytes"
	"fmt"
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

	if err := checkFile(paths[0], expect); err != nil {
		t.Error(err)
	}
}

func TestLogcat2csv_Exec_File_Nil(t *testing.T) {
	fileName := "not_a_file"
	expect := "File open error: " + fileName + "\n"

	out := new(bytes.Buffer)
	params := cmdParams{
		paths: []string{fileName},
		error: out,
	}

	logcat2csv := logcat2csv{}
	logcat2csv.Exec(params)

	if out.String() != expect {
		t.Errorf("\n  result: %q\n  expect: %q", out.String(), expect)
	}
}

func TestLogcat2csv_Exec_File_Raw(t *testing.T) {
	fileName := "./test/logcat.raw.txt"
	expect := "Parse error. Conversion canceled: " + fileName + "\n"

	out := new(bytes.Buffer)
	params := cmdParams{
		paths: []string{fileName},
		error: out,
	}

	logcat2csv := logcat2csv{}
	logcat2csv.Exec(params)

	if out.String() != expect {
		t.Errorf("\n  result: %q\n  expect: %q", out.String(), expect)
	}
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

	if err := checkFile(paths[0], expect0); err != nil {
		t.Error(err)
	}
	if err := checkFile(paths[1], expect1); err != nil {
		t.Error(err)
	}
}

func checkFile(file string, expect []string) error {
	var out string
	fp, err := os.Open(file + ".csv")
	if err != nil {
		return err
	}
	defer fp.Close()
	defer os.Remove(file + ".csv")

	if expect != nil {
		scanner := bufio.NewScanner(fp)
		i := 0
		for range expect {
			if !scanner.Scan() {
				return fmt.Errorf("\n  result: %q\n  expect: %q", "", expect[i])
			}
			out = scanner.Text()
			//fmt.Printf("res[%d]: %s\n", i, out)
			//fmt.Printf("exp[%d]: %s\n", i, expect[i])
			if out != expect[i] {
				return fmt.Errorf("\n  result: %q\n  expect: %q", out, expect[i])
			}
			i = i + 1
		}
		if err := scanner.Err(); err != nil {
			return err
		}
	}

	return nil
}
