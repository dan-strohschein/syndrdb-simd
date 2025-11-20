package syndrdbsimd

import (
	"fmt"
	"sync/atomic"
)

// PatternType represents the type of SQL LIKE pattern for optimized matching.
type PatternType int

const (
	// PatternExact matches exact strings (no wildcards)
	// Example: "hello" matches only "hello"
	PatternExact PatternType = iota

	// PatternPrefix matches strings starting with a given prefix
	// Example: "hello%" matches "hello", "hello world", etc.
	PatternPrefix

	// PatternSuffix matches strings ending with a given suffix
	// Example: "%world" matches "world", "hello world", etc.
	PatternSuffix

	// PatternContains matches strings containing a substring
	// Example: "%test%" matches "testing", "contest", "test", etc.
	PatternContains

	// PatternWildcard matches complex patterns with % (any chars) and _ (single char)
	// Example: "h_llo%" matches "hello world", "hallo there", etc.
	PatternWildcard
)

// String returns a human-readable representation of the PatternType.
func (pt PatternType) String() string {
	switch pt {
	case PatternExact:
		return "Exact"
	case PatternPrefix:
		return "Prefix"
	case PatternSuffix:
		return "Suffix"
	case PatternContains:
		return "Contains"
	case PatternWildcard:
		return "Wildcard"
	default:
		return fmt.Sprintf("Unknown(%d)", int(pt))
	}
}

// thresholdConfig holds the SIMD threshold configuration for string operations.
type thresholdConfig struct {
	minStrings       int // Minimum number of strings to use SIMD
	avgByteThreshold int // Average string length threshold for adaptive SIMD
}

// Global atomic pointer to threshold configuration
var stringSIMDThreshold atomic.Pointer[thresholdConfig]

func init() {
	// Set default thresholds
	defaultConfig := &thresholdConfig{
		minStrings:       16,
		avgByteThreshold: 32,
	}
	stringSIMDThreshold.Store(defaultConfig)
}

// SetStringSIMDThreshold configures the adaptive SIMD thresholds for string operations.
//
// Parameters:
//   - minStrings: Minimum number of strings required to use SIMD (must be > 0)
//   - avgByteThreshold: Average string length threshold for SIMD selection (must be > 0)
//
// The function uses these parameters to decide when to use SIMD:
//   - If len(strings) < minStrings, use scalar operations
//   - If average string length > avgByteThreshold, may use lower minStrings threshold
//   - Default values: minStrings=16, avgByteThreshold=32
//
// This function is thread-safe and can be called at runtime, though it's recommended
// to set thresholds once at startup for best performance.
//
// Panics if minStrings <= 0 or avgByteThreshold <= 0.
//
// TODO: Return error instead of panic for library-friendly API in future.
func SetStringSIMDThreshold(minStrings, avgByteThreshold int) {
	if minStrings <= 0 {
		panic(fmt.Sprintf("minStrings must be > 0, got %d", minStrings))
	}
	if avgByteThreshold <= 0 {
		panic(fmt.Sprintf("avgByteThreshold must be > 0, got %d", avgByteThreshold))
	}

	config := &thresholdConfig{
		minStrings:       minStrings,
		avgByteThreshold: avgByteThreshold,
	}
	stringSIMDThreshold.Store(config)
}

// GetStringSIMDThreshold returns the current SIMD threshold configuration.
// Returns (minStrings, avgByteThreshold).
func GetStringSIMDThreshold() (int, int) {
	config := stringSIMDThreshold.Load()
	return config.minStrings, config.avgByteThreshold
}
