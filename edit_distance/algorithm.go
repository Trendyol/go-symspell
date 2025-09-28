package edit_distance

import (
	"github.com/hbollon/go-edlib"
)

const (
	kSpace = rune(' ')
)

// AlgorithmType can be used to represent an edit distance algorithm choice.
type AlgorithmType int

// Algorithm defines the interface for edit distance comparison.
type Algorithm interface {
	// Compare calculates the edit distance between strings a and b.
	Compare(a, b string, maxDistance int) (int, error)
}

// DamerauOSAFastAlgorithm is a dummy(!) implementation of the DamerauOSA fast algorithm.
type DamerauOSAFastAlgorithm struct {
}

func (d *DamerauOSAFastAlgorithm) Compare(a, b string, maxDistance int) (int, error) {
	if maxDistance <= 0 {
		if a == b {
			return 0, nil
		}
		return -1, nil
	}

	// Convert to runes once (UTF-8 safe).
	r1 := []rune(a)
	r2 := []rune(b)
	n1 := len(r1)
	n2 := len(r2)

	// If length difference already exceeds maxDistance, prune early.
	if diff := n2 - n1; diff > maxDistance || diff < -maxDistance {
		return -1, nil
	}

	// Ensure r1 is the shorter (or equal) one.
	if n1 > n2 {
		r1, r2 = r2, r1
		n1, n2 = n2, n1
	}

	// Trim common suffix
	for n1 != 0 && r1[n1-1] == r2[n2-1] {
		n1--
		n2--
	}
	if n1 == 0 {
		if n2 <= maxDistance {
			return n2, nil
		}
		return -1, nil
	}

	// Trim common prefix
	start := 0
	for start != n1 && r1[start] == r2[start] {
		start++
	}
	if start != 0 {
		n1 -= start
		n2 -= start
	}
	if n1 == 0 {
		if n2 <= maxDistance {
			return n2, nil
		}
		return -1, nil
	}

	// Choose bounded/unbounded DP based on maxDistance
	if maxDistance < n2 {
		return d.internalDistanceMax(r1, r2, n1, n2, start, maxDistance)
	}
	return d.internalDistance(r1, r2, n1, n2, start)
}

func (d *DamerauOSAFastAlgorithm) internalDistance(r1, r2 []rune, len1, len2, start int) (int, error) {
	prevCharCosts := make([]int, len2) // zero-initialized
	charCosts := make([]int, len2)
	for j := 0; j < len2; j++ {
		charCosts[j] = j + 1
	}

	char1 := kSpace
	current := 0

	for i := 0; i < len1; i++ {
		prevChar1 := char1
		char1 = r1[start+i]

		char2 := kSpace
		above := i
		left := i
		nextTrans := 0

		for j := 0; j < len2; j++ {
			thisTrans := nextTrans
			nextTrans = prevCharCosts[j]

			// Diagonal (substitution) base
			current = left
			prevCharCosts[j] = left

			// Left now equals previous charCosts (diagonal for next iter)
			left = charCosts[j]

			prevChar2 := char2
			char2 = r2[start+j]

			if char1 != char2 {
				// Substitution if neither insertion nor deletion cheaper
				if above < current {
					current = above
				}
				if left < current {
					current = left
				}
				current++

				// Transposition
				if i != 0 && j != 0 && char1 == prevChar2 && prevChar1 == char2 && thisTrans+1 < current {
					current = thisTrans + 1
				}
			}

			above = current
			charCosts[j] = current
		}
	}
	return current, nil
}

func (d *DamerauOSAFastAlgorithm) internalDistanceMax(r1, r2 []rune, len1, len2, start, maxDistance int) (int, error) {
	prevCharCosts := make([]int, len2) // zero-initialized
	charCosts := make([]int, len2)
	for j := 0; j < len2; j++ {
		charCosts[j] = j + 1
	}

	lenDiff := len2 - len1
	jStartOffset := maxDistance - lenDiff
	jStart := 0
	jEnd := maxDistance
	if jEnd > len2 {
		jEnd = len2
	}

	char1 := kSpace
	current := 0

	for i := 0; i < len1; i++ {
		prevChar1 := char1
		char1 = r1[start+i]

		char2 := kSpace
		above := i
		left := i
		nextTrans := 0

		// Slide window: lower-right diagonal is (i - lenDiff), upper-left is i.
		if i > jStartOffset {
			jStart++
		}
		if jEnd < len2 {
			jEnd++
		}

		// Clamp (safety)
		if jStart < 0 {
			jStart = 0
		}
		if jEnd > len2 {
			jEnd = len2
		}

		for j := jStart; j < jEnd; j++ {
			thisTrans := nextTrans
			nextTrans = prevCharCosts[j]

			// Diagonal (substitution) base
			current = left
			prevCharCosts[j] = left

			// Left becomes previous charCosts (diagonal for next)
			left = charCosts[j]

			prevChar2 := char2
			char2 = r2[start+j]

			if char1 != char2 {
				// Substitution if neither insertion nor deletion cheaper
				if above < current {
					current = above
				}
				if left < current {
					current = left
				}
				current++

				// Transposition
				if i != 0 && j != 0 && char1 == prevChar2 && prevChar1 == char2 && thisTrans+1 < current {
					current = thisTrans + 1
				}
			}

			above = current
			charCosts[j] = current
		}

		// Early-exit: the cell on the (i+lenDiff)-th column exceeds maxDistance
		idx := i + lenDiff
		if idx >= 0 && idx < len2 && charCosts[idx] > maxDistance {
			return -1, nil
		}
	}

	if current <= maxDistance {
		return current, nil
	}
	return -1, nil
}

type DamerauOSAAlgorithm struct {
}

func (d *DamerauOSAAlgorithm) Compare(a, b string, maxDistance int) (int, error) {
	return edlib.OSADamerauLevenshteinDistance(a, b), nil
}

type DamerauLevenshteinAlgorithm struct {
}

func (d *DamerauLevenshteinAlgorithm) Compare(a, b string, maxDistance int) (int, error) {
	return edlib.DamerauLevenshteinDistance(a, b), nil
}

type LevenshteinAlgorithm struct {
}

func (l *LevenshteinAlgorithm) Compare(a, b string, maxDistance int) (int, error) {
	return edlib.LevenshteinDistance(a, b), nil
}
