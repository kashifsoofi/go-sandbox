package stringcalculator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func Add(input string) (int, error) {
	numbers := parse(input)

	total := 0
	msg := ""
	for _, n := range numbers {
		if n < 0 {
			if msg == "" {
				msg = "negatives not allowed:"
			}
			msg += fmt.Sprintf(" %v", n)
		}
		if n > 1000 {
			continue
		}

		total += n
	}

	if msg != "" {
		return 0, errors.New(msg)
	}

	return total, nil
}

func parse(input string) []int {
	normalisedInput := normaliseInput(input)
	separated := strings.Split(normalisedInput, ",")
	numbers := []int{}
	for _, s := range separated {
		n, _ := strconv.Atoi(s)
		numbers = append(numbers, n)
	}

	return numbers
}

func normaliseInput(input string) string {
	delimeter := ","
	customDelimeters := []string{}
	if strings.Index(input, "//") == 0 {
		input = input[2:]
		idx := 0
		if strings.Index(input, "[") == 0 {
			for strings.Index(input, "[") == 0 {
				idx := strings.Index(input, "]")
				customDelimeter := input[1:idx]
				customDelimeters = append(customDelimeters, customDelimeter)
				input = input[idx+1:]
			}
			idx = strings.Index(input, "\n")
		} else {
			idx = strings.Index(input, "\n")
			customDelimeter := input[:idx]
			customDelimeters = append(customDelimeters, customDelimeter)
		}

		input = input[idx+1:]
	}

	for _, cd := range customDelimeters {
		input = strings.ReplaceAll(input, cd, delimeter)
	}

	if strings.ContainsRune(input, '\n') {
		return strings.ReplaceAll(input, "\n", delimeter)
	}

	return input
}
