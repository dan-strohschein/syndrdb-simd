//go:build arm64

package syndrdbsimd

// cmpGtFloat64NEON compares 2 float64 values against a threshold using NEON.
// Returns a bitmask where bit i is set if values[i] > threshold.
// NEON processes 2 float64 values at a time (128-bit registers).
// NaN comparisons return false per IEEE 754.
//
//go:noescape
func cmpGtFloat64NEON(values *float64, threshold float64) uint64

// cmpGeFloat64NEON compares 2 float64 values against a threshold using NEON.
// Returns a bitmask where bit i is set if values[i] >= threshold.
// NaN comparisons return false per IEEE 754.
//
//go:noescape
func cmpGeFloat64NEON(values *float64, threshold float64) uint64

// cmpLtFloat64NEON compares 2 float64 values against a threshold using NEON.
// Returns a bitmask where bit i is set if values[i] < threshold.
// NaN comparisons return false per IEEE 754.
//
//go:noescape
func cmpLtFloat64NEON(values *float64, threshold float64) uint64

// cmpLeFloat64NEON compares 2 float64 values against a threshold using NEON.
// Returns a bitmask where bit i is set if values[i] <= threshold.
// NaN comparisons return false per IEEE 754.
//
//go:noescape
func cmpLeFloat64NEON(values *float64, threshold float64) uint64

// cmpEqFloat64NEON compares 2 float64 values for equality against a threshold using NEON.
// Returns a bitmask where bit i is set if values[i] == threshold.
// NaN comparisons return false per IEEE 754.
//
//go:noescape
func cmpEqFloat64NEON(values *float64, threshold float64) uint64

// cmpNeFloat64NEON compares 2 float64 values for inequality against a threshold using NEON.
// Returns a bitmask where bit i is set if values[i] != threshold.
// NaN != x returns true for all x per IEEE 754.
//
//go:noescape
func cmpNeFloat64NEON(values *float64, threshold float64) uint64
