package common

import (
	"reflect"
	"testing"
)

func TestSplitN(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		sep      string
		n        int
		expected []string
	}{
		{
			name:     "normal split with comma",
			s:        "apple,banana,cherry",
			sep:      ",",
			n:        3,
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "split with limit",
			s:        "apple,banana,cherry,date",
			sep:      ",",
			n:        3,
			expected: []string{"apple", "banana", "cherry,date"},
		},
		{
			name:     "empty separator",
			s:        "hello",
			sep:      "",
			n:        3,
			expected: []string{"hello"},
		},
		{
			name:     "empty string",
			s:        "",
			sep:      ",",
			n:        3,
			expected: []string{""},
		},
		{
			name:     "separator not found",
			s:        "hello",
			sep:      ",",
			n:        3,
			expected: []string{"hello"},
		},
		{
			name:     "split with space",
			s:        "hello world test",
			sep:      " ",
			n:        2,
			expected: []string{"hello", "world test"},
		},
		{
			name:     "split with newline",
			s:        "line1\nline2\nline3",
			sep:      "\n",
			n:        2,
			expected: []string{"line1", "line2\nline3"},
		},
		{
			name:     "n is 1",
			s:        "apple,banana,cherry",
			sep:      ",",
			n:        1,
			expected: []string{"apple,banana,cherry"},
		},
		{
			name:     "n is 0",
			s:        "apple,banana,cherry",
			sep:      ",",
			n:        0,
			expected: nil,
		},
		{
			name:     "n is negative",
			s:        "apple,banana,cherry",
			sep:      ",",
			n:        -1,
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "multiple character separator",
			s:        "apple::banana::cherry",
			sep:      "::",
			n:        2,
			expected: []string{"apple", "banana::cherry"},
		},
		{
			name:     "unicode characters",
			s:        "héllo•wörld•tëst",
			sep:      "•",
			n:        3,
			expected: []string{"héllo", "wörld", "tëst"},
		},
		{
			name:     "consecutive separators",
			s:        "apple,,banana,,cherry",
			sep:      ",",
			n:        5,
			expected: []string{"apple", "", "banana", "", "cherry"},
		},
		{
			name:     "string is just separator",
			s:        ",",
			sep:      ",",
			n:        3,
			expected: []string{"", ""},
		},
		{
			name:     "string starts with separator",
			s:        ",apple,banana",
			sep:      ",",
			n:        3,
			expected: []string{"", "apple", "banana"},
		},
		{
			name:     "string ends with separator",
			s:        "apple,banana,",
			sep:      ",",
			n:        3,
			expected: []string{"apple", "banana", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitN(tt.s, tt.sep, tt.n)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SplitN(%q, %q, %d) = %v, want %v", tt.s, tt.sep, tt.n, result, tt.expected)
			}
		})
	}
}

func TestIsAcronym(t *testing.T) {
	tests := []struct {
		name                 string
		word                 string
		ignoreTermWithDigits bool
		expected             bool
	}{
		{
			name:                 "simple acronym",
			word:                 "FBI",
			ignoreTermWithDigits: false,
			expected:             true,
		},
		{
			name:                 "mixed case not acronym",
			word:                 "FbI",
			ignoreTermWithDigits: false,
			expected:             false,
		},
		{
			name:                 "lowercase not acronym",
			word:                 "fbi",
			ignoreTermWithDigits: false,
			expected:             false,
		},
		{
			name:                 "single uppercase letter",
			word:                 "A",
			ignoreTermWithDigits: false,
			expected:             true,
		},
		{
			name:                 "single lowercase letter",
			word:                 "a",
			ignoreTermWithDigits: false,
			expected:             false,
		},
		{
			name:                 "empty string",
			word:                 "",
			ignoreTermWithDigits: false,
			expected:             false,
		},
		{
			name:                 "acronym with digits",
			word:                 "FBI2",
			ignoreTermWithDigits: false,
			expected:             true,
		},
		{
			name:                 "acronym with digits ignore false",
			word:                 "FBI2",
			ignoreTermWithDigits: true,
			expected:             true,
		},
		{
			name:                 "only digits",
			word:                 "123",
			ignoreTermWithDigits: false,
			expected:             false,
		},
		{
			name:                 "only digits ignore true",
			word:                 "123",
			ignoreTermWithDigits: true,
			expected:             true,
		},
		{
			name:                 "mixed digits and uppercase",
			word:                 "A1B2C3",
			ignoreTermWithDigits: false,
			expected:             true,
		},
		{
			name:                 "mixed digits and uppercase ignore true",
			word:                 "A1B2C3",
			ignoreTermWithDigits: true,
			expected:             true,
		},
		{
			name:                 "mixed digits and lowercase",
			word:                 "A1b2C3",
			ignoreTermWithDigits: false,
			expected:             false,
		},
		{
			name:                 "mixed digits and lowercase ignore true",
			word:                 "A1b2C3",
			ignoreTermWithDigits: true,
			expected:             true,
		},
		{
			name:                 "special characters with uppercase",
			word:                 "A-B-C",
			ignoreTermWithDigits: false,
			expected:             false,
		},
		{
			name:                 "unicode uppercase",
			word:                 "ÄÖÜ",
			ignoreTermWithDigits: false,
			expected:             true,
		},
		{
			name:                 "unicode mixed case",
			word:                 "Äöü",
			ignoreTermWithDigits: false,
			expected:             false,
		},
		{
			name:                 "has letter check - no letters",
			word:                 "123!@#",
			ignoreTermWithDigits: false,
			expected:             false,
		},
		{
			name:                 "has letter check - no letters ignore true",
			word:                 "123!@#",
			ignoreTermWithDigits: true,
			expected:             true,
		},
		{
			name:                 "single digit",
			word:                 "1",
			ignoreTermWithDigits: false,
			expected:             false,
		},
		{
			name:                 "single digit ignore true",
			word:                 "1",
			ignoreTermWithDigits: true,
			expected:             true,
		},
		{
			name:                 "symbols only",
			word:                 "!@#",
			ignoreTermWithDigits: false,
			expected:             false,
		},
		{
			name:                 "symbols only ignore true",
			word:                 "!@#",
			ignoreTermWithDigits: true,
			expected:             false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsAcronym(tt.word, tt.ignoreTermWithDigits)
			if result != tt.expected {
				t.Errorf("IsAcronym(%q, %t) = %t, want %t", tt.word, tt.ignoreTermWithDigits, result, tt.expected)
			}
		})
	}
}

func TestCaseTransferSimilar(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		target   string
		expected string
	}{
		{
			name:     "all uppercase source",
			source:   "HELLO",
			target:   "world",
			expected: "WORLD",
		},
		{
			name:     "title case source",
			source:   "Hello",
			target:   "world",
			expected: "World",
		},
		{
			name:     "lowercase source",
			source:   "hello",
			target:   "world",
			expected: "world",
		},
		{
			name:     "mixed case source",
			source:   "HeLLo",
			target:   "world",
			expected: "World",
		},
		{
			name:     "empty source",
			source:   "",
			target:   "world",
			expected: "",
		},
		{
			name:     "empty target",
			source:   "HELLO",
			target:   "",
			expected: "",
		},
		{
			name:     "both empty",
			source:   "",
			target:   "",
			expected: "",
		},
		{
			name:     "single character uppercase source",
			source:   "A",
			target:   "world",
			expected: "WORLD",
		},
		{
			name:     "single character lowercase source",
			source:   "a",
			target:   "world",
			expected: "world",
		},
		{
			name:     "single character target",
			source:   "HELLO",
			target:   "w",
			expected: "W",
		},
		{
			name:     "title case single character target",
			source:   "Hello",
			target:   "w",
			expected: "W",
		},
		{
			name:     "unicode uppercase source",
			source:   "HÉLLO",
			target:   "wörld",
			expected: "WÖRLD",
		},
		{
			name:     "unicode title case source",
			source:   "Héllo",
			target:   "wörld",
			expected: "Wörld",
		},
		{
			name:     "unicode lowercase source",
			source:   "héllo",
			target:   "wörld",
			expected: "wörld",
		},
		{
			name:     "unicode mixed case source",
			source:   "HéLLo",
			target:   "wörld",
			expected: "Wörld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CaseTransferSimilar(tt.source, tt.target)
			if result != tt.expected {
				t.Errorf("CaseTransferSimilar(%q, %q) = %q, want %q", tt.source, tt.target, result, tt.expected)
			}
		})
	}
}
