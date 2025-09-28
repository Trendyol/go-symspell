package common

import "testing"

func TestHasSeen(t *testing.T) {
	tests := []struct {
		name     string
		set      map[string]struct{}
		key      string
		expected bool
	}{
		{
			name:     "key exists in set",
			set:      map[string]struct{}{"hello": {}, "world": {}},
			key:      "hello",
			expected: true,
		},
		{
			name:     "key does not exist in set",
			set:      map[string]struct{}{"hello": {}, "world": {}},
			key:      "test",
			expected: false,
		},
		{
			name:     "empty set",
			set:      map[string]struct{}{},
			key:      "test",
			expected: false,
		},
		{
			name:     "nil set",
			set:      nil,
			key:      "test",
			expected: false,
		},
		{
			name:     "empty key in set",
			set:      map[string]struct{}{"": {}, "hello": {}},
			key:      "",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasSeen(tt.set, tt.key)
			if result != tt.expected {
				t.Errorf("HasSeen() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHasSeenInt(t *testing.T) {
	tests := []struct {
		name     string
		set      map[int]struct{}
		key      int
		expected bool
	}{
		{
			name:     "int key exists",
			set:      map[int]struct{}{1: {}, 2: {}, 3: {}},
			key:      2,
			expected: true,
		},
		{
			name:     "int key does not exist",
			set:      map[int]struct{}{1: {}, 2: {}, 3: {}},
			key:      5,
			expected: false,
		},
		{
			name:     "zero key exists",
			set:      map[int]struct{}{0: {}, 1: {}},
			key:      0,
			expected: true,
		},
		{
			name:     "negative key",
			set:      map[int]struct{}{-1: {}, 0: {}, 1: {}},
			key:      -1,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasSeen(tt.set, tt.key)
			if result != tt.expected {
				t.Errorf("HasSeen() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestContainsRune(t *testing.T) {
	tests := []struct {
		name     string
		runes    []rune
		target   rune
		expected bool
	}{
		{
			name:     "rune exists at beginning",
			runes:    []rune{'h', 'e', 'l', 'l', 'o'},
			target:   'h',
			expected: true,
		},
		{
			name:     "rune exists in middle",
			runes:    []rune{'h', 'e', 'l', 'l', 'o'},
			target:   'l',
			expected: true,
		},
		{
			name:     "rune exists at end",
			runes:    []rune{'h', 'e', 'l', 'l', 'o'},
			target:   'o',
			expected: true,
		},
		{
			name:     "rune does not exist",
			runes:    []rune{'h', 'e', 'l', 'l', 'o'},
			target:   'x',
			expected: false,
		},
		{
			name:     "empty slice",
			runes:    []rune{},
			target:   'h',
			expected: false,
		},
		{
			name:     "nil slice",
			runes:    nil,
			target:   'h',
			expected: false,
		},
		{
			name:     "unicode rune exists",
			runes:    []rune{'h', 'é', 'l', 'l', 'ö'},
			target:   'é',
			expected: true,
		},
		{
			name:     "unicode rune does not exist",
			runes:    []rune{'h', 'e', 'l', 'l', 'o'},
			target:   'é',
			expected: false,
		},
		{
			name:     "emoji rune exists",
			runes:    []rune{'🙂', 'h', 'e', 'l', 'l', 'o'},
			target:   '🙂',
			expected: true,
		},
		{
			name:     "single rune slice match",
			runes:    []rune{'a'},
			target:   'a',
			expected: true,
		},
		{
			name:     "single rune slice no match",
			runes:    []rune{'a'},
			target:   'b',
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainsRune(tt.runes, tt.target)
			if result != tt.expected {
				t.Errorf("ContainsRune() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRunesEqual(t *testing.T) {
	tests := []struct {
		name     string
		a        []rune
		b        []rune
		expected bool
	}{
		{
			name:     "equal rune slices",
			a:        []rune{'h', 'e', 'l', 'l', 'o'},
			b:        []rune{'h', 'e', 'l', 'l', 'o'},
			expected: true,
		},
		{
			name:     "different rune slices",
			a:        []rune{'h', 'e', 'l', 'l', 'o'},
			b:        []rune{'w', 'o', 'r', 'l', 'd'},
			expected: false,
		},
		{
			name:     "different lengths",
			a:        []rune{'h', 'e', 'l', 'l', 'o'},
			b:        []rune{'h', 'e', 'l', 'l'},
			expected: false,
		},
		{
			name:     "both empty",
			a:        []rune{},
			b:        []rune{},
			expected: true,
		},
		{
			name:     "both nil",
			a:        nil,
			b:        nil,
			expected: true,
		},
		{
			name:     "one empty one nil",
			a:        []rune{},
			b:        nil,
			expected: true,
		},
		{
			name:     "one empty one with content",
			a:        []rune{},
			b:        []rune{'a'},
			expected: false,
		},
		{
			name:     "one nil one with content",
			a:        nil,
			b:        []rune{'a'},
			expected: false,
		},
		{
			name:     "single rune equal",
			a:        []rune{'a'},
			b:        []rune{'a'},
			expected: true,
		},
		{
			name:     "single rune different",
			a:        []rune{'a'},
			b:        []rune{'b'},
			expected: false,
		},
		{
			name:     "unicode runes equal",
			a:        []rune{'h', 'é', 'l', 'l', 'ö'},
			b:        []rune{'h', 'é', 'l', 'l', 'ö'},
			expected: true,
		},
		{
			name:     "unicode runes different",
			a:        []rune{'h', 'é', 'l', 'l', 'ö'},
			b:        []rune{'h', 'e', 'l', 'l', 'o'},
			expected: false,
		},
		{
			name:     "emoji runes equal",
			a:        []rune{'🙂', '👍'},
			b:        []rune{'🙂', '👍'},
			expected: true,
		},
		{
			name:     "same content different order",
			a:        []rune{'a', 'b', 'c'},
			b:        []rune{'c', 'b', 'a'},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RunesEqual(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("RunesEqual() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRuneIndexEqual(t *testing.T) {
	tests := []struct {
		name     string
		a        []rune
		i        int
		b        []rune
		j        int
		expected bool
	}{
		{
			name:     "equal runes at valid indices",
			a:        []rune{'h', 'e', 'l', 'l', 'o'},
			i:        0,
			b:        []rune{'h', 'a', 'l', 'l', 'o'},
			j:        0,
			expected: true,
		},
		{
			name:     "different runes at valid indices",
			a:        []rune{'h', 'e', 'l', 'l', 'o'},
			i:        0,
			b:        []rune{'w', 'a', 'l', 'l', 'o'},
			j:        0,
			expected: false,
		},
		{
			name:     "negative index i",
			a:        []rune{'h', 'e', 'l', 'l', 'o'},
			i:        -1,
			b:        []rune{'h', 'a', 'l', 'l', 'o'},
			j:        0,
			expected: false,
		},
		{
			name:     "negative index j",
			a:        []rune{'h', 'e', 'l', 'l', 'o'},
			i:        0,
			b:        []rune{'h', 'a', 'l', 'l', 'o'},
			j:        -1,
			expected: false,
		},
		{
			name:     "index i out of bounds",
			a:        []rune{'h', 'e', 'l', 'l', 'o'},
			i:        5,
			b:        []rune{'h', 'a', 'l', 'l', 'o'},
			j:        0,
			expected: false,
		},
		{
			name:     "index j out of bounds",
			a:        []rune{'h', 'e', 'l', 'l', 'o'},
			i:        0,
			b:        []rune{'h', 'a', 'l', 'l', 'o'},
			j:        5,
			expected: false,
		},
		{
			name:     "both indices out of bounds",
			a:        []rune{'h', 'e', 'l', 'l', 'o'},
			i:        10,
			b:        []rune{'h', 'a', 'l', 'l', 'o'},
			j:        10,
			expected: false,
		},
		{
			name:     "empty slice a",
			a:        []rune{},
			i:        0,
			b:        []rune{'h', 'a', 'l', 'l', 'o'},
			j:        0,
			expected: false,
		},
		{
			name:     "empty slice b",
			a:        []rune{'h', 'e', 'l', 'l', 'o'},
			i:        0,
			b:        []rune{},
			j:        0,
			expected: false,
		},
		{
			name:     "both slices empty",
			a:        []rune{},
			i:        0,
			b:        []rune{},
			j:        0,
			expected: false,
		},
		{
			name:     "nil slice a",
			a:        nil,
			i:        0,
			b:        []rune{'h', 'a', 'l', 'l', 'o'},
			j:        0,
			expected: false,
		},
		{
			name:     "nil slice b",
			a:        []rune{'h', 'e', 'l', 'l', 'o'},
			i:        0,
			b:        nil,
			j:        0,
			expected: false,
		},
		{
			name:     "unicode runes equal",
			a:        []rune{'h', 'é', 'l', 'l', 'ö'},
			i:        1,
			b:        []rune{'w', 'é', 'l', 'l', 'ö'},
			j:        1,
			expected: true,
		},
		{
			name:     "emoji runes equal",
			a:        []rune{'🙂', '👍'},
			i:        0,
			b:        []rune{'🙂', '😀'},
			j:        0,
			expected: true,
		},
		{
			name:     "boundary valid indices",
			a:        []rune{'a', 'b', 'c'},
			i:        2,
			b:        []rune{'x', 'y', 'c'},
			j:        2,
			expected: true,
		},
		{
			name:     "zero indices equal",
			a:        []rune{'a'},
			i:        0,
			b:        []rune{'a'},
			j:        0,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RuneIndexEqual(tt.a, tt.i, tt.b, tt.j)
			if result != tt.expected {
				t.Errorf("RuneIndexEqual() = %v, want %v", result, tt.expected)
			}
		})
	}
}
