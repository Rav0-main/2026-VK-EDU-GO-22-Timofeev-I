package eval

import (
	"errors"
	"math"
	"testing"
)

func TestCalculatePositive(t *testing.T) {
	const EPS float64 = 1e-9

	tests := []struct {
		name   string
		expr   string
		answer float64
	}{
		{"Сложение двух целых чисел", "2 + 3", 5},
		{"Умножение двух целых чисел", "3 * 4", 12},
		{"Деление двух целых чисел c остатком", "5 / 2", 2.5},
		{"Деление двух целых чисел без остатка", "70 / 7", 10},
		{"Вычитание двух целых чисел", "3 - 5", -2},
		{"Скобки в начале с целыми числами", "(3 + 5 - 1) * 10 - 9", 61},
		{"Скобки в середине с целыми числами", "10 + 2 * (30 - 7) - 10", 46},
		{"Скобки в конце с целыми числами", "30 + 23 + 10 * (5 - 5)", 53},
		{"Унарный для целого числа минус без скобок", "-5 + 10", 5},
		{"Унарный для целого числа минус в скобках", "10 + (-50 + 5)", -35},
		{"Вложенность скобок с целыми числами", "10 + (50 * (20 + 10) - 30 * (25 + (35 - 10))) - 10", 0},
		{"Умножение отрицательных целых чисел", "-10 * (-10) * (-3) * (-2)", 600},
		{"Сложение двух вещественных чисел", "2.1 + 3.5", 5.6},
		{"Сложение вещественного и целого чисел", "5.2 + 6.7", 11.9},
		{"Умножение двух вещественных чисел", "3.5 * 4.5", 15.75},
		{"Умножение целого и вещественного числа", "3 * 0.5", 1.5},
		{"Деление двух вещественных чисел", "60.5 / 2.5", 24.2},
		{"Деление вещественного числа на целое", "10.5 / 2", 5.25},
		{"Вычитание двух вещественных чисел", "3.1 - 10.5", -7.4},
		{"Одно число", "+10", 10},
		{"Много пробелов", "    10 +          30\t+20               +(10-10*0)-         10", 60},
		{"Нет пробелов", "+1+2+3+4-10", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Calculate(tt.expr)
			if err != nil {
				t.Errorf("RE.\nExpected success\nbut given: %s", err)
				return
			}

			if math.Abs(tt.answer-res) > EPS {
				t.Errorf("WA.\nExpected: %f\nbut given: %f", tt.answer, res)
			}
		})
	}
}

func TestCalculateNegative(t *testing.T) {
	tests := []struct {
		name        string
		expr        string
		expectedErr error
	}{
		{"Пустое выражение", "", ErrInvalidExpression},
		{"Выражение только из пробелов", "   ", ErrInvalidExpression},
		{"Неизвестный символ", "2 + a", ErrInvalidExpression},
		{"Неизвестный спецсимвол", "10 $ 5", ErrInvalidExpression},
		{"Деление на ноль", "5 / 0", ErrDivisionByZero},
		{"Деление на ноль в скобках", "10 / (5 - 5)", ErrDivisionByZero},
		{"Отсутствует закрывающая скобка", "(2 + 3", ErrMismatchedParentheses},
		{"Отсутствует открывающая скобка", "2 + 3)", ErrMismatchedParentheses},
		{"Перепутан порядок скобок", ")2 + 3(", ErrMismatchedParentheses},
		{"Несколько точек в числе", "2.5.3 + 1", ErrInvalidExpression},
		{"Одиночная точка вместо числа", "1 + .", ErrInvalidExpression},
		{"Два бинарных оператора подряд", "2 * * 3", ErrInvalidExpression},
		{"Оператор в конце выражения", "2 + 3 +", ErrInvalidExpression},
		{"Оператор умножения в начале выражения", "* 2 + 3", ErrInvalidExpression},
		{"Пустые скобки", "()", ErrInvalidExpression},
		{"Пропущен оператор между числами", "2 3 + 4", ErrInvalidExpression},
		{"Две точки в вещественном числе", "2....5 + 3.1", ErrInvalidExpression},
		{"Унарный оператор после бинарного", "1 / +2", ErrInvalidExpression},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Calculate(tt.expr)
			if err == nil {
				t.Errorf("WA.\nExpected error matching %v\nbut got nil", tt.expectedErr)
				return
			}

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("WA.\nExpected error wrapping: %v\nbut given: %v", tt.expectedErr, err)
			}
		})
	}
}
