package symspell

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/Trendyol/go-symspell/edit_distance"
)

func TestDefaultOptions(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		expected interface{}
		actual   interface{}
	}{
		{
			name:     "InitialCapacity",
			field:    "InitialCapacity",
			expected: 16,
			actual:   DefaultOptions.InitialCapacity,
		},
		{
			name:     "MaxDictionaryEditDistance",
			field:    "MaxDictionaryEditDistance",
			expected: 2,
			actual:   DefaultOptions.MaxDictionaryEditDistance,
		},
		{
			name:     "PrefixLength",
			field:    "PrefixLength",
			expected: 7,
			actual:   DefaultOptions.PrefixLength,
		},
		{
			name:     "CountThreshold",
			field:    "CountThreshold",
			expected: 1,
			actual:   DefaultOptions.CountThreshold,
		},
		{
			name:     "DistanceAlgorithm",
			field:    "DistanceAlgorithm",
			expected: edit_distance.DamerauOSAFast,
			actual:   DefaultOptions.DistanceAlgorithm,
		},
		{
			name:     "IncludeUnknown",
			field:    "IncludeUnknown",
			expected: false,
			actual:   DefaultOptions.IncludeUnknown,
		},
		{
			name:     "TransferCasing",
			field:    "TransferCasing",
			expected: false,
			actual:   DefaultOptions.TransferCasing,
		},
		{
			name:     "IgnoreToken",
			field:    "IgnoreToken",
			expected: (*regexp.Regexp)(nil),
			actual:   DefaultOptions.IgnoreToken,
		},
		{
			name:     "IgnoreNonWords",
			field:    "IgnoreNonWords",
			expected: false,
			actual:   DefaultOptions.IgnoreNonWords,
		},
		{
			name:     "IgnoreTermWithDigits",
			field:    "IgnoreTermWithDigits",
			expected: false,
			actual:   DefaultOptions.IgnoreTermWithDigits,
		},
		{
			name:     "SplitBySpace",
			field:    "SplitBySpace",
			expected: false,
			actual:   DefaultOptions.SplitBySpace,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.actual, tt.expected) {
				t.Errorf("DefaultOptions.%s = %v, want %v", tt.field, tt.actual, tt.expected)
			}
		})
	}
}

func TestOptionsStruct(t *testing.T) {
	opts := options{
		InitialCapacity:           1000,
		MaxDictionaryEditDistance: 3,
		PrefixLength:              10,
		CountThreshold:            5,
		DistanceAlgorithm:         edit_distance.Levenshtein,
		IncludeUnknown:            true,
		TransferCasing:            true,
		IgnoreToken:               regexp.MustCompile(`test`),
		IgnoreNonWords:            true,
		IgnoreTermWithDigits:      true,
		SplitBySpace:              true,
	}

	if opts.InitialCapacity != 1000 {
		t.Errorf("InitialCapacity = %d, want 1000", opts.InitialCapacity)
	}
	if opts.MaxDictionaryEditDistance != 3 {
		t.Errorf("MaxDictionaryEditDistance = %d, want 3", opts.MaxDictionaryEditDistance)
	}
	if opts.PrefixLength != 10 {
		t.Errorf("PrefixLength = %d, want 10", opts.PrefixLength)
	}
	if opts.CountThreshold != 5 {
		t.Errorf("CountThreshold = %d, want 5", opts.CountThreshold)
	}
	if opts.DistanceAlgorithm != edit_distance.Levenshtein {
		t.Errorf("DistanceAlgorithm = %v, want %v", opts.DistanceAlgorithm, edit_distance.Levenshtein)
	}
	if opts.IncludeUnknown != true {
		t.Errorf("IncludeUnknown = %v, want true", opts.IncludeUnknown)
	}
	if opts.TransferCasing != true {
		t.Errorf("TransferCasing = %v, want true", opts.TransferCasing)
	}
	if opts.IgnoreToken == nil {
		t.Error("IgnoreToken should not be nil")
	}
	if opts.IgnoreNonWords != true {
		t.Errorf("IgnoreNonWords = %v, want true", opts.IgnoreNonWords)
	}
	if opts.IgnoreTermWithDigits != true {
		t.Errorf("IgnoreTermWithDigits = %v, want true", opts.IgnoreTermWithDigits)
	}
	if opts.SplitBySpace != true {
		t.Errorf("SplitBySpace = %v, want true", opts.SplitBySpace)
	}
}

func TestFuncConfig_Apply(t *testing.T) {
	tests := []struct {
		name     string
		modifier func(*options)
		check    func(*options) bool
		desc     string
	}{
		{
			name: "modify InitialCapacity",
			modifier: func(o *options) {
				o.InitialCapacity = 500
			},
			check: func(o *options) bool {
				return o.InitialCapacity == 500
			},
			desc: "InitialCapacity should be 500",
		},
		{
			name: "modify MaxDictionaryEditDistance",
			modifier: func(o *options) {
				o.MaxDictionaryEditDistance = 4
			},
			check: func(o *options) bool {
				return o.MaxDictionaryEditDistance == 4
			},
			desc: "MaxDictionaryEditDistance should be 4",
		},
		{
			name: "modify PrefixLength",
			modifier: func(o *options) {
				o.PrefixLength = 15
			},
			check: func(o *options) bool {
				return o.PrefixLength == 15
			},
			desc: "PrefixLength should be 15",
		},
		{
			name: "modify IncludeUnknown",
			modifier: func(o *options) {
				o.IncludeUnknown = true
			},
			check: func(o *options) bool {
				return o.IncludeUnknown == true
			},
			desc: "IncludeUnknown should be true",
		},
		{
			name: "modify multiple fields",
			modifier: func(o *options) {
				o.InitialCapacity = 100
				o.TransferCasing = true
				o.IgnoreNonWords = true
			},
			check: func(o *options) bool {
				return o.InitialCapacity == 100 && o.TransferCasing == true && o.IgnoreNonWords == true
			},
			desc: "Multiple fields should be modified correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := FuncConfig{options: tt.modifier}
			opts := &options{}

			config.Apply(opts)

			if !tt.check(opts) {
				t.Errorf("FuncConfig.Apply() failed: %s", tt.desc)
			}
		})
	}
}

func TestNewFuncOption(t *testing.T) {
	tests := []struct {
		name     string
		modifier func(*options)
		expected interface{}
		getter   func(*options) interface{}
		desc     string
	}{
		{
			name: "create option for InitialCapacity",
			modifier: func(o *options) {
				o.InitialCapacity = 777
			},
			expected: 777,
			getter: func(o *options) interface{} {
				return o.InitialCapacity
			},
			desc: "InitialCapacity",
		},
		{
			name: "create option for CountThreshold",
			modifier: func(o *options) {
				o.CountThreshold = 25
			},
			expected: 25,
			getter: func(o *options) interface{} {
				return o.CountThreshold
			},
			desc: "CountThreshold",
		},
		{
			name: "create option for TransferCasing",
			modifier: func(o *options) {
				o.TransferCasing = true
			},
			expected: true,
			getter: func(o *options) interface{} {
				return o.TransferCasing
			},
			desc: "TransferCasing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := NewFuncOption(tt.modifier)

			if option == nil {
				t.Fatal("NewFuncOption returned nil")
			}

			var _ Options = option

			opts := &options{}
			option.Apply(opts)

			actual := tt.getter(opts)
			if actual != tt.expected {
				t.Errorf("NewFuncOption for %s: got %v, want %v", tt.desc, actual, tt.expected)
			}
		})
	}
}

func TestWithInitialCapacity(t *testing.T) {
	tests := []struct {
		name     string
		capacity int
		expected int
	}{
		{
			name:     "positive capacity",
			capacity: 1000,
			expected: 1000,
		},
		{
			name:     "zero capacity",
			capacity: 0,
			expected: 0,
		},
		{
			name:     "negative capacity",
			capacity: -100,
			expected: -100,
		},
		{
			name:     "large capacity",
			capacity: 1000000,
			expected: 1000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := WithInitialCapacity(tt.capacity)
			opts := &options{}
			option.Apply(opts)

			if opts.InitialCapacity != tt.expected {
				t.Errorf("WithInitialCapacity(%d) = %d, want %d", tt.capacity, opts.InitialCapacity, tt.expected)
			}
		})
	}
}

func TestWithMaxDictionaryEditDistance(t *testing.T) {
	tests := []struct {
		name     string
		distance int
		expected int
	}{
		{
			name:     "distance 1",
			distance: 1,
			expected: 1,
		},
		{
			name:     "distance 2",
			distance: 2,
			expected: 2,
		},
		{
			name:     "distance 3",
			distance: 3,
			expected: 3,
		},
		{
			name:     "zero distance",
			distance: 0,
			expected: 0,
		},
		{
			name:     "negative distance",
			distance: -1,
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := WithMaxDictionaryEditDistance(tt.distance)
			opts := &options{}
			option.Apply(opts)

			if opts.MaxDictionaryEditDistance != tt.expected {
				t.Errorf("WithMaxDictionaryEditDistance(%d) = %d, want %d", tt.distance, opts.MaxDictionaryEditDistance, tt.expected)
			}
		})
	}
}

func TestWithPrefixLength(t *testing.T) {
	tests := []struct {
		name     string
		length   int
		expected int
	}{
		{
			name:     "length 5",
			length:   5,
			expected: 5,
		},
		{
			name:     "length 7",
			length:   7,
			expected: 7,
		},
		{
			name:     "length 12",
			length:   12,
			expected: 12,
		},
		{
			name:     "zero length",
			length:   0,
			expected: 0,
		},
		{
			name:     "negative length",
			length:   -5,
			expected: -5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := WithPrefixLength(tt.length)
			opts := &options{}
			option.Apply(opts)

			if opts.PrefixLength != tt.expected {
				t.Errorf("WithPrefixLength(%d) = %d, want %d", tt.length, opts.PrefixLength, tt.expected)
			}
		})
	}
}

func TestWithCountThreshold(t *testing.T) {
	tests := []struct {
		name      string
		threshold int
		expected  int
	}{
		{
			name:      "threshold 1",
			threshold: 1,
			expected:  1,
		},
		{
			name:      "threshold 10",
			threshold: 10,
			expected:  10,
		},
		{
			name:      "threshold 100",
			threshold: 100,
			expected:  100,
		},
		{
			name:      "zero threshold",
			threshold: 0,
			expected:  0,
		},
		{
			name:      "negative threshold",
			threshold: -10,
			expected:  -10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := WithCountThreshold(tt.threshold)
			opts := &options{}
			option.Apply(opts)

			if opts.CountThreshold != tt.expected {
				t.Errorf("WithCountThreshold(%d) = %d, want %d", tt.threshold, opts.CountThreshold, tt.expected)
			}
		})
	}
}

func TestWithDistanceAlgorithm(t *testing.T) {
	tests := []struct {
		name      string
		algorithm edit_distance.AlgorithmType
		expected  edit_distance.AlgorithmType
	}{
		{
			name:      "DamerauOSAFast",
			algorithm: edit_distance.DamerauOSAFast,
			expected:  edit_distance.DamerauOSAFast,
		},
		{
			name:      "DamerauOSA",
			algorithm: edit_distance.DamerauOSA,
			expected:  edit_distance.DamerauOSA,
		},
		{
			name:      "DamerauLevenshtein",
			algorithm: edit_distance.DamerauLevenshtein,
			expected:  edit_distance.DamerauLevenshtein,
		},
		{
			name:      "Levenshtein",
			algorithm: edit_distance.Levenshtein,
			expected:  edit_distance.Levenshtein,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := WithDistanceAlgorithm(tt.algorithm)
			opts := &options{}
			option.Apply(opts)

			if opts.DistanceAlgorithm != tt.expected {
				t.Errorf("WithDistanceAlgorithm(%v) = %v, want %v", tt.algorithm, opts.DistanceAlgorithm, tt.expected)
			}
		})
	}
}

func TestWithIncludeUnknown(t *testing.T) {
	tests := []struct {
		name     string
		value    bool
		expected bool
	}{
		{
			name:     "include unknown true",
			value:    true,
			expected: true,
		},
		{
			name:     "include unknown false",
			value:    false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := WithIncludeUnknown(tt.value)
			opts := &options{}
			option.Apply(opts)

			if opts.IncludeUnknown != tt.expected {
				t.Errorf("WithIncludeUnknown(%v) = %v, want %v", tt.value, opts.IncludeUnknown, tt.expected)
			}
		})
	}
}

func TestWithTransferCasing(t *testing.T) {
	tests := []struct {
		name     string
		value    bool
		expected bool
	}{
		{
			name:     "transfer casing true",
			value:    true,
			expected: true,
		},
		{
			name:     "transfer casing false",
			value:    false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := WithTransferCasing(tt.value)
			opts := &options{}
			option.Apply(opts)

			if opts.TransferCasing != tt.expected {
				t.Errorf("WithTransferCasing(%v) = %v, want %v", tt.value, opts.TransferCasing, tt.expected)
			}
		})
	}
}

func TestWithIgnoreToken(t *testing.T) {
	tests := []struct {
		name     string
		pattern  *regexp.Regexp
		expected *regexp.Regexp
	}{
		{
			name:     "nil regexp",
			pattern:  nil,
			expected: nil,
		},
		{
			name:     "valid regexp",
			pattern:  regexp.MustCompile(`test`),
			expected: regexp.MustCompile(`test`),
		},
		{
			name:     "complex regexp",
			pattern:  regexp.MustCompile(`^(https?://|[^@]+@[^@]+\.[^@]+)`),
			expected: regexp.MustCompile(`^(https?://|[^@]+@[^@]+\.[^@]+)`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := WithIgnoreToken(tt.pattern)
			opts := &options{}
			option.Apply(opts)

			if tt.expected == nil {
				if opts.IgnoreToken != nil {
					t.Errorf("WithIgnoreToken(nil) should set IgnoreToken to nil")
				}
			} else {
				if opts.IgnoreToken == nil {
					t.Errorf("WithIgnoreToken should not set IgnoreToken to nil when given valid regexp")
				} else if opts.IgnoreToken.String() != tt.expected.String() {
					t.Errorf("WithIgnoreToken regexp pattern = %s, want %s", opts.IgnoreToken.String(), tt.expected.String())
				}
			}
		})
	}
}

func TestWithIgnoreNonWords(t *testing.T) {
	tests := []struct {
		name     string
		value    bool
		expected bool
	}{
		{
			name:     "ignore non-words true",
			value:    true,
			expected: true,
		},
		{
			name:     "ignore non-words false",
			value:    false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := WithIgnoreNonWords(tt.value)
			opts := &options{}
			option.Apply(opts)

			if opts.IgnoreNonWords != tt.expected {
				t.Errorf("WithIgnoreNonWords(%v) = %v, want %v", tt.value, opts.IgnoreNonWords, tt.expected)
			}
		})
	}
}

func TestWithIgnoreTermWithDigits(t *testing.T) {
	tests := []struct {
		name     string
		value    bool
		expected bool
	}{
		{
			name:     "ignore terms with digits true",
			value:    true,
			expected: true,
		},
		{
			name:     "ignore terms with digits false",
			value:    false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := WithIgnoreTermWithDigits(tt.value)
			opts := &options{}
			option.Apply(opts)

			if opts.IgnoreTermWithDigits != tt.expected {
				t.Errorf("WithIgnoreTermWithDigits(%v) = %v, want %v", tt.value, opts.IgnoreTermWithDigits, tt.expected)
			}
		})
	}
}

func TestWithSplitBySpace(t *testing.T) {
	tests := []struct {
		name     string
		value    bool
		expected bool
	}{
		{
			name:     "split by space true",
			value:    true,
			expected: true,
		},
		{
			name:     "split by space false",
			value:    false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := WithSplitBySpace(tt.value)
			opts := &options{}
			option.Apply(opts)

			if opts.SplitBySpace != tt.expected {
				t.Errorf("WithSplitBySpace(%v) = %v, want %v", tt.value, opts.SplitBySpace, tt.expected)
			}
		})
	}
}

func TestOptionsInterface(t *testing.T) {
	var _ Options = &FuncConfig{}
	var _ Options = NewFuncOption(func(*options) {})

	option := WithInitialCapacity(500)
	var _ Options = option

	opts := &options{}
	option.Apply(opts)

	if opts.InitialCapacity != 500 {
		t.Errorf("Options interface implementation failed")
	}
}

func TestFunctionalOptionsPattern(t *testing.T) {
	tests := []struct {
		name           string
		options        []Options
		expectedValues map[string]interface{}
	}{
		{
			name: "single option",
			options: []Options{
				WithInitialCapacity(100),
			},
			expectedValues: map[string]interface{}{
				"InitialCapacity": 100,
			},
		},
		{
			name: "multiple options",
			options: []Options{
				WithInitialCapacity(200),
				WithMaxDictionaryEditDistance(3),
				WithPrefixLength(10),
			},
			expectedValues: map[string]interface{}{
				"InitialCapacity":           200,
				"MaxDictionaryEditDistance": 3,
				"PrefixLength":              10,
			},
		},
		{
			name: "boolean options",
			options: []Options{
				WithIncludeUnknown(true),
				WithTransferCasing(true),
				WithIgnoreNonWords(false),
			},
			expectedValues: map[string]interface{}{
				"IncludeUnknown": true,
				"TransferCasing": true,
				"IgnoreNonWords": false,
			},
		},
		{
			name: "all option types",
			options: []Options{
				WithInitialCapacity(1000),
				WithMaxDictionaryEditDistance(2),
				WithPrefixLength(8),
				WithCountThreshold(5),
				WithDistanceAlgorithm(edit_distance.Levenshtein),
				WithIncludeUnknown(true),
				WithTransferCasing(false),
				WithIgnoreNonWords(true),
				WithIgnoreTermWithDigits(false),
				WithSplitBySpace(true),
			},
			expectedValues: map[string]interface{}{
				"InitialCapacity":           1000,
				"MaxDictionaryEditDistance": 2,
				"PrefixLength":              8,
				"CountThreshold":            5,
				"DistanceAlgorithm":         edit_distance.Levenshtein,
				"IncludeUnknown":            true,
				"TransferCasing":            false,
				"IgnoreNonWords":            true,
				"IgnoreTermWithDigits":      false,
				"SplitBySpace":              true,
			},
		},
		{
			name: "overriding options (last wins)",
			options: []Options{
				WithInitialCapacity(100),
				WithInitialCapacity(200),
				WithInitialCapacity(300),
			},
			expectedValues: map[string]interface{}{
				"InitialCapacity": 300,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &options{}

			for _, option := range tt.options {
				option.Apply(opts)
			}

			for field, expected := range tt.expectedValues {
				var actual interface{}
				switch field {
				case "InitialCapacity":
					actual = opts.InitialCapacity
				case "MaxDictionaryEditDistance":
					actual = opts.MaxDictionaryEditDistance
				case "PrefixLength":
					actual = opts.PrefixLength
				case "CountThreshold":
					actual = opts.CountThreshold
				case "DistanceAlgorithm":
					actual = opts.DistanceAlgorithm
				case "IncludeUnknown":
					actual = opts.IncludeUnknown
				case "TransferCasing":
					actual = opts.TransferCasing
				case "IgnoreNonWords":
					actual = opts.IgnoreNonWords
				case "IgnoreTermWithDigits":
					actual = opts.IgnoreTermWithDigits
				case "SplitBySpace":
					actual = opts.SplitBySpace
				default:
					t.Errorf("Unknown field in test: %s", field)
					continue
				}

				if !reflect.DeepEqual(actual, expected) {
					t.Errorf("Field %s = %v, want %v", field, actual, expected)
				}
			}
		})
	}
}

func TestOptionsWithRegexp(t *testing.T) {
	pattern := regexp.MustCompile(`^https?://`)
	option := WithIgnoreToken(pattern)

	opts := &options{}
	option.Apply(opts)

	if opts.IgnoreToken == nil {
		t.Fatal("IgnoreToken should not be nil")
	}

	if !opts.IgnoreToken.MatchString("https://example.com") {
		t.Error("Regexp should match https URL")
	}

	if opts.IgnoreToken.MatchString("ftp://example.com") {
		t.Error("Regexp should not match ftp URL")
	}
}

func TestOptionsChaining(t *testing.T) {
	opts := &options{}

	WithInitialCapacity(500).Apply(opts)
	WithMaxDictionaryEditDistance(3).Apply(opts)
	WithIncludeUnknown(true).Apply(opts)

	if opts.InitialCapacity != 500 {
		t.Errorf("InitialCapacity = %d, want 500", opts.InitialCapacity)
	}
	if opts.MaxDictionaryEditDistance != 3 {
		t.Errorf("MaxDictionaryEditDistance = %d, want 3", opts.MaxDictionaryEditDistance)
	}
	if opts.IncludeUnknown != true {
		t.Errorf("IncludeUnknown = %v, want true", opts.IncludeUnknown)
	}
}

func TestEmptyOptionsApplication(t *testing.T) {
	opts := &options{}

	emptyOption := NewFuncOption(func(*options) {})
	emptyOption.Apply(opts)

	if opts.InitialCapacity != 0 {
		t.Errorf("Empty option should not modify InitialCapacity")
	}
}

func TestOptionsStressTest(t *testing.T) {
	opts := &options{}

	for i := 0; i < 1000; i++ {
		WithInitialCapacity(i).Apply(opts)
	}

	if opts.InitialCapacity != 999 {
		t.Errorf("Final InitialCapacity = %d, want 999", opts.InitialCapacity)
	}
}
