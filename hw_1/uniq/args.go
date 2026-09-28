package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ParseArgsToConfig проверяет os.Args и возвращает заполненную структуру Config или error.
func ParseArgsToConfig(args []string) (*Config, error) {
	config := &Config{}
	var positionals []string

	for i := 0; i < len(args); {
		arg := args[i]

		if !strings.HasPrefix(arg, "-") || arg == "-" {
			positionals = append(positionals, args[i])
			i++
			continue
		}

		switch arg {
		case "-c", "-d", "-u":
			mode := strings.TrimPrefix(arg, "-")
			if config.Mode != "" {
				return nil, fmt.Errorf("флаги -c, -d и -u взаимоисключающие (уже задан -%s)", config.Mode)
			}
			config.Mode = mode

		case "-i":
			if config.IgnoreCase {
				return nil, errors.New("флаг -i дублируется")
			}
			config.IgnoreCase = true

		case "-f":
			if config.FValue != 0 {
				return nil, errors.New("флаг -f дублируется")
			}
			if i+1 >= len(args) {
				return nil, errors.New("флаг -f требует числового значения")
			}
			num, err := strconv.Atoi(args[i+1])
			if err != nil || num <= 0 {
				return nil, fmt.Errorf("некорректное значение для -f: %s (ожидается целое положительное число)", args[i+1])
			}
			config.FValue = num
			i++

		case "-s":
			if config.SValue != 0 {
				return nil, errors.New("флаг -s дублируется")
			}
			if i+1 >= len(args) {
				return nil, errors.New("флаг -s требует числового значения")
			}
			num, err := strconv.Atoi(args[i+1])
			if err != nil || num < 0 {
				return nil, fmt.Errorf("некорректное значение для -s: %s (ожидается целое неотрицательное число)", args[i+1])
			}
			config.SValue = num
			i++

		default:
			return nil, fmt.Errorf("неизвестный флаг: %s", arg)
		}

		i++
	}

	if len(positionals) > 2 {
		return nil, errors.New("слишком много аргументов: требуется [input_file [output_file]]")
	}

	if len(positionals) >= 1 {
		config.InputFile = positionals[0]
	}
	if len(positionals) == 2 {
		config.OutputFile = positionals[1]
	}

	return config, nil
}
