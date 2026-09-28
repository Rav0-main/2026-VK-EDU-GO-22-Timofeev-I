package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"uniq/views"
)

func TestUniq(t *testing.T) {
	type TestCase struct {
		Name     string
		Input    io.Reader
		LineView views.View
		Config   Config
	}
	testcases := []TestCase{
		{"Проверка без параметров", strings.NewReader("Hello\nHello"), views.Std, Config{}},
		{"Проверка -i", strings.NewReader("Hello\nHeLLO\nhello\n"), views.IgnoreCase, Config{IgnoreCase: true}},
		{"Проверка -c", strings.NewReader("Hello\nHello\nVAR\nVAR\nVAR\nV"), views.Std, Config{Mode: "c"}},
		{"Проверка -i, -c", strings.NewReader("Hello\nHELLO\nv\nV"), views.IgnoreCase, Config{Mode: "c", IgnoreCase: true}},
		{"Проверка -d", strings.NewReader("Hello\nHello\nXR\nxr\nSolo\n"), views.Std, Config{Mode: "d"}},
		{"Проверка -u", strings.NewReader("Hello\nXR\nXR\nXR\nHello"), views.Std, Config{Mode: "u"}},
		{"Проверка -f 1", strings.NewReader("Hello world\nLINE world\nPRIVET world\nWORLD woRLD"), views.SkipFirstFFields, Config{FValue: 1}},
		{"Проверка -s 1", strings.NewReader("HELLO\nhELLO\nNELLO\nXELLO\nPRDELLO"), views.SkipFirstSRunes, Config{SValue: 1}},
		{"Проверка -s 3 -f 1", strings.NewReader("SKIP VK.RU\nskipXCZAD cd.RU\nsKIp ch.ru\nSKIP__chtoto__ cmnru"), views.Compose(views.SkipFirstFFields, views.SkipFirstSRunes), Config{FValue: 1, SValue: 3}},
		{"Проверка -i -s 3 -f 1", strings.NewReader("SKIP VK.RU\nFIRST_WORD zxcVK\nSECOND_WORD zxcvk\nNICE vk.vk"), views.Compose(views.SkipFirstFFields, views.SkipFirstSRunes, views.IgnoreCase), Config{FValue: 1, SValue: 3, IgnoreCase: true}},
		{"Проверка -u -s 3 -f 1", strings.NewReader("ASDADS XXXP\nBPPDAOdsad      XRPP\nXA CPRX"), views.Compose(views.SkipFirstFFields, views.SkipFirstSRunes, views.IgnoreCase), Config{Mode: "u", SValue: 3, FValue: 1, IgnoreCase: true}},
	}
	tests := []struct {
		TestCase
		Answer string
	}{
		{testcases[0], "Hello\n"},
		{testcases[1], "Hello\n"},
		{testcases[2], "2 Hello\n3 VAR\n1 V\n"},
		{testcases[3], "2 Hello\n2 v\n"},
		{testcases[4], "Hello\n"},
		{testcases[5], "Hello\nHello\n"},
		{testcases[6], "Hello world\nWORLD woRLD\n"},
		{testcases[7], "HELLO\nPRDELLO\n"},
		{testcases[8], "SKIP VK.RU\nsKIp ch.ru\n"},
		{testcases[9], "SKIP VK.RU\nFIRST_WORD zxcVK\n"},
		{testcases[10], "XA CPRX\n"},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			var output bytes.Buffer
			err := Uniq(tt.Input, &output, tt.LineView, &tt.Config)

			if err != nil {
				t.Errorf("RE.\nОжидается успешное выполнение\nно получено: %s", err)
			}

			if output.String() != tt.Answer {
				t.Errorf("WA.\nОжидается: %+v\nно получено: %+v", tt.Answer, output.String())
			}
		})
	}
}
