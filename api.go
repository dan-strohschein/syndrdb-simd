// Package syndrdbsimd provides high-performance SIMD operations for SyndrDB.
//
// This package implements database-critical operations using SIMD instructions
// (AVX2 on x86-64, NEON on ARM64) with automatic fallback to scalar implementations
// when SIMD is not available.
//
// Phase 1 includes:
//   - Comparison operations (==, !=, >, <, >=, <=) on int64 arrays
//   - Bitmap operations (AND, OR, XOR, NOT, PopCount)
//   - Both bool slice and bitmask return variants
//
// Performance characteristics:
//   - AVX2: ~4x speedup on comparisons, ~4x on bitmap operations
//   - NEON: ~2x speedup on comparisons, ~2x on bitmap operations
//   - Automatic CPU feature detection and optimal implementation selection
package syndrdbsimd

// CmpEqInt64 compares int64 values for equality against a threshold.
// Returns a slice of booleans where result[i] == true if values[i] == threshold.
//
// This function automatically selects the best implementation:
//   - AVX2 on x86-64 processors (4 elements per operation)
//   - NEON on ARM64 processors (2 elements per operation)
//   - Scalar fallback on other architectures
func CmpEqInt64(values []int64, threshold int64) []bool {
	if len(values) == 0 {
		return []bool{}
	}

	// TODO: I could add a heuristic to skip SIMD for very small arrays (< 8 elements)
	// where the overhead might not be worth it

	return cmpEqInt64Impl(values, threshold)
}

// CmpEqInt64Mask compares int64 values for equality and returns a bitmask.
// Returns a slice of uint64 where bit i in result[j] is set if values[j*64+i] == threshold.
// This is more memory-efficient than CmpEqInt64 for large datasets.
func CmpEqInt64Mask(values []int64, threshold int64) []uint64 {
	if len(values) == 0 {
		return []uint64{}
	}

	return cmpEqInt64MaskImpl(values, threshold)
}

// CmpNeInt64 compares int64 values for inequality against a threshold.
// Returns a slice of booleans where result[i] == true if values[i] != threshold.
func CmpNeInt64(values []int64, threshold int64) []bool {
	if len(values) == 0 {
		return []bool{}
	}

	return cmpNeInt64Impl(values, threshold)
}

// CmpNeInt64Mask compares int64 values for inequality and returns a bitmask.
func CmpNeInt64Mask(values []int64, threshold int64) []uint64 {
	if len(values) == 0 {
		return []uint64{}
	}

	return cmpNeInt64MaskImpl(values, threshold)
}

// CmpGtInt64 compares int64 values for greater-than against a threshold.
// Returns a slice of booleans where result[i] == true if values[i] > threshold.
func CmpGtInt64(values []int64, threshold int64) []bool {
	if len(values) == 0 {
		return []bool{}
	}

	return cmpGtInt64Impl(values, threshold)
}

// CmpGtInt64Mask compares int64 values for greater-than and returns a bitmask.
func CmpGtInt64Mask(values []int64, threshold int64) []uint64 {
	if len(values) == 0 {
		return []uint64{}
	}

	return cmpGtInt64MaskImpl(values, threshold)
}

// CmpLtInt64 compares int64 values for less-than against a threshold.
// Returns a slice of booleans where result[i] == true if values[i] < threshold.
func CmpLtInt64(values []int64, threshold int64) []bool {
	if len(values) == 0 {
		return []bool{}
	}

	return cmpLtInt64Impl(values, threshold)
}

// CmpLtInt64Mask compares int64 values for less-than and returns a bitmask.
func CmpLtInt64Mask(values []int64, threshold int64) []uint64 {
	if len(values) == 0 {
		return []uint64{}
	}

	return cmpLtInt64MaskImpl(values, threshold)
}

// CmpGeInt64 compares int64 values for greater-than-or-equal against a threshold.
// Returns a slice of booleans where result[i] == true if values[i] >= threshold.
func CmpGeInt64(values []int64, threshold int64) []bool {
	if len(values) == 0 {
		return []bool{}
	}

	return cmpGeInt64Impl(values, threshold)
}

// CmpGeInt64Mask compares int64 values for greater-than-or-equal and returns a bitmask.
func CmpGeInt64Mask(values []int64, threshold int64) []uint64 {
	if len(values) == 0 {
		return []uint64{}
	}

	return cmpGeInt64MaskImpl(values, threshold)
}

// CmpLeInt64 compares int64 values for less-than-or-equal against a threshold.
// Returns a slice of booleans where result[i] == true if values[i] <= threshold.
func CmpLeInt64(values []int64, threshold int64) []bool {
	if len(values) == 0 {
		return []bool{}
	}

	return cmpLeInt64Impl(values, threshold)
}

// CmpLeInt64Mask compares int64 values for less-than-or-equal and returns a bitmask.
func CmpLeInt64Mask(values []int64, threshold int64) []uint64 {
	if len(values) == 0 {
		return []uint64{}
	}

	return cmpLeInt64MaskImpl(values, threshold)
}

// AndBitmap performs bitwise AND on two uint64 bitmaps.
// Returns a new bitmap where result[i] = a[i] & b[i].
// Uses SIMD when available for ~4x (AVX2) or ~2x (NEON) speedup.
func AndBitmap(a, b []uint64) []uint64 {
	if len(a) == 0 || len(b) == 0 {
		return []uint64{}
	}

	return andBitmapImpl(a, b)
}

// OrBitmap performs bitwise OR on two uint64 bitmaps.
// Returns a new bitmap where result[i] = a[i] | b[i].
func OrBitmap(a, b []uint64) []uint64 {
	if len(a) == 0 || len(b) == 0 {
		return []uint64{}
	}

	return orBitmapImpl(a, b)
}

// XorBitmap performs bitwise XOR on two uint64 bitmaps.
// Returns a new bitmap where result[i] = a[i] ^ b[i].
func XorBitmap(a, b []uint64) []uint64 {
	if len(a) == 0 || len(b) == 0 {
		return []uint64{}
	}

	return xorBitmapImpl(a, b)
}

// NotBitmap performs bitwise NOT on a uint64 bitmap.
// Returns a new bitmap where result[i] = ^a[i].
func NotBitmap(a []uint64) []uint64 {
	if len(a) == 0 {
		return []uint64{}
	}

	return notBitmapImpl(a)
}

// PopCount counts the number of set bits (1s) in a bitmap.
// This is extremely useful for counting matching rows in database queries.
// Uses SIMD when available for significant speedup.
func PopCount(bitmap []uint64) int {
	if len(bitmap) == 0 {
		return 0
	}

	return popCountImpl(bitmap)
}

// BoolsToBitmask converts a boolean slice to a compact bitmask representation.
// Each uint64 in the result holds 64 boolean values as bits.
// This is useful for memory-efficient storage of boolean arrays.
func BoolsToBitmask(bools []bool) []uint64 {
	return boolsToBitmask(bools)
}

// ============================================================================
// Float64 Comparisons
// ============================================================================

// CmpGtFloat64 compares float64 values for greater-than against a threshold.
// Returns a slice of booleans where result[i] == true if values[i] > threshold.
//
// NaN comparisons always return false per IEEE 754.
// Infinity values are compared normally (e.g., +Inf > any finite number is true).
//
// This function automatically selects the best implementation:
//   - AVX2 on x86-64 processors (4 elements per operation)
//   - NEON on ARM64 processors (2 elements per operation)
//   - Scalar fallback on other architectures
func CmpGtFloat64(values []float64, threshold float64) []bool {
	if len(values) == 0 {
		return []bool{}
	}

	return cmpGtFloat64Impl(values, threshold)
}

// CmpGtFloat64Mask compares float64 values for greater-than and returns a bitmask.
// Returns a slice of uint64 where bit i in result[j] is set if values[j*64+i] > threshold.
// This is more memory-efficient than CmpGtFloat64 for large datasets.
//
// NaN comparisons always return false per IEEE 754.
func CmpGtFloat64Mask(values []float64, threshold float64) []uint64 {
	if len(values) == 0 {
		return []uint64{}
	}

	return cmpGtFloat64MaskImpl(values, threshold)
}

// CmpGeFloat64 compares float64 values for greater-than-or-equal against a threshold.
// Returns a slice of booleans where result[i] == true if values[i] >= threshold.
//
// NaN comparisons always return false per IEEE 754.
func CmpGeFloat64(values []float64, threshold float64) []bool {
	if len(values) == 0 {
		return []bool{}
	}

	return cmpGeFloat64Impl(values, threshold)
}

// CmpGeFloat64Mask compares float64 values for greater-than-or-equal and returns a bitmask.
func CmpGeFloat64Mask(values []float64, threshold float64) []uint64 {
	if len(values) == 0 {
		return []uint64{}
	}

	return cmpGeFloat64MaskImpl(values, threshold)
}

// CmpLtFloat64 compares float64 values for less-than against a threshold.
// Returns a slice of booleans where result[i] == true if values[i] < threshold.
//
// NaN comparisons always return false per IEEE 754.
func CmpLtFloat64(values []float64, threshold float64) []bool {
	if len(values) == 0 {
		return []bool{}
	}

	return cmpLtFloat64Impl(values, threshold)
}

// CmpLtFloat64Mask compares float64 values for less-than and returns a bitmask.
func CmpLtFloat64Mask(values []float64, threshold float64) []uint64 {
	if len(values) == 0 {
		return []uint64{}
	}

	return cmpLtFloat64MaskImpl(values, threshold)
}

// CmpLeFloat64 compares float64 values for less-than-or-equal against a threshold.
// Returns a slice of booleans where result[i] == true if values[i] <= threshold.
//
// NaN comparisons always return false per IEEE 754.
func CmpLeFloat64(values []float64, threshold float64) []bool {
	if len(values) == 0 {
		return []bool{}
	}

	return cmpLeFloat64Impl(values, threshold)
}

// CmpLeFloat64Mask compares float64 values for less-than-or-equal and returns a bitmask.
func CmpLeFloat64Mask(values []float64, threshold float64) []uint64 {
	if len(values) == 0 {
		return []uint64{}
	}

	return cmpLeFloat64MaskImpl(values, threshold)
}

// CmpEqFloat64 compares float64 values for equality against a threshold.
// Returns a slice of booleans where result[i] == true if values[i] == threshold.
//
// NaN comparisons always return false per IEEE 754 (even NaN == NaN is false).
func CmpEqFloat64(values []float64, threshold float64) []bool {
	if len(values) == 0 {
		return []bool{}
	}

	return cmpEqFloat64Impl(values, threshold)
}

// CmpEqFloat64Mask compares float64 values for equality and returns a bitmask.
func CmpEqFloat64Mask(values []float64, threshold float64) []uint64 {
	if len(values) == 0 {
		return []uint64{}
	}

	return cmpEqFloat64MaskImpl(values, threshold)
}

// CmpNeFloat64 compares float64 values for inequality against a threshold.
// Returns a slice of booleans where result[i] == true if values[i] != threshold.
//
// NaN != x returns true for all x per IEEE 754 (including NaN != NaN).
func CmpNeFloat64(values []float64, threshold float64) []bool {
	if len(values) == 0 {
		return []bool{}
	}

	return cmpNeFloat64Impl(values, threshold)
}

// CmpNeFloat64Mask compares float64 values for inequality and returns a bitmask.
func CmpNeFloat64Mask(values []float64, threshold float64) []uint64 {
	if len(values) == 0 {
		return []uint64{}
	}

	return cmpNeFloat64MaskImpl(values, threshold)
}

// BitmaskToBools converts a bitmask back to a boolean slice.
// The length parameter specifies how many booleans to extract.
func BitmaskToBools(mask []uint64, length int) []bool {
	return bitmaskToBools(mask, length)
}

// ============================================================================
// String Comparisons
// ============================================================================

// CmpEqString compares string values for equality against a threshold.
// Returns a slice of booleans where result[i] == true if values[i] == threshold.
//
// Uses adaptive SIMD based on string count and average length. Configure thresholds
// via SetStringSIMDThreshold.
func CmpEqString(values []string, threshold string) []bool {
	if len(values) == 0 {
		return []bool{}
	}

	byteValues := stringsToBytes(values)
	byteThreshold := []byte(threshold)
	return cmpEqStringImpl(byteValues, byteThreshold)
}

// CmpEqStringMask compares string values for equality and returns a bitmask.
func CmpEqStringMask(values []string, threshold string) []uint64 {
	if len(values) == 0 {
		return []uint64{}
	}

	byteValues := stringsToBytes(values)
	byteThreshold := []byte(threshold)
	return cmpEqStringMaskImpl(byteValues, byteThreshold)
}

// CmpNeString compares string values for inequality against a threshold.
// Returns a slice of booleans where result[i] == true if values[i] != threshold.
func CmpNeString(values []string, threshold string) []bool {
	if len(values) == 0 {
		return []bool{}
	}

	byteValues := stringsToBytes(values)
	byteThreshold := []byte(threshold)
	return cmpNeStringImpl(byteValues, byteThreshold)
}

// CmpNeStringMask compares string values for inequality and returns a bitmask.
func CmpNeStringMask(values []string, threshold string) []uint64 {
	if len(values) == 0 {
		return []uint64{}
	}

	byteValues := stringsToBytes(values)
	byteThreshold := []byte(threshold)
	return cmpNeStringMaskImpl(byteValues, byteThreshold)
}

// CmpHasPrefixString checks if string values start with a given prefix.
// Returns a slice of booleans where result[i] == true if values[i] starts with prefix.
func CmpHasPrefixString(values []string, prefix string) []bool {
	if len(values) == 0 {
		return []bool{}
	}

	byteValues := stringsToBytes(values)
	bytePrefix := []byte(prefix)
	return cmpHasPrefixStringImpl(byteValues, bytePrefix)
}

// CmpHasPrefixStringMask checks for prefix match and returns a bitmask.
func CmpHasPrefixStringMask(values []string, prefix string) []uint64 {
	if len(values) == 0 {
		return []uint64{}
	}

	byteValues := stringsToBytes(values)
	bytePrefix := []byte(prefix)
	return cmpHasPrefixStringMaskImpl(byteValues, bytePrefix)
}

// CmpHasSuffixString checks if string values end with a given suffix.
// Returns a slice of booleans where result[i] == true if values[i] ends with suffix.
func CmpHasSuffixString(values []string, suffix string) []bool {
	if len(values) == 0 {
		return []bool{}
	}

	byteValues := stringsToBytes(values)
	byteSuffix := []byte(suffix)
	return cmpHasSuffixStringImpl(byteValues, byteSuffix)
}

// CmpHasSuffixStringMask checks for suffix match and returns a bitmask.
func CmpHasSuffixStringMask(values []string, suffix string) []uint64 {
	if len(values) == 0 {
		return []uint64{}
	}

	byteValues := stringsToBytes(values)
	byteSuffix := []byte(suffix)
	return cmpHasSuffixStringMaskImpl(byteValues, byteSuffix)
}

// CmpContainsString checks if string values contain a substring.
// Returns a slice of booleans where result[i] == true if values[i] contains substr.
func CmpContainsString(values []string, substr string) []bool {
	if len(values) == 0 {
		return []bool{}
	}

	byteValues := stringsToBytes(values)
	byteSubstr := []byte(substr)
	return cmpContainsStringImpl(byteValues, byteSubstr)
}

// CmpContainsStringMask checks for substring match and returns a bitmask.
func CmpContainsStringMask(values []string, substr string) []uint64 {
	if len(values) == 0 {
		return []uint64{}
	}

	byteValues := stringsToBytes(values)
	byteSubstr := []byte(substr)
	return cmpContainsStringMaskImpl(byteValues, byteSubstr)
}

// CmpEqStringIgnoreCase compares string values for equality (case-insensitive ASCII).
// Returns a slice of booleans where result[i] == true if values[i] equals threshold
// ignoring case for ASCII characters (A-Z, a-z). Non-ASCII bytes are compared directly.
func CmpEqStringIgnoreCase(values []string, threshold string) []bool {
	if len(values) == 0 {
		return []bool{}
	}

	byteValues := stringsToBytes(values)
	byteThreshold := []byte(threshold)
	return cmpEqStringIgnoreCaseImpl(byteValues, byteThreshold)
}

// CmpEqStringIgnoreCaseMask compares strings case-insensitively and returns a bitmask.
func CmpEqStringIgnoreCaseMask(values []string, threshold string) []uint64 {
	if len(values) == 0 {
		return []uint64{}
	}

	byteValues := stringsToBytes(values)
	byteThreshold := []byte(threshold)
	return cmpEqStringIgnoreCaseMaskImpl(byteValues, byteThreshold)
}

// CmpLikeString performs SQL LIKE pattern matching on string values.
// Supports % (any chars) and _ (single char) wildcards.
// Returns a slice of booleans where result[i] == true if values[i] matches pattern.
//
// Pattern is compiled and optimized based on structure:
//   - "exact" → exact equality check
//   - "prefix%" → prefix match
//   - "%suffix" → suffix match
//   - "%contains%" → substring match
//   - complex patterns with % and _ → wildcard matching with backtracking
//
// For repeated pattern usage, consider using CmpLikeStringCompiled for better performance.
func CmpLikeString(values []string, pattern string) []bool {
	if len(values) == 0 {
		return []bool{}
	}

	compiled, err := CompilePatternAuto(pattern)
	if err != nil {
		// Return all false on invalid pattern
		results := make([]bool, len(values))
		return results
	}

	return cmpLikeStringCompiled(values, compiled)
}

// CmpLikeStringMask performs SQL LIKE pattern matching and returns a bitmask.
func CmpLikeStringMask(values []string, pattern string) []uint64 {
	bools := CmpLikeString(values, pattern)
	return boolsToBitmask(bools)
}

// CmpLikeStringCompiled performs SQL LIKE pattern matching using a pre-compiled pattern.
// This is more efficient when matching the same pattern against multiple value sets.
//
// Example:
//
//	pattern, err := CompilePattern("hello%")
//	if err != nil { ... }
//	results1 := CmpLikeStringCompiled(values1, pattern)
//	results2 := CmpLikeStringCompiled(values2, pattern)
func CmpLikeStringCompiled(values []string, pattern *CompiledPattern) []bool {
	if len(values) == 0 {
		return []bool{}
	}

	return cmpLikeStringCompiled(values, pattern)
}

// CmpLikeStringCompiledMask performs compiled LIKE matching and returns a bitmask.
func CmpLikeStringCompiledMask(values []string, pattern *CompiledPattern) []uint64 {
	bools := CmpLikeStringCompiled(values, pattern)
	return boolsToBitmask(bools)
}

// cmpLikeStringCompiled is the internal implementation of LIKE matching.
func cmpLikeStringCompiled(values []string, pattern *CompiledPattern) []bool {
	byteValues := stringsToBytes(values)

	switch pattern.Type {
	case PatternExact:
		// Exact match
		return cmpEqStringImpl(byteValues, pattern.Segments[0])
	case PatternPrefix:
		// Prefix match
		return cmpHasPrefixStringImpl(byteValues, pattern.Segments[0])
	case PatternSuffix:
		// Suffix match
		return cmpHasSuffixStringImpl(byteValues, pattern.Segments[0])
	case PatternContains:
		// Substring match
		return cmpContainsStringImpl(byteValues, pattern.Segments[0])
	case PatternWildcard:
		// Complex wildcard matching
		return cmpMatchWildcardImpl(byteValues, []byte(pattern.OriginalPattern))
	default:
		// Unknown pattern type - return all false
		results := make([]bool, len(values))
		return results
	}
}

// ============================================================================
// Phase 2: Aggregation Operations
// ============================================================================

// SumInt64 computes the sum of all int64 values in the array.
// Returns 0 for empty arrays.
//
// This function automatically selects the best implementation:
//   - AVX2 on x86-64 processors (4 elements per operation)
//   - NEON on ARM64 processors (2 elements per operation)
//   - Scalar fallback on other architectures
//
// Performance: ~4-6x speedup with SIMD on large arrays.
func SumInt64(values []int64) int64 {
	if len(values) == 0 {
		return 0
	}

	return sumInt64Impl(values)
}

// MinInt64 finds the minimum int64 value in the array.
// Returns math.MaxInt64 for empty arrays.
//
// This function automatically selects the best implementation:
//   - AVX2 on x86-64 processors (4 elements per operation)
//   - NEON on ARM64 processors (2 elements per operation)
//   - Scalar fallback on other architectures
//
// Performance: ~4-6x speedup with SIMD on large arrays.
func MinInt64(values []int64) int64 {
	if len(values) == 0 {
		return 1<<63 - 1 // math.MaxInt64
	}

	return minInt64Impl(values)
}

// MaxInt64 finds the maximum int64 value in the array.
// Returns math.MinInt64 for empty arrays.
//
// This function automatically selects the best implementation:
//   - AVX2 on x86-64 processors (4 elements per operation)
//   - NEON on ARM64 processors (2 elements per operation)
//   - Scalar fallback on other architectures
//
// Performance: ~4-6x speedup with SIMD on large arrays.
func MaxInt64(values []int64) int64 {
	if len(values) == 0 {
		return -1 << 63 // math.MinInt64
	}

	return maxInt64Impl(values)
}

// CountNonNull counts the number of non-null values in the array.
// The nullBitmap parameter indicates which values are null (bit 1 = null, bit 0 = not null).
// If nullBitmap is nil or empty, all values are considered non-null.
//
// This function automatically selects the best implementation:
//   - AVX2 on x86-64 processors
//   - NEON on ARM64 processors
//   - Scalar fallback on other architectures
//
// Performance: ~2-4x speedup with SIMD on large arrays.
func CountNonNull(values []int64, nullBitmap []uint64) int64 {
	if len(values) == 0 {
		return 0
	}

	return countNonNullImpl(values, nullBitmap)
}

// AvgInt64 computes the average of all int64 values in the array.
// Returns 0.0 for empty arrays.
//
// Note: This uses SumInt64 internally, so it benefits from SIMD acceleration.
func AvgInt64(values []int64) float64 {
	if len(values) == 0 {
		return 0
	}

	sum := SumInt64(values)
	return float64(sum) / float64(len(values))
}

// ========================================
// Phase 3: Hashing Operations
// ========================================

// HashInt64 computes FNV-1a hashes for int64 values.
// The output slice must be pre-allocated with the same length as values.
//
// This function automatically selects the best implementation:
//   - AVX2 on x86-64 processors (4 elements per operation)
//   - NEON on ARM64 processors (2 elements per operation)
//   - Scalar fallback on other architectures
//
// FNV-1a is a fast, simple hash suitable for hash table operations.
// Performance: ~2-4x speedup with SIMD on large arrays.
func HashInt64(values []int64, output []uint64) {
	if len(values) == 0 || len(output) != len(values) {
		return
	}

	hashInt64Impl(values, output)
}

// CRC32 computes the CRC32 checksum of a byte slice using the IEEE polynomial.
// This uses the standard library's hash/crc32 implementation.
func CRC32(data []byte) uint32 {
	return crc32Generic(data)
}

// CRC32Int64 computes CRC32 checksums for int64 values.
// The output slice must be pre-allocated with the same length as values.
//
// This function automatically selects the best implementation:
//   - Hardware CRC32C on x86-64 processors (when available)
//   - Hardware CRC32C on ARM64 processors (when available)
//   - Scalar fallback otherwise
//
// Performance: ~2-3x speedup with hardware CRC32C instructions.
func CRC32Int64(values []int64, output []uint32) {
	if len(values) == 0 || len(output) != len(values) {
		return
	}

	crc32Int64Impl(values, output)
}

// XXHash64 computes XXHash64 hashes for int64 values.
// The output slice must be pre-allocated with the same length as values.
//
// This function automatically selects the best implementation:
//   - AVX2 on x86-64 processors (4 elements per operation)
//   - NEON on ARM64 processors (2 elements per operation)
//   - Scalar fallback on other architectures
//
// XXHash64 is a fast, high-quality non-cryptographic hash.
// Performance: ~3-5x speedup with SIMD on large arrays.
func XXHash64(values []int64, output []uint64) {
	if len(values) == 0 || len(output) != len(values) {
		return
	}

	xxhash64Impl(values, output)
}

// XXHash64Bytes computes the XXHash64 hash of a byte slice.
// This is useful for hashing variable-length keys in hash tables.
func XXHash64Bytes(data []byte) uint64 {
	return xxhash64BytesGeneric(data)
}

// ========================================
// Phase 4: String Operations
// ========================================

// StrCmp compares two byte slices lexicographically.
// Returns:
//
//	-1 if a < b
//	 0 if a == b
//	+1 if a > b
//
// This is equivalent to bytes.Compare but can use SIMD for acceleration.
func StrCmp(a, b []byte) int {
	return strCmpImpl(a, b)
}

// StrLen returns the length of a byte slice.
// This is a simple wrapper for len() but provided for API consistency.
func StrLen(s []byte) int {
	return len(s)
}

// StrPrefixCmp checks if str starts with prefix.
// Returns true if str has the prefix, false otherwise.
//
// This function automatically selects the best implementation:
//   - AVX2 on x86-64 processors (32 bytes per operation)
//   - NEON on ARM64 processors (16 bytes per operation)
//   - Scalar fallback on other architectures
//
// Equivalent to bytes.HasPrefix but can use SIMD for acceleration.
// Performance: ~2-4x speedup with SIMD on long prefixes.
func StrPrefixCmp(str, prefix []byte) bool {
	return strPrefixCmpImpl(str, prefix)
}

// StrContains checks if str contains substr.
// Returns true if substr is found in str, false otherwise.
//
// Equivalent to bytes.Contains.
func StrContains(str, substr []byte) bool {
	return strContainsGeneric(str, substr)
}

// StrEq checks if two byte slices are equal.
// Returns true if equal, false otherwise.
//
// This function automatically selects the best implementation:
//   - AVX2 on x86-64 processors (32 bytes per operation)
//   - NEON on ARM64 processors (16 bytes per operation)
//   - Scalar fallback on other architectures
//
// Equivalent to bytes.Equal but can use SIMD for acceleration.
// Performance: ~3-5x speedup with SIMD on long strings.
func StrEq(a, b []byte) bool {
	return strEqImpl(a, b)
}

// StrToLower converts a byte slice to lowercase (ASCII only).
// The input slice is modified in-place.
//
// This function automatically selects the best implementation:
//   - AVX2 on x86-64 processors (32 bytes per operation)
//   - NEON on ARM64 processors (16 bytes per operation)
//   - Scalar fallback on other architectures
//
// Performance: ~4-6x speedup with SIMD on long strings.
func StrToLower(s []byte) {
	strToLowerImpl(s)
}

// StrToUpper converts a byte slice to uppercase (ASCII only).
// The input slice is modified in-place.
//
// This function automatically selects the best implementation:
//   - AVX2 on x86-64 processors (32 bytes per operation)
//   - NEON on ARM64 processors (16 bytes per operation)
//   - Scalar fallback on other architectures
//
// Performance: ~4-6x speedup with SIMD on long strings.
func StrToUpper(s []byte) {
	strToUpperImpl(s)
}

// StrEqIgnoreCase checks if two byte slices are equal, ignoring case (ASCII only).
// Returns true if equal (case-insensitive), false otherwise.
func StrEqIgnoreCase(a, b []byte) bool {
	return strEqIgnoreCaseGeneric(a, b)
}
