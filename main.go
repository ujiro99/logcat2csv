package main

import "os"

func main() {
	cli := &CLI{inStream: os.Stdin, outStream: os.Stdout, errStream: os.Stderr}
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		cli.inStream = nil // There is no Stdin
	}
	os.Exit(cli.Run(os.Args))
}
