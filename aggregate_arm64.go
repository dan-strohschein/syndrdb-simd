// +build arm64

package syndrdbsimd

// sumInt64NEON computes sum using NEON SIMD instructions.
// Processes 2 int64 values per iteration (128-bit NEON register).
func sumInt64NEON(values *int64, length int) int64

// minInt64NEON finds minimum using NEON SIMD instructions.
// Processes 2 int64 values per iteration.
func minInt64NEON(values *int64, length int) int64

// maxInt64NEON finds maximum using NEON SIMD instructions.
// Processes 2 int64 values per iteration.
func maxInt64NEON(values *int64, length int) int64

// countNonNullNEON counts non-null values using NEON SIMD instructions.
// The nullBitmap uses bit i to indicate if values[i] is null (1 = null, 0 = not null).
func countNonNullNEON(values *int64, nullBitmap *uint64, length int) int64
