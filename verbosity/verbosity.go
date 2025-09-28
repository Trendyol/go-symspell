package verbosity

// Verbosity controls the closeness/quantity of returned spelling suggestions.
type Verbosity int

const (
	// Top returns the single best suggestion:
	// the one with the smallest edit distance, and among those,
	// the highest term frequency.
	Top Verbosity = iota
	// Closest returns all suggestions with the smallest edit distance,
	// sorted by term frequency.
	Closest
	// All returns all suggestions within maxEditDistance, sorted first by
	// edit distance and then by term frequency (slower, no early termination).
	All
)
