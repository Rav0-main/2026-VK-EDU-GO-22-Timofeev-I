package main

import (
	"bufio"
	"fmt"
	"io"
	"uniq/views"
)

type Config struct {
	Mode       string // сделано string для удобства парсинга
	IgnoreCase bool
	InputFile  string
	OutputFile string
	views.Params
}

func uniqFprintln(out io.Writer, count uint, line string, config *Config) error {
	var err error
	switch config.Mode {
	case "":
		_, err = fmt.Fprintln(out, line)
	case "c":
		_, err = fmt.Fprintf(out, "%d %s\n", count, line)
	case "d":
		if count > 1 {
			_, err = fmt.Fprintln(out, line)
		}
	case "u":
		if count == 1 {
			_, err = fmt.Fprintln(out, line)
		}
	default:
	}
	return err
}

func Uniq(istream io.Reader, ostream io.Writer, lineView views.View, config *Config) error {
	// Передача входных данных и выходных данных сделана из io, так как
	// это позволяет не считывать весь файл в память.

	if config == nil {
		panic("Аргумент 'config' не должен быть nil")
	}

	in := bufio.NewScanner(istream)
	var lineCount uint = 0
	var prevLineView string
	var prevLine string
	if in.Scan() {
		prevLine = in.Text()
		prevLineView = lineView(prevLine, &config.Params)
		lineCount++
	}

	for in.Scan() {
		currentLine := in.Text()
		currentLineView := lineView(currentLine, &config.Params)

		if currentLineView != prevLineView {
			err := uniqFprintln(ostream, lineCount, prevLine, config)
			if err != nil {
				return err
			}

			lineCount = 1
			prevLineView = currentLineView
			prevLine = currentLine
		} else {
			lineCount++
		}
	}
	err := uniqFprintln(ostream, lineCount, prevLine, config)
	if err != nil {
		return err
	}

	return in.Err()
}
