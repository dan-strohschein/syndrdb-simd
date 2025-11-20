//go:build amd64

package syndrdbsimd

// cmpGtFloat64AVX2 compares 4 float64 values against a threshold using AVX2.
// Returns a bitmask where bit i is set if values[i] > threshold.
// This function processes exactly 4 float64 values (32 bytes).
// NaN comparisons return false per IEEE 754.
//
//go:noescape
func cmpGtFloat64AVX2(values *float64, threshold float64) uint64

// cmpGeFloat64AVX2 compares 4 float64 values against a threshold using AVX2.
// Returns a bitmask where bit i is set if values[i] >= threshold.
// NaN comparisons return false per IEEE 754.
//
//go:noescape
func cmpGeFloat64AVX2(values *float64, threshold float64) uint64

// cmpLtFloat64AVX2 compares 4 float64 values against a threshold using AVX2.
// Returns a bitmask where bit i is set if values[i] < threshold.
// NaN comparisons return false per IEEE 754.
//
//go:noescape
func cmpLtFloat64AVX2(values *float64, threshold float64) uint64

// cmpLeFloat64AVX2 compares 4 float64 values against a threshold using AVX2.
// Returns a bitmask where bit i is set if values[i] <= threshold.
// NaN comparisons return false per IEEE 754.
//
//go:noescape
func cmpLeFloat64AVX2(values *float64, threshold float64) uint64

// cmpEqFloat64AVX2 compares 4 float64 values for equality against a threshold using AVX2.
// Returns a bitmask where bit i is set if values[i] == threshold.
// NaN comparisons return false per IEEE 754.
//
//go:noescape
func cmpEqFloat64AVX2(values *float64, threshold float64) uint64

// cmpNeFloat64AVX2 compares 4 float64 values for inequality against a threshold using AVX2.
// Returns a bitmask where bit i is set if values[i] != threshold.
// NaN != x returns true for all x per IEEE 754.
//
//go:noescape
func cmpNeFloat64AVX2(values *float64, threshold float64) uint64
