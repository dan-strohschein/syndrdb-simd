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
