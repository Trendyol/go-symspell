package common

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// SplitN splits s by sep into at most n parts (like strings.SplitN),
// but returns nil if sep is empty to avoid undefined behavior.
func SplitN(s, sep string, n int) []string {
	if sep == "" {
		// Fall back to whole string if separator is invalid.
		return []string{s}
	}
	return strings.SplitN(s, sep, n)
}

// IsAcronym checks if a word is an acronym (all uppercase).
func IsAcronym(word string, ignoreTermWithDigits bool) bool {
	if ignoreTermWithDigits {
		for _, r := range word {
			if unicode.IsDigit(r) {
				return true
			}
		}
	}
	hasLetter := false
	for _, r := range word {
		if !unicode.IsUpper(r) && !unicode.IsDigit(r) {
			return false
		}
		if unicode.IsLetter(r) {
			hasLetter = true
		}
	}
	return hasLetter
}

// CaseTransferSimilar transfers the case of a string from one to another.
func CaseTransferSimilar(source, target string) string {
	if len(source) == 0 {
		return ""
	}
	if len(target) == 0 {
		return target
	}
	isSourceTitle := true
	isSourceUpper := true
	for i, r := range source {
		if !unicode.IsUpper(r) {
			isSourceUpper = false
		}
		if i == 0 && !unicode.IsUpper(r) {
			isSourceTitle = false
		}
		if !isSourceUpper && !isSourceTitle {
			break
		}
	}

	if isSourceUpper {
		return strings.ToUpper(target)
	}
	if isSourceTitle {
		r, size := utf8.DecodeRuneInString(target)
		return string(unicode.ToUpper(r)) + target[size:]
	}
	return target
}
