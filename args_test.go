package main

import (
	"testing"
)

func TestParseArgsToConfigPositive(t *testing.T) {
	tests := []struct {
		Name   string
		Args   []string
		Answer Config
	}{
		{"Нет аргументов командной строки", []string{}, Config{}},
		{"Передан только input_file", []string{"in"}, Config{InputFile: "in"}},
		{"Передан input_file и output_file", []string{"in", "out"}, Config{InputFile: "in", OutputFile: "out"}},
		{"Передан мод и аргумент s", []string{"-c", "-s", "5"}, Config{Mode: "c", SValue: 5}},
		{"Передан мод и аргументы s, f", []string{"-c", "-s", "7", "-f", "10"}, Config{Mode: "c", SValue: 7, FValue: 10}},
		{"Передан только аргумент f", []string{"-f", "102"}, Config{FValue: 102}},
		{"Передан мод и input_file, output_file", []string{"-u", "in", "out"}, Config{InputFile: "in", OutputFile: "out", Mode: "u"}},
		{"Передан только мод", []string{"-d"}, Config{Mode: "d"}},
		{"Передан мод, аргумент s и input_file, output_file", []string{"-c", "-s", "7", "in", "out"}, Config{Mode: "c", SValue: 7, InputFile: "in", OutputFile: "out"}},
		{"Передан мод, аргумент f и input_file", []string{"-u", "-f", "8", "in"}, Config{Mode: "u", FValue: 8, InputFile: "in"}},
		{"Переданы все возможные аргументы", []string{"-c", "-i", "-s", "7", "-f", "1", "in", "out"}, Config{Mode: "c", IgnoreCase: true, SValue: 7, FValue: 1, InputFile: "in", OutputFile: "out"}},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			parsed, err := ParseArgsToConfig(tt.Args)
			if err != nil {
				t.Errorf("RE.\nMust be success (err == nil)\nbut given: %s", err)
			} else if *parsed != tt.Answer {
				t.Errorf("WA.\nExpected: %+v\nbut given: %+v", tt.Answer, *parsed)
			}
		})
	}
}

func TestParseArgsConfigNegative(t *testing.T) {
	tests := []struct {
		Name string
		Args []string
	}{
		{"Количество файлов > 2", []string{"in", "in", "out", "bc"}},
		{"Два взаимоисключающие аргументы s и d", []string{"-c", "-d"}},
		{"Два взаимоисключающие аргументы u и d", []string{"-u", "-d"}},
		{"Два повторяющихся аргумента рядом", []string{"-c", "-c"}},
		{"Два повторяющихся аргумента через другие аргументы", []string{"-d", "-s", "7", "-d"}},
		{"Неправильное значение аргумента s", []string{"-d", "-s", "six-seven"}},
		{"Пропущено значение f", []string{"-f", "-s", "7"}},
		{"Повторяющиеся аргументы -i", []string{"-c", "-i", "-s", "8", "-i", "in"}},
		{"Неизвестный аргумент", []string{"-i", "in", "out", "-c", "8"}},
		{"Отрицательное значение аргумента s", []string{"-s", "-150"}},
		{"Неизвестный аргумент с -", []string{"-s", "-XY"}},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := ParseArgsToConfig(tt.Args)
			if err == nil {
				t.Errorf("WA.\nMust be error (error != nil)")
			}
		})
	}
}
