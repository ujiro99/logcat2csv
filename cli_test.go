package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func TestRun_versionFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("./logcat2csv -version", " ")

	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}

	expected := fmt.Sprintf("logcat2csv version %s", Version)
	if !strings.Contains(errStream.String(), expected) {
		t.Errorf("expected %q to eq %q", errStream.String(), expected)
	}
}

func TestRun_No_Args(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{inStream: nil, outStream: outStream, errStream: errStream}
	args := strings.Split("./logcat2csv", " ")

	status := cli.Run(args)
	if status != ExitCodeError {
		t.Errorf("expected %d to eq %d", status, ExitCodeError)
	}
}

func TestRun_Not_File(t *testing.T) {
	file_name := "not_a_file"
	expect := "File does not exist: " + file_name + "\nPlease Enter to continue...\n"
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{inStream: nil, outStream: outStream, errStream: errStream}
	args := []string{"logcat2csv", file_name}

	status := cli.Run(args)
	if status != ExitCodeError {
		t.Errorf("expected %d to eq %d", status, ExitCodeError)
	}
	if errStream.String() != expect {
		t.Errorf("\n  result: %q\n  expect: %q", errStream.String(), expect)
	}
}

func TestRun_Exec_Stdio(t *testing.T) {
	expect := "01-01 00:00:00.000,930,931,I,tag_value,message_value\n"

	inStream := strings.NewReader("01-01 00:00:00.000   930   931 I tag_value  : message_value")
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{inStream: inStream, outStream: outStream, errStream: errStream}
	args := strings.Split("./logcat2csv", " ")

	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}
	if outStream.String() != expect {
		t.Errorf("\n  result: %q\n  expect: %q", outStream.String(), expect)
	}
}

func TestRun_encodeFlag(t *testing.T) {
	expect := []string{
		toShiftJis("01-01 00:00:00.000,930,931,I,tag_value,message_value_1"),
		toShiftJis("01-01 00:00:01.000,930,931,I,tag_value,message_value_あ亜Ａア￥凜熙♪堯"),
	}
	cli := &CLI{inStream: nil}
	args := strings.Split("./logcat2csv --encode shift-jis test/logcat_kanji.txt", " ")

	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}
	checkFile(args[3], expect, t)
}

func toShiftJis(str string) string {
	buf := new(bytes.Buffer)
	w := transform.NewWriter(buf, japanese.ShiftJIS.NewEncoder())
	w.Write([]byte(str))
	return buf.String()
}

func TestRun_Exec_Multiple_File(t *testing.T) {
	expect0 := []string{
		"01-01 00:00:00.000,930,931,I,tag_value,message_value_1",
		"01-01 00:00:01.000,930,931,I,tag_value,message_value_2",
	}
	expect1 := []string{
		"01-01 00:00:00.000,930,931,I,tag_value,message_value_3",
		"01-01 00:00:01.000,930,931,I,tag_value,message_value_4",
	}

	cli := &CLI{inStream: nil}
	args := strings.Split("./logcat2csv test/logcat.txt test/logcat2.txt", " ")

	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}
	checkFile(args[1], expect0, t)
	checkFile(args[2], expect1, t)
}

func TestRun_Exec_File_Not_File(t *testing.T) {
	expect0 := []string{
		"01-01 00:00:00.000,930,931,I,tag_value,message_value_1",
		"01-01 00:00:01.000,930,931,I,tag_value,message_value_2",
	}

	file_name := "not_a_file"
	expect := "File does not exist: " + file_name + "\nPlease Enter to continue...\n"
	errStream := new(bytes.Buffer)
	cli := &CLI{inStream: nil, errStream: errStream}
	args := []string{"logcat2csv", "test/logcat.txt", file_name}

	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}
	checkFile(args[1], expect0, t)

	if errStream.String() != expect {
		t.Errorf("\n  result: %q\n  expect: %q", errStream.String(), expect)
	}
}
