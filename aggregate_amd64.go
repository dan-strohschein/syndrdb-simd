// +build amd64

package syndrdbsimd

// sumInt64AVX2 computes sum using AVX2 SIMD instructions.
// Processes 4 int64 values per iteration (256-bit AVX2 register).
func sumInt64AVX2(values *int64, length int) int64

// minInt64AVX2 finds minimum using AVX2 SIMD instructions.
// Processes 4 int64 values per iteration.
func minInt64AVX2(values *int64, length int) int64

// maxInt64AVX2 finds maximum using AVX2 SIMD instructions.
// Processes 4 int64 values per iteration.
func maxInt64AVX2(values *int64, length int) int64

// countNonNullAVX2 counts non-null values using AVX2 SIMD instructions.
// The nullBitmap uses bit i to indicate if values[i] is null (1 = null, 0 = not null).
func countNonNullAVX2(values *int64, nullBitmap *uint64, length int) int64
