// Package views содержит функции для изменения строки
package views

import (
	"strings"
	"unicode"
)

type Params struct {
	FValue int
	SValue int
}

type View func(str string, params *Params) string

// Compose возвращает View, которая является композицией views...
// Порядок композиции: вначале views0, далее views1 и т.д.
func Compose(views ...View) View {
	if len(views) == 0 {
		panic("Аргумент 'views' пустой")
	}

	return func(str string, params *Params) string {
		for _, view := range views {
			str = view(str, params)
		}
		return str
	}
}

// Std возвращает саму строку
func Std(str string, params *Params) string {
	return str
}

// IgnoreCase возвращает строку в нижним регистре
func IgnoreCase(str string, params *Params) string {
	return strings.ToLower(str)
}

// SkipFirstSRunes возвращает строку без первых s символов.
// Если строка меньше, то возвращается пустая строка
func SkipFirstSRunes(str string, params *Params) string {
	if params.SValue < 0 {
		panic("Аргумент 'SValue' должен быть int >= 0")
	}

	count := 0
	for i := range str {
		if count == params.SValue {
			return str[i:]
		}
		count++
	}
	return ""
}

// SkipFirstFFields возвращает строку без первых f слов.
// Если строка содержит слов меньше, то возвращает пустую строку
func SkipFirstFFields(str string, params *Params) string {
	if params.FValue < 0 {
		panic("Аргумент 'FValue' должен быть int >= 0")
	}

	runes := []rune(str)
	i := 0

	for range params.FValue {
		for i < len(runes) && unicode.IsSpace(runes[i]) {
			i++
		}

		for i < len(runes) && !unicode.IsSpace(runes[i]) {
			i++
		}

		if i == len(runes) {
			return ""
		}
	}

	for i < len(runes) && unicode.IsSpace(runes[i]) {
		i++
	}

	return string(runes[i:])
}
