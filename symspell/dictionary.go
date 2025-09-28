package symspell

import (
	"bufio"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

type DictItem struct {
	text  string
	runes []rune
	count int
}

func (di *DictItem) Len() int {
	return len(di.runes)
}

func NewDictItem(word string, count int) *DictItem {
	return &DictItem{
		text:  word,
		runes: []rune(word),
		count: count,
	}
}

type Dictionary struct {
	unigrams       map[string]*DictItem
	bigramCounts   map[string]int
	deletes        map[string][]*DictItem
	belowThreshold map[string]int

	// maxLength tracks the longest word among all valid dictionary entries.
	// Used for optimization - words longer than maxLength + maxEditDistance
	// cannot be corrections for any dictionary word.
	maxLength int

	// bigramCountMin stores the minimum frequency among all bigrams.
	// Used for probability calculations and filtering low-frequency bigrams.
	bigramCountMin int

	// prefixLength limits the length of prefixes used for delete key generation.
	// Longer words are truncated to this length before generating delete variants.
	// This balances memory usage with correction quality.
	prefixLength int

	// maxDictionaryEditDistance defines the maximum edit distance for precomputed deletes.
	// All possible character deletions up to this distance are indexed.
	// Higher values provide better recall but use more memory.
	maxDictionaryEditDistance int

	// countThreshold sets the minimum frequency required for a word to be considered valid.
	// Words with frequency below this threshold are stored in belowThreshold map
	// until they accumulate sufficient frequency to be promoted.
	countThreshold int
}

// EntryCount returns the total number of delete variant entries in the dictionary.
//
// This represents the number of unique delete keys indexed for spell correction,
// not the number of unique words. Each word may generate multiple delete variants
// based on the maxDictionaryEditDistance setting.
//
// Returns:
//
//	The number of delete variant entries stored in the deletes map.
func (d *Dictionary) EntryCount() int {
	return len(d.deletes)
}

// MaxLength returns the maximum rune length among all valid dictionary words.
//
// This value is automatically maintained as words are added and removed.
// It's used for optimization during spell checking - candidate words that
// are too long relative to this value can be quickly filtered out.
//
// Returns:
//
//	The length in runes of the longest valid word in the dictionary.
func (d *Dictionary) MaxLength() int { return d.maxLength }

// PrefixLength returns the configured prefix length for delete key generation.
//
// This is the maximum length of word prefixes used when generating delete
// variants. Words longer than this are truncated before processing.
//
// Returns:
//
//	The prefix length setting used for delete key generation.
func (d *Dictionary) PrefixLength() int { return d.prefixLength }

// MaxDictionaryEditDistance returns the maximum edit distance for precomputed deletes.
//
// This is the maximum number of character deletions that are precomputed
// and indexed for each dictionary word. Higher values provide better recall
// but require more memory.
//
// Returns:
//
//	The maximum edit distance setting for the dictionary.
func (d *Dictionary) MaxDictionaryEditDistance() int { return d.maxDictionaryEditDistance }

// CountThreshold returns the minimum frequency required for word validity.
//
// Words with corpus frequency below this threshold are not considered valid
// for spell checking until their accumulated frequency reaches this value.
//
// Returns:
//
//	The count threshold setting for word validity.
func (d *Dictionary) CountThreshold() int { return d.countThreshold }

// WordCount returns the total number of unique words in the dictionary.
//
// This includes all words that have been added to the dictionary, regardless
// of whether they meet the count threshold for being considered valid.
// Use LookupWord to check if a specific word is considered valid.
//
// Returns:
//
//	The total number of unique words stored in the word pool.
func (d *Dictionary) WordCount() int { return len(d.unigrams) }

func NewDictionary(initialCapacity, prefixLength, maxDictEdit, countThreshold int) *Dictionary {
	return &Dictionary{
		unigrams:       make(map[string]*DictItem, initialCapacity),
		bigramCounts:   make(map[string]int, initialCapacity),
		deletes:        make(map[string][]*DictItem, initialCapacity),
		belowThreshold: make(map[string]int, initialCapacity/10+1),

		maxLength:                 0,
		bigramCountMin:            math.MaxInt,
		prefixLength:              prefixLength,
		maxDictionaryEditDistance: maxDictEdit,
		countThreshold:            countThreshold,
	}
}

func (d *Dictionary) LookupWord(word string) *DictItem {
	return d.unigrams[word]
}

func (d *Dictionary) LookupDelete(text string) ([]*DictItem, bool) {
	items, ok := d.deletes[text]
	return items, ok
}

// -------- Building / Loading --------

// CreateDictionaryStream loads multiple dictionary words from a stream containing plain text.
//
// This method reads text from an io.Reader line by line and extracts words using
// the provided tokenization function. Each extracted word is added to the dictionary
// with a frequency count of 1. This method is useful for building dictionaries
// from large text corpora.
//
// **NOTE**: Merges with any dictionary data already loaded.
//
// Args:
//
//	r: The io.Reader containing the text corpus to process.
//	tokenize: A function that takes a line of text and returns a slice of words.
//	          This allows for custom tokenization logic (e.g., handling punctuation,
//	          case normalization, filtering stop words).
//
// Returns:
//
//	True if the stream was processed successfully, false if an I/O error occurred.
//	Any scanning error encountered while reading the stream.
//
// Example:
//
//	tokenizer := func(line string) []string { return strings.Fields(strings.ToLower(line)) }
//	success, err := dict.CreateDictionaryStream(file, tokenizer)
func (d *Dictionary) CreateDictionaryStream(r io.Reader, tokenize func(string) []string) (bool, error) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		for _, key := range tokenize(sc.Text()) {
			d.CreateDictionaryEntry(key, 1)
		}
	}
	return true, sc.Err()
}

// CreateDictionary loads multiple dictionary words from a file containing plain text.
//
// This method opens a text file and processes it using CreateDictionaryStream.
// It's a convenience wrapper for loading dictionary data from files on disk.
//
// **NOTE**: Merges with any dictionary data already loaded.
//
// Args:
//
//	path: The path and filename of the text corpus file to load.
//	tokenize: A function that takes a line of text and returns a slice of words.
//	          This allows for custom tokenization logic.
//
// Returns:
//
//	True if the file was loaded successfully, false if the file couldn't be opened
//	or an I/O error occurred. Any file opening or scanning error encountered.
//
// Example:
//
//	tokenizer := func(line string) []string { return strings.Fields(strings.ToLower(line)) }
//	success, err := dict.CreateDictionary("corpus.txt", tokenizer)
func (d *Dictionary) CreateDictionary(path string, tokenize func(string) []string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()
	return d.CreateDictionaryStream(f, tokenize)
}

// LoadDictionary loads word-frequency pairs from a structured file.
//
// This method loads dictionary data from a file where each line contains
// structured data with word and frequency information separated by a delimiter.
// This is commonly used for loading pre-computed frequency dictionaries.
//
// **NOTE**: Merges with any dictionary data already loaded.
//
// Args:
//
//	path: The path and filename of the structured dictionary file.
//	termIndex: The zero-based column index containing the word terms.
//	countIndex: The zero-based column index containing the frequency counts.
//	sep: The field separator string (e.g., " ", "\t", ",", "|").
//
// Returns:
//
//	True if the file was loaded successfully, false if the file couldn't be opened
//	or an I/O error occurred. Any file opening or scanning error encountered.
//
// Example:
//
//	// Load from tab-separated file where word is column 0, count is column 1
//	success, err := dict.LoadDictionary("frequency.txt", 0, 1, "\t")
func (d *Dictionary) LoadDictionary(path string, termIndex, countIndex int, sep string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()
	return d.LoadDictionaryStream(f, termIndex, countIndex, sep)
}

// BigramCountMin returns the minimum frequency among all stored bigrams.
//
// This value is automatically maintained and used for probability calculations
// and filtering during context-aware spell checking.
//
// Returns:
//
//	The minimum bigram frequency, or math.MaxInt if no bigrams are stored.
func (d *Dictionary) BigramCountMin() int { return d.bigramCountMin }

// LoadDictionaryStream loads word-frequency pairs from a structured stream.
//
// This method processes structured data from an io.Reader where each line
// contains delimited fields with word and frequency information. Lines that
// don't have enough fields or contain invalid frequency values are skipped.
//
// **NOTE**: Merges with any dictionary data already loaded.
//
// Args:
//
//	r: The io.Reader containing the structured dictionary data.
//	termIndex: The zero-based column index containing the word terms.
//	countIndex: The zero-based column index containing the frequency counts.
//	sep: The field separator string used to split each line.
//
// Returns:
//
//	True if the stream was processed successfully, false if an I/O error occurred.
//	Any scanning error encountered while reading the stream.
//
// Example:
//
//	// Process comma-separated data: "word,123,other_data"
//	success, err := dict.LoadDictionaryStream(reader, 0, 1, ",")
func (d *Dictionary) LoadDictionaryStream(r io.Reader, termIndex, countIndex int, sep string) (bool, error) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		parts := strings.Split(sc.Text(), sep)
		if len(parts) <= termIndex || len(parts) <= countIndex {
			continue
		}
		key := parts[termIndex]
		c, err := strconv.Atoi(parts[countIndex])
		if err != nil {
			continue
		}
		d.CreateDictionaryEntry(key, c)
	}
	return true, sc.Err()
}

// LoadBigramDictionary loads bigram frequency data from a structured file.
//
// This method loads bigram (two-word sequence) frequency data from a file
// where each line contains bigram terms and their corpus frequencies.
// Bigrams are used for context-aware spell checking and improved correction ranking.
//
// **NOTE**: Merges with any bigram data already loaded.
//
// Args:
//
//	path: The path and filename of the bigram frequency file.
//	termIndex: The zero-based column index containing the bigram terms (e.g., "word1 word2").
//	countIndex: The zero-based column index containing the frequency counts.
//	sep: The field separator string used to split each line.
//
// Returns:
//
//	True if the file was loaded successfully, false if the file couldn't be opened
//	or an I/O error occurred. Any file opening or scanning error encountered.
//
// Example:
//
//	// Load from tab-separated file: "the cat\t45\tother_data"
//	success, err := dict.LoadBigramDictionary("bigrams.txt", 0, 1, "\t")
func (d *Dictionary) LoadBigramDictionary(path string, termIndex, countIndex int, sep string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()
	return d.LoadBigramDictionaryStream(f, termIndex, countIndex, sep)
}

// LoadBigramDictionaryStream loads bigram frequency data from a structured stream.
//
// This method processes bigram data from an io.Reader where each line contains
// delimited fields with bigram terms and frequency information. Lines that
// don't have enough fields or contain invalid frequency values are skipped.
//
// **NOTE**: Merges with any bigram data already loaded.
//
// Args:
//
//	r: The io.Reader containing the structured bigram data.
//	termIndex: The zero-based column index containing the bigram terms.
//	countIndex: The zero-based column index containing the frequency counts.
//	sep: The field separator string used to split each line.
//
// Returns:
//
//	True if the stream was processed successfully, false if an I/O error occurred.
//	Any scanning error encountered while reading the stream.
func (d *Dictionary) LoadBigramDictionaryStream(r io.Reader, termIndex, countIndex int, sep string) (bool, error) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		parts := strings.Split(sc.Text(), sep)
		if len(parts) <= termIndex || len(parts) <= countIndex {
			continue
		}
		key := parts[termIndex]
		count, err := strconv.Atoi(parts[countIndex])
		if err != nil {
			continue
		}
		d.SetBigramCount(key, count)
	}
	return true, sc.Err()
}

func (d *Dictionary) GetBigramCount(text string) int {
	count, ok := d.bigramCounts[text]
	if ok {
		return count
	}
	return 0
}

func (d *Dictionary) getOrCreateUnigram(text string, count int) (*DictItem, bool) {
	existing, ok := d.unigrams[text]
	if ok {
		return existing, false
	}
	newItem := NewDictItem(text, count)
	d.unigrams[text] = newItem
	return newItem, true
}

// SetBigramCount adds or updates a bigram frequency entry in the dictionary.
//
// This method parses a bigram string (typically "word1 word2") and stores
// the frequency count in the packed bigram map. It automatically creates
// WordIDs for the constituent words if they don't exist. The method also
// maintains the minimum bigram count for probability calculations.
//
// Args:
//
//	term: A string containing the bigram, typically space-separated words (e.g., "the cat").
//	count: The corpus frequency of this bigram.
//
// Note:
//
//	If the term contains fewer than 2 words, the operation is silently ignored.
//	Only the first two words are used if the term contains more than 2 words.
//
// Example:
//
//	dict.SetBigramCount("the cat", 1250)  // Adds bigram with frequency 1250
func (d *Dictionary) SetBigramCount(term string, count int) {
	ws := strings.Fields(term)
	if len(ws) < 2 {
		return
	}
	_, ok := d.bigramCounts[term]
	if !ok {
		d.bigramCounts[term] = count
	}
	if count < d.bigramCountMin {
		d.bigramCountMin = count
	}
}

// CreateDictionaryEntry adds or updates a word entry with frequency accumulation and threshold promotion.
//
// This is the core method for adding words to the dictionary. It handles frequency
// accumulation, count threshold checking, and automatic promotion to valid status.
// When a word reaches the count threshold, it becomes available for spell checking
// and its delete variants are indexed for fast lookup.
//
// The method implements several behaviors:
// - Accumulates frequency counts for existing words
// - Stores below-threshold words separately until they qualify for promotion
// - Automatically indexes delete variants when words become valid
// - Updates the maximum word length cache
//
// Args:
//
//	word: The word string to add or update in the dictionary.
//	addCount: The frequency count to add. Must be > 0 when countThreshold > 0.
//
// Note:
//
//	If addCount <= 0 and countThreshold > 0, the operation is silently ignored.
//	Words are only indexed for spell checking once they reach the count threshold.
//
// Example:
//
//	dict.CreateDictionaryEntry("example", 5)  // Add "example" with frequency 5
//	dict.CreateDictionaryEntry("example", 3)  // Accumulate to frequency 8
func (d *Dictionary) CreateDictionaryEntry(word string, addCount int) {
	// Early return if count is zero, as it can't change anything
	if addCount <= 0 && d.countThreshold > 0 {
		return
	}

	// Look first in below threshold words, update count, and allow promotion
	// to correct spelling word if count reaches threshold threshold must be
	// >1 for there to be the possibility of low threshold words
	if d.countThreshold > 1 {
		if prev, ok := d.belowThreshold[word]; ok {
			addCount += prev
			delete(d.belowThreshold, word)
		}
	}

	// Still below threshold
	if addCount < d.countThreshold {
		d.belowThreshold[word] = addCount
		return
	}

	item, isNew := d.getOrCreateUnigram(word, addCount)

	if item.Len() > d.maxLength {
		d.maxLength = item.Len()
	}

	// Already Valid and Existing
	if item.count > 0 && !isNew {
		item.count += addCount
		return
	}

	d.indexDeletesForID(item)
}

// indexDeletesForID generates and indexes all delete variants for a specific word.
//
// This method implements the core of the symmetric delete algorithm by generating
// all possible character deletion variants of a word up to the maximum edit distance
// and indexing them for fast spell correction lookup. Each delete variant hash
// points to this word as a potential correction.
//
// The indexing process includes:
// - The original word (if short enough)
// - The word truncated to prefix length
// - All possible single and multi-character deletions up to maxDictionaryEditDistance
//
// Args:
//
//	id: The WordID of the word to index delete variants for.
//
// Note:
//
//	This method is automatically called when a word is promoted to valid status.
//	It's computationally expensive for long words and high edit distances,
//	but enables very fast spell checking at query time.
func (d *Dictionary) indexDeletesForID(item *DictItem) {
	rs := item.runes
	if len(rs) <= d.maxDictionaryEditDistance {
		d.deletes[item.text] = appendIfMissing(d.deletes[item.text], item)
	}
	if len(rs) > d.prefixLength {
		rs = rs[:d.prefixLength]
	}
	// include the (possibly truncated) prefix itself
	{
		str := string(rs)
		d.deletes[str] = appendIfMissing(d.deletes[str], item)
	}
	// deeper deletes
	d.walkDeletes(rs, func(key string) {
		d.deletes[key] = appendIfMissing(d.deletes[key], item)
	})
}

// walkDeletes systematically generates all possible character deletion variants up to maxDictionaryEditDistance.
//
// This method implements a depth-first traversal using a stack to generate all
// possible combinations of character deletions from the input rune slice.
// It's the core algorithm for building the symmetric delete index that enables
// fast spell checking lookups.
//
// The algorithm works by:
// 1. Starting with the original rune sequence
// 2. For each position, creating a variant with that character deleted
// 3. Recursively applying deletions to each variant
// 4. Stopping when the maximum edit distance is reached
// 5. Calling the provided function with the hash of each deletion variant
//
// Args:
//
//	cur: The current rune slice to generate deletions from.
//	fn: A callback function that receives the hash of each generated delete variant.
//
// Note:
//
//	This is a computationally intensive operation. The number of variants grows
//	exponentially with edit distance: O(n^d) where n is word length and d is edit distance.
//	For typical values (word length ~10, edit distance 2), this generates hundreds of variants.
func (d *Dictionary) walkDeletes(cur []rune, fn func(deleteKey string)) {
	type node struct {
		runes []rune
		dist  int
		start int
	}
	stack := []node{{runes: append([]rune(nil), cur...), dist: 0, start: 0}}

	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if n.dist >= d.maxDictionaryEditDistance {
			continue
		}
		rl := len(n.runes)
		if rl == 0 {
			continue
		}

		for i := n.start; i < rl; i++ {
			buf := make([]rune, rl-1)
			copy(buf[:i], n.runes[:i])
			copy(buf[i:], n.runes[i+1:])
			str := string(buf)
			fn(str)

			next := make([]rune, rl-1)
			copy(next, buf)
			stack = append(stack, node{runes: next, dist: n.dist + 1, start: i})
		}
	}
}

func appendIfMissing(a []*DictItem, x *DictItem) []*DictItem {
	for _, v := range a {
		if v.text == x.text {
			return a
		}
	}
	return append(a, x)
}
