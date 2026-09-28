// Package eval вычисляет выражение на основе предоставленной строки
package eval

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

var (
	ErrInvalidExpression     = errors.New("invalid expression")
	ErrDivisionByZero        = errors.New("division by zero")
	ErrMismatchedParentheses = errors.New("mismatched parentheses")
)

// TokenType определяет тип токена в выражении.
type tokenType int

// TODO: добавить tokenUnaryMinus и tokenUnaryPlus, заменить tokenOperator на tokenBinaryOperator
const (
	tokenNumber tokenType = iota
	tokenOperator
	tokenLeftParen
	tokenRightParen
)

type token struct {
	typ tokenType

	// TODO: может быть стоит оптимизировать память: хранить (start_i, end_i)
	value string
}

// Calculate принимает математическое выражение в виде строки и возвращает результат его вычисления.
func Calculate(expr string) (float64, error) {
	tokens, err := tokenize(expr)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrInvalidExpression, err)
	}

	if len(tokens) == 0 {
		return 0, fmt.Errorf("%w: empty expression", ErrInvalidExpression)
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

// tokenize разбивает строку на токены с поддержкой унарных минусов.
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
						return nil, fmt.Errorf("multiple dots in number at index %d", i)
					}
					hasDot = true
				}
				i++
			}
			numStr := string(runes[start:i])
			if numStr == "." {
				return nil, fmt.Errorf("invalid number '.' at index %d", start)
			}
			tokens = append(tokens, token{typ: tokenNumber, value: numStr})
			continue
		}

		// Операторы
		if isOperator(ch) {
			// Проверка на унарный минус:
			// Минус является унарным, если он идет первым или сразу после '(' или другого оператора.
			if ch == '-' {
				isUnary := false
				if len(tokens) == 0 {
					isUnary = true
				} else {
					prevToken := tokens[len(tokens)-1]
					if prevToken.typ == tokenLeftParen || prevToken.typ == tokenOperator {
						isUnary = true
					}
				}

				if isUnary {
					tokens = append(tokens, token{typ: tokenOperator, value: "u-"})
					i++
					continue
				}
			}

			// TODO: сделать унарный плюс

			tokens = append(tokens, token{typ: tokenOperator, value: string(ch)})
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

		return nil, fmt.Errorf("unexpected character '%c' at index %d", ch, i)
	}

	return tokens, nil
}

func isOperator(ch rune) bool {
	return ch == '+' || ch == '-' || ch == '*' || ch == '/'
}

func precedence(op string) int {
	// TODO: может быть стоит это в map?
	switch op {
	case "u-":
		return 3 // Самый высокий приоритет для унарного минуса
	case "*", "/":
		return 2
	case "+", "-":
		return 1
	default:
		return 0
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

		case tokenOperator:
			for len(stack) > 0 && stack[len(stack)-1].typ == tokenOperator {
				top := stack[len(stack)-1]
				// Унарный минус правоассоциативен, бинарные - левоассоциативны
				if (t.value != "u-" && precedence(t.value) <= precedence(top.value)) ||
					(t.value == "u-" && precedence(t.value) < precedence(top.value)) {
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
				return 0, fmt.Errorf("%w: invalid float format", ErrInvalidExpression)
			}
			stack = append(stack, val)

		case tokenOperator:
			if t.value == "u-" {
				if len(stack) < 1 {
					return 0, ErrInvalidExpression
				}
				stack[len(stack)-1] = -stack[len(stack)-1]
				continue
			}

			if len(stack) < 2 {
				return 0, ErrInvalidExpression
			}

			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var res float64
			switch t.value {
			case "+":
				res = a + b
			case "-":
				res = a - b
			case "*":
				res = a * b
			case "/":
				if b == 0 {
					return 0, ErrDivisionByZero
				}
				res = a / b
			}
			stack = append(stack, res)
		}
	}

	if len(stack) != 1 {
		return 0, ErrInvalidExpression
	}

	return stack[0], nil
}
