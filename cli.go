package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int    = 0
	ExitCodeError int    = 1 + iota
	Name          string = "logcat2csv"
)

var (
	// Version represents this version.
	Version = "0.1.0"
)

// CLI is the command line object
type CLI struct {
	inStream             io.Reader
	outStream, errStream io.Writer
}

type cmdParams struct {
	reader         io.Reader
	writer, error  io.Writer
	encode, osName string
	paths          []string
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string, osName string) int {
	var (
		encode  string
		version bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)
	flags.Usage = func() { fmt.Fprintf(cli.outStream, helpText) }
	flags.StringVar(&encode, "encode", "", "charactor encoding of output file")
	flags.StringVar(&encode, "e", "", "charactor encoding of output file(Short)")
	flags.BoolVar(&version, "version", false, "Print version information and quit.")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	// Show version
	if version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	// Parse arguments
	params := cmdParams{
		error:  cli.errStream,
		encode: encode,
		osName: osName,
	}
	if cli.inStream != nil {
		params.reader = cli.inStream
		params.writer = cli.outStream
	} else {
		params.paths = cli.expandArgs(flags.Args())
		if len(params.paths) <= 0 {
			return ExitCodeError
		}
	}

	// Execute
	logcat2csv := logcat2csv{}
	return logcat2csv.Exec(params)
}

func (cli *CLI) waitForKey() {
	fmt.Fprintf(cli.errStream, "Please Enter to continue...\n")
	var buf [1]byte
	os.Stdin.Read(buf[:])
}

func (cli *CLI) listFiles(dirName string) []string {
	fileInfos, err := ioutil.ReadDir(dirName)
	if err != nil {
		return []string{}
	}
	files := make([]string, len(fileInfos))
	i := 0
	for _, fileInfo := range fileInfos {
		// don't list recursively
		if !fileInfo.IsDir() {
			files[i] = filepath.Join(dirName, fileInfo.Name())
			i = i + 1
		}
	}
	return files[:i]
}

func (cli *CLI) expandArgs(args []string) []string {
	if len(args) <= 0 {
		fmt.Fprintf(cli.errStream, "Please specify a file, or drag & drop to icon.\n")
		cli.waitForKey()
		return []string{}
	}

	// list and validate filepaths.
	files := make([]string, len(args))
	pathMap := map[string][]string{
		"": files,
	}
	total := 0 // count of total files.
	count := 0 // count of files specified directlly.
	for _, path := range args {
		if isDir(path) {
			pathMap[path] = cli.listFiles(path)
			total = total + len(pathMap[path])
		} else if isFile(path) {
			files[count] = path
			total = total + 1
			count = count + 1
		} else {
			fmt.Fprintf(cli.errStream, "File does not exist: %s\n", path)
		}
	}
	pathMap[""] = files[:count]
	filePaths := make([]string, total)
	i := 0
	for _, paths := range pathMap {
		for _, p := range paths {
			filePaths[i] = p
			i = i + 1
		}
	}
	return filePaths
}

func isFile(file string) bool {
	if s, err := os.Stat(file); err != nil || s.IsDir() {
		return false
	}
	return true
}

func isDir(file string) bool {
	if s, err := os.Stat(file); err == nil && s.IsDir() {
		return true
	}
	return false
}

var helpText = `logcat2csv is tool for convert logcat to csv.

https://github.com/ujiro99/logcat2csv

Usage:
  logcat2csv [options] PATH ...

Options:
  --encode, -e   Charactor encoding of output file.
  --version      Show version.
  --help         Show this help.
`
