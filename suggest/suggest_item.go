package suggest

import (
	"fmt"
	"math"
	"sort"
)

// SuggestItem represents a spelling suggestion.
type SuggestItem struct {
	// Suggested word.
	Term string
	// Edit distance from the searched term.
	Distance int
	// Frequency in the dictionary or a Naive Bayes probability score.
	Count int
}

// NewSuggestItem creates a new SuggestItem with the given term, distance, and count.
func NewSuggestItem(term string, distance int, count int) SuggestItem {
	return SuggestItem{
		Term:     term,
		Distance: distance,
		Count:    count,
	}
}

// NewWithProbability creates a SuggestItem using a Naive Bayes–style probability
// as the count, based on the term length.
func NewWithProbability(term string, distance int) SuggestItem {
	count := 10 / int(math.Pow(10, float64(len(term))))
	return SuggestItem{
		Term:     term,
		Distance: distance,
		Count:    count,
	}
}

// String returns the SuggestItem as "term, distance, count".
func (s SuggestItem) String() string {
	return fmt.Sprintf("%s, %d, %d", s.Term, s.Distance, s.Count)
}

// SuggestItems is a slice of SuggestItem.
type SuggestItems []SuggestItem

// Len returns the number of items in the slice.
func (si SuggestItems) Len() int {
	return len(si)
}

// Swap exchanges two items in the slice.
func (si SuggestItems) Swap(i, j int) {
	si[i], si[j] = si[j], si[i]
}

// Less compares two items, ordering by ascending distance,
// then by descending count when distances are equal.
func (si SuggestItems) Less(i, j int) bool {
	if si[i].Distance == si[j].Distance {
		return si[i].Count > si[j].Count
	}
	return si[i].Distance < si[j].Distance
}

// SortSuggestItems sorts a slice of SuggestItem in place.
func SortSuggestItems(items []SuggestItem) {
	sort.Sort(SuggestItems(items))
}

// NewSuggestItemWithProbability is an alternative constructor that
// calculates the count as a probability based on term length.
func NewSuggestItemWithProbability(term string, distance int) SuggestItem {
	count := 10 / int(math.Pow(10, float64(len(term))))
	return SuggestItem{Term: term, Distance: distance, Count: count}
}
