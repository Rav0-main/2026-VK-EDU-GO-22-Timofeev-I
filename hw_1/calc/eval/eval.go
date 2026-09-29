// Package eval вычисляет выражение на основе предоставленной строки
package eval

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

var (
	ErrInvalidExpression     = errors.New("некорректное выражение")
	ErrDivisionByZero        = errors.New("деление на ноль")
	ErrMismatchedParentheses = errors.New("утеряны открывающие/закрывающие скобки")
)

type tokenType int

const (
	tokenNumber tokenType = iota
	tokenBinaryOperator
	tokenUnaryOperator
	tokenLeftParen
	tokenRightParen
)

type unaryOperator struct {
	precedence int
	action     func(num float64) (float64, error)
}

type binaryOperator struct {
	precedence int
	action     func(num1 float64, num2 float64) (float64, error)
}

var unaryOperators map[rune]unaryOperator = map[rune]unaryOperator{
	'+': {3, func(num float64) (float64, error) {
		return num, nil
	}},
	'-': {3, func(num float64) (float64, error) {
		return -num, nil
	}},
}

var binaryOperators map[rune]binaryOperator = map[rune]binaryOperator{
	'*': {2, func(num1 float64, num2 float64) (float64, error) {
		return num1 * num2, nil
	}},
	'/': {2, func(num1 float64, num2 float64) (float64, error) {
		if num2 == 0 {
			return 0, ErrDivisionByZero
		}
		return num1 / num2, nil
	}},
	'+': {1, func(num1 float64, num2 float64) (float64, error) {
		return num1 + num2, nil
	}},
	'-': {1, func(num1 float64, num2 float64) (float64, error) {
		return num1 - num2, nil
	}},
}

type token struct {
	typ tokenType

	// Указатели на token, т.е. token = expr[startIndex:endIndex]
	startIndex int
	endIndex   int
}

// Calculate принимает математическое выражение в виде строки и возвращает результат его вычисления.
func Calculate(expr string) (float64, error) {
	runes := []rune(expr)
	tokens, err := tokenize(runes)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrInvalidExpression, err)
	}

	if len(tokens) == 0 {
		return 0, fmt.Errorf("%w: пустое выражение", ErrInvalidExpression)
	}

	rpn, err := shuntingYard(runes, tokens)
	if err != nil {
		return 0, err
	}

	result, err := calculateRPN(runes, rpn)
	if err != nil {
		return 0, err
	}

	return result, nil
}

// tokenize разбивает строку на токены с поддержкой унарных плюсов и минусов.
func tokenize(runes []rune) ([]token, error) {
	var tokens []token
	n := len(runes)
	i := 0

	for i < n {
		ch := runes[i]

		if unicode.IsSpace(ch) {
			i++

		} else if unicode.IsDigit(ch) || ch == '.' { // Числа (целые или с плавающей точкой)
			start := i
			hasDot := false
			for i < n && (unicode.IsDigit(runes[i]) || runes[i] == '.') {
				if runes[i] == '.' {
					if hasDot {
						return nil, fmt.Errorf("множество точек в числе с начальным индексом: %d", i)
					}
					hasDot = true
				}
				i++
			}
			if i-start == 1 && runes[start] == '.' {
				return nil, fmt.Errorf("неправильное число '.' с индексом: %d", start)
			}
			tokens = append(tokens, token{typ: tokenNumber, startIndex: start, endIndex: i})

		} else if isOperator(ch) {
			if isUnaryOperator(ch) {
				isUnary := false
				if len(tokens) == 0 {
					isUnary = true
				} else {
					prevToken := tokens[len(tokens)-1]
					if prevToken.typ == tokenLeftParen {
						isUnary = true
					}
				}

				if isUnary {
					tokens = append(tokens, token{typ: tokenUnaryOperator, startIndex: i, endIndex: i + 1})
					i++
					continue
				}
			}

			tokens = append(tokens, token{typ: tokenBinaryOperator, startIndex: i, endIndex: i + 1})
			i++

		} else if ch == '(' {
			tokens = append(tokens, token{typ: tokenLeftParen, startIndex: i, endIndex: i + 1})
			i++

		} else if ch == ')' {
			tokens = append(tokens, token{typ: tokenRightParen, startIndex: i, endIndex: i + 1})
			i++

		} else {
			return nil, fmt.Errorf("неизвестный символ '%c' с индексом: %d", ch, i)
		}
	}

	return tokens, nil
}

func isOperator(ch rune) bool {
	return isUnaryOperator(ch) || isBinaryOperator(ch)
}

func isUnaryOperator(ch rune) bool {
	_, ok := unaryOperators[ch]
	return ok
}

func isBinaryOperator(ch rune) bool {
	_, ok := binaryOperators[ch]
	return ok
}

func precedence(runes []rune, t token) int {
	switch t.typ {
	case tokenUnaryOperator:
		return unaryOperators[runes[t.startIndex]].precedence

	case tokenBinaryOperator:
		return binaryOperators[runes[t.startIndex]].precedence

	default:
		return -1
	}
}

// shuntingYard преобразует список токенов в инфиксной нотации в ОПЗ (RPN).
func shuntingYard(runes []rune, tokens []token) ([]token, error) {
	var output []token
	var stack []token

	for _, t := range tokens {
		switch t.typ {
		case tokenNumber:
			output = append(output, t)

		case tokenUnaryOperator:
			stack = append(stack, t)

		case tokenBinaryOperator:
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				if top.typ == tokenLeftParen {
					break
				}

				if precedence(runes, top) >= precedence(runes, t) {
					output = append(output, top)
					stack = stack[:len(stack)-1]
				} else {
					break
				}
			}
			stack = append(stack, t)

		case tokenLeftParen:
			stack = append(stack, t)

		case tokenRightParen:
			foundLeft := false
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if top.typ == tokenLeftParen {
					foundLeft = true
					break
				}
				output = append(output, top)
			}
			if !foundLeft {
				return nil, ErrMismatchedParentheses
			}
		}
	}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if top.typ == tokenLeftParen || top.typ == tokenRightParen {
			return nil, ErrMismatchedParentheses
		}
		output = append(output, top)
	}

	return output, nil
}

// calculateRPN вычисляет значение выражения в ОПЗ (RPN).
func calculateRPN(runes []rune, tokens []token) (float64, error) {
	var stack []float64

	for _, t := range tokens {
		switch t.typ {
		case tokenNumber:
			val, err := strconv.ParseFloat(string(runes[t.startIndex:t.endIndex]), 64)
			if err != nil {
				return 0, fmt.Errorf("%w: неправильный формат вещественного числа", ErrInvalidExpression)
			}
			stack = append(stack, val)

		case tokenUnaryOperator:
			if len(stack) < 1 {
				return 0, ErrInvalidExpression
			}
			res, err := unaryOperators[runes[t.startIndex]].action(stack[len(stack)-1])
			if err != nil {
				return 0, fmt.Errorf("%w: ошибка унарного оператора %s", err, string(runes[t.startIndex:t.endIndex]))
			}
			stack[len(stack)-1] = res

		case tokenBinaryOperator:
			if len(stack) < 2 {
				return 0, ErrInvalidExpression
			}

			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			res, err := binaryOperators[runes[t.startIndex]].action(a, b)
			if err != nil {
				return 0, fmt.Errorf("%w: ошибка бинарного оператора %s", err, string(runes[t.startIndex:t.endIndex]))
			}
			stack = append(stack, res)
		}
	}

	if len(stack) != 1 {
		return 0, ErrInvalidExpression
	}

	return stack[0], nil
}
