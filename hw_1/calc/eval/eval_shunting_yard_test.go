package eval

import (
	"reflect"
	"testing"
)

func TestShuntingYardPositive(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		tokens []token
		answer []token
	}{
		{
			"Простая арифметика с учетом приоритета (* превосходит +)",
			"2 + 3 * 4",
			[]token{
				{tokenNumber, 0, 1},
				{tokenBinaryOperator, 2, 3},
				{tokenNumber, 4, 5},
				{tokenBinaryOperator, 6, 7},
				{tokenNumber, 8, 9},
			},
			[]token{
				{tokenNumber, 0, 1},
				{tokenNumber, 4, 5},
				{tokenNumber, 8, 9},
				{tokenBinaryOperator, 6, 7},
				{tokenBinaryOperator, 2, 3},
			},
		},
		{
			"Изменение приоритета скобками",
			"(2 + 3) * 4",
			[]token{
				{tokenLeftParen, 0, 1},
				{tokenNumber, 1, 2},
				{tokenBinaryOperator, 3, 4},
				{tokenNumber, 5, 6},
				{tokenRightParen, 6, 7},
				{tokenBinaryOperator, 8, 9},
				{tokenNumber, 10, 11},
			},
			[]token{
				{tokenNumber, 1, 2},
				{tokenNumber, 5, 6},
				{tokenBinaryOperator, 3, 4},
				{tokenNumber, 10, 11},
				{tokenBinaryOperator, 8, 9},
			},
		},
		{
			"Унарный оператор в начале выражения",
			"-3 + 5",
			[]token{
				{tokenUnaryOperator, 0, 1},
				{tokenNumber, 1, 2},
				{tokenBinaryOperator, 3, 4},
				{tokenNumber, 5, 6},
			},
			[]token{
				{tokenNumber, 1, 2},
				{tokenUnaryOperator, 0, 1},
				{tokenNumber, 5, 6},
				{tokenBinaryOperator, 3, 4},
			},
		},
		{
			"Унарный оператор перед скобкой",
			"-(3 + 5)",
			[]token{
				{tokenUnaryOperator, 0, 1},
				{tokenLeftParen, 1, 2},
				{tokenNumber, 2, 3},
				{tokenBinaryOperator, 4, 5},
				{tokenNumber, 6, 7},
				{tokenRightParen, 7, 8},
			},
			[]token{
				{tokenNumber, 2, 3},
				{tokenNumber, 6, 7},
				{tokenBinaryOperator, 4, 5},
				{tokenUnaryOperator, 0, 1},
			},
		},
		{
			"Унарный оператор внутри скобок",
			"(-3)",
			[]token{
				{tokenLeftParen, 0, 1},
				{tokenUnaryOperator, 1, 2},
				{tokenNumber, 2, 3},
				{tokenRightParen, 3, 4},
			},
			[]token{
				{tokenNumber, 2, 3},
				{tokenUnaryOperator, 1, 2},
			},
		},
		{
			"Вложенные скобки",
			"((1 + 2) * 3)",
			[]token{
				{tokenLeftParen, 0, 1},
				{tokenLeftParen, 1, 2},
				{tokenNumber, 2, 3},
				{tokenBinaryOperator, 4, 5},
				{tokenNumber, 6, 7},
				{tokenRightParen, 7, 8},
				{tokenBinaryOperator, 9, 10},
				{tokenNumber, 11, 12},
				{tokenRightParen, 12, 13},
			},
			[]token{
				{tokenNumber, 2, 3},
				{tokenNumber, 6, 7},
				{tokenBinaryOperator, 4, 5},
				{tokenNumber, 11, 12},
				{tokenBinaryOperator, 9, 10},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := shuntingYard([]rune(tt.input), tt.tokens)

			if err != nil {
				t.Errorf("RE.\nОжидается успешное выполнение\nно получено = %+v", err)
				return
			}

			if !reflect.DeepEqual(res, tt.answer) {
				t.Errorf("WA.\nОжидается: %+v\nно получено: %+v", tt.answer, res)
			}
		})
	}
}

func TestShuntingYardNegative(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		tokens []token
	}{
		{
			"Лишняя закрывающая скобка",
			"1 + 2)",
			[]token{
				{tokenNumber, 0, 1},
				{tokenBinaryOperator, 2, 3},
				{tokenNumber, 4, 5},
				{tokenRightParen, 5, 6},
			},
		},
		{
			"Незакрытая открывающая скобка",
			"(1 + 2",
			[]token{
				{tokenLeftParen, 0, 1},
				{tokenNumber, 1, 2},
				{tokenBinaryOperator, 3, 4},
				{tokenNumber, 5, 6},
			},
		},
		{
			"Перепутан порядок скобок",
			")1 + 2(",
			[]token{
				{tokenRightParen, 0, 1},
				{tokenNumber, 1, 2},
				{tokenBinaryOperator, 3, 4},
				{tokenNumber, 5, 6},
				{tokenLeftParen, 6, 7},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := shuntingYard([]rune(tt.input), tt.tokens)

			if err == nil {
				t.Errorf("WA.\nОжидается ошибочное выполнение\nно получено успешное выполнение")
			}
		})
	}
}
