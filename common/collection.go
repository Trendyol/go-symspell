package common

// HasSeen checks whether k is in set.
func HasSeen[T comparable](set map[T]struct{}, k T) bool {
	_, ok := set[k]
	return ok
}

// ContainsRune checks if r exists in slice (avoids strings.ContainsRune + alloc).
func ContainsRune(rs []rune, r rune) bool {
	for _, v := range rs {
		if v == r {
			return true
		}
	}
	return false
}

// RunesEqual compares two rune slices for equality without allocations.
func RunesEqual(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// RuneIndexEqual compares a[i] == b[j], guarding bounds.
func RuneIndexEqual(a []rune, i int, b []rune, j int) bool {
	if i < 0 || j < 0 || i >= len(a) || j >= len(b) {
		return false
	}
	return a[i] == b[j]
}
