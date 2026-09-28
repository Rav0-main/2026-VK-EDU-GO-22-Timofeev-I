package main

import (
	"calc/eval"
	"fmt"
)

func main() {
	s := "1 + 2 / 2 + (-5 * 1)"

	r, _ := eval.Calculate(s)

	fmt.Println(r)
}
