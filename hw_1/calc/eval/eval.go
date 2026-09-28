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
	'+': unaryOperator{3, func(num float64) (float64, error) {
		return num, nil
	}},
	'-': unaryOperator{3, func(num float64) (float64, error) {
		return -num, nil
	}},
}

var binaryOperators map[rune]binaryOperator = map[rune]binaryOperator{
	'*': binaryOperator{2, func(num1 float64, num2 float64) (float64, error) {
		return num1 * num2, nil
	}},
	'/': binaryOperator{2, func(num1 float64, num2 float64) (float64, error) {
		if num2 == 0 {
			return 0, ErrDivisionByZero
		}
		return num1 / num2, nil
	}},
	'+': binaryOperator{1, func(num1 float64, num2 float64) (float64, error) {
		return num1 + num2, nil
	}},
	'-': binaryOperator{1, func(num1 float64, num2 float64) (float64, error) {
		return num1 - num2, nil
	}},
}

type token struct {
	typ tokenType

	// TODO: оптимизировать двумя указателями: (start_i, end_i)
	value string
}

// Calculate принимает математическое выражение в виде строки и возвращает результат его вычисления.
func Calculate(expr string) (float64, error) {
	tokens, err := tokenize(expr)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrInvalidExpression, err)
	}

	if len(tokens) == 0 {
		return 0, fmt.Errorf("%w: пустое выражение", ErrInvalidExpression)
	}

	rpn, err := shuntingYard(tokens)
	if err != nil {
		return 0, err
	}

	result, err := evalRPN(rpn)
	if err != nil {
		return 0, err
	}

	return result, nil
}

// tokenize разбивает строку на токены с поддержкой унарных плюсов и минусов.
func tokenize(expr string) ([]token, error) {
	var tokens []token
	runes := []rune(expr)
	n := len(runes)
	i := 0

	for i < n {
		ch := runes[i]

		if unicode.IsSpace(ch) {
			i++
			continue
		}

		// Числа (целые или с плавающей точкой)
		if unicode.IsDigit(ch) || ch == '.' {
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
			numStr := string(runes[start:i])
			if numStr == "." {
				return nil, fmt.Errorf("неправильное число '.' с индексом: %d", start)
			}
			tokens = append(tokens, token{typ: tokenNumber, value: numStr})
			continue
		}

		// Операторы
		if isOperator(ch) {
			// Проверка на унарный плюс/минус:
			// Оператор унарный, если он идет первым, сразу после '('.
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
					tokens = append(tokens, token{typ: tokenUnaryOperator, value: string(ch)})
					i++
					continue
				}
			}

			tokens = append(tokens, token{typ: tokenBinaryOperator, value: string(ch)})
			i++
			continue
		}

		// Скобки
		if ch == '(' {
			tokens = append(tokens, token{typ: tokenLeftParen, value: "("})
			i++
			continue
		}
		if ch == ')' {
			tokens = append(tokens, token{typ: tokenRightParen, value: ")"})
			i++
			continue
		}

		return nil, fmt.Errorf("неизвестный символ '%c' с индексом: %d", ch, i)
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

func precedence(t token) int {
	switch t.typ {
	case tokenUnaryOperator:
		return unaryOperators[rune(t.value[0])].precedence

	case tokenBinaryOperator:
		return binaryOperators[rune(t.value[0])].precedence

	default:
		return -1
	}
}

// shuntingYard преобразует список токенов в инфиксной нотации в ОПЗ (RPN).
func shuntingYard(tokens []token) ([]token, error) {
	var output []token
	var stack []token

	for _, t := range tokens {
		switch t.typ {
		case tokenNumber:
			output = append(output, t)

		case tokenUnaryOperator:
			// Унарный оператор правоассоциативен — просто помещаем его в стек
			stack = append(stack, t)

		case tokenBinaryOperator:
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				if top.typ == tokenLeftParen {
					break
				}

				// Выталкиваем из стека операторы с большим или равным приоритетом
				if precedence(top) >= precedence(t) {
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

// evalRPN вычисляет значение выражения в ОПЗ (RPN).
func evalRPN(tokens []token) (float64, error) {
	var stack []float64

	for _, t := range tokens {
		switch t.typ {
		case tokenNumber:
			val, err := strconv.ParseFloat(t.value, 64)
			if err != nil {
				return 0, fmt.Errorf("%w: неправильный формат вещественного числа", ErrInvalidExpression)
			}
			stack = append(stack, val)

		case tokenUnaryOperator:
			if len(stack) < 1 {
				return 0, ErrInvalidExpression
			}
			res, err := unaryOperators[rune(t.value[0])].action(stack[len(stack)-1])
			if err != nil {
				return 0, fmt.Errorf("%w: ошибка унарного оператора %s", err, t.value)
			}
			stack[len(stack)-1] = res

		case tokenBinaryOperator:
			if len(stack) < 2 {
				return 0, ErrInvalidExpression
			}

			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			res, err := binaryOperators[rune(t.value[0])].action(a, b)
			if err != nil {
				return 0, fmt.Errorf("%w: ошибка бинарного оператора %s", err, t.value)
			}
			stack = append(stack, res)
		}
	}

	if len(stack) != 1 {
		return 0, ErrInvalidExpression
	}

	return stack[0], nil
}
