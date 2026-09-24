package main

import (
	"fmt"
	"os"
	"uniq/views"
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
	config, err := ParseArgsToConfig(os.Args[1:])
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
	lineView := views.Std
	if config.IgnoreCase {
		lineView = views.Compose(lineView, views.IgnoreCase)
	}
	if config.FValue != 0 {
		lineView = views.Compose(lineView, views.SkipFirstFFields)
	}
	if config.SValue != 0 {
		lineView = views.Compose(lineView, views.SkipFirstSRunes)
	}

	if err = Uniq(fin, out, lineView, config); err != nil {
		printErr(err)
		os.Exit(int(ErrRuntime))
	}
}
