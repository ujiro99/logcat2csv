package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var pathsInDir = []string{
	"test/logcat.txt",
	"test/logcat2.txt",
	"test/logcat.threadtime.txt",
	"test/logcat_kanji.txt",
	"test/logcat_not_shiftjis.txt",
}

func convertTo(str string, encode string) string {
	buf := new(bytes.Buffer)
	var w io.Writer
	switch encode {
	case ShiftJIS:
		w = transform.NewWriter(buf, japanese.ShiftJIS.NewEncoder())
	case EUCJP:
		w = transform.NewWriter(buf, japanese.EUCJP.NewEncoder())
	case ISO2022JP:
		w = transform.NewWriter(buf, japanese.ISO2022JP.NewEncoder())
	default:
		w = buf
	}
	w.Write([]byte(str))
	return buf.String()
}

func TestRun_versionFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("./logcat2csv -version", " ")

	status := cli.Run(args, "")
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}

	expected := fmt.Sprintf("logcat2csv version %s", Version)
	if !strings.Contains(errStream.String(), expected) {
		t.Errorf("expected %q to eq %q", errStream.String(), expected)
	}
}

func TestRun_No_Args(t *testing.T) {
	outStream := new(bytes.Buffer)
	inStream := strings.NewReader("\n")
	cli := &CLI{inStream: inStream, outStream: outStream}
	args := strings.Split("./logcat2csv", " ")

	status := cli.Run(args, "")
	if status != ExitCodeError {
		t.Errorf("expected %d to eq %d", status, ExitCodeError)
	}
}

func TestRun_Not_File(t *testing.T) {
	fileName := "not_a_file"
	expect := "File does not exist: " + fileName + "\nTarget not found.\n"
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{inStream: nil, outStream: outStream, errStream: errStream}
	args := []string{"logcat2csv", fileName}

	status := cli.Run(args, "")
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
	outStream := new(bytes.Buffer)
	cli := &CLI{inStream: inStream, outStream: outStream}
	args := strings.Split("./logcat2csv", " ")

	status := cli.Run(args, "")
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}
	if outStream.String() != expect {
		t.Errorf("\n  result: %q\n  expect: %q", outStream.String(), expect)
	}
}

func TestRun_encodeFlag(t *testing.T) {
	expect := []string{
		convertTo("01-01 00:00:00.000,930,931,I,tag_value,message_value_1", ShiftJIS),
		convertTo("01-01 00:00:01.000,930,931,I,tag_value,message_value_あ亜Ａア￥凜熙♪堯", ShiftJIS),
	}
	cli := &CLI{inStream: nil}
	args := strings.Split("./logcat2csv --encode shift-jis test/logcat_kanji.txt", " ")

	status := cli.Run(args, "")
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}
	err := checkFile(args[3], expect)
	if err != nil {
		t.Error(err)
	}
}

func TestRun_encodeFlag_output_with_utf8_if_encoding_failed(t *testing.T) {
	expect := []string{
		convertTo("01-01 00:00:01.000,930,931,I,tag_value,message_value_あ亜Ａア￥凜熙♪堯", ShiftJIS),
		convertTo("01-01 00:00:01.000,930,931,I,tag_value,\"AddressBook Labels [en-US]: [, A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S, T, U, V, W, X, Y, Z, Α, Β, Γ, Δ, Ε, Ζ, Η, Θ, Ι, Κ, Λ, Μ, Ν, Ξ, Ο, Π, Ρ, Σ, Τ, Υ, Φ, Χ, Ψ, Ω, , А, Б, В, Г, Д, Ђ, Е, Є, Ж, З, И, І, Й, Ј, К, Л, Љ, М, Н, Њ, О, П, Р, С, Т, Ћ, У, Ф, Х, Ц, Ч, Џ, Ш, Щ, Ю, Я, , א, ב, ג, ד, ה, ו, ז, ח, ט, י, כ, ל, מ, נ, ס, ע, פ, צ, ק, ר, ש, ת, , ا, ب, ت, ث, ج, ح, خ, د, ذ, ر, ز, س, ش, ص, ض, ط, ظ, ع, غ, ف, ق, ك, ل, م, ن, ه, و, ي, , ก, ข, ฃ, ค, ฅ, ฆ, ง, จ, ฉ, ช, ซ, ฌ, ญ, ฎ, ฏ, ฐ, ฑ, ฒ, ณ, ด, ต, ถ, ท, ธ, น, บ, ป, ผ, ฝ, พ, ฟ, ภ, ม, ย, ร, ฤ, ล, ฦ, ว, ศ, ษ, ส, ห, ฬ, อ, ฮ, , ㄱ, ㄴ, ㄷ, ㄹ, ㅁ, ㅂ, ㅅ, ㅇ, ㅈ, ㅊ, ㅋ, ㅌ, ㅍ, ㅎ, , あ, か, さ, た, な, は, ま, や, ら, わ, #, ]\"", UTF8),
		convertTo("01-01 00:00:01.000,930,931,I,tag_value,\"AddressBook Labels [en-US]: [, A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S, T, U, V, W, X, Y, Z]\"", ShiftJIS),
	}
	cli := &CLI{inStream: nil}
	args := strings.Split("./logcat2csv --encode shift-jis test/logcat_not_shiftjis.txt", " ")

	status := cli.Run(args, "")
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}
	err := checkFile(args[3], expect)
	if err != nil {
		t.Error(err)
	}
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

	status := cli.Run(args, "")
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}
	if err := checkFile(args[1], expect0); err != nil {
		t.Error(err)
	}
	if err := checkFile(args[2], expect1); err != nil {
		t.Error(err)
	}
}

func TestRun_Exec_File_Not_File(t *testing.T) {
	expect0 := []string{
		"01-01 00:00:00.000,930,931,I,tag_value,message_value_1",
		"01-01 00:00:01.000,930,931,I,tag_value,message_value_2",
	}

	fileName := "not_a_file"
	expect := "File does not exist: " + fileName + "\n"
	errStream := new(bytes.Buffer)
	cli := &CLI{inStream: nil, errStream: errStream}
	args := []string{"logcat2csv", "test/logcat.txt", fileName}

	status := cli.Run(args, "")
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}
	if err := checkFile(args[1], expect0); err != nil {
		t.Error(err)
	}

	if errStream.String() != expect {
		t.Errorf("\n  result: %q\n  expect: %q", errStream.String(), expect)
	}
}

func TestRun_Exec_Directory(t *testing.T) {
	cli := &CLI{inStream: nil}
	args := []string{"logcat2csv", "./test"}

	status := cli.Run(args, "")
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}

	for _, path := range pathsInDir {
		if err := checkFile(path, nil); err != nil {
			t.Error(err)
		}
	}
}

func TestRun_Exec_Directory_Ignore_if_not_logcat_file(t *testing.T) {
	ignorePaths := []string{
		"test/logcat.raw.txt",
		"test/notLogcat.txt",
	}
	cli := &CLI{inStream: nil}
	args := []string{"logcat2csv", "./test"}

	status := cli.Run(args, "")
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}

	for _, path := range ignorePaths {
		if s, err := os.Stat(path + ".csv"); err == nil && !s.IsDir() {
			t.Error(path + ".csv is created.")
		}
		continue
	}

	// clean generated files
	for _, path := range pathsInDir {
		checkFile(path, nil)
	}
}

func TestRun_Exec_Directory_ignore_if_csv_already_exists(t *testing.T) {
	ignorePaths := []string{
		"test/ignore.txt",
	}
	cli := &CLI{inStream: nil}
	args := []string{"logcat2csv", "./test"}

	status := cli.Run(args, "")
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}

	for _, path := range ignorePaths {
		if _, err := os.Stat(path + ".csv"); err != nil {
			t.Error(path + ".csv is deleted.")
		}
		continue
	}

	// clean generated files
	for _, path := range pathsInDir {
		checkFile(path, nil)
	}
}
