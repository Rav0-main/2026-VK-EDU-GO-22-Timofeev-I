package views_test

import (
	"testing"
	"uniq/views"
)

type TestArgs struct {
	Name   string
	Str    string
	Params views.Params
	View   views.View
	Answer string
}

func runTests(tests []TestArgs, t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			result := tt.View(tt.Str, &tt.Params)
			if result != tt.Answer {
				t.Errorf("WA.\nExpected: %s,\nbut given: %s\n", tt.Answer, result)
			}
		})
	}
}

func TestIgnoreCase(t *testing.T) {
	tests := []TestArgs{
		{"Только буквы", "HeLLO", views.Params{}, views.IgnoreCase, "hello"},
		{"Только цифры", "123", views.Params{}, views.IgnoreCase, "123"},
		{"Цифры, буквы, спецсимволы", "MegaTesT1!", views.Params{}, views.IgnoreCase, "megatest1!"},
		{"Не только ascii", "ПрИВЕТ", views.Params{}, views.IgnoreCase, "привет"},
	}

	runTests(tests, t)
}

func TestSkipFirstSRunes(t *testing.T) {
	tests := []TestArgs{
		{"Без начальных пробелов", "Привет мир", views.Params{0, 3}, views.SkipFirstSRunes, "вет мир"},
		{"С начальными пробелами, удаляя все пробелы", "   Пробелы", views.Params{0, 4}, views.SkipFirstSRunes, "робелы"},
		{"С начальными пробелами, оставляя несколько", "  Пробелы", views.Params{0, 3}, views.SkipFirstSRunes, "робелы"},
		{"Удаление всей строки", "Мало", views.Params{0, 100}, views.SkipFirstSRunes, ""},
		{"Только ascii символы", "Hello", views.Params{0, 2}, views.SkipFirstSRunes, "llo"},
		{"Без удаления начальных символов", "Full", views.Params{0, 0}, views.SkipFirstSRunes, "Full"},
	}

	runTests(tests, t)
}

func TestSkipFirstFFields(t *testing.T) {
	tests := []TestArgs{
		{"Полей меньше, чем f", "1_слово 2_ое", views.Params{100, 0}, views.SkipFirstFFields, ""},
		{"Одно поле", "Поле какой-то текст", views.Params{1, 0}, views.SkipFirstFFields, "какой-то текст"},
		{"Ascii строка", "Field1 field2 text field3", views.Params{2, 0}, views.SkipFirstFFields, "text field3"},
		{"Удаление ведущих пробелов", "RM_ME               text", views.Params{1, 0}, views.SkipFirstFFields, "text"},
		{"Без удаления полей", "FULL", views.Params{0, 0}, views.SkipFirstFFields, "FULL"},
	}

	runTests(tests, t)
}

func TestCompose(t *testing.T) {
	tests := []TestArgs{
		{"Compose(Std, SkipFirstSRunes)", "Пропуск s символов", views.Params{0, 2}, views.Compose(views.Std, views.SkipFirstSRunes), "опуск s символов"},
		{"Compose(Std, SkipFirstFFields)", "Поле какое-то текст", views.Params{1, 0}, views.Compose(views.Std, views.SkipFirstFFields), "какое-то текст"},
		{"Compose(IgnoreCase, SkipFirstSRunes)", "ПРИВЕТ мир", views.Params{0, 3}, views.Compose(views.IgnoreCase, views.SkipFirstSRunes), "вет мир"},
		{"Compose(IgnoreCase, SkipFirstFFields)", "ПОЛЕ ТЕКСТ ЕЩЁ", views.Params{1, 0}, views.Compose(views.IgnoreCase, views.SkipFirstFFields), "текст ещё"},
		{"Compose трёх функций", "ПРИВЕТ МИР ТЕСТ", views.Params{1, 0}, views.Compose(views.IgnoreCase, views.SkipFirstFFields, views.Std), "мир тест"},
	}

	runTests(tests, t)
}
