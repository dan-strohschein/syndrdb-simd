//go:build arm64

package syndrdbsimd

// cmpGtInt64NEON compares 2 int64 values against a threshold using NEON.
// Returns a bitmask where bit i is set if values[i] > threshold.
// NEON processes 2 int64 values at a time (128-bit registers).
//
//go:noescape
func cmpGtInt64NEON(values *int64, threshold int64) uint64

// cmpEqInt64NEON compares 2 int64 values for equality against a threshold using NEON.
// Returns a bitmask where bit i is set if values[i] == threshold.
//
//go:noescape
func cmpEqInt64NEON(values *int64, threshold int64) uint64

// cmpLtInt64NEON compares 2 int64 values against a threshold using NEON.
// Returns a bitmask where bit i is set if values[i] < threshold.
//
//go:noescape
func cmpLtInt64NEON(values *int64, threshold int64) uint64

// cmpGeInt64NEON compares 2 int64 values against a threshold using NEON.
// Returns a bitmask where bit i is set if values[i] >= threshold.
//
//go:noescape
func cmpGeInt64NEON(values *int64, threshold int64) uint64

// cmpLeInt64NEON compares 2 int64 values against a threshold using NEON.
// Returns a bitmask where bit i is set if values[i] <= threshold.
//
//go:noescape
func cmpLeInt64NEON(values *int64, threshold int64) uint64

// cmpNeInt64NEON compares 2 int64 values for inequality against a threshold using NEON.
// Returns a bitmask where bit i is set if values[i] != threshold.
//
//go:noescape
func cmpNeInt64NEON(values *int64, threshold int64) uint64
