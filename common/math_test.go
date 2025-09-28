package common

import (
	"math"
	"testing"
)

func TestAbs(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{
			name:     "positive number",
			input:    5,
			expected: 5,
		},
		{
			name:     "negative number",
			input:    -5,
			expected: 5,
		},
		{
			name:     "zero",
			input:    0,
			expected: 0,
		},
		{
			name:     "negative zero",
			input:    -0,
			expected: 0,
		},
		{
			name:     "one",
			input:    1,
			expected: 1,
		},
		{
			name:     "negative one",
			input:    -1,
			expected: 1,
		},
		{
			name:     "large positive number",
			input:    1000000,
			expected: 1000000,
		},
		{
			name:     "large negative number",
			input:    -1000000,
			expected: 1000000,
		},
		{
			name:     "max int",
			input:    math.MaxInt,
			expected: math.MaxInt,
		},
		{
			name:     "min int plus one",
			input:    math.MinInt + 1,
			expected: math.MaxInt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Abs(tt.input)
			if result != tt.expected {
				t.Errorf("Abs(%d) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAbsEdgeCases(t *testing.T) {
	result := Abs(math.MinInt)
	if result != math.MinInt {
		t.Errorf("Abs(math.MinInt) = %d, but due to overflow it remains math.MinInt", result)
	}
}

func TestAbsBenchmark(t *testing.T) {
	testCases := []int{-1000, -1, 0, 1, 1000}

	for _, tc := range testCases {
		result := Abs(tc)
		if tc >= 0 && result != tc {
			t.Errorf("Abs(%d) = %d, want %d", tc, result, tc)
		}
		if tc < 0 && result != -tc {
			t.Errorf("Abs(%d) = %d, want %d", tc, result, -tc)
		}
	}
}
