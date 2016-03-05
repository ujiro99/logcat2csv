package main

import (
	"os"
	"runtime"
)

func main() {
	cli := &CLI{inStream: os.Stdin, outStream: os.Stdout, errStream: os.Stderr}
	stat, err := os.Stdin.Stat()
	if (err != nil) || (stat.Mode() & os.ModeCharDevice) != 0 {
		cli.inStream = nil // There is no Stdin
	}
	os.Exit(cli.Run(os.Args, runtime.GOOS))
}
