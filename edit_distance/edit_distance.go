package edit_distance

import (
	"errors"
)

// Supported algorithms
const (
	DamerauOSAFast AlgorithmType = iota
	DamerauOSA
	DamerauLevenshtein
	Levenshtein
)

// DistanceComparer compares two strings with a chosen edit-distance algorithm.
type DistanceComparer interface {
	// Compare returns the distance between a and b.
	// If maxDistance >= 0 and the distance exceeds it, returns -1 with no error.
	Compare(a, b string, maxDistance int) (int, error)
}

// distanceComparer holds the selected algorithm.
type distanceComparer struct {
	algorithm Algorithm
}

// NewDistanceComparer creates a comparer for the given algorithm.
func NewDistanceComparer(algorithmType AlgorithmType) (DistanceComparer, error) {
	var algorithm Algorithm
	switch algorithmType {
	case DamerauOSAFast:
		algorithm = &DamerauOSAFastAlgorithm{}
	case DamerauOSA:
		algorithm = &DamerauOSAAlgorithm{}
	case DamerauLevenshtein:
		algorithm = &DamerauLevenshteinAlgorithm{}
	case Levenshtein:
		algorithm = &LevenshteinAlgorithm{}
	default:
		return nil, errors.New("invalid algorithm")
	}

	return &distanceComparer{algorithm: algorithm}, nil
}

// Compare computes the edit distance using the configured algorithm.
func (d *distanceComparer) Compare(a, b string, maxDistance int) (int, error) {
	return d.algorithm.Compare(a, b, maxDistance)
}
