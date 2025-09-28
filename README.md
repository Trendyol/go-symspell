# Go SymSpell

[![Go Version](https://img.shields.io/badge/go-1.18+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/Trendyol/go-symspell)](https://goreportcard.com/report/github.com/Trendyol/go-symspell)
[![GoDoc](https://godoc.org/github.com/Trendyol/go-symspell?status.svg)](https://godoc.org/github.com/Trendyol/go-symspell)

## Overview

Go SymSpell is a fast and efficient spell-checking and correction library for Go. It implements the SymSpell algorithm with the “symmetric delete” approach, enabling both speed and accuracy. Unlike traditional spell checkers that generate variations of the input word, SymSpell precomputes all possible deletions of dictionary words up to a given edit distance. This allows very quick lookups while keeping correction quality high.

## Installation

```bash
go mod init your-project
go get github.com/Trendyol/go-symspell
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "github.com/Trendyol/go-symspell/symspell"
    "github.com/Trendyol/go-symspell/verbosity"
)

func main() {
    // Create a new SymSpell instance with default settings
    ss, err := symspell.NewSymSpell()
    if err != nil {
        log.Fatal("Failed to create spell checker:", err)
    }

    // Load dictionary (word frequency_count format)
    _, err = ss.LoadDictionary("dictionary.txt", 0, 1, " ")
    if err != nil {
        log.Fatal("Failed to load dictionary:", err)
    }

    // Get spelling suggestions
    suggestions, err := ss.Lookup("speling", verbosity.Top, 2)
    if err != nil {
        log.Fatal("Lookup failed:", err)
    }

    // Print results
    for _, suggestion := range suggestions {
        fmt.Printf("Suggestion: %s (Distance: %d, Frequency: %d)\n", 
            suggestion.Term, suggestion.Distance, suggestion.Count)
    }
}
```

## Configuration Options

Create a SymSpell instance with custom configuration:

```go
ss, err := symspell.NewSymSpell(
    symspell.WithMaxDictionaryEditDistance(2),  // Maximum edit distance for dictionary
    symspell.WithPrefixLength(7),               // Prefix length for optimization
    symspell.WithIncludeUnknown(true),          // Include unknown words in results
    symspell.WithTransferCasing(true),          // Transfer original casing
    symspell.WithIgnoreNonWords(true),          // Skip non-word tokens
    symspell.WithSplitBySpace(true),            // Enable compound word splitting
)
```

### Available Options

| Option | Default | Description |
|--------|---------|-------------|
| `MaxDictionaryEditDistance` | 2 | Maximum edit distance for precomputed deletes |
| `PrefixLength` | 7 | Length of word prefixes for optimization |
| `InitialCapacity` | 16 | Initial dictionary capacity |
| `CountThreshold` | 1 | Minimum frequency threshold for words |
| `DistanceAlgorithm` | DamerauOSAFast | Edit distance algorithm |
| `IncludeUnknown` | false | Include input word even if not in dictionary |
| `TransferCasing` | false | Apply original casing to suggestions |
| `IgnoreNonWords` | false | Skip tokens that aren't words |
| `IgnoreTermWithDigits` | false | Skip words containing digits |
| `SplitBySpace` | false | Split compound words automatically |

## Dictionary Format

Dictionary files should contain words with their frequencies:

```
the 1061396
of 593677
to 416629
and 411764
a 409757
```

## Verbosity Levels

Control the detail level of suggestions:

```go
import "github.com/Trendyol/go-symspell/verbosity"

// Top suggestion only
suggestions, _ := ss.Lookup("word", verbosity.Top, 2)

// Closest matches within edit distance
suggestions, _ := ss.Lookup("word", verbosity.Closest, 2)

// All suggestions within edit distance
suggestions, _ := ss.Lookup("word", verbosity.All, 2)
```