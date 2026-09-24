package main

import (
	"testing"
)

func TestParseArgsToConfigPositive(t *testing.T) {
	tests := []struct {
		Args   []string
		Answer Config
	}{
		{[]string{}, Config{}},
		{[]string{"in"}, Config{InputFile: "in"}},
		{[]string{"in", "out"}, Config{InputFile: "in", OutputFile: "out"}},
		{[]string{"-c", "-s", "5"}, Config{Mode: "c", SValue: 5}},
		{[]string{"-c", "-s", "7", "-f", "10"}, Config{Mode: "c", SValue: 7, FValue: 10}},
		{[]string{"-f", "102"}, Config{FValue: 102}},
		{[]string{"-u", "in", "out"}, Config{InputFile: "in", OutputFile: "out", Mode: "u"}},
		{[]string{"-d"}, Config{Mode: "d"}},
		{[]string{"-c", "-s", "7", "in", "out"}, Config{Mode: "c", SValue: 7, InputFile: "in", OutputFile: "out"}},
		{[]string{"-c", "-s", "8", "in"}, Config{Mode: "c", SValue: 8, InputFile: "in"}},
	}

	for _, tt := range tests {
		parsed, err := ParseArgsToConfig(tt.Args)
		if err != nil {
			t.Errorf("RE.\nMust be success (err == nil)\nbut given: %s", err)
			continue
		}
		if *parsed != tt.Answer {
			t.Errorf("WA.\nExpected: %+v\nbut given: %+v", tt.Answer, *parsed)
		}
	}
}
