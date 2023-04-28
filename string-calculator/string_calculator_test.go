package stringcalculator

import (
	"testing"
)

func TestAdd(t *testing.T) {
	cases := []struct {
		name           string
		input          string
		expectedResult int
		expectedMsg    string
	}{
		{
			name:           "empty string should return 0",
			input:          "",
			expectedResult: 0,
			expectedMsg:    "",
		},
		{
			name:           "one number should return number",
			input:          "42",
			expectedResult: 42,
			expectedMsg:    "",
		},
		{
			name:           "two numbers with comma delimiter numbers should return sum",
			input:          "1,2",
			expectedResult: 3,
			expectedMsg:    "",
		},
		{
			name:           "multiple numbers with comma delimiter should return sum",
			input:          "1,2,3",
			expectedResult: 6,
			expectedMsg:    "",
		},
		{
			name:           "two numbers with newline delimiter should return sum",
			input:          "1\n2",
			expectedResult: 3,
			expectedMsg:    "",
		},
		{
			name:           "multiple numbers with newline delimiter should return sum",
			input:          "1\n2\n3",
			expectedResult: 6,
			expectedMsg:    "",
		},
		{
			name:           "multiple numbers with comma and newline delimiter should return sum",
			input:          "1\n2,3",
			expectedResult: 6,
			expectedMsg:    "",
		},
		{
			name:           "two numbers with custom delimiter should return sum",
			input:          "//;\n1;2",
			expectedResult: 3,
			expectedMsg:    "",
		},
		{
			name:           "negative number should return error",
			input:          "-1",
			expectedResult: 0,
			expectedMsg:    "negatives not allowed: -1",
		},
		{
			name:           "multiple negative numbers should return error with all negative numbers",
			input:          "-1,2,-3",
			expectedResult: 0,
			expectedMsg:    "negatives not allowed: -1 -3",
		},
		{
			name:           "numbers bigger than 1000 should be ignored",
			input:          "2,1001",
			expectedResult: 2,
			expectedMsg:    "",
		},
		{
			name:           "multiple numbers with multichar delimiter should return sum",
			input:          "//[***]\n1***2***3",
			expectedResult: 6,
			expectedMsg:    "",
		},
		{
			name:           "multiple numbers with multiple custom delimiter should return sum",
			input:          "//[*][%]\n1*2%3",
			expectedResult: 6,
			expectedMsg:    "",
		},
		{
			name:           "multiple numbers with multiple multichar custom delimiter should return sum",
			input:          "//[*%*][%*%]\n1*%*2%*%3",
			expectedResult: 6,
			expectedMsg:    "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result, err := Add(c.input)

			if c.expectedMsg != "" {
				if err == nil {
					t.Errorf("got: nil, want: %s", c.expectedMsg)
				}

				errMsg := err.Error()
				if errMsg != c.expectedMsg {
					t.Errorf("got: %s, want: %s", errMsg, c.expectedMsg)
				}
			}
			if result != c.expectedResult {
				t.Errorf("got: %d, want: %d", result, c.expectedResult)
			}
		})
	}
}
