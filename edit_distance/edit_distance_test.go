package edit_distance

import (
	"testing"
)

func TestAlgorithmTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant AlgorithmType
		expected int
	}{
		{
			name:     "DamerauOSAFast should be 0",
			constant: DamerauOSAFast,
			expected: 0,
		},
		{
			name:     "DamerauOSA should be 1",
			constant: DamerauOSA,
			expected: 1,
		},
		{
			name:     "DamerauLevenshtein should be 2",
			constant: DamerauLevenshtein,
			expected: 2,
		},
		{
			name:     "Levenshtein should be 3",
			constant: Levenshtein,
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.constant) != tt.expected {
				t.Errorf("Algorithm constant %v = %d, want %d", tt.constant, int(tt.constant), tt.expected)
			}
		})
	}
}

func TestNewDistanceComparer(t *testing.T) {
	tests := []struct {
		name          string
		algorithmType AlgorithmType
		expectError   bool
		errorMessage  string
	}{
		{
			name:          "DamerauOSAFast algorithm",
			algorithmType: DamerauOSAFast,
			expectError:   false,
		},
		{
			name:          "DamerauOSA algorithm",
			algorithmType: DamerauOSA,
			expectError:   false,
		},
		{
			name:          "DamerauLevenshtein algorithm",
			algorithmType: DamerauLevenshtein,
			expectError:   false,
		},
		{
			name:          "Levenshtein algorithm",
			algorithmType: Levenshtein,
			expectError:   false,
		},
		{
			name:          "invalid algorithm type",
			algorithmType: AlgorithmType(999),
			expectError:   true,
			errorMessage:  "invalid algorithm",
		},
		{
			name:          "negative algorithm type",
			algorithmType: AlgorithmType(-1),
			expectError:   true,
			errorMessage:  "invalid algorithm",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comparer, err := NewDistanceComparer(tt.algorithmType)

			if tt.expectError {
				if err == nil {
					t.Errorf("NewDistanceComparer() expected error but got none")
				}
				if err != nil && err.Error() != tt.errorMessage {
					t.Errorf("NewDistanceComparer() error = %v, want %v", err.Error(), tt.errorMessage)
				}
				if comparer != nil {
					t.Errorf("NewDistanceComparer() should return nil when error occurs")
				}
			} else {
				if err != nil {
					t.Errorf("NewDistanceComparer() unexpected error = %v", err)
				}
				if comparer == nil {
					t.Errorf("NewDistanceComparer() should not return nil for valid algorithm")
				}

				dc, ok := comparer.(*distanceComparer)
				if !ok {
					t.Errorf("NewDistanceComparer() should return *distanceComparer")
				}
				if dc.algorithm == nil {
					t.Errorf("NewDistanceComparer() should set algorithm field")
				}
			}
		})
	}
}

func TestDistanceComparerAlgorithmTypes(t *testing.T) {
	algorithms := []AlgorithmType{
		DamerauOSAFast,
		DamerauOSA,
		DamerauLevenshtein,
		Levenshtein,
	}

	for _, alg := range algorithms {
		t.Run("algorithm_"+string(rune(alg+'0')), func(t *testing.T) {
			comparer, err := NewDistanceComparer(alg)
			if err != nil {
				t.Fatalf("NewDistanceComparer() error = %v", err)
			}

			dc := comparer.(*distanceComparer)
			switch alg {
			case DamerauOSAFast:
				if _, ok := dc.algorithm.(*DamerauOSAFastAlgorithm); !ok {
					t.Errorf("Expected DamerauOSAFastAlgorithm for DamerauOSAFast")
				}
			case DamerauOSA:
				if _, ok := dc.algorithm.(*DamerauOSAAlgorithm); !ok {
					t.Errorf("Expected DamerauOSAAlgorithm for DamerauOSA")
				}
			case DamerauLevenshtein:
				if _, ok := dc.algorithm.(*DamerauLevenshteinAlgorithm); !ok {
					t.Errorf("Expected DamerauLevenshteinAlgorithm for DamerauLevenshtein")
				}
			case Levenshtein:
				if _, ok := dc.algorithm.(*LevenshteinAlgorithm); !ok {
					t.Errorf("Expected LevenshteinAlgorithm for Levenshtein")
				}
			}
		})
	}
}

func TestDistanceComparerCompare(t *testing.T) {
	tests := []struct {
		name        string
		algorithm   AlgorithmType
		a           string
		b           string
		maxDistance int
	}{
		{
			name:        "identical strings",
			algorithm:   DamerauOSAFast,
			a:           "hello",
			b:           "hello",
			maxDistance: 5,
		},
		{
			name:        "empty strings",
			algorithm:   DamerauOSAFast,
			a:           "",
			b:           "",
			maxDistance: 0,
		},
		{
			name:        "one empty string",
			algorithm:   DamerauOSAFast,
			a:           "hello",
			b:           "",
			maxDistance: 10,
		},
		{
			name:        "different strings",
			algorithm:   DamerauOSAFast,
			a:           "kitten",
			b:           "sitting",
			maxDistance: 10,
		},
		{
			name:        "unicode strings",
			algorithm:   DamerauOSAFast,
			a:           "héllo",
			b:           "hëllo",
			maxDistance: 5,
		},
		{
			name:        "max distance zero",
			algorithm:   DamerauOSAFast,
			a:           "hello",
			b:           "world",
			maxDistance: 0,
		},
		{
			name:        "negative max distance",
			algorithm:   DamerauOSAFast,
			a:           "hello",
			b:           "world",
			maxDistance: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comparer, err := NewDistanceComparer(tt.algorithm)
			if err != nil {
				t.Fatalf("NewDistanceComparer() error = %v", err)
			}

			distance, err := comparer.Compare(tt.a, tt.b, tt.maxDistance)
			if err != nil {
				t.Errorf("Compare() error = %v", err)
			}

			if tt.a == tt.b && distance != 0 {
				t.Errorf("Compare() distance = %d for identical strings, want 0", distance)
			}

			if distance < -1 {
				t.Errorf("Compare() distance = %d, should not be less than -1", distance)
			}
		})
	}
}

func TestDistanceComparerInterface(t *testing.T) {
	comparer, err := NewDistanceComparer(DamerauOSAFast)
	if err != nil {
		t.Fatalf("NewDistanceComparer() error = %v", err)
	}

	var _ DistanceComparer = comparer

	distance, err := comparer.Compare("test", "best", 5)
	if err != nil {
		t.Errorf("Compare() error = %v", err)
	}

	if distance < 0 && distance != -1 {
		t.Errorf("Compare() distance = %d, invalid value", distance)
	}
}

func TestDistanceComparerAllAlgorithms(t *testing.T) {
	algorithms := []AlgorithmType{
		DamerauOSAFast,
		DamerauOSA,
		DamerauLevenshtein,
		Levenshtein,
	}

	testCases := []struct {
		a string
		b string
	}{
		{"", ""},
		{"a", "a"},
		{"abc", "abc"},
		{"", "a"},
		{"a", ""},
		{"abc", "def"},
		{"kitten", "sitting"},
		{"saturday", "sunday"},
	}

	for _, alg := range algorithms {
		t.Run("algorithm_"+string(rune(alg+'A')), func(t *testing.T) {
			comparer, err := NewDistanceComparer(alg)
			if err != nil {
				t.Fatalf("NewDistanceComparer() error = %v", err)
			}

			for _, tc := range testCases {
				distance, err := comparer.Compare(tc.a, tc.b, 10)
				if err != nil {
					t.Errorf("Compare(%q, %q) error = %v", tc.a, tc.b, err)
				}

				if tc.a == tc.b && distance != 0 {
					t.Errorf("Compare(%q, %q) = %d, want 0 for identical strings", tc.a, tc.b, distance)
				}

				if distance < -1 {
					t.Errorf("Compare(%q, %q) = %d, should not be less than -1", tc.a, tc.b, distance)
				}
			}
		})
	}
}

func TestDistanceComparerMaxDistanceHandling(t *testing.T) {
	comparer, err := NewDistanceComparer(DamerauOSAFast)
	if err != nil {
		t.Fatalf("NewDistanceComparer() error = %v", err)
	}

	tests := []struct {
		name         string
		a            string
		b            string
		maxDistance  int
		expectExceed bool
	}{
		{
			name:         "within max distance",
			a:            "cat",
			b:            "bat",
			maxDistance:  5,
			expectExceed: false,
		},
		{
			name:         "exact max distance",
			a:            "a",
			b:            "bb",
			maxDistance:  2,
			expectExceed: false,
		},
		{
			name:         "exceed max distance",
			a:            "hello",
			b:            "world",
			maxDistance:  1,
			expectExceed: true,
		},
		{
			name:         "zero max distance different strings",
			a:            "a",
			b:            "b",
			maxDistance:  0,
			expectExceed: true,
		},
		{
			name:         "zero max distance same strings",
			a:            "a",
			b:            "a",
			maxDistance:  0,
			expectExceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			distance, err := comparer.Compare(tt.a, tt.b, tt.maxDistance)
			if err != nil {
				t.Errorf("Compare() error = %v", err)
			}

			if tt.expectExceed && distance != -1 {
				t.Errorf("Compare() = %d, expected -1 when exceeding maxDistance", distance)
			}

			if !tt.expectExceed && distance == -1 {
				t.Errorf("Compare() = -1, expected valid distance within maxDistance")
			}

			if !tt.expectExceed && distance > tt.maxDistance {
				t.Errorf("Compare() = %d, should not exceed maxDistance %d", distance, tt.maxDistance)
			}
		})
	}
}
