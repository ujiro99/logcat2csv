package main

import (
	"flag"
	"fmt"
	"io"
	"os"
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
	reader io.Reader
	writer io.Writer
	encode string
	paths  []string
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
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
	params := cmdParams{}
	params.encode = encode
	if cli.inStream != nil {
		params.reader = cli.inStream
		params.writer = cli.outStream
	} else {
		if len(flags.Args()) <= 0 {
			fmt.Fprintf(cli.errStream, "You must specify a file!\n")
			cli.waitForKey()
			return ExitCodeError
		}
		paths := make([]string, len(flags.Args()))
		i := 0
		for _, path := range flags.Args() {
			if !isFile(path) {
				fmt.Fprintf(cli.errStream, "File does not exist: %s\n", path)
				cli.waitForKey()
			} else {
				paths[i] = path
				i = i + 1
			}
		}
		if i == 0 {
			return ExitCodeError
		}
		params.paths = paths[:i]
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

func isFile(file string) bool {
	if s, err := os.Stat(file); err != nil || s.IsDir() {
		return false
	}
	return true
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
