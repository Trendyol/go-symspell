package symspell

import (
	"regexp"

	"github.com/Trendyol/go-symspell/edit_distance"
)

// DefaultOptions provides a sensible default configuration for SymSpell spell checking.
//
// These default values are optimized for general-purpose spell checking with
// good performance characteristics. They provide a balance between correction
// quality and computational efficiency suitable for most applications.
//
// Default configuration:
//   - InitialCapacity: 16 - Minimal initial memory allocation
//   - MaxDictionaryEditDistance: 2 - Handles most common typos (insertions, deletions, substitutions)
//   - PrefixLength: 7 - Good balance between memory usage and correction quality
//   - CountThreshold: 1 - Accept all words immediately (no frequency threshold)
//   - DistanceAlgorithm: DamerauOSAFast - Fast algorithm including transpositions
//   - IncludeUnknown: false - Don't include unknown words in suggestions
//   - TransferCasing: false - Don't transfer original word casing to suggestions
//   - IgnoreToken: nil - No regex-based token filtering
//   - IgnoreNonWords: false - Process all input tokens
//   - IgnoreTermWithDigits: false - Process words containing digits
//   - SplitBySpace: false - Don't automatically split compound words
//
// Example:
//
//	checker := NewSymSpell()  // Use default configuration
var DefaultOptions = options{
	InitialCapacity:           16,
	MaxDictionaryEditDistance: 2,
	PrefixLength:              7,
	CountThreshold:            1,
	DistanceAlgorithm:         edit_distance.DamerauOSAFast,
	IncludeUnknown:            false,
	TransferCasing:            false,
	IgnoreToken:               nil,
	IgnoreNonWords:            false,
	IgnoreTermWithDigits:      false,
	SplitBySpace:              false,
}

// options holds all configuration parameters for SymSpell spell checking behavior.
//
// This struct contains all the configurable parameters that control how the
// SymSpell algorithm operates, from dictionary initialization through spell
// checking and suggestion generation. Each field affects different aspects
// of the spell checking process.
type options struct {
	// InitialCapacity sets the expected number of dictionary words for memory pre-allocation.
	// Higher values reduce memory reallocations during dictionary loading but use more
	// initial memory. Set to 0 if dictionary size is unknown.
	// Typical values: 16 (small), 1000 (medium), 100000 (large dictionaries).
	InitialCapacity int
	// MaxDictionaryEditDistance defines the maximum edit distance for precomputed delete variants.
	// This is the key parameter controlling correction quality vs memory usage.
	// Higher values find more corrections but exponentially increase memory requirements.
	// Typical values: 1 (fast, minimal memory), 2 (balanced), 3 (high quality, more memory).
	MaxDictionaryEditDistance int
	// PrefixLength limits the length of word prefixes used for delete key generation.
	// Longer prefixes improve correction quality but increase memory usage.
	// Words longer than this are truncated before generating delete variants.
	// Typical values: 5-7 (balanced), 10+ (high quality, more memory).
	PrefixLength int
	// CountThreshold sets the minimum frequency required for words to be considered valid.
	// Words below this threshold are accumulated until they reach the minimum count.
	// Higher values filter out rare/misspelled words but may reject valid low-frequency terms.
	// Values: 1 (accept all), 5-10 (filter rare words), 100+ (only common words).
	CountThreshold int
	// DistanceAlgorithm specifies which edit distance algorithm to use for ranking suggestions.
	// Different algorithms have different performance characteristics and correction types:
	// - Levenshtein: Insertions, deletions, substitutions
	// - DamerauOSA: Adds single-character transpositions
	// - DamerauOSAFast: Optimized version of DamerauOSA (recommended)
	DistanceAlgorithm edit_distance.AlgorithmType
	// IncludeUnknown determines whether to include the original (unknown) word in suggestions.
	// When true, the input word is always included as a suggestion even if not in dictionary.
	// Useful for preserving proper nouns, technical terms, or intentional misspellings.
	IncludeUnknown bool
	// TransferCasing controls whether to preserve the original word's casing in suggestions.
	// When true, suggestions adopt the capitalization pattern of the input word.
	// Examples: "HELLO" -> "WORLD" becomes "WORLD", "Hello" -> "world" becomes "World".
	TransferCasing bool
	// IgnoreToken is a regex pattern for tokens that should be ignored during spell checking.
	// Tokens matching this pattern are passed through unchanged without correction attempts.
	// Useful for URLs, email addresses, special codes, or other non-word tokens.
	// Set to nil to disable token filtering.
	IgnoreToken *regexp.Regexp
	// IgnoreNonWords determines whether to skip tokens that don't look like words.
	// When true, tokens containing only punctuation, numbers, or special characters
	// are passed through without spell checking. Helps performance with mixed content.
	IgnoreNonWords bool
	// IgnoreTermWithDigits controls whether to skip words containing numeric digits.
	// When true, terms like "abc123", "3rd", or "IPv4" are ignored during spell checking.
	// Useful for technical content where alphanumeric terms are intentional.
	IgnoreTermWithDigits bool
	// SplitBySpace enables automatic splitting of compound words at space boundaries.
	// When true, the spell checker attempts to split unknown words and check parts separately.
	// Useful for languages with compound words or handling concatenated terms.
	// Example: "spellchecker" might be split into "spell checker".
	SplitBySpace bool
}

// Options defines the interface for configuring SymSpell spell checker behavior.
//
// This interface implements the functional options pattern, allowing for flexible
// and extensible configuration of SymSpell instances. Options can be combined
// and chained together to create custom configurations.
//
// Example:
//
//	checker := NewSymSpell(
//	    WithMaxDictionaryEditDistance(3),
//	    WithPrefixLength(10),
//	    WithIncludeUnknown(true),
//	)
type Options interface {
	// Apply modifies the provided options configuration with this option's settings.
	//
	// This method is called internally during SymSpell initialization to apply
	// each provided option to the base configuration. Options are applied in
	// the order they are provided, allowing later options to override earlier ones.
	//
	// Args:
	//     options: Pointer to the options struct to modify.
	Apply(options *options)
}

// FuncConfig implements the Options interface using a function-based approach.
type FuncConfig struct {
	// options is the configuration function that modifies an options struct.
	// This function encapsulates the specific configuration change this option represents.
	options func(options *options)
}

// Apply implements the Options interface by calling the wrapped configuration function.
//
// Args:
//
//	conf: Pointer to the options struct to modify.
func (w FuncConfig) Apply(conf *options) {
	w.options(conf)
}

// NewFuncOption creates a new FuncConfig option from a configuration function.
//
// Args:
//
//	f: A function that takes an options pointer and modifies it.
//
// Returns:
//
//	A new FuncConfig instance that implements the Options interface.
//
// Example:
//
//	customOption := NewFuncOption(func(o *options) { o.InitialCapacity = 1000 })
func NewFuncOption(f func(options *options)) *FuncConfig {
	return &FuncConfig{options: f}
}

// WithInitialCapacity sets the expected number of dictionary words for memory optimization.
//
// This option controls the initial capacity hint for internal data structures
// during dictionary creation. Setting an appropriate value reduces memory
// reallocations and improves loading performance when the dictionary size is known.
//
// Args:
//
//	initialCapacity: Expected number of unique words in the dictionary.
//	                Use 0 if unknown. Negative values are treated as 0.
//
// Returns:
//
//	An Options instance that sets the InitialCapacity configuration.
//
// Example:
//
//	// For a small dictionary (< 1000 words)
//	checker := NewSymSpell(WithInitialCapacity(1000))
//
//	// For a large dictionary (100k words)
//	checker := NewSymSpell(WithInitialCapacity(100000))
func WithInitialCapacity(initialCapacity int) Options {
	return NewFuncOption(func(options *options) {
		options.InitialCapacity = initialCapacity
	})
}

// WithMaxDictionaryEditDistance sets the maximum edit distance for precomputed delete variants.
//
// This is the most important parameter affecting both correction quality and memory usage.
// It determines how many character deletions are precomputed and indexed for each
// dictionary word. Higher values provide better correction recall but exponentially
// increase memory requirements and initialization time.
//
// Args:
//
//	maxDictionaryEditDistance: Maximum edit distance for delete variant indexing.
//	                          Typical values: 1 (minimal), 2 (balanced), 3 (comprehensive).
//
// Returns:
//
//	An Options instance that sets the MaxDictionaryEditDistance configuration.
//
// Example:
//
//	// Minimal memory usage, basic corrections
//	checker := NewSymSpell(WithMaxDictionaryEditDistance(1))
//
//	// Balanced approach (recommended for most use cases)
//	checker := NewSymSpell(WithMaxDictionaryEditDistance(2))
//
//	// High-quality corrections, more memory
//	checker := NewSymSpell(WithMaxDictionaryEditDistance(3))
func WithMaxDictionaryEditDistance(maxDictionaryEditDistance int) Options {
	return NewFuncOption(func(options *options) {
		options.MaxDictionaryEditDistance = maxDictionaryEditDistance
	})
}

// WithPrefixLength sets the maximum length of word prefixes used for delete key generation.
//
// This parameter controls how much of each word is used when generating delete
// variants. Longer prefixes improve correction quality but increase memory usage.
// Words longer than this value are truncated before delete variant generation.
//
// Args:
//
//	prefixLength: Maximum prefix length for delete key generation.
//	             Typical values: 5 (minimal), 7 (balanced), 10+ (comprehensive).
//
// Returns:
//
//	An Options instance that sets the PrefixLength configuration.
//
// Example:
//
//	// Memory-efficient for short words
//	checker := NewSymSpell(WithPrefixLength(5))
//
//	// Balanced approach (recommended)
//	checker := NewSymSpell(WithPrefixLength(7))
//
//	// Better for long words, more memory
//	checker := NewSymSpell(WithPrefixLength(12))
func WithPrefixLength(prefixLength int) Options {
	return NewFuncOption(func(options *options) {
		options.PrefixLength = prefixLength
	})
}

// WithCountThreshold sets the minimum frequency required for words to be considered valid.
//
// This parameter filters out low-frequency words that might be misspellings or
// rare terms. Words below the threshold are accumulated until they reach the
// minimum count. Higher values improve dictionary quality but may reject
// valid low-frequency terms.
//
// Args:
//
//	countThreshold: Minimum frequency count for word validation.
//	               Values: 1 (accept all), 5-10 (filter rare), 100+ (common only).
//
// Returns:
//
//	An Options instance that sets the CountThreshold configuration.
//
// Example:
//
//	// Accept all words (no frequency filtering)
//	checker := NewSymSpell(WithCountThreshold(1))
//
//	// Filter out rare words and likely typos
//	checker := NewSymSpell(WithCountThreshold(10))
//
//	// Only accept well-established words
//	checker := NewSymSpell(WithCountThreshold(100))
func WithCountThreshold(countThreshold int) Options {
	return NewFuncOption(func(options *options) {
		options.CountThreshold = countThreshold
	})
}

// WithDistanceAlgorithm sets the edit distance algorithm used for ranking suggestions.
//
// Different algorithms support different types of edits and have different
// performance characteristics. The choice affects both the accuracy of
// suggestion ranking and the computational cost.
//
// Args:
//
//	algorithmType: The edit distance algorithm to use for suggestion ranking.
//	              Available options from edit_distance package:
//	              - Levenshtein: Insertions, deletions, substitutions
//	              - DamerauOSA: Adds single-character transpositions
//	              - DamerauOSAFast: Optimized version (recommended)
//
// Returns:
//
//	An Options instance that sets the DistanceAlgorithm configuration.
//
// Example:
//
//	// Basic but fast algorithm
//	checker := NewSymSpell(WithDistanceAlgorithm(edit_distance.Levenshtein))
//
//	// Recommended: fast and handles transpositions
//	checker := NewSymSpell(WithDistanceAlgorithm(edit_distance.DamerauOSAFast))
//
//	// Most comprehensive (slower)
//	checker := NewSymSpell(WithDistanceAlgorithm(edit_distance.DamerauOSA))
func WithDistanceAlgorithm(algorithmType edit_distance.AlgorithmType) Options {
	return NewFuncOption(func(options *options) {
		options.DistanceAlgorithm = algorithmType
	})
}

// WithIncludeUnknown controls whether to include the original word in suggestions.
//
// When enabled, the original input word is always included as a suggestion
// even if it's not found in the dictionary. This is useful for preserving
// proper nouns, technical terms, or intentional non-dictionary words.
//
// Args:
//
//	includeUnknown: Whether to include unknown words in suggestion results.
//	               true = always include original word
//	               false = only include dictionary-based suggestions
//
// Returns:
//
//	An Options instance that sets the IncludeUnknown configuration.
//
// Example:
//
//	// Preserve original words even if unknown
//	checker := NewSymSpell(WithIncludeUnknown(true))
//
//	// Only suggest dictionary words
//	checker := NewSymSpell(WithIncludeUnknown(false))
func WithIncludeUnknown(includeUnknown bool) Options {
	return NewFuncOption(func(options *options) {
		options.IncludeUnknown = includeUnknown
	})
}

// WithTransferCasing controls whether to preserve original word casing in suggestions.
//
// When enabled, suggestions adopt the capitalization pattern of the input word.
// This provides more natural-looking corrections that maintain the user's
// intended casing style.
//
// Args:
//
//	transferCasing: Whether to transfer casing from input to suggestions.
//	               true = apply input casing pattern to suggestions
//	               false = return suggestions with dictionary casing
//
// Returns:
//
//	An Options instance that sets the TransferCasing configuration.
//
// Example:
//
//	// Preserve user's capitalization style
//	checker := NewSymSpell(WithTransferCasing(true))
//
//	// Use dictionary's original casing
//	checker := NewSymSpell(WithTransferCasing(false))
func WithTransferCasing(transferCasing bool) Options {
	return NewFuncOption(func(options *options) {
		options.TransferCasing = transferCasing
	})
}

// WithIgnoreToken sets a regex pattern for tokens that should be ignored during spell checking.
//
// Tokens matching the provided regular expression are passed through unchanged
// without any spell checking attempts. This is useful for excluding URLs,
// email addresses, codes, or other structured text that shouldn't be corrected.
//
// Args:
//
//	ignoreToken: A compiled regular expression pattern for tokens to ignore.
//	            Set to nil to disable token filtering.
//	            The pattern is matched against the entire token.
//
// Returns:
//
//	An Options instance that sets the IgnoreToken configuration.
//
// Example:
//
//	// Ignore URLs and email addresses
//	urlPattern := regexp.MustCompile(`^(https?://|[^@]+@[^@]+\.[^@]+)`)
//	checker := NewSymSpell(WithIgnoreToken(urlPattern))
//
//	// Disable token filtering
//	checker := NewSymSpell(WithIgnoreToken(nil))
func WithIgnoreToken(ignoreToken *regexp.Regexp) Options {
	return NewFuncOption(func(options *options) {
		options.IgnoreToken = ignoreToken
	})
}

// WithIgnoreNonWords controls whether to skip tokens that don't look like words.
//
// When enabled, tokens containing only punctuation, numbers, or special
// characters are passed through without spell checking. This improves
// performance and prevents unnecessary processing of non-word content.
//
// Args:
//
//	ignoreNonWords: Whether to skip non-word tokens during spell checking.
//	               true = skip punctuation, numbers, symbols
//	               false = attempt to spell check all tokens
//
// Returns:
//
//	An Options instance that sets the IgnoreNonWords configuration.
//
// Example:
//
//	// Skip obvious non-words for better performance
//	checker := NewSymSpell(WithIgnoreNonWords(true))
//
//	// Check all tokens regardless of content
//	checker := NewSymSpell(WithIgnoreNonWords(false))
func WithIgnoreNonWords(ignoreNonWords bool) Options {
	return NewFuncOption(func(options *options) {
		options.IgnoreNonWords = ignoreNonWords
	})
}

// WithIgnoreTermWithDigits controls whether to skip words containing numeric digits.
//
// When enabled, words that contain any numeric digits are passed through
// without spell checking. This is useful for technical content where
// alphanumeric terms (model numbers, versions, codes) are intentional.
//
// Args:
//
//	ignoreTermWithDigits: Whether to skip digit-containing words.
//	                     true = ignore words with any digits
//	                     false = spell check all words regardless of digits
//
// Returns:
//
//	An Options instance that sets the IgnoreTermWithDigits configuration.
//
// Example:
//
//	// Ignore technical terms with digits
//	checker := NewSymSpell(WithIgnoreTermWithDigits(true))
//
//	// Check all words including those with digits
//	checker := NewSymSpell(WithIgnoreTermWithDigits(false))
func WithIgnoreTermWithDigits(ignoreTermWithDigits bool) Options {
	return NewFuncOption(func(options *options) {
		options.IgnoreTermWithDigits = ignoreTermWithDigits
	})
}

// WithSplitBySpace enables automatic splitting of compound words for separate checking.
//
// When enabled, unknown words are split at space boundaries and each part
// is checked separately. This helps with compound words, concatenated terms,
// or cases where multiple words were accidentally joined together.
//
// Args:
//
//	splitBySpace: Whether to attempt word splitting during spell checking.
//	             true = split unknown words and check parts separately
//	             false = treat each token as a single unit
//
// Returns:
//
//	An Options instance that sets the SplitBySpace configuration.
//
// Example:
//
//	// Enable compound word splitting
//	checker := NewSymSpell(WithSplitBySpace(true))
//
//	// Treat all tokens as indivisible units
//	checker := NewSymSpell(WithSplitBySpace(false))
func WithSplitBySpace(splitBySpace bool) Options {
	return NewFuncOption(func(options *options) {
		options.SplitBySpace = splitBySpace
	})
}
