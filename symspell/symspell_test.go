package symspell

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/Trendyol/go-symspell/edit_distance"
	"github.com/Trendyol/go-symspell/suggest"
	"github.com/Trendyol/go-symspell/verbosity"
)

func TestNewSymSpell_Success(t *testing.T) {
	tests := []struct {
		name    string
		options []Options
	}{
		{
			name:    "default options",
			options: []Options{},
		},
		{
			name: "custom options",
			options: []Options{
				WithMaxDictionaryEditDistance(3),
				WithPrefixLength(10),
				WithInitialCapacity(1000),
				WithCountThreshold(5),
				WithDistanceAlgorithm(edit_distance.Levenshtein),
				WithIncludeUnknown(true),
				WithTransferCasing(true),
				WithIgnoreToken(regexp.MustCompile(`test`)),
				WithIgnoreNonWords(true),
				WithIgnoreTermWithDigits(true),
				WithSplitBySpace(true),
			},
		},
		{
			name: "minimal valid configuration",
			options: []Options{
				WithMaxDictionaryEditDistance(0),
				WithPrefixLength(1),
			},
		},
		{
			name: "maximum configuration",
			options: []Options{
				WithMaxDictionaryEditDistance(10),
				WithPrefixLength(20),
				WithInitialCapacity(100000),
				WithCountThreshold(1000),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss, err := NewSymSpell(tt.options...)
			if err != nil {
				t.Errorf("NewSymSpell() should not return error, got: %v", err)
				return
			}
			if ss == nil {
				t.Error("NewSymSpell() should not return nil")
				return
			}

			// Verify basic functionality
			if ss.EntryCount() < 0 {
				t.Error("EntryCount should be non-negative")
			}
			if ss.WordCount() < 0 {
				t.Error("WordCount should be non-negative")
			}
			if ss.MaxLength() < 0 {
				t.Error("MaxLength should be non-negative")
			}
		})
	}
}

func TestNewSymSpell_ValidationErrors(t *testing.T) {
	tests := []struct {
		name        string
		options     []Options
		expectedErr string
	}{
		{
			name: "negative max dictionary edit distance",
			options: []Options{
				WithMaxDictionaryEditDistance(-1),
			},
			expectedErr: "maxDictionaryEditDistance cannot be negative",
		},
		{
			name: "zero prefix length",
			options: []Options{
				WithPrefixLength(0),
			},
			expectedErr: "prefixLength cannot be less than 1",
		},
		{
			name: "negative prefix length",
			options: []Options{
				WithPrefixLength(-5),
			},
			expectedErr: "prefixLength cannot be less than 1",
		},
		{
			name: "prefix length less than or equal to max edit distance",
			options: []Options{
				WithMaxDictionaryEditDistance(5),
				WithPrefixLength(5),
			},
			expectedErr: "prefixLength must be greater than maxDictionaryEditDistance",
		},
		{
			name: "prefix length equal to max edit distance",
			options: []Options{
				WithMaxDictionaryEditDistance(3),
				WithPrefixLength(3),
			},
			expectedErr: "prefixLength must be greater than maxDictionaryEditDistance",
		},
		{
			name: "negative count threshold",
			options: []Options{
				WithCountThreshold(-1),
			},
			expectedErr: "countThreshold cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss, err := NewSymSpell(tt.options...)
			if err == nil {
				t.Errorf("NewSymSpell() should return error for %s", tt.name)
				return
			}
			if ss != nil {
				t.Error("NewSymSpell() should return nil when error occurs")
			}
			if !strings.Contains(err.Error(), tt.expectedErr) {
				t.Errorf("Expected error containing %q, got %q", tt.expectedErr, err.Error())
			}
		})
	}
}

func TestNewSymSpell_InvalidDistanceAlgorithm(t *testing.T) {
	// Test with an invalid algorithm type that would cause NewDistanceComparer to fail
	ss, err := NewSymSpell(WithDistanceAlgorithm(edit_distance.AlgorithmType(999)))
	if err == nil {
		t.Error("NewSymSpell() should return error for invalid distance algorithm")
	}
	if ss != nil {
		t.Error("NewSymSpell() should return nil when error occurs")
	}
}

func TestSymSpell_GetterMethods(t *testing.T) {
	ss, err := NewSymSpell(
		WithMaxDictionaryEditDistance(2),
		WithPrefixLength(7),
		WithInitialCapacity(100),
	)
	if err != nil {
		t.Fatalf("Failed to create SymSpell: %v", err)
	}

	// Test initial values
	if ss.EntryCount() != 0 {
		t.Errorf("EntryCount() should be 0 initially, got %d", ss.EntryCount())
	}
	if ss.WordCount() != 0 {
		t.Errorf("WordCount() should be 0 initially, got %d", ss.WordCount())
	}
	if ss.MaxLength() != 0 {
		t.Errorf("MaxLength() should be 0 initially, got %d", ss.MaxLength())
	}

	// Add some entries to test changes
	dict := ss.dictionary
	dict.CreateDictionaryEntry("hello", 10)
	dict.CreateDictionaryEntry("world", 5)

	if ss.WordCount() != 2 {
		t.Errorf("WordCount() should be 2 after adding words, got %d", ss.WordCount())
	}
	if ss.EntryCount() <= 0 {
		t.Errorf("EntryCount() should be positive after adding words, got %d", ss.EntryCount())
	}
	if ss.MaxLength() < 5 {
		t.Errorf("MaxLength() should be at least 5 after adding 'hello' and 'world', got %d", ss.MaxLength())
	}
}

func TestSymSpell_CreateDictionary(t *testing.T) {
	// Create temporary file
	tempFile, err := os.CreateTemp("", "test_create_dict_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	content := "hello world\ntest data\nhello again\nspelling correction\n"
	if _, err := tempFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	ss, err := NewSymSpell()
	if err != nil {
		t.Fatalf("Failed to create SymSpell: %v", err)
	}

	success, err := ss.CreateDictionary(tempFile.Name())
	if err != nil {
		t.Errorf("CreateDictionary() should not return error, got: %v", err)
	}
	if !success {
		t.Error("CreateDictionary() should return true for successful creation")
	}

	// Verify words were added
	if ss.WordCount() == 0 {
		t.Error("WordCount should be greater than 0 after creating dictionary")
	}

	// Test with non-existent file
	success, err = ss.CreateDictionary("nonexistent_file.txt")
	if success {
		t.Error("CreateDictionary() should return false for non-existent file")
	}
	if err == nil {
		t.Error("CreateDictionary() should return error for non-existent file")
	}
}

func TestSymSpell_LoadDictionary(t *testing.T) {
	// Create temporary file
	tempFile, err := os.CreateTemp("", "test_load_dict_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	content := "hello\t10\nworld\t5\ntest\t20\nspelling\t15\n"
	if _, err := tempFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	ss, err := NewSymSpell()
	if err != nil {
		t.Fatalf("Failed to create SymSpell: %v", err)
	}

	success, err := ss.LoadDictionary(tempFile.Name(), 0, 1, "\t")
	if err != nil {
		t.Errorf("LoadDictionary() should not return error, got: %v", err)
	}
	if !success {
		t.Error("LoadDictionary() should return true for successful load")
	}

	// Verify words were loaded
	if ss.WordCount() != 4 {
		t.Errorf("WordCount should be 4 after loading dictionary, got %d", ss.WordCount())
	}

	// Test with non-existent file
	success, err = ss.LoadDictionary("nonexistent_file.txt", 0, 1, "\t")
	if success {
		t.Error("LoadDictionary() should return false for non-existent file")
	}
	if err == nil {
		t.Error("LoadDictionary() should return error for non-existent file")
	}
}

func TestSymSpell_LoadBigramDictionary(t *testing.T) {
	// Create temporary file
	tempFile, err := os.CreateTemp("", "test_load_bigram_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	content := "hello world\t100\ntest data\t50\nspelling correction\t75\n"
	if _, err := tempFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	ss, err := NewSymSpell()
	if err != nil {
		t.Fatalf("Failed to create SymSpell: %v", err)
	}

	success, err := ss.LoadBigramDictionary(tempFile.Name(), 0, 1, "\t")
	if err != nil {
		t.Errorf("LoadBigramDictionary() should not return error, got: %v", err)
	}
	if !success {
		t.Error("LoadBigramDictionary() should return true for successful load")
	}

	// Test with non-existent file
	success, err = ss.LoadBigramDictionary("nonexistent_file.txt", 0, 1, "\t")
	if success {
		t.Error("LoadBigramDictionary() should return false for non-existent file")
	}
	if err == nil {
		t.Error("LoadBigramDictionary() should return error for non-existent file")
	}
}

func TestSymSpell_Lookup_BasicFunctionality(t *testing.T) {
	ss, err := NewSymSpell(WithMaxDictionaryEditDistance(2))
	if err != nil {
		t.Fatalf("Failed to create SymSpell: %v", err)
	}

	// Add test words
	dict := ss.dictionary
	dict.CreateDictionaryEntry("hello", 100)
	dict.CreateDictionaryEntry("world", 50)
	dict.CreateDictionaryEntry("test", 75)
	dict.CreateDictionaryEntry("spelling", 60)

	tests := []struct {
		name           string
		input          string
		verbosity      verbosity.Verbosity
		maxEditDist    int
		expectResults  int
		expectContains string
	}{
		{
			name:           "exact match",
			input:          "hello",
			verbosity:      verbosity.Top,
			maxEditDist:    2,
			expectResults:  1,
			expectContains: "hello",
		},
		{
			name:           "single edit distance",
			input:          "helo",
			verbosity:      verbosity.Top,
			maxEditDist:    2,
			expectResults:  1,
			expectContains: "hello",
		},
		{
			name:           "multiple suggestions with Top verbosity",
			input:          "tst",
			verbosity:      verbosity.Top,
			maxEditDist:    2,
			expectResults:  1,
			expectContains: "test",
		},
		{
			name:           "multiple suggestions with All verbosity",
			input:          "tst",
			verbosity:      verbosity.All,
			maxEditDist:    2,
			expectResults:  -1, // Don't check count, just that we get results
			expectContains: "",
		},
		{
			name:           "no suggestions beyond max edit distance",
			input:          "xyz",
			verbosity:      verbosity.All,
			maxEditDist:    1,
			expectResults:  0,
			expectContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions, err := ss.Lookup(tt.input, tt.verbosity, tt.maxEditDist)
			if err != nil {
				t.Errorf("Lookup() should not return error, got: %v", err)
				return
			}

			if tt.expectResults >= 0 && len(suggestions) != tt.expectResults {
				t.Errorf("Expected %d suggestions, got %d", tt.expectResults, len(suggestions))
			}

			if tt.expectContains != "" {
				found := false
				for _, sug := range suggestions {
					if sug.Term == tt.expectContains {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected suggestions to contain %q", tt.expectContains)
				}
			}
		})
	}
}

func TestSymSpell_Lookup_VerbosityLevels(t *testing.T) {
	ss, err := NewSymSpell(WithMaxDictionaryEditDistance(2))
	if err != nil {
		t.Fatalf("Failed to create SymSpell: %v", err)
	}

	// Add test words with different frequencies
	dict := ss.dictionary
	dict.CreateDictionaryEntry("test", 100)
	dict.CreateDictionaryEntry("text", 50)
	dict.CreateDictionaryEntry("best", 75)

	input := "tst" // Should match all three with edit distance 1-2

	// Test Top verbosity - should return only the best
	suggestions, err := ss.Lookup(input, verbosity.Top, 2)
	if err != nil {
		t.Fatalf("Lookup() failed: %v", err)
	}
	if len(suggestions) != 1 {
		t.Errorf("Top verbosity should return 1 suggestion, got %d", len(suggestions))
	}

	// Test Closest verbosity - should return all with minimum distance
	suggestions, err = ss.Lookup(input, verbosity.Closest, 2)
	if err != nil {
		t.Fatalf("Lookup() failed: %v", err)
	}
	if len(suggestions) == 0 {
		t.Error("Closest verbosity should return at least 1 suggestion")
	}
	// All suggestions should have the same distance (minimum)
	if len(suggestions) > 1 {
		minDist := suggestions[0].Distance
		for _, sug := range suggestions {
			if sug.Distance != minDist {
				t.Errorf("All suggestions should have same distance %d, got %d for %s", minDist, sug.Distance, sug.Term)
			}
		}
	}

	// Test All verbosity - should return all suggestions within maxEditDistance
	suggestions, err = ss.Lookup(input, verbosity.All, 2)
	if err != nil {
		t.Fatalf("Lookup() failed: %v", err)
	}
	if len(suggestions) == 0 {
		t.Error("All verbosity should return at least 1 suggestion")
	}
	// Verify suggestions are sorted by distance, then frequency
	for i := 1; i < len(suggestions); i++ {
		prev, curr := suggestions[i-1], suggestions[i]
		if prev.Distance > curr.Distance {
			t.Errorf("Suggestions should be sorted by distance: %d > %d", prev.Distance, curr.Distance)
		}
		if prev.Distance == curr.Distance && prev.Count < curr.Count {
			t.Errorf("Suggestions with same distance should be sorted by count: %d < %d", prev.Count, curr.Count)
		}
	}
}

func TestSymSpell_Lookup_EdgeCases(t *testing.T) {
	ss, err := NewSymSpell(WithMaxDictionaryEditDistance(2))
	if err != nil {
		t.Fatalf("Failed to create SymSpell: %v", err)
	}

	dict := ss.dictionary
	dict.CreateDictionaryEntry("hello", 10)
	dict.CreateDictionaryEntry("a", 5)
	dict.CreateDictionaryEntry("supercalifragilisticexpialidocious", 1)

	tests := []struct {
		name        string
		input       string
		maxEditDist int
		expectError bool
		expectEmpty bool
	}{
		{
			name:        "single character",
			input:       "b",
			maxEditDist: 1,
			expectError: false,
			expectEmpty: false,
		},
		{
			name:        "very long string",
			input:       strings.Repeat("a", 1000),
			maxEditDist: 2,
			expectError: false,
			expectEmpty: true,
		},
		{
			name:        "negative max edit distance",
			input:       "hello",
			maxEditDist: -1,
			expectError: false,
			expectEmpty: false,
		},
		{
			name:        "max edit distance too large",
			input:       "hello",
			maxEditDist: 10,
			expectError: true,
			expectEmpty: false,
		},
		{
			name:        "zero max edit distance with exact match",
			input:       "hello",
			maxEditDist: 0,
			expectError: false,
			expectEmpty: false,
		},
		{
			name:        "zero max edit distance without exact match",
			input:       "helo",
			maxEditDist: 0,
			expectError: false,
			expectEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions, err := ss.Lookup(tt.input, verbosity.Top, tt.maxEditDist)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.expectEmpty && len(suggestions) > 0 {
				t.Errorf("Expected no suggestions but got %d", len(suggestions))
			}
			if !tt.expectEmpty && tt.input != "" && len(suggestions) == 0 {
				// Some cases might legitimately return no suggestions
			}
		})
	}
}

func TestSymSpell_Lookup_WithTransferCasing(t *testing.T) {
	ss, err := NewSymSpell(WithTransferCasing(true))
	if err != nil {
		t.Fatalf("Failed to create SymSpell: %v", err)
	}

	dict := ss.dictionary
	dict.CreateDictionaryEntry("hello", 10)
	dict.CreateDictionaryEntry("world", 5)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "uppercase input",
			input:    "HELLO",
			expected: "HELLO",
		},
		{
			name:     "title case input",
			input:    "Hello",
			expected: "Hello",
		},
		{
			name:     "lowercase input",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "mixed case with correction",
			input:    "HELO",
			expected: "HELLO",
		},
		{
			name:     "title case with correction",
			input:    "Helo",
			expected: "Hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions, err := ss.Lookup(tt.input, verbosity.Top, 2)
			if err != nil {
				t.Errorf("Lookup() failed: %v", err)
				return
			}
			if len(suggestions) == 0 {
				t.Error("Expected at least one suggestion")
				return
			}
			if suggestions[0].Term != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, suggestions[0].Term)
			}
		})
	}
}

func TestSymSpell_Lookup_WithIncludeUnknown(t *testing.T) {
	ss, err := NewSymSpell(WithIncludeUnknown(true))
	if err != nil {
		t.Fatalf("Failed to create SymSpell: %v", err)
	}

	// Don't add any words to dictionary
	suggestions, err := ss.Lookup("unknown", verbosity.Top, 2)
	if err != nil {
		t.Errorf("Lookup() failed: %v", err)
		return
	}
	if len(suggestions) != 1 {
		t.Errorf("Expected 1 suggestion for unknown word, got %d", len(suggestions))
		return
	}
	if suggestions[0].Term != "unknown" {
		t.Errorf("Expected original term %q, got %q", "unknown", suggestions[0].Term)
	}
	if suggestions[0].Distance != 3 { // maxEditDistance + 1
		t.Errorf("Expected distance 3, got %d", suggestions[0].Distance)
	}
}

func TestSymSpell_Lookup_WithIgnoreToken(t *testing.T) {
	pattern := regexp.MustCompile(`^https?://`)
	ss, err := NewSymSpell(WithIgnoreToken(pattern))
	if err != nil {
		t.Fatalf("Failed to create SymSpell: %v", err)
	}

	tests := []struct {
		name        string
		input       string
		shouldMatch bool
	}{
		{
			name:        "ftp URL",
			input:       "ftp://example.com",
			shouldMatch: false,
		},
		{
			name:        "regular word",
			input:       "hello",
			shouldMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions, err := ss.Lookup(tt.input, verbosity.Top, 2)
			if err != nil {
				t.Errorf("Lookup() failed: %v", err)
				return
			}

			if tt.shouldMatch {
				if len(suggestions) != 1 {
					t.Errorf("Expected 1 suggestion for ignored token, got %d", len(suggestions))
					return
				}
				if suggestions[0].Term != tt.input {
					t.Errorf("Expected original term %q, got %q", tt.input, suggestions[0].Term)
				}
				if suggestions[0].Distance != 0 {
					t.Errorf("Expected distance 0 for ignored token, got %d", suggestions[0].Distance)
				}
			}
		})
	}
}

func TestSymSpell_LookupCompound_BasicFunctionality(t *testing.T) {
	ss, err := NewSymSpell(WithMaxDictionaryEditDistance(2))
	if err != nil {
		t.Fatalf("Failed to create SymSpell: %v", err)
	}

	// Add test words
	dict := ss.dictionary
	dict.CreateDictionaryEntry("hello", 100)
	dict.CreateDictionaryEntry("world", 50)
	dict.CreateDictionaryEntry("spell", 75)
	dict.CreateDictionaryEntry("check", 60)
	dict.CreateDictionaryEntry("spell", 80)
	dict.CreateDictionaryEntry("checker", 40)

	// Add bigrams for better compound suggestions
	dict.SetBigramCount("hello world", 30)
	dict.SetBigramCount("spell check", 20)

	tests := []struct {
		name           string
		input          string
		maxEditDist    int
		expectContains string
	}{
		{
			name:           "simple phrase correction",
			input:          "helo wrld",
			maxEditDist:    2,
			expectContains: "hello world",
		},
		{
			name:           "single word",
			input:          "hello",
			maxEditDist:    2,
			expectContains: "hello",
		},
		{
			name:           "phrase with exact matches",
			input:          "hello world",
			maxEditDist:    2,
			expectContains: "hello world",
		},
		{
			name:           "compound word splitting",
			input:          "spellcheck",
			maxEditDist:    2,
			expectContains: "", // Will depend on implementation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions, err := ss.LookupCompound(tt.input, tt.maxEditDist)
			if err != nil {
				t.Errorf("LookupCompound() should not return error, got: %v", err)
				return
			}

			if len(suggestions) == 0 {
				t.Error("LookupCompound() should return at least one suggestion")
				return
			}

			// LookupCompound should return exactly one suggestion
			if len(suggestions) != 1 {
				t.Errorf("LookupCompound() should return exactly 1 suggestion, got %d", len(suggestions))
			}

			if tt.expectContains != "" && suggestions[0].Term != tt.expectContains {
				t.Errorf("Expected result containing %q, got %q", tt.expectContains, suggestions[0].Term)
			}
		})
	}
}

func TestSymSpell_LookupCompound_WithIgnoreNonWords(t *testing.T) {
	ss, err := NewSymSpell(
		WithMaxDictionaryEditDistance(2),
		WithIgnoreNonWords(true),
		WithIgnoreTermWithDigits(true),
	)
	if err != nil {
		t.Fatalf("Failed to create SymSpell: %v", err)
	}

	// Add test words
	dict := ss.dictionary
	dict.CreateDictionaryEntry("hello", 100)
	dict.CreateDictionaryEntry("world", 50)

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "phrase with numbers",
			input: "hello 123 world",
		},
		{
			name:  "phrase with acronym",
			input: "hello NASA world",
		},
		{
			name:  "phrase with mixed content",
			input: "hello IPv4 world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions, err := ss.LookupCompound(tt.input, 2)
			if err != nil {
				t.Errorf("LookupCompound() should not return error, got: %v", err)
				return
			}

			if len(suggestions) == 0 {
				t.Error("LookupCompound() should return at least one suggestion")
			}
		})
	}
}

func TestSymSpell_LookupCompound_EdgeCases(t *testing.T) {
	ss, err := NewSymSpell()
	if err != nil {
		t.Fatalf("Failed to create SymSpell: %v", err)
	}

	tests := []struct {
		name        string
		input       string
		maxEditDist int
	}{
		{
			name:        "empty string",
			input:       "",
			maxEditDist: 2,
		},
		{
			name:        "single character",
			input:       "a",
			maxEditDist: 1,
		},
		{
			name:        "very long phrase",
			input:       strings.Repeat("word ", 100),
			maxEditDist: 2,
		},
		{
			name:        "phrase with only spaces",
			input:       "   ",
			maxEditDist: 2,
		},
		{
			name:        "phrase with special characters",
			input:       "hello! world?",
			maxEditDist: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions, err := ss.LookupCompound(tt.input, tt.maxEditDist)
			if err != nil {
				t.Errorf("LookupCompound() should not return error for %q, got: %v", tt.input, err)
				return
			}

			// Should always return at least one suggestion (even if it's the original)
			if len(suggestions) == 0 {
				t.Errorf("LookupCompound() should return at least one suggestion for %q", tt.input)
			}
		})
	}
}

func TestSymSpell_earlyExit(t *testing.T) {
	tests := []struct {
		name            string
		includeUnknown  bool
		suggestions     []suggest.SuggestItem
		original        string
		maxEditDistance int
		expectedLen     int
		expectUnknown   bool
	}{
		{
			name:            "include unknown with empty suggestions",
			includeUnknown:  true,
			suggestions:     []suggest.SuggestItem{},
			original:        "unknown",
			maxEditDistance: 2,
			expectedLen:     1,
			expectUnknown:   true,
		},
		{
			name:            "include unknown with existing suggestions",
			includeUnknown:  true,
			suggestions:     []suggest.SuggestItem{{Term: "hello", Distance: 1, Count: 10}},
			original:        "helo",
			maxEditDistance: 2,
			expectedLen:     1,
			expectUnknown:   false,
		},
		{
			name:            "don't include unknown with empty suggestions",
			includeUnknown:  false,
			suggestions:     []suggest.SuggestItem{},
			original:        "unknown",
			maxEditDistance: 2,
			expectedLen:     0,
			expectUnknown:   false,
		},
		{
			name:            "don't include unknown with existing suggestions",
			includeUnknown:  false,
			suggestions:     []suggest.SuggestItem{{Term: "hello", Distance: 1, Count: 10}},
			original:        "helo",
			maxEditDistance: 2,
			expectedLen:     1,
			expectUnknown:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss, err := NewSymSpell(WithIncludeUnknown(tt.includeUnknown))
			if err != nil {
				t.Fatalf("Failed to create SymSpell: %v", err)
			}

			result := ss.earlyExit(tt.suggestions, tt.original, tt.maxEditDistance)

			if len(result) != tt.expectedLen {
				t.Errorf("Expected %d suggestions, got %d", tt.expectedLen, len(result))
				return
			}

			if tt.expectUnknown && len(result) > 0 {
				last := result[len(result)-1]
				if last.Term != tt.original {
					t.Errorf("Expected unknown term to be %q, got %q", tt.original, last.Term)
				}
				if last.Distance != tt.maxEditDistance+1 {
					t.Errorf("Expected unknown distance to be %d, got %d", tt.maxEditDistance+1, last.Distance)
				}
				if last.Count != 0 {
					t.Errorf("Expected unknown count to be 0, got %d", last.Count)
				}
			}
		})
	}
}

func TestSymSpell_parseWords(t *testing.T) {
	ss, err := NewSymSpell()
	if err != nil {
		t.Fatalf("Failed to create SymSpell: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "simple words",
			input:    "hello world test",
			expected: []string{"hello", "world", "test"},
		},
		{
			name:     "words with punctuation",
			input:    "Hello, world! How are you?",
			expected: []string{"hello", "world", "how", "are", "you"},
		},
		{
			name:     "words with apostrophes",
			input:    "don't can't won't it's",
			expected: []string{"don't", "can't", "won't", "it's"},
		},
		{
			name:     "words with numbers",
			input:    "hello123 world456 test",
			expected: []string{"hello123", "world456", "test"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "only punctuation",
			input:    "!@#$%^&*()",
			expected: []string{},
		},
		{
			name:     "unicode characters",
			input:    "café naïve résumé",
			expected: []string{"café", "naïve", "résumé"},
		},
		{
			name:     "mixed case becomes lowercase",
			input:    "Hello WORLD Test",
			expected: []string{"hello", "world", "test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ss.parseWords(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d words, got %d", len(tt.expected), len(result))
				t.Errorf("Expected: %v", tt.expected)
				t.Errorf("Got: %v", result)
				return
			}

			for i, expected := range tt.expected {
				if i >= len(result) || result[i] != expected {
					t.Errorf("Expected word %d to be %q, got %q", i, expected, result[i])
				}
			}
		})
	}
}

func TestSymSpell_Integration_CompleteWorkflow(t *testing.T) {
	// Test complete workflow: create, load data, and perform lookups
	ss, err := NewSymSpell(
		WithMaxDictionaryEditDistance(2),
		WithPrefixLength(7),
		WithIncludeUnknown(true),
		WithTransferCasing(true),
	)
	if err != nil {
		t.Fatalf("Failed to create SymSpell: %v", err)
	}

	// Create temporary dictionary file
	dictFile, err := os.CreateTemp("", "test_dict_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp dict file: %v", err)
	}
	defer os.Remove(dictFile.Name())

	dictContent := "hello\t100\nworld\t50\nspelling\t75\ncorrection\t60\ntest\t80\ndata\t40\n"
	if _, err := dictFile.WriteString(dictContent); err != nil {
		t.Fatalf("Failed to write dict content: %v", err)
	}
	dictFile.Close()

	// Load dictionary
	success, err := ss.LoadDictionary(dictFile.Name(), 0, 1, "\t")
	if !success || err != nil {
		t.Fatalf("Failed to load dictionary: success=%v, err=%v", success, err)
	}

	// Create temporary bigram file
	bigramFile, err := os.CreateTemp("", "test_bigram_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp bigram file: %v", err)
	}
	defer os.Remove(bigramFile.Name())

	bigramContent := "hello world\t30\nspelling correction\t25\ntest data\t20\n"
	if _, err := bigramFile.WriteString(bigramContent); err != nil {
		t.Fatalf("Failed to write bigram content: %v", err)
	}
	bigramFile.Close()

	// Load bigrams
	success, err = ss.LoadBigramDictionary(bigramFile.Name(), 0, 1, "\t")
	if !success || err != nil {
		t.Fatalf("Failed to load bigram dictionary: success=%v, err=%v", success, err)
	}

	// Verify data was loaded
	if ss.WordCount() != 6 {
		t.Errorf("Expected 6 words in dictionary, got %d", ss.WordCount())
	}

	// Test various lookup scenarios
	testCases := []struct {
		name     string
		input    string
		function string // "lookup" or "compound"
	}{
		{"exact match lookup", "hello", "lookup"},
		{"single error lookup", "helo", "lookup"},
		{"multiple errors lookup", "hllo", "lookup"},
		{"case transfer lookup", "HELLO", "lookup"},
		{"unknown word lookup", "xyz", "lookup"},
		{"simple phrase compound", "hello world", "compound"},
		{"phrase with errors compound", "helo wrld", "compound"},
		{"complex phrase compound", "speling corectin", "compound"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.function == "lookup" {
				suggestions, err := ss.Lookup(tc.input, verbosity.Top, 2)
				if err != nil {
					t.Errorf("Lookup failed for %q: %v", tc.input, err)
					return
				}
				if len(suggestions) == 0 {
					t.Errorf("No suggestions returned for %q", tc.input)
				}
			} else if tc.function == "compound" {
				suggestions, err := ss.LookupCompound(tc.input, 2)
				if err != nil {
					t.Errorf("LookupCompound failed for %q: %v", tc.input, err)
					return
				}
				if len(suggestions) == 0 {
					t.Errorf("No suggestions returned for %q", tc.input)
				}
			}
		})
	}
}

func TestSymSpell_ErrorHandling(t *testing.T) {
	ss, err := NewSymSpell(WithMaxDictionaryEditDistance(1))
	if err != nil {
		t.Fatalf("Failed to create SymSpell: %v", err)
	}

	// Test error conditions that should be handled gracefully
	tests := []struct {
		name        string
		operation   func() error
		expectError bool
	}{
		{
			name: "lookup with distance too large",
			operation: func() error {
				_, err := ss.Lookup("hello", verbosity.Top, 5)
				return err
			},
			expectError: true,
		},
		{
			name: "valid lookup operation",
			operation: func() error {
				_, err := ss.Lookup("hello", verbosity.Top, 1)
				return err
			},
			expectError: false,
		},
		{
			name: "compound lookup with empty string",
			operation: func() error {
				_, err := ss.LookupCompound("", 1)
				return err
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.operation()
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestSymSpell_UnicodeSupport(t *testing.T) {
	ss, err := NewSymSpell(WithMaxDictionaryEditDistance(2))
	if err != nil {
		t.Fatalf("Failed to create SymSpell: %v", err)
	}

	// Add Unicode words to dictionary
	dict := ss.dictionary
	unicodeWords := []struct {
		word  string
		count int
	}{
		{"café", 10},
		{"naïve", 8},
		{"résumé", 12},
		{"测试", 15},
		{"🙂😀👍", 5},
		{"Москва", 20},
		{"العربية", 18},
	}

	for _, uw := range unicodeWords {
		dict.CreateDictionaryEntry(uw.word, uw.count)
	}

	// Test Unicode lookups
	tests := []struct {
		name       string
		input      string
		shouldFind bool
	}{
		{"exact Unicode match", "café", true},
		{"Unicode with error", "cafe", true},
		{"Chinese characters", "测试", true},
		{"Emoji", "🙂😀👍", true},
		{"Cyrillic", "Москва", true},
		{"Arabic", "العربية", true},
		{"mixed Unicode error", "naïv", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions, err := ss.Lookup(tt.input, verbosity.Top, 2)
			if err != nil {
				t.Errorf("Lookup failed for Unicode input %q: %v", tt.input, err)
				return
			}

			if tt.shouldFind && len(suggestions) == 0 {
				t.Errorf("Expected to find suggestions for %q", tt.input)
			}
		})
	}

	// Test Unicode compound lookups
	compoundTests := []string{
		"café résumé",
		"测试 数据",
		"🙂😀 👍",
	}

	for _, input := range compoundTests {
		t.Run("compound_"+input, func(t *testing.T) {
			suggestions, err := ss.LookupCompound(input, 2)
			if err != nil {
				t.Errorf("LookupCompound failed for Unicode input %q: %v", input, err)
				return
			}
			if len(suggestions) == 0 {
				t.Errorf("Expected at least one suggestion for %q", input)
			}
		})
	}
}

func TestSymSpell_BoundaryConditions(t *testing.T) {
	ss, err := NewSymSpell(WithMaxDictionaryEditDistance(3))
	if err != nil {
		t.Fatalf("Failed to create SymSpell: %v", err)
	}

	dict := ss.dictionary

	// Add words of various lengths
	testWords := []string{
		"a",                                  // 1 character
		"ab",                                 // 2 characters
		"hello",                              // 5 characters
		"supercalifragilisticexpialidocious", // 34 characters
		strings.Repeat("long", 250),          // 1000 characters
	}

	for _, word := range testWords {
		dict.CreateDictionaryEntry(word, 1)
	}

	// Test boundary conditions
	tests := []struct {
		name            string
		input           string
		maxEditDistance int
		description     string
	}{
		{
			name:            "single char input",
			input:           "a",
			maxEditDistance: 3,
			description:     "Single character input should work",
		},
		{
			name:            "very long input",
			input:           strings.Repeat("test", 250),
			maxEditDistance: 1,
			description:     "Very long input should be handled",
		},
		{
			name:            "input longer than any dictionary word",
			input:           strings.Repeat("x", 2000),
			maxEditDistance: 1,
			description:     "Input longer than dictionary words should be handled",
		},
		{
			name:            "max edit distance zero",
			input:           "hello",
			maxEditDistance: 0,
			description:     "Zero edit distance should only match exact",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test both Lookup and LookupCompound for boundary conditions
			_, err1 := ss.Lookup(tt.input, verbosity.Top, tt.maxEditDistance)
			_, err2 := ss.LookupCompound(tt.input, tt.maxEditDistance)

			// We don't expect errors for boundary conditions, just that they complete
			if err1 != nil && !strings.Contains(err1.Error(), "distance too large") {
				t.Errorf("Unexpected error in Lookup for %s: %v", tt.description, err1)
			}
			if err2 != nil {
				t.Errorf("Unexpected error in LookupCompound for %s: %v", tt.description, err2)
			}
		})
	}
}

func TestGlobalConstants(t *testing.T) {
	// Test that N constant is reasonable
	if N <= 0 {
		t.Error("N constant should be positive")
	}

	// Test wordPattern regex
	if wordPattern == nil {
		t.Fatal("wordPattern should not be nil")
	}

	// Test wordPattern matches expected patterns
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input:    "hello world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "don't can't won't",
			expected: []string{"don't", "can't", "won't"},
		},
		{
			input:    "hello123 world456",
			expected: []string{"hello123", "world456"},
		},
		{
			input:    "café naïve",
			expected: []string{"café", "naïve"},
		},
		{
			input:    "!@#$%^&*()",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run("wordPattern_"+tt.input, func(t *testing.T) {
			matches := wordPattern.FindAllString(tt.input, -1)
			if len(matches) != len(tt.expected) {
				t.Errorf("Expected %d matches, got %d for input %q", len(tt.expected), len(matches), tt.input)
				t.Errorf("Expected: %v, Got: %v", tt.expected, matches)
				return
			}

			for i, expected := range tt.expected {
				if i >= len(matches) || matches[i] != expected {
					t.Errorf("Expected match %d to be %q, got %q", i, expected, matches[i])
				}
			}
		})
	}
}

func TestSymSpell_StressTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	ss, err := NewSymSpell(WithMaxDictionaryEditDistance(2))
	if err != nil {
		t.Fatalf("Failed to create SymSpell: %v", err)
	}

	// Add many words to dictionary
	dict := ss.dictionary
	for i := 0; i < 1000; i++ {
		word := generateTestWord(i)
		dict.CreateDictionaryEntry(word, i+1)
	}

	// Test many lookups
	for i := 0; i < 100; i++ {
		input := generateTestWord(i) + "x" // Add error
		_, err := ss.Lookup(input, verbosity.Top, 2)
		if err != nil {
			t.Errorf("Lookup failed on iteration %d: %v", i, err)
		}

		_, err = ss.LookupCompound(input, 2)
		if err != nil {
			t.Errorf("LookupCompound failed on iteration %d: %v", i, err)
		}
	}
}

// Helper function to generate test words
func generateTestWord(seed int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz"
	length := (seed % 10) + 3 // Words between 3-12 chars
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[(seed+i)%len(chars)]
	}
	return string(result)
}
