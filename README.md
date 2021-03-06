# logcat2csv [![Build Status](https://travis-ci.org/ujiro99/logcat2csv.svg?branch=master)](https://travis-ci.org/ujiro99/logcat2csv)  [![Coverage Status](https://coveralls.io/repos/github/ujiro99/logcat2csv/badge.svg?branch=master)](https://coveralls.io/github/ujiro99/logcat2csv?branch=master)

Command line tool for convert logcat to csv.

## SYNOPSIS

```
logcat2csv is tool for convert logcat to csv.

Usage:
  logcat2csv [options] PATH ...

Options:
  --encode, -e   Charactor encoding of output file.
  --version      Show version.
  --help         Show this help.
```


## Install

You can get binary from github release page.

[-> Release Page](https://github.com/ujiro99/logcat2csv/releases)

or, use `go get`:

```bash
$ go get -d github.com/ujiro99/logcat2csv
```

## Contribution

1. Fork ([https://github.com/ujiro99/logcat2csv/fork](https://github.com/ujiro99/logcat2csv/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[ujiro99](https://github.com/ujiro99)
