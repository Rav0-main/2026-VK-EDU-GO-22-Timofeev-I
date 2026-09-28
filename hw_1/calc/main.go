package main

import (
	"bufio"
	"calc/eval"
	"fmt"
	"os"
)

func printErr(err error) {
	fmt.Fprintf(os.Stderr, "Ошибка: %s\n", err)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	if scanner.Scan() {
		input := scanner.Text()
		res, err := eval.Calculate(input)

		if err != nil {
			printErr(err)
			os.Exit(1)
		}

		fmt.Println(res)
	}

	if err := scanner.Err(); err != nil {
		printErr(err)
		os.Exit(1)
	}
}
