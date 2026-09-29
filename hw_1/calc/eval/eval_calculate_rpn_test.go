package eval

import (
	"math"
	"testing"
)

func TestCalculateRPNPositive(t *testing.T) {
	const EPS float64 = 1e-9

	tests := []struct {
		name   string
		input  string
		tokens []token
		answer float64
	}{
		{
			"Сложение и умножение: 2 + 3 * 4 -> [2, 3, 4, *, +]",
			"2 + 3 * 4",
			[]token{
				{tokenNumber, 0, 1},
				{tokenNumber, 4, 5},
				{tokenNumber, 8, 9},
				{tokenBinaryOperator, 6, 7},
				{tokenBinaryOperator, 2, 3},
			},
			14,
		},
		{
			"Приоритет скобок: (2 + 3) * 4 -> [2, 3, +, 4, *]",
			"(2 + 3) * 4",
			[]token{
				{tokenNumber, 1, 2},
				{tokenNumber, 5, 6},
				{tokenBinaryOperator, 3, 4},
				{tokenNumber, 10, 11},
				{tokenBinaryOperator, 8, 9},
			},
			20,
		},
		{
			"Унарный минус перед скобками: -(3 + 5) -> [3, 5, +, -]",
			"-(3 + 5)",
			[]token{
				{tokenNumber, 2, 3},
				{tokenNumber, 6, 7},
				{tokenBinaryOperator, 4, 5},
				{tokenUnaryOperator, 0, 1},
			},
			-8,
		},
		{
			"Числа с плавающей точкой: 3.5 / 2 -> [3.5, 2, /]",
			"3.5 / 2",
			[]token{
				{tokenNumber, 0, 3},
				{tokenNumber, 6, 7},
				{tokenBinaryOperator, 4, 5},
			},
			1.75,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := calculateRPN([]rune(tt.input), tt.tokens)

			if err != nil {
				t.Errorf("RE.\nОжидается успешное выполнение\nно получено: %v", err)
				return
			}

			if math.Abs(res-tt.answer) > EPS {
				t.Errorf("WA.\nОжидается: %+v\nно получено: %+v", tt.answer, res)
			}
		})
	}
}

func TestCalculateRPNNegative(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		tokens []token
	}{
		{
			"Недостаточно операндов для бинарного оператора",
			"+",
			[]token{
				{tokenBinaryOperator, 0, 1},
			},
		},
		{
			"Недостаточно операндов для унарного оператора",
			"-",
			[]token{
				{tokenUnaryOperator, 0, 1},
			},
		},
		{
			"Лишнее число в стеке (не хватило оператора)",
			"2 3",
			[]token{
				{tokenNumber, 0, 1},
				{tokenNumber, 2, 3},
			},
		},
		{
			"Пустой список токенов",
			"",
			[]token{},
		},
		{
			"Ошибка парсинга числа",
			"abc",
			[]token{
				{tokenNumber, 0, 3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := calculateRPN([]rune(tt.input), tt.tokens)

			if err == nil {
				t.Errorf("WA.\nОжидается ошибочное выполнение\nно получено успешное выполнение")
			}
		})
	}
}
