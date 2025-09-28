package suggest

import (
	"testing"
)

func TestNewSuggestItem(t *testing.T) {
	tests := []struct {
		name     string
		term     string
		distance int
		count    int
		expected SuggestItem
	}{
		{
			name:     "basic creation",
			term:     "hello",
			distance: 1,
			count:    100,
			expected: SuggestItem{Term: "hello", Distance: 1, Count: 100},
		},
		{
			name:     "empty term",
			term:     "",
			distance: 0,
			count:    0,
			expected: SuggestItem{Term: "", Distance: 0, Count: 0},
		},
		{
			name:     "negative values",
			term:     "test",
			distance: -1,
			count:    -5,
			expected: SuggestItem{Term: "test", Distance: -1, Count: -5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewSuggestItem(tt.term, tt.distance, tt.count)
			if result != tt.expected {
				t.Errorf("NewSuggestItem() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewWithProbability(t *testing.T) {
	tests := []struct {
		name     string
		term     string
		distance int
		expected SuggestItem
	}{
		{
			name:     "single character",
			term:     "a",
			distance: 1,
			expected: SuggestItem{Term: "a", Distance: 1, Count: 1},
		},
		{
			name:     "two characters",
			term:     "ab",
			distance: 2,
			expected: SuggestItem{Term: "ab", Distance: 2, Count: 0},
		},
		{
			name:     "three characters",
			term:     "abc",
			distance: 0,
			expected: SuggestItem{Term: "abc", Distance: 0, Count: 0},
		},
		{
			name:     "empty string",
			term:     "",
			distance: 0,
			expected: SuggestItem{Term: "", Distance: 0, Count: 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewWithProbability(tt.term, tt.distance)
			if result != tt.expected {
				t.Errorf("NewWithProbability() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewSuggestItemWithProbability(t *testing.T) {
	tests := []struct {
		name     string
		term     string
		distance int
		expected SuggestItem
	}{
		{
			name:     "single character",
			term:     "x",
			distance: 1,
			expected: SuggestItem{Term: "x", Distance: 1, Count: 1},
		},
		{
			name:     "two characters",
			term:     "xy",
			distance: 2,
			expected: SuggestItem{Term: "xy", Distance: 2, Count: 0},
		},
		{
			name:     "empty string",
			term:     "",
			distance: 0,
			expected: SuggestItem{Term: "", Distance: 0, Count: 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewSuggestItemWithProbability(tt.term, tt.distance)
			if result != tt.expected {
				t.Errorf("NewSuggestItemWithProbability() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSuggestItem_String(t *testing.T) {
	tests := []struct {
		name     string
		item     SuggestItem
		expected string
	}{
		{
			name:     "normal item",
			item:     SuggestItem{Term: "hello", Distance: 1, Count: 100},
			expected: "hello, 1, 100",
		},
		{
			name:     "empty term",
			item:     SuggestItem{Term: "", Distance: 0, Count: 0},
			expected: ", 0, 0",
		},
		{
			name:     "negative values",
			item:     SuggestItem{Term: "test", Distance: -1, Count: -5},
			expected: "test, -1, -5",
		},
		{
			name:     "special characters",
			item:     SuggestItem{Term: "test,item", Distance: 2, Count: 50},
			expected: "test,item, 2, 50",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.item.String()
			if result != tt.expected {
				t.Errorf("String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSuggestItems_Len(t *testing.T) {
	tests := []struct {
		name     string
		items    SuggestItems
		expected int
	}{
		{
			name:     "empty slice",
			items:    SuggestItems{},
			expected: 0,
		},
		{
			name:     "single item",
			items:    SuggestItems{{Term: "test", Distance: 1, Count: 10}},
			expected: 1,
		},
		{
			name: "multiple items",
			items: SuggestItems{
				{Term: "test1", Distance: 1, Count: 10},
				{Term: "test2", Distance: 2, Count: 20},
				{Term: "test3", Distance: 3, Count: 30},
			},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.items.Len()
			if result != tt.expected {
				t.Errorf("Len() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestSuggestItems_Swap(t *testing.T) {
	tests := []struct {
		name     string
		items    SuggestItems
		i        int
		j        int
		expected SuggestItems
	}{
		{
			name: "swap first and second",
			items: SuggestItems{
				{Term: "first", Distance: 1, Count: 10},
				{Term: "second", Distance: 2, Count: 20},
			},
			i: 0,
			j: 1,
			expected: SuggestItems{
				{Term: "second", Distance: 2, Count: 20},
				{Term: "first", Distance: 1, Count: 10},
			},
		},
		{
			name: "swap same index",
			items: SuggestItems{
				{Term: "same", Distance: 1, Count: 10},
			},
			i: 0,
			j: 0,
			expected: SuggestItems{
				{Term: "same", Distance: 1, Count: 10},
			},
		},
		{
			name: "swap middle items",
			items: SuggestItems{
				{Term: "first", Distance: 1, Count: 10},
				{Term: "second", Distance: 2, Count: 20},
				{Term: "third", Distance: 3, Count: 30},
			},
			i: 0,
			j: 2,
			expected: SuggestItems{
				{Term: "third", Distance: 3, Count: 30},
				{Term: "second", Distance: 2, Count: 20},
				{Term: "first", Distance: 1, Count: 10},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items := make(SuggestItems, len(tt.items))
			copy(items, tt.items)
			items.Swap(tt.i, tt.j)

			for k := 0; k < len(items); k++ {
				if items[k] != tt.expected[k] {
					t.Errorf("After Swap(%d, %d), items[%d] = %v, want %v", tt.i, tt.j, k, items[k], tt.expected[k])
				}
			}
		})
	}
}

func TestSuggestItems_Less(t *testing.T) {
	tests := []struct {
		name     string
		items    SuggestItems
		i        int
		j        int
		expected bool
	}{
		{
			name: "different distances - first less",
			items: SuggestItems{
				{Term: "first", Distance: 1, Count: 10},
				{Term: "second", Distance: 2, Count: 20},
			},
			i:        0,
			j:        1,
			expected: true,
		},
		{
			name: "different distances - first greater",
			items: SuggestItems{
				{Term: "first", Distance: 3, Count: 10},
				{Term: "second", Distance: 1, Count: 20},
			},
			i:        0,
			j:        1,
			expected: false,
		},
		{
			name: "same distance - first higher count",
			items: SuggestItems{
				{Term: "first", Distance: 1, Count: 30},
				{Term: "second", Distance: 1, Count: 20},
			},
			i:        0,
			j:        1,
			expected: true,
		},
		{
			name: "same distance - first lower count",
			items: SuggestItems{
				{Term: "first", Distance: 1, Count: 10},
				{Term: "second", Distance: 1, Count: 20},
			},
			i:        0,
			j:        1,
			expected: false,
		},
		{
			name: "same distance and count",
			items: SuggestItems{
				{Term: "first", Distance: 1, Count: 20},
				{Term: "second", Distance: 1, Count: 20},
			},
			i:        0,
			j:        1,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.items.Less(tt.i, tt.j)
			if result != tt.expected {
				t.Errorf("Less(%d, %d) = %t, want %t", tt.i, tt.j, result, tt.expected)
			}
		})
	}
}

func TestSortSuggestItems(t *testing.T) {
	tests := []struct {
		name     string
		items    []SuggestItem
		expected []SuggestItem
	}{
		{
			name:     "empty slice",
			items:    []SuggestItem{},
			expected: []SuggestItem{},
		},
		{
			name: "single item",
			items: []SuggestItem{
				{Term: "test", Distance: 1, Count: 10},
			},
			expected: []SuggestItem{
				{Term: "test", Distance: 1, Count: 10},
			},
		},
		{
			name: "sort by distance ascending",
			items: []SuggestItem{
				{Term: "third", Distance: 3, Count: 10},
				{Term: "first", Distance: 1, Count: 10},
				{Term: "second", Distance: 2, Count: 10},
			},
			expected: []SuggestItem{
				{Term: "first", Distance: 1, Count: 10},
				{Term: "second", Distance: 2, Count: 10},
				{Term: "third", Distance: 3, Count: 10},
			},
		},
		{
			name: "sort by count descending when distances equal",
			items: []SuggestItem{
				{Term: "low", Distance: 1, Count: 10},
				{Term: "high", Distance: 1, Count: 30},
				{Term: "medium", Distance: 1, Count: 20},
			},
			expected: []SuggestItem{
				{Term: "high", Distance: 1, Count: 30},
				{Term: "medium", Distance: 1, Count: 20},
				{Term: "low", Distance: 1, Count: 10},
			},
		},
		{
			name: "complex sorting",
			items: []SuggestItem{
				{Term: "d2c10", Distance: 2, Count: 10},
				{Term: "d1c20", Distance: 1, Count: 20},
				{Term: "d1c30", Distance: 1, Count: 30},
				{Term: "d3c50", Distance: 3, Count: 50},
				{Term: "d2c40", Distance: 2, Count: 40},
			},
			expected: []SuggestItem{
				{Term: "d1c30", Distance: 1, Count: 30},
				{Term: "d1c20", Distance: 1, Count: 20},
				{Term: "d2c40", Distance: 2, Count: 40},
				{Term: "d2c10", Distance: 2, Count: 10},
				{Term: "d3c50", Distance: 3, Count: 50},
			},
		},
		{
			name: "already sorted",
			items: []SuggestItem{
				{Term: "first", Distance: 1, Count: 30},
				{Term: "second", Distance: 1, Count: 20},
				{Term: "third", Distance: 2, Count: 10},
			},
			expected: []SuggestItem{
				{Term: "first", Distance: 1, Count: 30},
				{Term: "second", Distance: 1, Count: 20},
				{Term: "third", Distance: 2, Count: 10},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items := make([]SuggestItem, len(tt.items))
			copy(items, tt.items)
			SortSuggestItems(items)

			if len(items) != len(tt.expected) {
				t.Errorf("SortSuggestItems() result length = %d, want %d", len(items), len(tt.expected))
				return
			}

			for i := 0; i < len(items); i++ {
				if items[i] != tt.expected[i] {
					t.Errorf("SortSuggestItems() items[%d] = %v, want %v", i, items[i], tt.expected[i])
				}
			}
		})
	}
}

func TestSuggestItemsInterface(t *testing.T) {
	items := SuggestItems{
		{Term: "third", Distance: 3, Count: 10},
		{Term: "first", Distance: 1, Count: 30},
		{Term: "second", Distance: 1, Count: 20},
	}

	if items.Len() != 3 {
		t.Errorf("Expected length 3, got %d", items.Len())
	}

	if !items.Less(1, 2) {
		t.Error("Expected first item (index 1) to be less than second item (index 2)")
	}

	items.Swap(0, 1)
	expected := SuggestItems{
		{Term: "first", Distance: 1, Count: 30},
		{Term: "third", Distance: 3, Count: 10},
		{Term: "second", Distance: 1, Count: 20},
	}

	if items[0] != expected[0] || items[1] != expected[1] {
		t.Errorf("Swap failed: got %v, want %v", items, expected)
	}
}
