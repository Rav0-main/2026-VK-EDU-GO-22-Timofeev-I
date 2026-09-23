package views_test

import (
	"testing"
	"uniq/views"
)

func TestIgnoreCase(t *testing.T) {
	tests := []struct {
		name   string
		str    string
		answer string
	}{
		{"Только буквы", "HeLLO", "hello"},
		{"Только цифры", "123", "123"},
		{"Цифры, буквы, спецсимволы", "MegaTesT1!", "megatest1!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := views.IgnoreCase(tt.str, nil)
			if result != tt.answer {
				t.Errorf("WA.\nExpected: %s,\nbut given: %s\n", tt.answer, result)
			}
		})
	}
}

func TestSkipFirstSRunes(t *testing.T) {
	tests := []struct {
		name   string
		str    string
		params views.Params
		answer string
	}{
		{"Без начальных пробелов", "Привет мир", views.Params{0, 3}, "вет мир"},
		{"С начальными пробелами, удаляя все пробелы", "   Пробелы", views.Params{0, 4}, "робелы"},
		{"С начальными пробелами, оставляя несколько", "  Пробелы", views.Params{0, 3}, "робелы"},
		{"Удаление всей строки", "Мало", views.Params{0, 100}, ""},
		{"Только ascii символы", "Hello", views.Params{0, 2}, "llo"},
		{"Без удаления начальных символов", "Full", views.Params{0, 0}, "Full"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := views.SkipFirstSRunes(tt.str, &tt.params)
			if result != tt.answer {
				t.Errorf("WA.\nExpected: %s,\nbut given: %s\n", tt.answer, result)
			}
		})
	}
}

func TestSkipFirstFFields(t *testing.T) {
	tests := []struct {
		name   string
		str    string
		params views.Params
		answer string
	}{
		{"Полей меньше, чем f", "1_слово 2_ое", views.Params{100, 0}, ""},
		{"Одно поле", "Поле какой-то текст", views.Params{1, 0}, "какой-то текст"},
		{"Ascii строка", "Field1 field2 text field3", views.Params{2, 0}, "text field3"},
		{"Удаление ведущих пробелов", "RM_ME               text", views.Params{1, 0}, "text"},
		{"Без удаления полей", "FULL", views.Params{0, 0}, "FULL"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := views.SkipFirstFFields(tt.str, &tt.params)
			if result != tt.answer {
				t.Errorf("WA.\nExpected: %s,\nbut given: %s\n", tt.answer, result)
			}
		})
	}
}
