package edit_distance

import (
	"testing"
)

func TestDamerauOSAFastAlgorithm_Compare(t *testing.T) {
	algorithm := &DamerauOSAFastAlgorithm{}

	tests := []struct {
		name        string
		a           string
		b           string
		maxDistance int
		expected    int
	}{
		{
			name:        "identical strings",
			a:           "hello",
			b:           "hello",
			maxDistance: 5,
			expected:    0,
		},
		{
			name:        "empty strings",
			a:           "",
			b:           "",
			maxDistance: 5,
			expected:    0,
		},
		{
			name:        "first string empty",
			a:           "",
			b:           "abc",
			maxDistance: 5,
			expected:    3,
		},
		{
			name:        "second string empty",
			a:           "abc",
			b:           "",
			maxDistance: 5,
			expected:    3,
		},
		{
			name:        "single character difference",
			a:           "cat",
			b:           "bat",
			maxDistance: 5,
			expected:    1,
		},
		{
			name:        "insertion needed",
			a:           "cat",
			b:           "cats",
			maxDistance: 5,
			expected:    1,
		},
		{
			name:        "deletion needed",
			a:           "cats",
			b:           "cat",
			maxDistance: 5,
			expected:    1,
		},
		{
			name:        "transposition",
			a:           "ab",
			b:           "ba",
			maxDistance: 5,
			expected:    1,
		},
		{
			name:        "multiple operations",
			a:           "kitten",
			b:           "sitting",
			maxDistance: 10,
			expected:    3,
		},
		{
			name:        "unicode characters",
			a:           "héllo",
			b:           "hëllo",
			maxDistance: 5,
			expected:    1,
		},
		{
			name:        "emoji characters",
			a:           "🙂hello",
			b:           "👍hello",
			maxDistance: 5,
			expected:    1,
		},
		{
			name:        "common prefix",
			a:           "prefixABC",
			b:           "prefixDEF",
			maxDistance: 5,
			expected:    3,
		},
		{
			name:        "common suffix",
			a:           "ABCsuffix",
			b:           "DEFsuffix",
			maxDistance: 5,
			expected:    3,
		},
		{
			name:        "common prefix and suffix",
			a:           "prefixABCsuffix",
			b:           "prefixDEFsuffix",
			maxDistance: 5,
			expected:    3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := algorithm.Compare(tt.a, tt.b, tt.maxDistance)
			if err != nil {
				t.Errorf("Compare() error = %v", err)
			}
			if result != tt.expected {
				t.Errorf("Compare(%q, %q, %d) = %d, want %d", tt.a, tt.b, tt.maxDistance, result, tt.expected)
			}
		})
	}
}

func TestDamerauOSAFastAlgorithm_MaxDistanceHandling(t *testing.T) {
	algorithm := &DamerauOSAFastAlgorithm{}

	tests := []struct {
		name        string
		a           string
		b           string
		maxDistance int
		expected    int
	}{
		{
			name:        "max distance zero with identical strings",
			a:           "hello",
			b:           "hello",
			maxDistance: 0,
			expected:    0,
		},
		{
			name:        "max distance zero with different strings",
			a:           "hello",
			b:           "world",
			maxDistance: 0,
			expected:    -1,
		},
		{
			name:        "negative max distance with identical strings",
			a:           "hello",
			b:           "hello",
			maxDistance: -1,
			expected:    0,
		},
		{
			name:        "negative max distance with different strings",
			a:           "hello",
			b:           "world",
			maxDistance: -1,
			expected:    -1,
		},
		{
			name:        "exceed max distance by length difference",
			a:           "a",
			b:           "abcdefgh",
			maxDistance: 5,
			expected:    -1,
		},
		{
			name:        "within max distance",
			a:           "abc",
			b:           "def",
			maxDistance: 5,
			expected:    3,
		},
		{
			name:        "exactly at max distance",
			a:           "abc",
			b:           "defg",
			maxDistance: 4,
			expected:    4,
		},
		{
			name:        "exceed max distance",
			a:           "hello",
			b:           "world",
			maxDistance: 2,
			expected:    -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := algorithm.Compare(tt.a, tt.b, tt.maxDistance)
			if err != nil {
				t.Errorf("Compare() error = %v", err)
			}
			if result != tt.expected {
				t.Errorf("Compare(%q, %q, %d) = %d, want %d", tt.a, tt.b, tt.maxDistance, result, tt.expected)
			}
		})
	}
}

func TestDamerauOSAFastAlgorithm_EdgeCases(t *testing.T) {
	algorithm := &DamerauOSAFastAlgorithm{}

	tests := []struct {
		name        string
		a           string
		b           string
		maxDistance int
		expected    int
	}{
		{
			name:        "both strings empty",
			a:           "",
			b:           "",
			maxDistance: 1,
			expected:    0,
		},
		{
			name:        "single character strings identical",
			a:           "a",
			b:           "a",
			maxDistance: 1,
			expected:    0,
		},
		{
			name:        "single character strings different",
			a:           "a",
			b:           "b",
			maxDistance: 1,
			expected:    1,
		},
		{
			name:        "first string much longer",
			a:           "verylongstring",
			b:           "short",
			maxDistance: 20,
			expected:    12,
		},
		{
			name:        "second string much longer",
			a:           "short",
			b:           "verylongstring",
			maxDistance: 20,
			expected:    12,
		},
		{
			name:        "strings with spaces",
			a:           "hello world",
			b:           "hello  world",
			maxDistance: 5,
			expected:    1,
		},
		{
			name:        "strings with special characters",
			a:           "hello!@#",
			b:           "hello$%^",
			maxDistance: 5,
			expected:    3,
		},
		{
			name:        "unicode normalization",
			a:           "café",
			b:           "café",
			maxDistance: 5,
			expected:    0,
		},
		{
			name:        "mixed case",
			a:           "Hello",
			b:           "HELLO",
			maxDistance: 5,
			expected:    4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := algorithm.Compare(tt.a, tt.b, tt.maxDistance)
			if err != nil {
				t.Errorf("Compare() error = %v", err)
			}
			if result != tt.expected {
				t.Errorf("Compare(%q, %q, %d) = %d, want %d", tt.a, tt.b, tt.maxDistance, result, tt.expected)
			}
		})
	}
}

func TestDamerauOSAFastAlgorithm_TranspositionCases(t *testing.T) {
	algorithm := &DamerauOSAFastAlgorithm{}

	tests := []struct {
		name        string
		a           string
		b           string
		maxDistance int
		expected    int
	}{
		{
			name:        "simple transposition",
			a:           "ab",
			b:           "ba",
			maxDistance: 5,
			expected:    1,
		},
		{
			name:        "transposition in middle",
			a:           "abcd",
			b:           "acbd",
			maxDistance: 5,
			expected:    1,
		},
		{
			name:        "multiple transpositions",
			a:           "abcd",
			b:           "badc",
			maxDistance: 5,
			expected:    2,
		},
		{
			name:        "transposition with unicode",
			a:           "αβ",
			b:           "βα",
			maxDistance: 5,
			expected:    1,
		},
		{
			name:        "no transposition possible",
			a:           "abc",
			b:           "def",
			maxDistance: 5,
			expected:    3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := algorithm.Compare(tt.a, tt.b, tt.maxDistance)
			if err != nil {
				t.Errorf("Compare() error = %v", err)
			}
			if result != tt.expected {
				t.Errorf("Compare(%q, %q, %d) = %d, want %d", tt.a, tt.b, tt.maxDistance, result, tt.expected)
			}
		})
	}
}

func TestDamerauOSAFastAlgorithm_PrefixSuffixOptimization(t *testing.T) {
	algorithm := &DamerauOSAFastAlgorithm{}

	tests := []struct {
		name        string
		a           string
		b           string
		maxDistance int
		expected    int
	}{
		{
			name:        "identical after prefix removal",
			a:           "prefixSAME",
			b:           "prefixSAME",
			maxDistance: 5,
			expected:    0,
		},
		{
			name:        "identical after suffix removal",
			a:           "SAMEsuffix",
			b:           "SAMEsuffix",
			maxDistance: 5,
			expected:    0,
		},
		{
			name:        "identical after both prefix and suffix removal",
			a:           "prefixSAMEsuffix",
			b:           "prefixSAMEsuffix",
			maxDistance: 5,
			expected:    0,
		},
		{
			name:        "different only in middle after trimming",
			a:           "prefixABCsuffix",
			b:           "prefixXYZsuffix",
			maxDistance: 5,
			expected:    3,
		},
		{
			name:        "empty after prefix trimming",
			a:           "commonpart",
			b:           "commonpartextra",
			maxDistance: 10,
			expected:    5,
		},
		{
			name:        "empty after suffix trimming",
			a:           "extracommonpart",
			b:           "commonpart",
			maxDistance: 10,
			expected:    5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := algorithm.Compare(tt.a, tt.b, tt.maxDistance)
			if err != nil {
				t.Errorf("Compare() error = %v", err)
			}
			if result != tt.expected {
				t.Errorf("Compare(%q, %q, %d) = %d, want %d", tt.a, tt.b, tt.maxDistance, result, tt.expected)
			}
		})
	}
}

func TestDamerauOSAAlgorithm_Compare(t *testing.T) {
	algorithm := &DamerauOSAAlgorithm{}

	tests := []struct {
		name        string
		a           string
		b           string
		maxDistance int
	}{
		{
			name:        "identical strings",
			a:           "hello",
			b:           "hello",
			maxDistance: 5,
		},
		{
			name:        "empty strings",
			a:           "",
			b:           "",
			maxDistance: 5,
		},
		{
			name:        "single character difference",
			a:           "cat",
			b:           "bat",
			maxDistance: 5,
		},
		{
			name:        "transposition",
			a:           "ab",
			b:           "ba",
			maxDistance: 5,
		},
		{
			name:        "unicode strings",
			a:           "héllo",
			b:           "hëllo",
			maxDistance: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := algorithm.Compare(tt.a, tt.b, tt.maxDistance)
			if err != nil {
				t.Errorf("Compare() error = %v", err)
			}
			if result < 0 {
				t.Errorf("Compare() = %d, should not be negative", result)
			}
			if tt.a == tt.b && result != 0 {
				t.Errorf("Compare() = %d for identical strings, want 0", result)
			}
		})
	}
}

func TestDamerauLevenshteinAlgorithm_Compare(t *testing.T) {
	algorithm := &DamerauLevenshteinAlgorithm{}

	tests := []struct {
		name        string
		a           string
		b           string
		maxDistance int
	}{
		{
			name:        "identical strings",
			a:           "hello",
			b:           "hello",
			maxDistance: 5,
		},
		{
			name:        "empty strings",
			a:           "",
			b:           "",
			maxDistance: 5,
		},
		{
			name:        "single character difference",
			a:           "cat",
			b:           "bat",
			maxDistance: 5,
		},
		{
			name:        "transposition",
			a:           "ab",
			b:           "ba",
			maxDistance: 5,
		},
		{
			name:        "unicode strings",
			a:           "héllo",
			b:           "hëllo",
			maxDistance: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := algorithm.Compare(tt.a, tt.b, tt.maxDistance)
			if err != nil {
				t.Errorf("Compare() error = %v", err)
			}
			if result < 0 {
				t.Errorf("Compare() = %d, should not be negative", result)
			}
			if tt.a == tt.b && result != 0 {
				t.Errorf("Compare() = %d for identical strings, want 0", result)
			}
		})
	}
}

func TestLevenshteinAlgorithm_Compare(t *testing.T) {
	algorithm := &LevenshteinAlgorithm{}

	tests := []struct {
		name        string
		a           string
		b           string
		maxDistance int
	}{
		{
			name:        "identical strings",
			a:           "hello",
			b:           "hello",
			maxDistance: 5,
		},
		{
			name:        "empty strings",
			a:           "",
			b:           "",
			maxDistance: 5,
		},
		{
			name:        "single character difference",
			a:           "cat",
			b:           "bat",
			maxDistance: 5,
		},
		{
			name:        "insertion",
			a:           "cat",
			b:           "cats",
			maxDistance: 5,
		},
		{
			name:        "deletion",
			a:           "cats",
			b:           "cat",
			maxDistance: 5,
		},
		{
			name:        "substitution",
			a:           "cat",
			b:           "bat",
			maxDistance: 5,
		},
		{
			name:        "unicode strings",
			a:           "héllo",
			b:           "hëllo",
			maxDistance: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := algorithm.Compare(tt.a, tt.b, tt.maxDistance)
			if err != nil {
				t.Errorf("Compare() error = %v", err)
			}
			if result < 0 {
				t.Errorf("Compare() = %d, should not be negative", result)
			}
			if tt.a == tt.b && result != 0 {
				t.Errorf("Compare() = %d for identical strings, want 0", result)
			}
		})
	}
}

func TestAlgorithmInterface(t *testing.T) {
	algorithms := []Algorithm{
		&DamerauOSAFastAlgorithm{},
		&DamerauOSAAlgorithm{},
		&DamerauLevenshteinAlgorithm{},
		&LevenshteinAlgorithm{},
	}

	for i, alg := range algorithms {
		t.Run("algorithm_interface_"+string(rune(i+'A')), func(t *testing.T) {
			var _ Algorithm = alg

			result, err := alg.Compare("test", "best", 5)
			if err != nil {
				t.Errorf("Compare() error = %v", err)
			}

			if result < -1 {
				t.Errorf("Compare() = %d, should not be less than -1", result)
			}
		})
	}
}

func TestDamerauOSAFastAlgorithm_LongStrings(t *testing.T) {
	algorithm := &DamerauOSAFastAlgorithm{}

	longString1 := "thisisaverylongstringfortestingpurposes"
	longString2 := "thisisaverylongstringfortestingpurposed"

	result, err := algorithm.Compare(longString1, longString2, 10)
	if err != nil {
		t.Errorf("Compare() error = %v", err)
	}

	if result < 0 {
		t.Errorf("Compare() = %d, should not be negative for long strings", result)
	}

	if result > 10 {
		t.Errorf("Compare() = %d, should not exceed reasonable bounds", result)
	}
}

func TestDamerauOSAFastAlgorithm_BoundedVsUnbounded(t *testing.T) {
	algorithm := &DamerauOSAFastAlgorithm{}

	testCases := []struct {
		a string
		b string
	}{
		{"hello", "world"},
		{"kitten", "sitting"},
		{"abc", "def"},
		{"", "abc"},
		{"abc", ""},
	}

	for _, tc := range testCases {
		t.Run("bounded_vs_unbounded", func(t *testing.T) {
			unbounded, err := algorithm.Compare(tc.a, tc.b, 100)
			if err != nil {
				t.Errorf("Unbounded Compare() error = %v", err)
			}

			if unbounded >= 0 {
				bounded, err := algorithm.Compare(tc.a, tc.b, unbounded+1)
				if err != nil {
					t.Errorf("Bounded Compare() error = %v", err)
				}

				if bounded != unbounded {
					t.Errorf("Bounded (%d) vs Unbounded (%d) should match when maxDistance is sufficient", bounded, unbounded)
				}
			}
		})
	}
}

func TestAlgorithmConstants(t *testing.T) {
	if kSpace != rune(' ') {
		t.Errorf("kSpace = %c, want space character", kSpace)
	}
}
