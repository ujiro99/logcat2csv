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
	path   string
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
		if len(os.Args) != 2 {
			fmt.Println("You must drag and drop a file!")
			waitForKey()
			return ExitCodeError
		} else if !isFile(os.Args[1]) {
			fmt.Println("File does not exist!")
			waitForKey()
			return ExitCodeError
		}
		params.path = os.Args[1]
	}

	// Execute
	logcat2csv := logcat2csv{}
	return logcat2csv.Exec(params)
}

func waitForKey() {
	fmt.Println("Please Enter to continue...")
	var buf [1]byte
	os.Stdin.Read(buf[:])
}

func isFile(file string) bool {
	if s, err := os.Stat(file); s.IsDir() || err != nil {
		return false
	}
	return true
}
