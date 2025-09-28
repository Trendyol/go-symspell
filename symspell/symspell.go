package symspell

import (
	"errors"
	"math"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/Trendyol/go-symspell/common"
	"github.com/Trendyol/go-symspell/edit_distance"
	"github.com/Trendyol/go-symspell/suggest"
	"github.com/Trendyol/go-symspell/verbosity"
)

var (
	// wordPattern matches a “word” token: letters/digits, optionally containing a single apostrophe.
	wordPattern = regexp.MustCompile(`([\p{L}\p{N}]+['’]*[\p{L}\p{N}]*)`)
)

const (
	// N is the total number of tokens in the corpus used to build the frequency dictionary.
	// It is used to convert counts (c) into probabilities (p) via: p=c/N.
	N = float64(1024908267229)
)

// SymSpell implements the Symmetric Delete spell-correction algorithm for fast and accurate spell checking.
//
// SymSpell is a high-performance spell checker that uses a novel approach called "Symmetric Delete"
// to achieve fast lookup times while maintaining excellent correction quality. Unlike traditional
// spell checkers that generate variants of the input word, SymSpell precomputes all possible
// character deletions of dictionary words up to a specified edit distance.
//
// Example:
//
//	symspell, err := NewSymSpell(WithMaxDictionaryEditDistance(2))
//	suggestions, err := symspell.Lookup("speling", verbosity.Top, 2)
type SymSpell struct {
	// dictionary holds the word database and precomputed delete variants for spell correction.
	// It contains all valid words, their frequencies, and the symmetric delete index that
	// enables fast spell correction lookups.
	dictionary *Dictionary
	// distanceComparer calculates edit distances between words for ranking suggestions.
	// Different algorithms (Levenshtein, Damerau-OSA) provide different trade-offs
	// between speed and accuracy in handling various types of spelling errors.
	distanceComparer edit_distance.DistanceComparer
	// includeUnknown determines whether unknown words are included in suggestion results.
	// When true, the original input word is always returned as a suggestion even if
	// not found in the dictionary, which helps preserve proper nouns and technical terms.
	includeUnknown bool
	// transferCasing controls whether the original word's capitalization is applied to suggestions.
	// When enabled, suggestions adopt the input word's casing pattern (all caps, title case, etc.)
	// to provide more natural-looking corrections.
	transferCasing bool
	// ignoreToken is a regex pattern for identifying tokens that should be ignored.
	// Tokens matching this pattern are passed through unchanged without spell checking,
	// useful for URLs, email addresses, or other structured text that shouldn't be corrected.
	ignoreToken *regexp.Regexp
	// ignoreNonWords determines whether non-word tokens are skipped during spell checking.
	// When enabled, tokens containing only punctuation, numbers, or symbols are ignored,
	// improving performance on mixed content with lots of non-word elements.
	ignoreNonWords bool
	// splitBySpace enables automatic compound word splitting for unknown terms.
	// When enabled, unknown words are split and each part checked separately,
	// helpful for compound words or accidentally concatenated terms.
	splitBySpace bool
	// ignoreTermWithDigits controls whether words containing digits are skipped.
	// When enabled, terms like "IPv4", "model123", or "3rd" are ignored during
	// spell checking, useful for technical content with intentional alphanumeric terms.
	ignoreTermWithDigits bool
}

// NewSymSpell creates and initializes a new SymSpell spell checker instance with the specified configuration.
//
// This constructor function creates a fully configured SymSpell instance using the functional
// options pattern. It applies all provided options to the default configuration, validates
// the resulting parameters for consistency and correctness, then initializes all internal
// components including the dictionary and edit distance comparer.
//
// The function performs comprehensive parameter validation to ensure the configuration
// is mathematically and algorithmically sound. Invalid configurations are rejected
// with descriptive error messages.
//
// Args:
//
//	options: Variable number of functional options to configure the spell symspell.
//	        Options are applied in order, with later options overriding earlier ones.
//	        If no options are provided, DefaultOptions configuration is used.
//
// Returns:
//
//	A pointer to a fully initialized SymSpell instance ready for dictionary loading
//	and spell checking operations. Returns nil and an error if the configuration
//	is invalid or initialization fails.
//
// Configuration validation rules:
//   - MaxDictionaryEditDistance >= 0 (must be non-negative)
//   - PrefixLength >= 1 (must be at least 1 character)
//   - PrefixLength > MaxDictionaryEditDistance (prefix must be longer than max edit distance)
//   - CountThreshold >= 0 (frequency threshold must be non-negative)
//   - DistanceAlgorithm must be a valid algorithm type
//
// Error conditions:
//   - Invalid parameter combinations (violate validation rules above)
//   - Unsupported or invalid edit distance algorithm
//   - Internal initialization failures
//
// Example:
//
//	// Create with default configuration
//	symspell, err := NewSymSpell()
//
//	// Create with custom configuration
//	symspell, err := NewSymSpell(
//	    WithMaxDictionaryEditDistance(2),
//	    WithPrefixLength(7),
//	    WithIncludeUnknown(true),
//	    WithTransferCasing(true),
//	)
//
//	// Handle initialization errors
//	if err != nil {
//	    log.Fatal("Failed to create spell checker:", err)
//	}
func NewSymSpell(options ...Options) (*SymSpell, error) {
	opts := DefaultOptions
	for _, opt := range options {
		opt.Apply(&opts)
	}
	if opts.MaxDictionaryEditDistance < 0 {
		return nil, errors.New("maxDictionaryEditDistance cannot be negative")
	}
	if opts.PrefixLength < 1 {
		return nil, errors.New("prefixLength cannot be less than 1")
	}
	if opts.PrefixLength <= opts.MaxDictionaryEditDistance {
		return nil, errors.New("prefixLength must be greater than maxDictionaryEditDistance")
	}
	if opts.CountThreshold < 0 {
		return nil, errors.New("countThreshold cannot be negative")
	}

	distanceComparer, err := edit_distance.NewDistanceComparer(opts.DistanceAlgorithm)
	if err != nil {
		return nil, err
	}

	return &SymSpell{
		dictionary:           NewDictionary(opts.InitialCapacity, opts.PrefixLength, opts.MaxDictionaryEditDistance, opts.CountThreshold),
		distanceComparer:     distanceComparer,
		includeUnknown:       opts.IncludeUnknown,
		transferCasing:       opts.TransferCasing,
		ignoreToken:          opts.IgnoreToken,
		ignoreNonWords:       opts.IgnoreNonWords,
		ignoreTermWithDigits: opts.IgnoreTermWithDigits,
		splitBySpace:         opts.SplitBySpace,
	}, nil
}

// EntryCount returns the total number of delete variant entries in the spell checker's dictionary.
//
// Example:
//
//	entries := symspell.EntryCount()
//	fmt.Printf("Dictionary has %d delete entries\n", entries)
func (s *SymSpell) EntryCount() int { return s.dictionary.EntryCount() }

// WordCount returns the total number of unique words stored in the spell checker's dictionary.
//
// Example:
//
//	words := symspell.WordCount()
//	fmt.Printf("Dictionary contains %d unique words\n", words)
func (s *SymSpell) WordCount() int { return s.dictionary.WordCount() }

// MaxLength returns the length in runes of the longest valid word in the dictionary.
//
// Example:
//
//	maxLen := symspell.MaxLength()
//	fmt.Printf("Longest word has %d characters\n", maxLen)
func (s *SymSpell) MaxLength() int { return s.dictionary.MaxLength() }

// CreateDictionary loads a dictionary from a plain text corpus file using the spell checker's word parsing rules.
//
// This method is a high-level wrapper around the dictionary's CreateDictionary method that
// automatically applies the spell checker's configured word parsing logic. It reads the
// text file line by line, extracts words using the wordPattern regex, and adds each
// word to the dictionary with a frequency count of 1.
//
// Args:
//
//	path: The file path to the plain text corpus to load.
//	     The file should contain natural language text, one sentence or paragraph per line.
//
// Returns:
//
//	True if the file was successfully processed, false if file opening failed.
//	An error if file I/O operations encounter problems.

// Example:
//
//	success, err := symspell.CreateDictionary("corpus.txt")
//	if err != nil {
//	    log.Fatal("Failed to load corpus:", err)
//	}
//	fmt.Printf("Dictionary loaded: %v, Word count: %d\n", success, symspell.WordCount())
//
// Note:
//
//	This method merges with any existing dictionary data. Words are added with frequency 1,
//	which may accumulate if the same word appears multiple times in the corpus.
func (s *SymSpell) CreateDictionary(path string) (bool, error) {
	return s.dictionary.CreateDictionary(path, s.parseWords)
}

// LoadDictionary loads word-frequency pairs from a structured dictionary file.
//
// This method loads pre-computed word frequency data from structured files. Each line should contain
// delimited fields where specific columns contain the word and its frequency count.
//
// Args:
//
//	path: The file path to the structured dictionary file.
//	termIndex: Zero-based column index containing the word terms.
//	countIndex: Zero-based column index containing the frequency counts.
//	sep: Field separator string used to split each line (e.g., " ", "\t", ",", "|").
//
// Returns:
//
//	True if the file was successfully processed, false if file opening failed.
//	An error if file I/O operations encounter problems.
//
// File format:
//
//	Each line: field1<sep>field2<sep>...fieldN
//	Where one field contains the word and another contains its frequency count.
//	Lines with insufficient fields or invalid counts are silently skipped.
//
// Example usage:
//
//	// Load tab-separated file: "word\t123\tother_data"
//	success, err := symspell.LoadDictionary("frequency.txt", 0, 1, "\t")
//	if err != nil {
//	    log.Fatal("Failed to load dictionary:", err)
//	}
//
// Note:
//
//	This method merges with existing dictionary data, accumulating frequency counts
//	for words that already exist in the dictionary.
func (s *SymSpell) LoadDictionary(path string, termIndex, countIndex int, sep string) (bool, error) {
	return s.dictionary.LoadDictionary(path, termIndex, countIndex, sep)
}

// LoadBigramDictionary loads bigram frequency data from a structured file to improve context-aware corrections.
//
// This method loads word pair (bigram) frequency data. Bigrams help the spell checker understand which
// word combinations are common, enabling better correction suggestions in context.
//
// Args:
//
//	path: The file path to the bigram frequency file.
//	termIndex: Zero-based column index containing the bigram terms (e.g., "word1 word2").
//	countIndex: Zero-based column index containing the bigram frequency counts.
//	sep: Field separator string used to split each line.
//
// Returns:
//
//	True if the file was successfully processed, false if file opening failed.
//	An error if file I/O operations encounter problems.
//
// File format:
//
//	Each line: field1<sep>field2<sep>...fieldN
//	Where one field contains space-separated word pairs and another contains frequency.
//	Example: "the cat\t1250\tother_data"
//
// Example usage:
//
//	// Load bigram data for context-aware corrections
//	success, err := symspell.LoadBigramDictionary("bigrams.txt", 0, 1, "\t")
//	if err != nil {
//	    log.Fatal("Failed to load bigrams:", err)
//	}
//
// Note:
//
//	Bigram data significantly improves correction quality but is optional.
//	The spell checker works without bigrams but provides better results with them.
func (s *SymSpell) LoadBigramDictionary(path string, termIndex, countIndex int, sep string) (bool, error) {
	return s.dictionary.LoadBigramDictionary(path, termIndex, countIndex, sep)
}

// Lookup finds spelling suggestions for a given word or phrase using the symmetric delete algorithm.
//
// This is the main spell checking method that implements the core SymSpell algorithm.
// It finds potential corrections for misspelled words by using precomputed delete variants
// and calculating edit distances to rank suggestions by likelihood and frequency.
//
// The algorithm works by:
//  1. Generating character deletions of the input word up to maxEditDistance
//  2. Looking up each deletion variant in the precomputed symmetric delete index
//  3. Calculating precise edit distances for candidate words found
//  4. Ranking suggestions by edit distance and word frequency
//  5. Applying configured filters and transformations to the results
//
// Args:
//
//	phrase: The input word or phrase to find spelling suggestions for.
//	       Single words work best; phrases are treated as single tokens.
//	vb: Verbosity level controlling how many suggestions to return:
//	    - verbosity.Top: Return only the best suggestion
//	    - verbosity.Closest: Return all suggestions with minimum edit distance
//	    - verbosity.All: Return all suggestions within maxEditDistance
//	maxEditDistance: Maximum edit distance to consider for suggestions.
//	                Must be <= dictionary's MaxDictionaryEditDistance.
//	                Use -1 to use dictionary's maximum edit distance.
//
// Returns:
//
//	A slice of SuggestItem structs containing suggested corrections, sorted by:
//	1. Edit distance (ascending) - closer matches first
//	2. Word frequency (descending) - more common words first
//	3. Alphabetical order (ascending) - consistent tie-breaking
//
//	Returns an error if maxEditDistance is invalid or algorithm fails.
//
// Example usage:
//
//	// Find the best correction
//	suggestions, err := symspell.Lookup("speling", verbosity.Top, 2)
//	if len(suggestions) > 0 {
//	    fmt.Printf("Did you mean: %s?\n", suggestions[0].Term)
//	}
//
//	// Find all close corrections
//	suggestions, err := symspell.Lookup("recieve", verbosity.Closest, 2)
//	for _, sug := range suggestions {
//	    fmt.Printf("Suggestion: %s (distance: %d, freq: %d)\n",
//	              sug.Term, sug.Distance, sug.Count)
//	}
func (s *SymSpell) Lookup(phrase string, vb verbosity.Verbosity, maxEditDistance int) ([]suggest.SuggestItem, error) {
	if maxEditDistance < 0 {
		maxEditDistance = s.dictionary.MaxDictionaryEditDistance()
	}
	if maxEditDistance > s.dictionary.MaxDictionaryEditDistance() {
		return nil, errors.New("distance too large")
	}

	orig := phrase
	if s.transferCasing {
		phrase = strings.ToLower(phrase)
	}
	phraseLen := utf8.RuneCountInString(phrase)
	prefixLen := s.dictionary.PrefixLength()

	suggestions := make([]suggest.SuggestItem, 0, 8)

	// Early exit, word is too big to possibly match any words
	if phraseLen-maxEditDistance > s.dictionary.MaxLength() {
		return s.earlyExit(suggestions, orig, maxEditDistance), nil
	}

	// Quick look for exact match
	if item := s.dictionary.LookupWord(phrase); item != nil {
		term := orig
		if !s.transferCasing {
			term = phrase
		}
		suggestions = append(suggestions, suggest.SuggestItem{Term: term, Distance: 0, Count: item.count})
		// Return exact match, unless client wants all matches
		if vb != verbosity.All {
			return s.earlyExit(suggestions, orig, maxEditDistance), nil
		}
	}
	if s.ignoreToken != nil && s.ignoreToken.MatchString(phrase) {
		suggestions = append(suggestions, suggest.SuggestItem{Term: phrase, Distance: 0, Count: 1})
		// Return exact match, unless client wants all matches
		if vb != verbosity.All {
			return s.earlyExit(suggestions, orig, maxEditDistance), nil
		}
	}

	// Early termination, if we only want to check if word in dictionary or
	// get its frequency.
	if maxEditDistance == 0 {
		return s.earlyExit(suggestions, orig, maxEditDistance), nil
	}

	consideredDeletes := make(map[string]struct{}, 64)
	// We considered the input phrase already in the Words
	consideredSuggestions := make(map[string]struct{}, 64)
	candidates := make([]candidate, 0, 64)
	candidatePointer := 0

	// Add original prefix
	phrasePrefixLen := phraseLen
	if phrasePrefixLen > prefixLen {
		phrasePrefixLen = prefixLen
	}
	firstRunes := make([]rune, phrasePrefixLen)
	phraseRunes := []rune(phrase)
	copy(firstRunes, phraseRunes[:phrasePrefixLen])
	candidates = append(candidates, candidate{text: string(firstRunes), runes: firstRunes})

	maxEditDistance2 := maxEditDistance

	// --- Main loop over delete candidates ---
	for candidatePointer < len(candidates) {
		c := candidates[candidatePointer]
		candidatePointer++

		// Early termination: if candidate distance is already higher than
		// suggestion distance, then there are no better suggestions to be
		// expected
		candidateLen := c.Len()
		lenDiff := phrasePrefixLen - candidateLen
		if lenDiff > maxEditDistance2 {
			// Skip to next candidate if Verbosity.ALL, look no further if Verbosity.TOP or CLOSEST
			// (candidates are ordered by delete distance, so none are closer than current)
			if vb == verbosity.All {
				// `max_edit_distance_2` only updated when verbosity != ALL.
				// New candidates are generated from deletes so it keeps getting shorter.
				// This should never be reached.
				continue
			}
			break
		}

		// Retrieve matching dictionary suggestions for this delete key
		if items, ok := s.dictionary.LookupDelete(c.text); ok {
			for _, item := range items {
				// Duplicate suggestion guard
				if _, seen := consideredSuggestions[item.text]; seen {
					continue
				}

				suggestion := item.text
				suggestionRunes := item.runes
				suggestionLen := item.Len()

				// Quick bounds check
				if common.Abs(suggestionLen-phraseLen) > maxEditDistance2 ||
					suggestionLen < candidateLen {
					continue
				}
				suggestionPrefixLen := min(suggestionLen, prefixLen)
				if suggestionPrefixLen > phrasePrefixLen && suggestionPrefixLen-candidateLen > maxEditDistance2 {
					continue
				}

				// True Damerau-Levenshtein Edit Distance: adjust distance, if both distances>0.
				// We allow simultaneous edits (deletes) of max_edit_distance on on both the dictionary and the phrase term.
				// For replaces and adjacent transposes the resulting edit distance stays <= max_edit_distance.
				// For inserts and deletes the resulting edit distance might exceed max_edit_distance.
				// To prevent suggestions of a higher edit distance, we need to calculate the resulting edit distance,
				// if there are simultaneous edits on both sides. Example: (bank==bnak and bank==bink, but bank!=kanb and bank!=xban and bank!=baxn for max_edit_distance=1).
				// Two deletes on each side of a pair makes them all equal, but the first two pairs have edit distance=1, the others edit distance=2.
				var distance int
				switch {
				case candidateLen == 0:
					// Suggestions which have no common chars with phrase
					distance = max(phraseLen, suggestionLen)
					if distance > maxEditDistance2 {
						continue
					}
				case suggestionLen == 1:
					// `suggestion` only gets added to `consideredSuggestions` when `suggestionLen>1`.
					// Given the maxDictionaryEditDistance and prefixLength restrictions, `distance` should never be >maxEditDistance2
					if !common.ContainsRune(phraseRunes, suggestionRunes[0]) {
						distance = phraseLen
					} else {
						distance = phraseLen - 1
					}
					if distance > maxEditDistance2 {
						continue
					}
				default:
					// Number of edits in prefix ==maxEditDistance AND no identical suffix, then editDistance>maxEditDistance and no need for Levenshtein calculation
					// (phraseLen >= prefixLength) && (suggestionLen >= prefixLength)

					minDistance := 0
					// Handles the shortcircuit of minDistance assignment when first boolean expression evaluates to False
					if prefixLen-maxEditDistance == candidateLen {
						minDistance = min(phraseLen, suggestionLen) - prefixLen
					}
					if prefixLen-maxEditDistance == candidateLen {
						if minDistance > 1 {
							startP := phraseLen + 1 - minDistance
							startS := suggestionLen + 1 - minDistance
							if !common.RunesEqual(phraseRunes[startP:], suggestionRunes[startS:]) {
								continue
							}
						}
						if minDistance > 0 {
							iP := phraseLen - minDistance
							iS := suggestionLen - minDistance
							sugRunes := suggestionRunes
							if !common.RuneIndexEqual(phraseRunes, iP, sugRunes, iS) &&
								(!common.RuneIndexEqual(phraseRunes, iP-1, sugRunes, iS) ||
									!common.RuneIndexEqual(phraseRunes, iP, sugRunes, iS-1)) {
								continue
							}
						}
					}

					// Mark seen before costly distance compare
					consideredSuggestions[item.text] = struct{}{}

					tmp, err := s.distanceComparer.Compare(phrase, suggestion, maxEditDistance2)
					if err != nil {
						return nil, err
					}
					if tmp < 0 {
						continue
					}
					distance = tmp
				}

				// Do not process higher distances than those already found, if verbosity<ALL
				// (note: maxEditDistance2 will always equal maxEditDistance when Verbosity.ALL)
				if distance <= maxEditDistance2 {
					sugCount := item.count
					item := suggest.SuggestItem{Term: suggestion, Distance: distance, Count: sugCount}

					if len(suggestions) > 0 {
						switch vb {
						case verbosity.Closest:
							// We will calculate distance only to the smallest found distance so far
							if distance < maxEditDistance2 {
								suggestions = suggestions[:0]
							}
						case verbosity.Top:
							if distance < maxEditDistance2 || sugCount > suggestions[0].Count {
								maxEditDistance2 = distance
								suggestions[0] = item
							}
							continue
						}
					}
					if vb != verbosity.All {
						maxEditDistance2 = distance
					}
					suggestions = append(suggestions, item)
				}
			}
		}

		// Add edits: derive edits (deletes) from candidate (phrase) and add them to candidates list.
		if lenDiff < maxEditDistance && candidateLen <= prefixLen {
			if vb != verbosity.All && lenDiff >= maxEditDistance2 {
				continue
			}
			if candidateLen > 0 {
				for i := 0; i < candidateLen; i++ {
					buf := make([]rune, candidateLen-1)
					copy(buf[:i], c.runes[:i])
					copy(buf[i:], c.runes[i+1:])
					str := string(buf)
					if _, ok := consideredDeletes[str]; ok {
						continue
					}
					consideredDeletes[str] = struct{}{}
					candidates = append(candidates, candidate{
						text:  str,
						runes: buf,
					})
				}
			}
		}
	}

	if len(suggestions) > 1 {
		sort.Sort(suggest.SuggestItems(suggestions))
	}

	if s.transferCasing {
		for i := range suggestions {
			suggestions[i].Term = common.CaseTransferSimilar(orig, suggestions[i].Term)
		}
	}
	suggestions = s.earlyExit(suggestions, orig, maxEditDistance)
	return suggestions, nil
}

// earlyExit handles the final processing of suggestion results with unknown word inclusion.
//
// This helper method applies the includeUnknown configuration setting to provide
// a fallback suggestion when no dictionary-based corrections are found. It ensures
// that user input is never completely rejected, which is useful for preserving
// proper nouns, technical terms, or intentionally non-dictionary words.
//
// Args:
//
//	suggestions: The current list of spelling suggestions found by the algorithm.
//	original: The original input word that was being spell checked.
//	maxEditDistance: The maximum edit distance used for the spell checking operation.
//
// Returns:
//
//	The modified suggestions list, potentially with the original word added
//	as a fallback suggestion if includeUnknown is enabled and no other
//	suggestions were found.
func (s *SymSpell) earlyExit(suggestions []suggest.SuggestItem, original string, maxEditDistance int) []suggest.SuggestItem {
	if s.includeUnknown && len(suggestions) == 0 {
		suggestions = append(suggestions, suggest.SuggestItem{
			Term:     original,
			Distance: maxEditDistance + 1,
			Count:    0,
		})
	}
	return suggestions
}

// LookupCompound performs advanced compound word spelling correction for multi-word phrases.
//
// This method provides sophisticated spell checking for phrases containing multiple words
// by analyzing each word individually and considering word combination possibilities.
// It handles word splitting, compound word detection, and context-aware correction
// using bigram statistics when available.
//
// Args:
//
//	phrase: The input phrase containing one or more words to spell check.
//	       Supports natural language text with multiple words separated by spaces.
//	maxEditDistance: Maximum edit distance for individual word corrections.
//	                Applied to each word separately during the correction process.
//
// Returns:
//
//	A slice containing a single SuggestItem representing the best correction
//	for the entire phrase. The suggestion combines optimal corrections for
//	all words with appropriate spacing and statistical ranking.
//
//	Returns an error if the algorithm fails or encounters invalid parameters.
//
// Example usage:
//
//	// Correct a multi-word phrase
//	suggestions, err := checker.LookupCompound("speling is importent", 2)
//	if err == nil && len(suggestions) > 0 {
//	    fmt.Printf("Corrected: %s\n", suggestions[0].Term)
//	    // Output: "spelling is important"
//	}
//
//	// Handle compound words and splitting
//	suggestions, err := checker.LookupCompound("spellchecker", 2)
//	if err == nil && len(suggestions) > 0 {
//	    fmt.Printf("Suggestion: %s\n", suggestions[0].Term)
//	    // May output: "spell checker" (if splitting improves the result)
//	}
func (s *SymSpell) LookupCompound(phrase string, maxEditDistance int) ([]suggest.SuggestItem, error) {
	terms1 := s.parseWords(phrase)
	var terms2 []string
	if s.ignoreNonWords {
		terms2 = terms1
	}

	suggestionParts := make([]suggest.SuggestItem, 0)
	isLastCombination := false

	// Translate every item to its best suggestion, otherwise it remains unchanged
	for i, term1 := range terms1 {
		if s.ignoreNonWords {
			if ok, _ := common.TryParseInt64(term1); ok {
				suggestionParts = append(suggestionParts, suggest.SuggestItem{
					Term:     term1,
					Distance: 0,
					Count:    int(N),
				})
				continue
			}
			if common.IsAcronym(terms2[i], s.ignoreTermWithDigits) {
				suggestionParts = append(suggestionParts, suggest.SuggestItem{
					Term:     terms2[i],
					Distance: 0,
					Count:    int(N),
				})
				continue
			}
		}

		suggestions, err := s.Lookup(term1, verbosity.Top, maxEditDistance)
		if err != nil {
			return nil, err
		}

		// Combination check, always before split
		if i > 0 && !isLastCombination {
			prevTerm := terms1[i-1]
			suggestionsCombi, err := s.Lookup(prevTerm+term1, verbosity.Top, maxEditDistance)
			if err != nil {
				return nil, err
			}

			if len(suggestionsCombi) > 0 {
				best1 := suggestionParts[len(suggestionParts)-1]
				var best2 suggest.SuggestItem
				if len(suggestions) > 0 {
					best2 = suggestions[0]
				} else {
					// Estimated word occurrence probability
					// P=10 / (N * 10^word length l)
					best2 = suggest.NewSuggestItemWithProbability(term1, maxEditDistance+1)
				}

				// distance1=edit distance between 2 split terms and their best corrections: also comparative value for the combination
				distance1 := best1.Distance + best2.Distance
				if distance1 >= 0 &&
					(suggestionsCombi[0].Distance+1 < distance1 ||
						(suggestionsCombi[0].Distance+1 == distance1 &&
							(float64(suggestionsCombi[0].Count) > float64(best1.Count)/float64(N)*float64(best2.Count)))) {
					suggestionsCombi[0].Distance++
					suggestionParts[len(suggestionParts)-1] = suggestionsCombi[0]
					isLastCombination = true
					continue
				}
			}
		}
		isLastCombination = false

		lenTerm1 := utf8.RuneCountInString(term1)
		// Always split terms without suggestion / never split terms with suggestion ed=0 / never split single char terms
		if len(suggestions) > 0 && (suggestions[0].Distance == 0 || lenTerm1 == 1) {
			suggestionParts = append(suggestionParts, suggestions[0])
		} else {
			// If no perfect suggestion, split word into pairs
			var suggestionSplitBest *suggest.SuggestItem
			if len(suggestions) > 0 {
				suggestionSplitBest = &suggestions[0]
			}

			if lenTerm1 > 1 {
				for j := 1; j < lenTerm1; j++ {
					term1Runes := []rune(term1)
					part1 := string(term1Runes[:j])
					part2 := string(term1Runes[j:])

					suggestions1, err := s.Lookup(part1, verbosity.Top, maxEditDistance)
					if err != nil || len(suggestions1) == 0 {
						continue
					}
					suggestions2, err := s.Lookup(part2, verbosity.Top, maxEditDistance)
					if err != nil || len(suggestions2) == 0 {
						continue
					}

					// Select best suggestion for split pair
					tmpTerm := suggestions1[0].Term + " " + suggestions2[0].Term
					tmpDistance, err := s.distanceComparer.Compare(term1, tmpTerm, maxEditDistance)
					if err != nil {
						return nil, err
					}
					if tmpDistance < 0 {
						tmpDistance = maxEditDistance + 1
					}

					if suggestionSplitBest != nil && tmpDistance > suggestionSplitBest.Distance {
						continue
					}
					if suggestionSplitBest != nil && tmpDistance < suggestionSplitBest.Distance {
						suggestionSplitBest = nil
					}

					tmpCount := 0
					item1 := s.dictionary.LookupWord(suggestions1[0].Term)
					item2 := s.dictionary.LookupWord(suggestions2[0].Term)
					if item1 != nil && item2 != nil {
						term := suggestions1[0].Term + " " + suggestions2[0].Term
						if count := s.dictionary.GetBigramCount(term); count > 0 {
							tmpCount = count
							if len(suggestions) > 0 {
								bestSi := suggestions[0]
								if suggestions1[0].Term+suggestions2[0].Term == term1 {
									tmpCount = max(tmpCount, bestSi.Count+2)
								} else if bestSi.Term == suggestions1[0].Term && bestSi.Term == suggestions2[0].Term {
									tmpCount = max(tmpCount, bestSi.Count+1)
								}
							} else if suggestions1[0].Term+suggestions2[0].Term == term1 {
								tmpCount = max(tmpCount, max(suggestions1[0].Count, suggestions2[0].Count)+2)
							}
						} else {
							// The Naive Bayes probability of the word combination is the product of the two word probabilities: P(AB)=P(A)*P(B)
							// Use it to estimate the frequency count of the combination, which then is used to rank/select the best splitting variant
							tmpCount = min(s.dictionary.BigramCountMin(),
								int(float64(suggestions1[0].Count)/N*float64(suggestions2[0].Count)))
						}
					}

					suggestionSplit := suggest.SuggestItem{Term: tmpTerm, Distance: tmpDistance, Count: tmpCount}
					if suggestionSplitBest == nil || suggestionSplit.Count > suggestionSplitBest.Count {
						suggestionSplitBest = &suggestionSplit
					}
				}

				if suggestionSplitBest != nil {
					// Select best suggestion for split pair
					suggestionParts = append(suggestionParts, *suggestionSplitBest)
				} else {
					item := suggest.NewSuggestItemWithProbability(term1, maxEditDistance+1)
					suggestionParts = append(suggestionParts, item)
				}
			} else {
				item := suggest.NewSuggestItemWithProbability(term1, maxEditDistance+1)
				suggestionParts = append(suggestionParts, item)
			}
		}
	}

	var sb strings.Builder
	joinedCount := N
	for _, item := range suggestionParts {
		sb.WriteString(item.Term)
		sb.WriteByte(' ')
		joinedCount *= float64(item.Count) / N
	}
	joinedTerm := strings.TrimSpace(sb.String())

	if s.transferCasing {
		joinedTerm = common.CaseTransferSimilar(phrase, joinedTerm)
	}

	finalDistance, err := s.distanceComparer.Compare(phrase, joinedTerm, math.MaxInt32)
	if err != nil {
		return nil, err
	}

	finalSuggestion := suggest.SuggestItem{Term: joinedTerm, Distance: finalDistance, Count: int(joinedCount)}
	return []suggest.SuggestItem{finalSuggestion}, nil
}

// parseWords extracts individual word tokens from input text using configured parsing rules.
func (s *SymSpell) parseWords(text string) []string {
	return wordPattern.FindAllString(strings.ToLower(text), -1)
}
