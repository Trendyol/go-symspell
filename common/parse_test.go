package common

import (
	"math"
	"testing"
)

func TestTryParseInt64(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedBool  bool
		expectedValue int64
	}{
		{
			name:          "valid positive number",
			input:         "123",
			expectedBool:  true,
			expectedValue: 123,
		},
		{
			name:          "valid negative number",
			input:         "-123",
			expectedBool:  true,
			expectedValue: -123,
		},
		{
			name:          "zero",
			input:         "0",
			expectedBool:  true,
			expectedValue: 0,
		},
		{
			name:          "negative zero",
			input:         "-0",
			expectedBool:  true,
			expectedValue: 0,
		},
		{
			name:          "max int64",
			input:         "9223372036854775807",
			expectedBool:  true,
			expectedValue: math.MaxInt64,
		},
		{
			name:          "min int64",
			input:         "-9223372036854775808",
			expectedBool:  true,
			expectedValue: math.MinInt64,
		},
		{
			name:          "large positive number",
			input:         "1000000000000",
			expectedBool:  true,
			expectedValue: 1000000000000,
		},
		{
			name:          "large negative number",
			input:         "-1000000000000",
			expectedBool:  true,
			expectedValue: -1000000000000,
		},
		{
			name:          "single digit",
			input:         "5",
			expectedBool:  true,
			expectedValue: 5,
		},
		{
			name:          "leading zeros",
			input:         "00123",
			expectedBool:  true,
			expectedValue: 123,
		},
		{
			name:          "leading zeros negative",
			input:         "-00123",
			expectedBool:  true,
			expectedValue: -123,
		},
		{
			name:          "invalid string",
			input:         "abc",
			expectedBool:  false,
			expectedValue: 0,
		},
		{
			name:          "empty string",
			input:         "",
			expectedBool:  false,
			expectedValue: 0,
		},
		{
			name:          "mixed alphanumeric",
			input:         "123abc",
			expectedBool:  false,
			expectedValue: 0,
		},
		{
			name:          "alphanumeric mixed",
			input:         "abc123",
			expectedBool:  false,
			expectedValue: 0,
		},
		{
			name:          "decimal number",
			input:         "123.45",
			expectedBool:  false,
			expectedValue: 0,
		},
		{
			name:          "scientific notation",
			input:         "1e5",
			expectedBool:  false,
			expectedValue: 0,
		},
		{
			name:          "hexadecimal",
			input:         "0xFF",
			expectedBool:  false,
			expectedValue: 0,
		},
		{
			name:          "octal",
			input:         "0123",
			expectedBool:  true,
			expectedValue: 123,
		},
		{
			name:          "binary",
			input:         "0b1010",
			expectedBool:  false,
			expectedValue: 0,
		},
		{
			name:          "with spaces",
			input:         " 123 ",
			expectedBool:  false,
			expectedValue: 0,
		},
		{
			name:          "with leading space",
			input:         " 123",
			expectedBool:  false,
			expectedValue: 0,
		},
		{
			name:          "with trailing space",
			input:         "123 ",
			expectedBool:  false,
			expectedValue: 0,
		},
		{
			name:          "plus sign",
			input:         "+123",
			expectedBool:  true,
			expectedValue: 123,
		},
		{
			name:          "double negative",
			input:         "--123",
			expectedBool:  false,
			expectedValue: 0,
		},
		{
			name:          "overflow positive",
			input:         "9223372036854775808",
			expectedBool:  false,
			expectedValue: 0,
		},
		{
			name:          "overflow negative",
			input:         "-9223372036854775809",
			expectedBool:  false,
			expectedValue: 0,
		},
		{
			name:          "very large number",
			input:         "99999999999999999999999999999999",
			expectedBool:  false,
			expectedValue: 0,
		},
		{
			name:          "unicode digits",
			input:         "১২৩",
			expectedBool:  false,
			expectedValue: 0,
		},
		{
			name:          "roman numerals",
			input:         "XVII",
			expectedBool:  false,
			expectedValue: 0,
		},
		{
			name:          "special characters",
			input:         "!@#",
			expectedBool:  false,
			expectedValue: 0,
		},
		{
			name:          "newline character",
			input:         "123\n",
			expectedBool:  false,
			expectedValue: 0,
		},
		{
			name:          "tab character",
			input:         "123\t",
			expectedBool:  false,
			expectedValue: 0,
		},
		{
			name:          "carriage return",
			input:         "123\r",
			expectedBool:  false,
			expectedValue: 0,
		},
		{
			name:          "comma separator",
			input:         "1,234",
			expectedBool:  false,
			expectedValue: 0,
		},
		{
			name:          "dot separator",
			input:         "1.234",
			expectedBool:  false,
			expectedValue: 0,
		},
		{
			name:          "underscore separator",
			input:         "1_234",
			expectedBool:  false,
			expectedValue: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultBool, resultValue := TryParseInt64(tt.input)

			if resultBool != tt.expectedBool {
				t.Errorf("TryParseInt64(%q) bool = %v, want %v", tt.input, resultBool, tt.expectedBool)
			}

			if resultValue != tt.expectedValue {
				t.Errorf("TryParseInt64(%q) value = %d, want %d", tt.input, resultValue, tt.expectedValue)
			}
		})
	}
}

func TestTryParseInt64EdgeCases(t *testing.T) {
	ok, val := TryParseInt64("0")
	if !ok || val != 0 {
		t.Errorf("TryParseInt64(\"0\") = (%v, %d), want (true, 0)", ok, val)
	}

	ok, val = TryParseInt64("-1")
	if !ok || val != -1 {
		t.Errorf("TryParseInt64(\"-1\") = (%v, %d), want (true, -1)", ok, val)
	}

	ok, val = TryParseInt64("1")
	if !ok || val != 1 {
		t.Errorf("TryParseInt64(\"1\") = (%v, %d), want (true, 1)", ok, val)
	}

	ok, val = TryParseInt64("")
	if ok || val != 0 {
		t.Errorf("TryParseInt64(\"\") = (%v, %d), want (false, 0)", ok, val)
	}
}

func TestTryParseInt64BoundaryValues(t *testing.T) {
	maxInt64Str := "9223372036854775807"
	ok, val := TryParseInt64(maxInt64Str)
	if !ok || val != math.MaxInt64 {
		t.Errorf("TryParseInt64(%q) = (%v, %d), want (true, %d)", maxInt64Str, ok, val, math.MaxInt64)
	}

	minInt64Str := "-9223372036854775808"
	ok, val = TryParseInt64(minInt64Str)
	if !ok || val != math.MinInt64 {
		t.Errorf("TryParseInt64(%q) = (%v, %d), want (true, %d)", minInt64Str, ok, val, math.MinInt64)
	}

	overflowStr := "9223372036854775808"
	ok, val = TryParseInt64(overflowStr)
	if ok {
		t.Errorf("TryParseInt64(%q) should fail due to overflow", overflowStr)
	}

	underflowStr := "-9223372036854775809"
	ok, val = TryParseInt64(underflowStr)
	if ok {
		t.Errorf("TryParseInt64(%q) should fail due to underflow", underflowStr)
	}
}

func TestTryParseInt64ConsistencyWithStandardLibrary(t *testing.T) {
	testInputs := []string{
		"123", "-123", "0", "+456", "9223372036854775807", "-9223372036854775808",
		"abc", "123abc", "", " 123", "123.45", "0xFF",
	}

	for _, input := range testInputs {
		t.Run("consistency_"+input, func(t *testing.T) {
			ok, _ := TryParseInt64(input)

			if input == "" || input == "abc" || input == "123abc" || input == " 123" || input == "123.45" || input == "0xFF" {
				if ok {
					t.Errorf("TryParseInt64(%q) should return false for invalid input", input)
				}
			} else if input == "123" || input == "-123" || input == "0" || input == "+456" || input == "9223372036854775807" || input == "-9223372036854775808" {
				if !ok {
					t.Errorf("TryParseInt64(%q) should return true for valid input", input)
				}
			}
		})
	}
}
