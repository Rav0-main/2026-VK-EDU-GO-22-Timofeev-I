package eval

import (
	"reflect"
	"testing"
)

func TestTokenizePositive(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		answer []token
	}{
		{
			"Простые целые числа и бинарный плюс",
			"12 + 345",
			[]token{
				{tokenNumber, 0, 2},
				{tokenBinaryOperator, 3, 4},
				{tokenNumber, 5, 8},
			},
		},
		{
			"Числа с плавающей точкой",
			"3.14 + .5",
			[]token{
				{tokenNumber, 0, 4},
				{tokenBinaryOperator, 5, 6},
				{tokenNumber, 7, 9},
			},
		},
		{
			"Унарный оператор в начале выражения",
			"-5",
			[]token{
				{tokenUnaryOperator, 0, 1},
				{tokenNumber, 1, 2},
			},
		},
		{
			"Унарный оператор после открывающей скобки",
			"(-5 + 2)",
			[]token{
				{tokenLeftParen, 0, 1},
				{tokenUnaryOperator, 1, 2},
				{tokenNumber, 2, 3},
				{tokenBinaryOperator, 4, 5},
				{tokenNumber, 6, 7},
				{tokenRightParen, 7, 8},
			},
		},
		{
			"Бинарный минус (не унарный)",
			"10 - 5",
			[]token{
				{tokenNumber, 0, 2},
				{tokenBinaryOperator, 3, 4},
				{tokenNumber, 5, 6},
			},
		},
		{
			"Пробелы и табуляции игнорируются",
			"  1   + \t 2 ",
			[]token{
				{tokenNumber, 2, 3},
				{tokenBinaryOperator, 6, 7},
				{tokenNumber, 10, 11},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := tokenize([]rune(tt.input))

			if err != nil {
				t.Errorf("RE.\nОжидается успешное выполнение\nно получено: %+v", err)
				return
			}

			if !reflect.DeepEqual(res, tt.answer) {
				t.Errorf("WA.\nОжидается: %v+\nно получено: %+v", tt.answer, res)
			}
		})
	}
}

func TestTokenizeNegative(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			"Несколько точек в одном числе",
			"1.2.3",
		},
		{
			"Одиночная точка без цифр",
			"5 + .",
		},
		{
			"Неизвестный символ",
			"12 @ 34",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tokenize([]rune(tt.input))

			if err == nil {
				t.Errorf("WA.\nОжидается ошибочное выполнение\nно получено успешное выполнение")
			}
		})
	}
}
