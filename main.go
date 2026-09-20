package main

import (
	"fmt"
	"os"
)

type ExitCode int

const (
	OK ExitCode = iota
	ErrWrongArgs
	ErrInputFilename
	ErrOutputFilename
	ErrRuntime
)

func printErr(err error) {
	fmt.Fprintf(os.Stderr, "Ошибка: %s\n", err)
}

func main() {
	config, err := GetUniqConfig(os.Args)
	if err != nil {
		printErr(err)
		os.Exit(int(ErrWrongArgs))
	}

	// поток ввода
	var fin *os.File
	if config.InputFile != "" {
		var err error
		fin, err = os.Open(config.InputFile)
		if err != nil {
			printErr(err)
			os.Exit(int(ErrInputFilename))
		}
	} else {
		fin = os.Stdin
	}

	// поток вывода
	var out *os.File
	if config.OutputFile != "" {
		var err error
		out, err = os.Create(config.OutputFile)
		if err != nil {
			printErr(err)
			os.Exit(int(ErrOutputFilename))
		}
	} else {
		out = os.Stdout
	}

	// основная логика uniq
	lineView := ViewStd
	if config.IgnoreCase {
		lineView = ViewCompose(lineView, ViewIgnoreCase)
	}
	if config.FValue != 0 {
		lineView = ViewCompose(lineView, ViewSkipFirstFFields)
	}
	if config.SValue != 0 {
		lineView = ViewCompose(lineView, ViewSkipFirstSRunes)
	}

	if err = Uniq(fin, out, lineView, config); err != nil {
		printErr(err)
		os.Exit(int(ErrRuntime))
	}
}
