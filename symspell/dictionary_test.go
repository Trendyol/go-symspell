package symspell

import (
	"fmt"
	"math"
	"os"
	"strings"
	"testing"
)

func TestNewDictItem(t *testing.T) {
	tests := []struct {
		name          string
		word          string
		count         int
		expectedText  string
		expectedCount int
		expectedLen   int
	}{
		{
			name:          "simple word",
			word:          "hello",
			count:         5,
			expectedText:  "hello",
			expectedCount: 5,
			expectedLen:   5,
		},
		{
			name:          "empty string",
			word:          "",
			count:         1,
			expectedText:  "",
			expectedCount: 1,
			expectedLen:   0,
		},
		{
			name:          "unicode word",
			word:          "café",
			count:         10,
			expectedText:  "café",
			expectedCount: 10,
			expectedLen:   4,
		},
		{
			name:          "emoji",
			word:          "🙂😀",
			count:         3,
			expectedText:  "🙂😀",
			expectedCount: 3,
			expectedLen:   2,
		},
		{
			name:          "zero count",
			word:          "test",
			count:         0,
			expectedText:  "test",
			expectedCount: 0,
			expectedLen:   4,
		},
		{
			name:          "negative count",
			word:          "negative",
			count:         -5,
			expectedText:  "negative",
			expectedCount: -5,
			expectedLen:   8,
		},
		{
			name:          "long word",
			word:          "supercalifragilisticexpialidocious",
			count:         1,
			expectedText:  "supercalifragilisticexpialidocious",
			expectedCount: 1,
			expectedLen:   34,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := NewDictItem(tt.word, tt.count)

			if item == nil {
				t.Fatal("NewDictItem returned nil")
			}

			if item.text != tt.expectedText {
				t.Errorf("DictItem.text = %q, want %q", item.text, tt.expectedText)
			}

			if item.count != tt.expectedCount {
				t.Errorf("DictItem.count = %d, want %d", item.count, tt.expectedCount)
			}

			if item.Len() != tt.expectedLen {
				t.Errorf("DictItem.Len() = %d, want %d", item.Len(), tt.expectedLen)
			}

			if len(item.runes) != tt.expectedLen {
				t.Errorf("len(DictItem.runes) = %d, want %d", len(item.runes), tt.expectedLen)
			}

			// Verify runes match the original text
			if string(item.runes) != tt.word {
				t.Errorf("string(DictItem.runes) = %q, want %q", string(item.runes), tt.word)
			}
		})
	}
}

func TestDictItem_Len(t *testing.T) {
	tests := []struct {
		name        string
		word        string
		expectedLen int
	}{
		{
			name:        "ascii word",
			word:        "hello",
			expectedLen: 5,
		},
		{
			name:        "empty string",
			word:        "",
			expectedLen: 0,
		},
		{
			name:        "single character",
			word:        "a",
			expectedLen: 1,
		},
		{
			name:        "unicode characters",
			word:        "naïve",
			expectedLen: 5,
		},
		{
			name:        "chinese characters",
			word:        "测试",
			expectedLen: 2,
		},
		{
			name:        "emoji",
			word:        "🙂",
			expectedLen: 1,
		},
		{
			name:        "mixed unicode",
			word:        "hello世界",
			expectedLen: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := NewDictItem(tt.word, 1)
			result := item.Len()

			if result != tt.expectedLen {
				t.Errorf("DictItem.Len() = %d, want %d for word %q", result, tt.expectedLen, tt.word)
			}

			// Verify Len() returns the same as len(runes)
			if result != len(item.runes) {
				t.Errorf("DictItem.Len() = %d, but len(runes) = %d", result, len(item.runes))
			}

			// Verify it matches the actual rune count
			actualRuneCount := len([]rune(tt.word))
			if result != actualRuneCount {
				t.Errorf("DictItem.Len() = %d, but actual rune count = %d", result, actualRuneCount)
			}
		})
	}
}

func TestNewDictionary(t *testing.T) {
	tests := []struct {
		name                    string
		initialCapacity         int
		prefixLength            int
		maxDictEdit             int
		countThreshold          int
		expectedInitialCapacity int
		expectedPrefixLength    int
		expectedMaxDictEdit     int
		expectedCountThreshold  int
		expectedMaxLength       int
		expectedBigramCountMin  int
	}{
		{
			name:                    "default values",
			initialCapacity:         16,
			prefixLength:            7,
			maxDictEdit:             2,
			countThreshold:          1,
			expectedInitialCapacity: 16,
			expectedPrefixLength:    7,
			expectedMaxDictEdit:     2,
			expectedCountThreshold:  1,
			expectedMaxLength:       0,
			expectedBigramCountMin:  math.MaxInt,
		},
		{
			name:                    "large values",
			initialCapacity:         10000,
			prefixLength:            15,
			maxDictEdit:             5,
			countThreshold:          100,
			expectedInitialCapacity: 10000,
			expectedPrefixLength:    15,
			expectedMaxDictEdit:     5,
			expectedCountThreshold:  100,
			expectedMaxLength:       0,
			expectedBigramCountMin:  math.MaxInt,
		},
		{
			name:                    "zero values",
			initialCapacity:         0,
			prefixLength:            0,
			maxDictEdit:             0,
			countThreshold:          0,
			expectedInitialCapacity: 0,
			expectedPrefixLength:    0,
			expectedMaxDictEdit:     0,
			expectedCountThreshold:  0,
			expectedMaxLength:       0,
			expectedBigramCountMin:  math.MaxInt,
		},
		{
			name:                    "negative values",
			initialCapacity:         -10,
			prefixLength:            -5,
			maxDictEdit:             -2,
			countThreshold:          -1,
			expectedInitialCapacity: -10,
			expectedPrefixLength:    -5,
			expectedMaxDictEdit:     -2,
			expectedCountThreshold:  -1,
			expectedMaxLength:       0,
			expectedBigramCountMin:  math.MaxInt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dict := NewDictionary(tt.initialCapacity, tt.prefixLength, tt.maxDictEdit, tt.countThreshold)

			if dict == nil {
				t.Fatal("NewDictionary returned nil")
			}

			// Test getter methods
			if dict.PrefixLength() != tt.expectedPrefixLength {
				t.Errorf("PrefixLength() = %d, want %d", dict.PrefixLength(), tt.expectedPrefixLength)
			}

			if dict.MaxDictionaryEditDistance() != tt.expectedMaxDictEdit {
				t.Errorf("MaxDictionaryEditDistance() = %d, want %d", dict.MaxDictionaryEditDistance(), tt.expectedMaxDictEdit)
			}

			if dict.CountThreshold() != tt.expectedCountThreshold {
				t.Errorf("CountThreshold() = %d, want %d", dict.CountThreshold(), tt.expectedCountThreshold)
			}

			if dict.MaxLength() != tt.expectedMaxLength {
				t.Errorf("MaxLength() = %d, want %d", dict.MaxLength(), tt.expectedMaxLength)
			}

			if dict.BigramCountMin() != tt.expectedBigramCountMin {
				t.Errorf("BigramCountMin() = %d, want %d", dict.BigramCountMin(), tt.expectedBigramCountMin)
			}

			// Test initial counts
			if dict.WordCount() != 0 {
				t.Errorf("WordCount() = %d, want 0", dict.WordCount())
			}

			if dict.EntryCount() != 0 {
				t.Errorf("EntryCount() = %d, want 0", dict.EntryCount())
			}

			// Test that maps are initialized
			if dict.unigrams == nil {
				t.Error("unigrams map should be initialized")
			}

			if dict.bigramCounts == nil {
				t.Error("bigramCounts map should be initialized")
			}

			if dict.deletes == nil {
				t.Error("deletes map should be initialized")
			}

			if dict.belowThreshold == nil {
				t.Error("belowThreshold map should be initialized")
			}
		})
	}
}

func TestDictionary_LookupWord(t *testing.T) {
	dict := NewDictionary(16, 7, 2, 1)

	// Test empty dictionary
	result := dict.LookupWord("hello")
	if result != nil {
		t.Errorf("LookupWord on empty dictionary should return nil, got %v", result)
	}

	// Add some words
	dict.CreateDictionaryEntry("hello", 10)
	dict.CreateDictionaryEntry("world", 5)
	dict.CreateDictionaryEntry("test", 20)

	tests := []struct {
		name        string
		word        string
		expectFound bool
		expectCount int
	}{
		{
			name:        "existing word hello",
			word:        "hello",
			expectFound: true,
			expectCount: 10,
		},
		{
			name:        "existing word world",
			word:        "world",
			expectFound: true,
			expectCount: 5,
		},
		{
			name:        "existing word test",
			word:        "test",
			expectFound: true,
			expectCount: 20,
		},
		{
			name:        "non-existing word",
			word:        "nonexistent",
			expectFound: false,
			expectCount: 0,
		},
		{
			name:        "empty string",
			word:        "",
			expectFound: false,
			expectCount: 0,
		},
		{
			name:        "case sensitive",
			word:        "Hello",
			expectFound: false,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := dict.LookupWord(tt.word)

			if tt.expectFound {
				if result == nil {
					t.Errorf("LookupWord(%q) should return item, got nil", tt.word)
					return
				}

				if result.text != tt.word {
					t.Errorf("LookupWord(%q).text = %q, want %q", tt.word, result.text, tt.word)
				}

				if result.count != tt.expectCount {
					t.Errorf("LookupWord(%q).count = %d, want %d", tt.word, result.count, tt.expectCount)
				}
			} else {
				if result != nil {
					t.Errorf("LookupWord(%q) should return nil, got %v", tt.word, result)
				}
			}
		})
	}
}

func TestDictionary_LookupDelete(t *testing.T) {
	dict := NewDictionary(16, 7, 2, 1)

	// Test empty dictionary
	items, found := dict.LookupDelete("hello")
	if found {
		t.Error("LookupDelete on empty dictionary should return false")
	}
	if items != nil {
		t.Error("LookupDelete on empty dictionary should return nil items")
	}

	// Add words to generate deletes
	dict.CreateDictionaryEntry("hello", 10)
	dict.CreateDictionaryEntry("world", 5)

	tests := []struct {
		name        string
		deleteKey   string
		expectFound bool
		minItems    int
	}{
		{
			name:        "original word hello",
			deleteKey:   "hello",
			expectFound: true,
			minItems:    1,
		},
		{
			name:        "original word world",
			deleteKey:   "world",
			expectFound: true,
			minItems:    1,
		},
		{
			name:        "delete variant hell",
			deleteKey:   "hell",
			expectFound: true,
			minItems:    1,
		},
		{
			name:        "delete variant worl",
			deleteKey:   "worl",
			expectFound: true,
			minItems:    1,
		},
		{
			name:        "non-existent delete",
			deleteKey:   "xyz",
			expectFound: false,
			minItems:    0,
		},
		{
			name:        "empty string",
			deleteKey:   "",
			expectFound: false,
			minItems:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, found := dict.LookupDelete(tt.deleteKey)

			if tt.expectFound {
				if !found {
					t.Errorf("LookupDelete(%q) should return true, got false", tt.deleteKey)
				}
				if len(items) < tt.minItems {
					t.Errorf("LookupDelete(%q) should return at least %d items, got %d", tt.deleteKey, tt.minItems, len(items))
				}
				if items == nil {
					t.Errorf("LookupDelete(%q) should return items, got nil", tt.deleteKey)
				}
			} else {
				if found {
					t.Errorf("LookupDelete(%q) should return false, got true", tt.deleteKey)
				}
			}
		})
	}
}

func TestDictionary_CreateDictionaryStream(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		tokenizer     func(string) []string
		expectedWords map[string]int
		expectSuccess bool
	}{
		{
			name:    "simple text",
			content: "hello world\nfoo bar\n",
			tokenizer: func(line string) []string {
				return strings.Fields(strings.ToLower(line))
			},
			expectedWords: map[string]int{
				"hello": 1,
				"world": 1,
				"foo":   1,
				"bar":   1,
			},
			expectSuccess: true,
		},
		{
			name:    "empty content",
			content: "",
			tokenizer: func(line string) []string {
				return strings.Fields(strings.ToLower(line))
			},
			expectedWords: map[string]int{},
			expectSuccess: true,
		},
		{
			name:    "single line",
			content: "test word count",
			tokenizer: func(line string) []string {
				return strings.Fields(strings.ToLower(line))
			},
			expectedWords: map[string]int{
				"test":  1,
				"word":  1,
				"count": 1,
			},
			expectSuccess: true,
		},
		{
			name:    "duplicate words",
			content: "hello world\nhello test\nworld again",
			tokenizer: func(line string) []string {
				return strings.Fields(strings.ToLower(line))
			},
			expectedWords: map[string]int{
				"hello": 2,
				"world": 2,
				"test":  1,
				"again": 1,
			},
			expectSuccess: true,
		},
		{
			name:    "custom tokenizer with punctuation",
			content: "hello, world! test.",
			tokenizer: func(line string) []string {
				words := strings.FieldsFunc(strings.ToLower(line), func(r rune) bool {
					return r == ' ' || r == ',' || r == '!' || r == '.'
				})
				var result []string
				for _, word := range words {
					if word != "" {
						result = append(result, word)
					}
				}
				return result
			},
			expectedWords: map[string]int{
				"hello": 1,
				"world": 1,
				"test":  1,
			},
			expectSuccess: true,
		},
		{
			name:    "lines with only whitespace",
			content: "hello\n   \n\t\nworld",
			tokenizer: func(line string) []string {
				return strings.Fields(strings.ToLower(line))
			},
			expectedWords: map[string]int{
				"hello": 1,
				"world": 1,
			},
			expectSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dict := NewDictionary(16, 7, 2, 1)
			reader := strings.NewReader(tt.content)

			success, err := dict.CreateDictionaryStream(reader, tt.tokenizer)

			if tt.expectSuccess {
				if !success {
					t.Errorf("CreateDictionaryStream should succeed, got false")
				}
				if err != nil {
					t.Errorf("CreateDictionaryStream should not return error, got: %v", err)
				}
			} else {
				if success {
					t.Errorf("CreateDictionaryStream should fail, got true")
				}
			}

			// Verify words were added correctly
			for expectedWord, expectedCount := range tt.expectedWords {
				item := dict.LookupWord(expectedWord)
				if item == nil {
					t.Errorf("Word %q should be in dictionary", expectedWord)
					continue
				}
				if item.count != expectedCount {
					t.Errorf("Word %q count = %d, want %d", expectedWord, item.count, expectedCount)
				}
			}

			// Verify no unexpected words were added
			if dict.WordCount() != len(tt.expectedWords) {
				t.Errorf("Dictionary should have %d words, got %d", len(tt.expectedWords), dict.WordCount())
			}
		})
	}
}

func TestDictionary_CreateDictionary(t *testing.T) {
	// Create temporary file
	tempFile, err := os.CreateTemp("", "test_create_dict_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	content := "hello world\ntest data\nhello again"
	if _, err := tempFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	dict := NewDictionary(16, 7, 2, 1)
	tokenizer := func(line string) []string {
		return strings.Fields(strings.ToLower(line))
	}

	success, err := dict.CreateDictionary(tempFile.Name(), tokenizer)

	if !success {
		t.Errorf("CreateDictionary should succeed")
	}
	if err != nil {
		t.Errorf("CreateDictionary should not return error, got: %v", err)
	}

	expectedWords := map[string]int{
		"hello": 2,
		"world": 1,
		"test":  1,
		"data":  1,
		"again": 1,
	}

	for expectedWord, expectedCount := range expectedWords {
		item := dict.LookupWord(expectedWord)
		if item == nil {
			t.Errorf("Word %q should be in dictionary", expectedWord)
			continue
		}
		if item.count != expectedCount {
			t.Errorf("Word %q count = %d, want %d", expectedWord, item.count, expectedCount)
		}
	}
}

func TestDictionary_CreateDictionary_FileNotFound(t *testing.T) {
	dict := NewDictionary(16, 7, 2, 1)
	tokenizer := func(line string) []string {
		return strings.Fields(strings.ToLower(line))
	}

	success, err := dict.CreateDictionary("nonexistent_file.txt", tokenizer)

	if success {
		t.Errorf("CreateDictionary should fail for non-existent file")
	}
	if err == nil {
		t.Errorf("CreateDictionary should return error for non-existent file")
	}
}

func TestDictionary_LoadDictionaryStream(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		termIndex     int
		countIndex    int
		separator     string
		expectedWords map[string]int
		expectSuccess bool
	}{
		{
			name:       "tab separated",
			content:    "hello\t10\nworld\t5\ntest\t15\n",
			termIndex:  0,
			countIndex: 1,
			separator:  "\t",
			expectedWords: map[string]int{
				"hello": 10,
				"world": 5,
				"test":  15,
			},
			expectSuccess: true,
		},
		{
			name:       "comma separated",
			content:    "word,100,extra\ntest,50,data\n",
			termIndex:  0,
			countIndex: 1,
			separator:  ",",
			expectedWords: map[string]int{
				"word": 100,
				"test": 50,
			},
			expectSuccess: true,
		},
		{
			name:       "space separated",
			content:    "hello 25\nworld 30\n",
			termIndex:  0,
			countIndex: 1,
			separator:  " ",
			expectedWords: map[string]int{
				"hello": 25,
				"world": 30,
			},
			expectSuccess: true,
		},
		{
			name:          "empty content",
			content:       "",
			termIndex:     0,
			countIndex:    1,
			separator:     "\t",
			expectedWords: map[string]int{},
			expectSuccess: true,
		},
		{
			name:       "mixed valid and invalid lines",
			content:    "valid\t10\ninvalid_line\nhello\t5\nbad_count\tabc\nworld\t20\n",
			termIndex:  0,
			countIndex: 1,
			separator:  "\t",
			expectedWords: map[string]int{
				"valid": 10,
				"hello": 5,
				"world": 20,
			},
			expectSuccess: true,
		},
		{
			name:       "reversed indices",
			content:    "10\thello\n5\tworld\n",
			termIndex:  1,
			countIndex: 0,
			separator:  "\t",
			expectedWords: map[string]int{
				"hello": 10,
				"world": 5,
			},
			expectSuccess: true,
		},
		{
			name:       "pipe separated",
			content:    "word1|100|extra|data\nword2|50|more|info\n",
			termIndex:  0,
			countIndex: 1,
			separator:  "|",
			expectedWords: map[string]int{
				"word1": 100,
				"word2": 50,
			},
			expectSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dict := NewDictionary(16, 7, 2, 1)
			reader := strings.NewReader(tt.content)

			success, err := dict.LoadDictionaryStream(reader, tt.termIndex, tt.countIndex, tt.separator)

			if tt.expectSuccess {
				if !success {
					t.Errorf("LoadDictionaryStream should succeed, got false")
				}
				if err != nil {
					t.Errorf("LoadDictionaryStream should not return error, got: %v", err)
				}
			}

			for expectedWord, expectedCount := range tt.expectedWords {
				item := dict.LookupWord(expectedWord)
				if item == nil {
					t.Errorf("Word %q should be in dictionary", expectedWord)
					continue
				}
				if item.count != expectedCount {
					t.Errorf("Word %q count = %d, want %d", expectedWord, item.count, expectedCount)
				}
			}
		})
	}
}

func TestDictionary_LoadDictionary(t *testing.T) {
	// Create temporary file
	tempFile, err := os.CreateTemp("", "test_load_dict_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	content := "hello\t10\nworld\t5\ntest\t20\n"
	if _, err := tempFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	dict := NewDictionary(16, 7, 2, 1)
	success, err := dict.LoadDictionary(tempFile.Name(), 0, 1, "\t")

	if !success {
		t.Errorf("LoadDictionary should succeed")
	}
	if err != nil {
		t.Errorf("LoadDictionary should not return error, got: %v", err)
	}

	expectedWords := map[string]int{
		"hello": 10,
		"world": 5,
		"test":  20,
	}

	for expectedWord, expectedCount := range expectedWords {
		item := dict.LookupWord(expectedWord)
		if item == nil {
			t.Errorf("Word %q should be in dictionary", expectedWord)
			continue
		}
		if item.count != expectedCount {
			t.Errorf("Word %q count = %d, want %d", expectedWord, item.count, expectedCount)
		}
	}
}

func TestDictionary_LoadDictionary_FileNotFound(t *testing.T) {
	dict := NewDictionary(16, 7, 2, 1)
	success, err := dict.LoadDictionary("nonexistent_file.txt", 0, 1, "\t")

	if success {
		t.Errorf("LoadDictionary should fail for non-existent file")
	}
	if err == nil {
		t.Errorf("LoadDictionary should return error for non-existent file")
	}
}

func TestDictionary_LoadBigramDictionaryStream(t *testing.T) {
	tests := []struct {
		name            string
		content         string
		termIndex       int
		countIndex      int
		separator       string
		expectedBigrams map[string]int
		expectSuccess   bool
	}{
		{
			name:       "simple bigrams",
			content:    "hello world\t100\nfoo bar\t50\n",
			termIndex:  0,
			countIndex: 1,
			separator:  "\t",
			expectedBigrams: map[string]int{
				"hello world": 100,
				"foo bar":     50,
			},
			expectSuccess: true,
		},
		{
			name:       "comma separated",
			content:    "the cat,75,extra\nquick brown,25,data\n",
			termIndex:  0,
			countIndex: 1,
			separator:  ",",
			expectedBigrams: map[string]int{
				"the cat":     75,
				"quick brown": 25,
			},
			expectSuccess: true,
		},
		{
			name:            "empty content",
			content:         "",
			termIndex:       0,
			countIndex:      1,
			separator:       "\t",
			expectedBigrams: map[string]int{},
			expectSuccess:   true,
		},
		{
			name:       "invalid lines",
			content:    "valid bigram\t10\ninvalid_line\nhello world\t5\nbad_count\tabc\n",
			termIndex:  0,
			countIndex: 1,
			separator:  "\t",
			expectedBigrams: map[string]int{
				"valid bigram": 10,
				"hello world":  5,
			},
			expectSuccess: true,
		},
		{
			name:       "zero count",
			content:    "zero count\t0\npositive count\t15\n",
			termIndex:  0,
			countIndex: 1,
			separator:  "\t",
			expectedBigrams: map[string]int{
				"zero count":     0,
				"positive count": 15,
			},
			expectSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dict := NewDictionary(16, 7, 2, 1)
			reader := strings.NewReader(tt.content)

			success, err := dict.LoadBigramDictionaryStream(reader, tt.termIndex, tt.countIndex, tt.separator)

			if tt.expectSuccess {
				if !success {
					t.Errorf("LoadBigramDictionaryStream should succeed, got false")
				}
				if err != nil {
					t.Errorf("LoadBigramDictionaryStream should not return error, got: %v", err)
				}
			}

			for expectedBigram, expectedCount := range tt.expectedBigrams {
				count := dict.GetBigramCount(expectedBigram)
				if count != expectedCount {
					t.Errorf("Bigram %q count = %d, want %d", expectedBigram, count, expectedCount)
				}
			}
		})
	}
}

func TestDictionary_LoadBigramDictionary(t *testing.T) {
	// Create temporary file
	tempFile, err := os.CreateTemp("", "test_load_bigram_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	content := "hello world\t100\nthe cat\t75\nquick brown\t50\n"
	if _, err := tempFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	dict := NewDictionary(16, 7, 2, 1)
	success, err := dict.LoadBigramDictionary(tempFile.Name(), 0, 1, "\t")

	if !success {
		t.Errorf("LoadBigramDictionary should succeed")
	}
	if err != nil {
		t.Errorf("LoadBigramDictionary should not return error, got: %v", err)
	}

	expectedBigrams := map[string]int{
		"hello world": 100,
		"the cat":     75,
		"quick brown": 50,
	}

	for expectedBigram, expectedCount := range expectedBigrams {
		count := dict.GetBigramCount(expectedBigram)
		if count != expectedCount {
			t.Errorf("Bigram %q count = %d, want %d", expectedBigram, count, expectedCount)
		}
	}
}

func TestDictionary_LoadBigramDictionary_FileNotFound(t *testing.T) {
	dict := NewDictionary(16, 7, 2, 1)
	success, err := dict.LoadBigramDictionary("nonexistent_file.txt", 0, 1, "\t")

	if success {
		t.Errorf("LoadBigramDictionary should fail for non-existent file")
	}
	if err == nil {
		t.Errorf("LoadBigramDictionary should return error for non-existent file")
	}
}

func TestDictionary_CreateDictionaryEntry(t *testing.T) {
	tests := []struct {
		name            string
		word            string
		addCount        int
		countThreshold  int
		expectedCount   int
		expectInDict    bool
		expectInDeletes bool
	}{
		{
			name:            "simple word above threshold",
			word:            "hello",
			addCount:        5,
			countThreshold:  1,
			expectedCount:   5,
			expectInDict:    true,
			expectInDeletes: true,
		},
		{
			name:            "word exceeding max length",
			word:            "hello-its-foo-from-test",
			addCount:        10,
			countThreshold:  1,
			expectedCount:   10,
			expectInDict:    true,
			expectInDeletes: false,
		},
		{
			name:            "zero count with zero threshold",
			word:            "test",
			addCount:        0,
			countThreshold:  0,
			expectedCount:   0,
			expectInDict:    true,
			expectInDeletes: true,
		},
		{
			name:            "negative count ignored with positive threshold",
			word:            "negative",
			addCount:        -5,
			countThreshold:  1,
			expectedCount:   0,
			expectInDict:    false,
			expectInDeletes: false,
		},
		{
			name:            "single character word",
			word:            "a",
			addCount:        10,
			countThreshold:  1,
			expectedCount:   10,
			expectInDict:    true,
			expectInDeletes: true,
		},
		{
			name:            "unicode word",
			word:            "café",
			addCount:        3,
			countThreshold:  1,
			expectedCount:   3,
			expectInDict:    true,
			expectInDeletes: true,
		},
		{
			name:            "emoji",
			word:            "🙂",
			addCount:        2,
			countThreshold:  1,
			expectedCount:   2,
			expectInDict:    true,
			expectInDeletes: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dict := NewDictionary(16, 7, 2, tt.countThreshold)

			// Record initial max length
			initialMaxLen := dict.MaxLength()

			dict.CreateDictionaryEntry(tt.word, tt.addCount)

			if tt.expectInDict {
				item := dict.LookupWord(tt.word)
				if item == nil {
					t.Errorf("Word %q should be in dictionary", tt.word)
				} else {
					if item.count != tt.expectedCount {
						t.Errorf("Word %q count = %d, want %d", tt.word, item.count, tt.expectedCount)
					}
					if item.text != tt.word {
						t.Errorf("Word text = %q, want %q", item.text, tt.word)
					}
				}

				// Check if max length was updated correctly
				expectedMaxLen := len([]rune(tt.word))
				if expectedMaxLen > initialMaxLen && dict.MaxLength() < expectedMaxLen {
					t.Errorf("MaxLength should be at least %d, got %d", expectedMaxLen, dict.MaxLength())
				}
			} else {
				item := dict.LookupWord(tt.word)
				if item != nil {
					t.Errorf("Word %q should not be in dictionary when below threshold", tt.word)
				}
			}

			if tt.expectInDeletes {
				_, found := dict.LookupDelete(tt.word)
				if !found {
					t.Errorf("Word %q should have delete variants indexed", tt.word)
				}
			}
		})
	}
}

func TestDictionary_CreateDictionaryEntry_Accumulation(t *testing.T) {
	dict := NewDictionary(16, 7, 2, 1)

	// Add word multiple times
	dict.CreateDictionaryEntry("hello", 5)
	dict.CreateDictionaryEntry("hello", 3)
	dict.CreateDictionaryEntry("hello", 2)

	item := dict.LookupWord("hello")
	if item == nil {
		t.Fatal("Word should be in dictionary")
	}

	expectedCount := 10
	if item.count != expectedCount {
		t.Errorf("Word count = %d, want %d", item.count, expectedCount)
	}
}

func TestDictionary_CreateDictionaryEntry_ThresholdPromotion(t *testing.T) {
	dict := NewDictionary(16, 7, 2, 5)

	// Add word below threshold
	dict.CreateDictionaryEntry("test", 2)

	item := dict.LookupWord("test")
	if item != nil {
		t.Fatal("Word should not be in unigrams map")
	}

	// Should not have delete variants yet
	_, found := dict.LookupDelete("test")
	if found {
		t.Errorf("Word should not have delete variants when below threshold")
	}

	// Add more to reach threshold
	dict.CreateDictionaryEntry("test", 4)

	// Should now have delete variants
	_, found = dict.LookupDelete("test")
	if !found {
		t.Errorf("Word should have delete variants after reaching threshold")
	}

	item = dict.LookupWord("test")
	if item == nil {
		t.Fatal("Word should be in unigrams map")
	}

	if item.count != 6 {
		t.Errorf("Word count = %d, want 6", item.count)
	}
}

func TestDictionary_SetBigramCount(t *testing.T) {
	tests := []struct {
		name            string
		term            string
		count           int
		expectedCount   int
		shouldStore     bool
		expectMinUpdate bool
	}{
		{
			name:            "valid bigram",
			term:            "hello world",
			count:           100,
			expectedCount:   100,
			shouldStore:     true,
			expectMinUpdate: true,
		},
		{
			name:            "single word ignored",
			term:            "hello",
			count:           50,
			expectedCount:   0,
			shouldStore:     false,
			expectMinUpdate: false,
		},
		{
			name:            "empty string ignored",
			term:            "",
			count:           25,
			expectedCount:   0,
			shouldStore:     false,
			expectMinUpdate: false,
		},
		{
			name:            "three words uses first two",
			term:            "the quick brown",
			count:           75,
			expectedCount:   75,
			shouldStore:     true,
			expectMinUpdate: true,
		},
		{
			name:            "zero count",
			term:            "zero count",
			count:           0,
			expectedCount:   0,
			shouldStore:     true,
			expectMinUpdate: true,
		},
		{
			name:            "negative count",
			term:            "negative count",
			count:           -10,
			expectedCount:   -10,
			shouldStore:     true,
			expectMinUpdate: true,
		},
		{
			name:            "extra whitespace",
			term:            "extra  spaces",
			count:           30,
			expectedCount:   30,
			shouldStore:     true,
			expectMinUpdate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dict := NewDictionary(16, 7, 2, 1)
			initialMin := dict.BigramCountMin()

			dict.SetBigramCount(tt.term, tt.count)

			count := dict.GetBigramCount(tt.term)
			if count != tt.expectedCount {
				t.Errorf("GetBigramCount(%q) = %d, want %d", tt.term, count, tt.expectedCount)
			}

			// Check if minimum was updated
			if tt.expectMinUpdate && tt.count < initialMin {
				if dict.BigramCountMin() != tt.count {
					t.Errorf("BigramCountMin should be %d, got %d", tt.count, dict.BigramCountMin())
				}
			}
		})
	}
}

func TestDictionary_SetBigramCount_NoDuplicate(t *testing.T) {
	dict := NewDictionary(16, 7, 2, 1)

	// Add same bigram twice - should not change count
	dict.SetBigramCount("hello world", 100)
	dict.SetBigramCount("hello world", 200)

	count := dict.GetBigramCount("hello world")
	if count != 100 {
		t.Errorf("Duplicate bigram should not change count, got %d, want 100", count)
	}
}

func TestDictionary_GetBigramCount(t *testing.T) {
	dict := NewDictionary(16, 7, 2, 1)

	// Test non-existent bigram
	count := dict.GetBigramCount("nonexistent bigram")
	if count != 0 {
		t.Errorf("Non-existent bigram should return 0, got %d", count)
	}

	// Add bigram and test
	dict.SetBigramCount("hello world", 50)
	count = dict.GetBigramCount("hello world")
	if count != 50 {
		t.Errorf("Bigram count should be 50, got %d", count)
	}

	// Test case sensitivity
	count = dict.GetBigramCount("Hello World")
	if count != 0 {
		t.Errorf("Case sensitive lookup should return 0, got %d", count)
	}
}

func TestDictionary_getOrCreateUnigram(t *testing.T) {
	dict := NewDictionary(16, 7, 2, 1)

	// Test creating new unigram
	item, isNew := dict.getOrCreateUnigram("hello", 5)
	if !isNew {
		t.Errorf("First call should create new unigram")
	}
	if item == nil {
		t.Fatal("Item should not be nil")
	}
	if item.text != "hello" {
		t.Errorf("Item text = %q, want 'hello'", item.text)
	}
	if item.count != 5 {
		t.Errorf("Item count = %d, want 5", item.count)
	}

	// Test getting existing unigram
	existingItem, isNew := dict.getOrCreateUnigram("hello", 10)
	if isNew {
		t.Errorf("Second call should return existing unigram")
	}
	if existingItem != item {
		t.Errorf("Should return same item instance")
	}
	if existingItem.count != 5 {
		t.Errorf("Count should remain unchanged at 5, got %d", existingItem.count)
	}
}

func TestDictionary_indexDeletesForID(t *testing.T) {
	dict := NewDictionary(16, 5, 1, 1) // prefix length 5, max edit distance 1

	// Test with short word (within prefix length and edit distance)
	item := NewDictItem("cat", 1)
	dict.indexDeletesForID(item)

	// Original word should be indexed
	items, found := dict.LookupDelete("cat")
	if !found {
		t.Error("Original word should be indexed")
	}
	if len(items) == 0 {
		t.Error("Should have at least one item")
	}

	// Delete variants should be indexed
	variants := []string{"ca", "ct", "at"}
	for _, variant := range variants {
		items, found := dict.LookupDelete(variant)
		if !found {
			t.Errorf("Delete variant %q should be indexed", variant)
		}
		if len(items) == 0 {
			t.Errorf("Delete variant %q should have items", variant)
		}
	}
}

func TestDictionary_indexDeletesForID_LongWord(t *testing.T) {
	dict := NewDictionary(16, 3, 1, 1) // prefix length 3

	// Test with long word (exceeds prefix length)
	item := NewDictItem("hello", 1)
	dict.indexDeletesForID(item)

	// Truncated prefix should be indexed
	items, found := dict.LookupDelete("hel")
	if !found {
		t.Error("Truncated prefix should be indexed")
	}
	if len(items) == 0 {
		t.Error("Should have at least one item")
	}
}

func TestDictionary_walkDeletes(t *testing.T) {
	dict := NewDictionary(16, 7, 2, 1) // max edit distance 2

	input := []rune("cat")
	var generated []string

	dict.walkDeletes(input, func(deleteKey string) {
		generated = append(generated, deleteKey)
	})

	// Should generate various delete combinations
	if len(generated) == 0 {
		t.Error("walkDeletes should generate delete variants")
	}

	// Check for expected single deletions
	expectedSingle := []string{"at", "ct", "ca"}
	for _, expected := range expectedSingle {
		found := false
		for _, gen := range generated {
			if gen == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected single deletion %q not found in generated variants", expected)
		}
	}

	// Should also generate double deletions (edit distance 2)
	expectedDouble := []string{"a", "c", "t"}
	for _, expected := range expectedDouble {
		found := false
		for _, gen := range generated {
			if gen == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected double deletion %q not found in generated variants", expected)
		}
	}
}

func TestDictionary_walkDeletes_EmptyInput(t *testing.T) {
	dict := NewDictionary(16, 7, 2, 1)

	input := []rune("")
	var generated []string

	dict.walkDeletes(input, func(deleteKey string) {
		generated = append(generated, deleteKey)
	})

	if len(generated) != 0 {
		t.Errorf("Empty input should not generate any deletes, got %d", len(generated))
	}
}

func TestDictionary_walkDeletes_SingleChar(t *testing.T) {
	dict := NewDictionary(16, 7, 2, 1)

	input := []rune("a")
	var generated []string

	dict.walkDeletes(input, func(deleteKey string) {
		generated = append(generated, deleteKey)
	})

	if len(generated) != 1 {
		t.Errorf("Single char should generate one delete (empty string), got %d", len(generated))
	}
	if len(generated) > 0 && generated[0] != "" {
		t.Errorf("Single char delete should be empty string, got %q", generated[0])
	}
}

func TestAppendIfMissing(t *testing.T) {
	item1 := NewDictItem("hello", 1)
	item2 := NewDictItem("world", 2)
	item3 := NewDictItem("hello", 3) // Same text as item1

	tests := []struct {
		name     string
		slice    []*DictItem
		item     *DictItem
		expected int
	}{
		{
			name:     "add to empty slice",
			slice:    []*DictItem{},
			item:     item1,
			expected: 1,
		},
		{
			name:     "add new item",
			slice:    []*DictItem{item1},
			item:     item2,
			expected: 2,
		},
		{
			name:     "add duplicate item (same text)",
			slice:    []*DictItem{item1},
			item:     item3,
			expected: 1,
		},
		{
			name:     "add to existing items",
			slice:    []*DictItem{item1, item2},
			item:     item3,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := appendIfMissing(tt.slice, tt.item)

			if len(result) != tt.expected {
				t.Errorf("Length = %d, want %d", len(result), tt.expected)
			}

			// Check if the item with same text already exists
			if tt.expected == len(tt.slice) {
				// Should not have added the item
				found := false
				for _, existing := range tt.slice {
					if existing.text == tt.item.text {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Item should have been found in original slice")
				}
			} else {
				// Should have added the item
				found := false
				for _, item := range result {
					if item == tt.item {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("New item should be in result slice")
				}
			}
		})
	}
}

func TestDictionary_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func() *Dictionary
		testFunc    func(*Dictionary) error
		description string
	}{
		{
			name: "very long word",
			setupFunc: func() *Dictionary {
				return NewDictionary(16, 7, 2, 1)
			},
			testFunc: func(d *Dictionary) error {
				longWord := strings.Repeat("abcdefghij", 100) // 1000 characters
				d.CreateDictionaryEntry(longWord, 1)

				item := d.LookupWord(longWord)
				if item == nil {
					return fmt.Errorf("Long word should be stored")
				}
				if item.Len() != 1000 {
					return fmt.Errorf("Long word length = %d, want 1000", item.Len())
				}
				return nil
			},
			description: "Handle very long words correctly",
		},
		{
			name: "unicode characters",
			setupFunc: func() *Dictionary {
				return NewDictionary(16, 7, 2, 1)
			},
			testFunc: func(d *Dictionary) error {
				unicodeWords := []string{"café", "naïve", "résumé", "测试", "🙂😀👍"}
				for _, word := range unicodeWords {
					d.CreateDictionaryEntry(word, 1)
				}

				for _, word := range unicodeWords {
					item := d.LookupWord(word)
					if item == nil {
						return fmt.Errorf("Unicode word %q should be stored", word)
					}
				}
				return nil
			},
			description: "Handle unicode characters correctly",
		},
		{
			name: "empty strings",
			setupFunc: func() *Dictionary {
				return NewDictionary(16, 7, 2, 1)
			},
			testFunc: func(d *Dictionary) error {
				d.CreateDictionaryEntry("", 5)

				item := d.LookupWord("")
				if item == nil {
					return fmt.Errorf("Empty string should be stored")
				}
				if item.Len() != 0 {
					return fmt.Errorf("Empty string length should be 0, got %d", item.Len())
				}
				return nil
			},
			description: "Handle empty strings",
		},
		{
			name: "zero parameters",
			setupFunc: func() *Dictionary {
				return NewDictionary(0, 0, 0, 0)
			},
			testFunc: func(d *Dictionary) error {
				if d.PrefixLength() != 0 {
					return fmt.Errorf("PrefixLength should be 0")
				}
				if d.MaxDictionaryEditDistance() != 0 {
					return fmt.Errorf("MaxDictionaryEditDistance should be 0")
				}
				if d.CountThreshold() != 0 {
					return fmt.Errorf("CountThreshold should be 0")
				}
				return nil
			},
			description: "Handle zero parameters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dict := tt.setupFunc()
			if err := tt.testFunc(dict); err != nil {
				t.Errorf("%s: %v", tt.description, err)
			}
		})
	}
}

func TestDictionary_BigramCountMin_Updates(t *testing.T) {
	dict := NewDictionary(16, 7, 2, 1)

	// Initial state
	if dict.BigramCountMin() != math.MaxInt {
		t.Errorf("Initial BigramCountMin should be MaxInt, got %d", dict.BigramCountMin())
	}

	// Add bigrams with decreasing counts
	bigrams := []struct {
		term  string
		count int
	}{
		{"hello world", 100},
		{"foo bar", 50},
		{"test case", 25},
		{"min value", 10},
	}

	for _, bigram := range bigrams {
		dict.SetBigramCount(bigram.term, bigram.count)

		// BigramCountMin should be updated to the smallest value seen so far
		expectedMin := bigram.count
		for _, prev := range bigrams {
			if prev.term == bigram.term {
				break
			}
			if prev.count < expectedMin {
				expectedMin = prev.count
			}
		}

		if dict.BigramCountMin() != expectedMin {
			t.Errorf("After adding %q, BigramCountMin = %d, want %d",
				bigram.term, dict.BigramCountMin(), expectedMin)
		}
	}
}

func TestDictionary_MaxLength_Updates(t *testing.T) {
	dict := NewDictionary(16, 7, 2, 1)

	// Test with progressively longer words
	words := []struct {
		text       string
		runeLength int
	}{
		{"a", 1},
		{"hi", 2},
		{"hello", 5},
		{"testing", 7},
		{"supercalifragilisticexpialidocious", 34},
		{"café", 4}, // Should not update max length
	}

	for _, word := range words {
		initialMax := dict.MaxLength()
		dict.CreateDictionaryEntry(word.text, 1)

		expectedMax := initialMax
		if word.runeLength > initialMax {
			expectedMax = word.runeLength
		}

		if dict.MaxLength() != expectedMax {
			t.Errorf("After adding %q (len=%d), MaxLength = %d, want %d",
				word.text, word.runeLength, dict.MaxLength(), expectedMax)
		}
	}
}

func TestDictionary_Integration_LoadAndQuery(t *testing.T) {
	dict := NewDictionary(16, 7, 2, 1)

	// Test integration of loading and querying
	content := "hello\t10\nworld\t5\ntest\t20\ndata\t15\n"
	reader := strings.NewReader(content)

	success, err := dict.LoadDictionaryStream(reader, 0, 1, "\t")
	if !success || err != nil {
		t.Fatalf("LoadDictionaryStream failed: success=%v, err=%v", success, err)
	}

	// Verify all words are loaded
	expectedWords := map[string]int{
		"hello": 10,
		"world": 5,
		"test":  20,
		"data":  15,
	}

	for word, expectedCount := range expectedWords {
		item := dict.LookupWord(word)
		if item == nil {
			t.Errorf("Word %q should be loaded", word)
			continue
		}
		if item.count != expectedCount {
			t.Errorf("Word %q count = %d, want %d", word, item.count, expectedCount)
		}

		// Verify delete variants are indexed
		_, found := dict.LookupDelete(word)
		if !found {
			t.Errorf("Word %q should have delete variants", word)
		}
	}

	// Test bigram loading
	bigramContent := "hello world\t100\ntest data\t50\n"
	bigramReader := strings.NewReader(bigramContent)

	success, err = dict.LoadBigramDictionaryStream(bigramReader, 0, 1, "\t")
	if !success || err != nil {
		t.Fatalf("LoadBigramDictionaryStream failed: success=%v, err=%v", success, err)
	}

	if dict.GetBigramCount("hello world") != 100 {
		t.Errorf("Bigram 'hello world' count should be 100")
	}
	if dict.GetBigramCount("test data") != 50 {
		t.Errorf("Bigram 'test data' count should be 50")
	}
}
