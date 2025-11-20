//go:build !amd64 && !arm64

package syndrdbsimd

// Generic implementation fallback for non-SIMD architectures

func cmpEqInt64Impl(values []int64, threshold int64) []bool {
	return cmpEqInt64Generic(values, threshold)
}

func cmpEqInt64MaskImpl(values []int64, threshold int64) []uint64 {
	return cmpEqInt64MaskGeneric(values, threshold)
}

func cmpNeInt64Impl(values []int64, threshold int64) []bool {
	return cmpNeInt64Generic(values, threshold)
}

func cmpNeInt64MaskImpl(values []int64, threshold int64) []uint64 {
	return cmpNeInt64MaskGeneric(values, threshold)
}

func cmpGtInt64Impl(values []int64, threshold int64) []bool {
	return cmpGtInt64Generic(values, threshold)
}

func cmpGtInt64MaskImpl(values []int64, threshold int64) []uint64 {
	return cmpGtInt64MaskGeneric(values, threshold)
}

func cmpLtInt64Impl(values []int64, threshold int64) []bool {
	return cmpLtInt64Generic(values, threshold)
}

func cmpLtInt64MaskImpl(values []int64, threshold int64) []uint64 {
	return cmpLtInt64MaskGeneric(values, threshold)
}

func cmpGeInt64Impl(values []int64, threshold int64) []bool {
	return cmpGeInt64Generic(values, threshold)
}

func cmpGeInt64MaskImpl(values []int64, threshold int64) []uint64 {
	return cmpGeInt64MaskGeneric(values, threshold)
}

func cmpLeInt64Impl(values []int64, threshold int64) []bool {
	return cmpLeInt64Generic(values, threshold)
}

func cmpLeInt64MaskImpl(values []int64, threshold int64) []uint64 {
	return cmpLeInt64MaskGeneric(values, threshold)
}

func andBitmapImpl(a, b []uint64) []uint64 {
	return andBitmapGeneric(a, b)
}

func orBitmapImpl(a, b []uint64) []uint64 {
	return orBitmapGeneric(a, b)
}

func xorBitmapImpl(a, b []uint64) []uint64 {
	return xorBitmapGeneric(a, b)
}

func notBitmapImpl(a []uint64) []uint64 {
	return notBitmapGeneric(a)
}

func popCountImpl(bitmap []uint64) int {
	return popCountGeneric(bitmap)
}

// ============================================================================
// Phase 2: Aggregation Operations
// ============================================================================

func sumInt64Impl(values []int64) int64 {
	return sumInt64Generic(values)
}

func minInt64Impl(values []int64) int64 {
	return minInt64Generic(values)
}

func maxInt64Impl(values []int64) int64 {
	return maxInt64Generic(values)
}

func countNonNullImpl(values []int64, nullBitmap []uint64) int64 {
	return countNonNullGeneric(values, nullBitmap)
}

// ============================================================================
// Phase 3: Hashing Operations
// ============================================================================

func hashInt64Impl(values []int64, output []uint64) {
	hashInt64SliceGeneric(values, output)
}

func crc32Int64Impl(values []int64, output []uint32) {
	crc32Int64SliceGeneric(values, output)
}

func xxhash64Impl(values []int64, output []uint64) {
	xxhash64SliceGeneric(values, output)
}

// ============================================================================
// Phase 4: String Operations
// ============================================================================

func strCmpImpl(a, b []byte) int {
	return strCmpGeneric(a, b)
}

func strPrefixCmpImpl(str, prefix []byte) bool {
	return strPrefixCmpGeneric(str, prefix)
}

func strEqImpl(a, b []byte) bool {
	return strEqGeneric(a, b)
}

func strToLowerImpl(s []byte) {
	strToLowerGeneric(s)
}

func strToUpperImpl(s []byte) {
	strToUpperGeneric(s)
}

// ============================================================================
// Float64 Comparisons
// ============================================================================

func cmpGtFloat64Impl(values []float64, threshold float64) []bool {
	return cmpGtFloat64Generic(values, threshold)
}

func cmpGtFloat64MaskImpl(values []float64, threshold float64) []uint64 {
	return cmpGtFloat64MaskGeneric(values, threshold)
}

func cmpGeFloat64Impl(values []float64, threshold float64) []bool {
	return cmpGeFloat64Generic(values, threshold)
}

func cmpGeFloat64MaskImpl(values []float64, threshold float64) []uint64 {
	return cmpGeFloat64MaskGeneric(values, threshold)
}

func cmpLtFloat64Impl(values []float64, threshold float64) []bool {
	return cmpLtFloat64Generic(values, threshold)
}

func cmpLtFloat64MaskImpl(values []float64, threshold float64) []uint64 {
	return cmpLtFloat64MaskGeneric(values, threshold)
}

func cmpLeFloat64Impl(values []float64, threshold float64) []bool {
	return cmpLeFloat64Generic(values, threshold)
}

func cmpLeFloat64MaskImpl(values []float64, threshold float64) []uint64 {
	return cmpLeFloat64MaskGeneric(values, threshold)
}

func cmpEqFloat64Impl(values []float64, threshold float64) []bool {
	return cmpEqFloat64Generic(values, threshold)
}

func cmpEqFloat64MaskImpl(values []float64, threshold float64) []uint64 {
	return cmpEqFloat64MaskGeneric(values, threshold)
}

func cmpNeFloat64Impl(values []float64, threshold float64) []bool {
	return cmpNeFloat64Generic(values, threshold)
}

func cmpNeFloat64MaskImpl(values []float64, threshold float64) []uint64 {
	return cmpNeFloat64MaskGeneric(values, threshold)
}

// ============================================================================
// String Comparisons
// ============================================================================

func cmpEqStringImpl(values [][]byte, threshold []byte) []bool {
	return cmpEqStringGeneric(values, threshold)
}

func cmpEqStringMaskImpl(values [][]byte, threshold []byte) []uint64 {
	return cmpEqStringMaskGeneric(values, threshold)
}

func cmpNeStringImpl(values [][]byte, threshold []byte) []bool {
	return cmpNeStringGeneric(values, threshold)
}

func cmpNeStringMaskImpl(values [][]byte, threshold []byte) []uint64 {
	return cmpNeStringMaskGeneric(values, threshold)
}

func cmpHasPrefixStringImpl(values [][]byte, prefix []byte) []bool {
	return cmpHasPrefixStringGeneric(values, prefix)
}

func cmpHasPrefixStringMaskImpl(values [][]byte, prefix []byte) []uint64 {
	return cmpHasPrefixStringMaskGeneric(values, prefix)
}

func cmpHasSuffixStringImpl(values [][]byte, suffix []byte) []bool {
	return cmpHasSuffixStringGeneric(values, suffix)
}

func cmpHasSuffixStringMaskImpl(values [][]byte, suffix []byte) []uint64 {
	return cmpHasSuffixStringMaskGeneric(values, suffix)
}

func cmpContainsStringImpl(values [][]byte, substr []byte) []bool {
	return cmpContainsStringGeneric(values, substr)
}

func cmpContainsStringMaskImpl(values [][]byte, substr []byte) []uint64 {
	return cmpContainsStringMaskGeneric(values, substr)
}

func cmpEqStringIgnoreCaseImpl(values [][]byte, threshold []byte) []bool {
	return cmpEqStringIgnoreCaseGeneric(values, threshold)
}

func cmpEqStringIgnoreCaseMaskImpl(values [][]byte, threshold []byte) []uint64 {
	return cmpEqStringIgnoreCaseMaskGeneric(values, threshold)
}

func cmpMatchWildcardImpl(values [][]byte, pattern []byte) []bool {
	return cmpMatchWildcardGeneric(values, pattern)
}

func cmpMatchWildcardMaskImpl(values [][]byte, pattern []byte) []uint64 {
	return cmpMatchWildcardMaskGeneric(values, pattern)
}
