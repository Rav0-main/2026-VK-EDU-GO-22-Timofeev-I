package main

import (
	"strings"
	"unicode"
)

type ViewParams struct {
	FValue int
	SValue int
}

type View func(str string, params *ViewParams) string

func ViewCompose(views ...View) View {
	if len(views) == 0 {
		panic("Аргумент 'views' пустой")
	}

	return func(str string, params *ViewParams) string {
		for _, view := range views {
			str = view(str, params)
		}
		return str
	}
}

// ViewStd возвращает саму строку
func ViewStd(str string, params *ViewParams) string {
	return str
}

// ViewIgnoreCase возвращает строку в нижним регистре
func ViewIgnoreCase(str string, params *ViewParams) string {
	return strings.ToLower(str)
}

// ViewSkipFirstSRunes возвращает строку без первых s символов.
// Если строка меньше, то возвращается пустая строка
func ViewSkipFirstSRunes(str string, params *ViewParams) string {
	if params.SValue < 0 {
		panic("Аргумент 'SValue' должен быть int >= 0")
	}

	if params.SValue >= len(str) {
		return ""
	}

	return str[params.SValue:]
}

// ViewSkipFirstFFields возвращает строку без первых f слов.
// Если строка содержит слов меньше, то возвращает пустую строку
func ViewSkipFirstFFields(str string, params *ViewParams) string {
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
