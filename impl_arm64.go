//go:build arm64

package syndrdbsimd

// cmpEqInt64Impl routes to NEON or generic implementation based on CPU capabilities
func cmpEqInt64Impl(values []int64, threshold int64) []bool {
	if !HasNEON() || len(values) < 8 {
		return cmpEqInt64Generic(values, threshold)
	}

	results := make([]bool, len(values))
	i := 0

	// Process 2 elements at a time with NEON
	for ; i+1 < len(values); i += 2 {
		mask := cmpEqInt64NEON(&values[i], threshold)
		results[i+0] = (mask & 0x1) != 0
		results[i+1] = (mask & 0x2) != 0
	}

	// Handle remainder with scalar
	for ; i < len(values); i++ {
		results[i] = values[i] == threshold
	}

	return results
}

func cmpEqInt64MaskImpl(values []int64, threshold int64) []uint64 {
	if !HasNEON() || len(values) < 8 {
		return cmpEqInt64MaskGeneric(values, threshold)
	}

	numWords := (len(values) + 63) / 64
	mask := make([]uint64, numWords)
	i := 0

	// Process 2 elements at a time
	for ; i+1 < len(values); i += 2 {
		cmpMask := cmpEqInt64NEON(&values[i], threshold)
		wordIdx := i / 64
		bitIdx := uint(i % 64)

		for lane := 0; lane < 2 && i+lane < len(values); lane++ {
			if (cmpMask & (1 << uint(lane))) != 0 {
				mask[wordIdx] |= 1 << (bitIdx + uint(lane))
			}
		}
	}

	// Handle remainder
	for ; i < len(values); i++ {
		if values[i] == threshold {
			wordIdx := i / 64
			bitIdx := uint(i % 64)
			mask[wordIdx] |= 1 << bitIdx
		}
	}

	return mask
}

func cmpNeInt64Impl(values []int64, threshold int64) []bool {
	if !HasNEON() || len(values) < 8 {
		return cmpNeInt64Generic(values, threshold)
	}

	results := make([]bool, len(values))
	i := 0

	for ; i+1 < len(values); i += 2 {
		mask := cmpNeInt64NEON(&values[i], threshold)
		results[i+0] = (mask & 0x1) != 0
		results[i+1] = (mask & 0x2) != 0
	}

	for ; i < len(values); i++ {
		results[i] = values[i] != threshold
	}

	return results
}

func cmpNeInt64MaskImpl(values []int64, threshold int64) []uint64 {
	if !HasNEON() || len(values) < 8 {
		return cmpNeInt64MaskGeneric(values, threshold)
	}

	numWords := (len(values) + 63) / 64
	mask := make([]uint64, numWords)
	i := 0

	for ; i+1 < len(values); i += 2 {
		cmpMask := cmpNeInt64NEON(&values[i], threshold)
		wordIdx := i / 64
		bitIdx := uint(i % 64)

		for lane := 0; lane < 2 && i+lane < len(values); lane++ {
			if (cmpMask & (1 << uint(lane))) != 0 {
				mask[wordIdx] |= 1 << (bitIdx + uint(lane))
			}
		}
	}

	for ; i < len(values); i++ {
		if values[i] != threshold {
			wordIdx := i / 64
			bitIdx := uint(i % 64)
			mask[wordIdx] |= 1 << bitIdx
		}
	}

	return mask
}

func cmpGtInt64Impl(values []int64, threshold int64) []bool {
	if !HasNEON() || len(values) < 8 {
		return cmpGtInt64Generic(values, threshold)
	}

	results := make([]bool, len(values))
	i := 0

	for ; i+1 < len(values); i += 2 {
		mask := cmpGtInt64NEON(&values[i], threshold)
		results[i+0] = (mask & 0x1) != 0
		results[i+1] = (mask & 0x2) != 0
	}

	for ; i < len(values); i++ {
		results[i] = values[i] > threshold
	}

	return results
}

func cmpGtInt64MaskImpl(values []int64, threshold int64) []uint64 {
	if !HasNEON() || len(values) < 8 {
		return cmpGtInt64MaskGeneric(values, threshold)
	}

	numWords := (len(values) + 63) / 64
	mask := make([]uint64, numWords)
	i := 0

	for ; i+1 < len(values); i += 2 {
		cmpMask := cmpGtInt64NEON(&values[i], threshold)
		wordIdx := i / 64
		bitIdx := uint(i % 64)

		for lane := 0; lane < 2 && i+lane < len(values); lane++ {
			if (cmpMask & (1 << uint(lane))) != 0 {
				mask[wordIdx] |= 1 << (bitIdx + uint(lane))
			}
		}
	}

	for ; i < len(values); i++ {
		if values[i] > threshold {
			wordIdx := i / 64
			bitIdx := uint(i % 64)
			mask[wordIdx] |= 1 << bitIdx
		}
	}

	return mask
}

func cmpLtInt64Impl(values []int64, threshold int64) []bool {
	if !HasNEON() || len(values) < 8 {
		return cmpLtInt64Generic(values, threshold)
	}

	results := make([]bool, len(values))
	i := 0

	for ; i+1 < len(values); i += 2 {
		mask := cmpLtInt64NEON(&values[i], threshold)
		results[i+0] = (mask & 0x1) != 0
		results[i+1] = (mask & 0x2) != 0
	}

	for ; i < len(values); i++ {
		results[i] = values[i] < threshold
	}

	return results
}

func cmpLtInt64MaskImpl(values []int64, threshold int64) []uint64 {
	if !HasNEON() || len(values) < 8 {
		return cmpLtInt64MaskGeneric(values, threshold)
	}

	numWords := (len(values) + 63) / 64
	mask := make([]uint64, numWords)
	i := 0

	for ; i+1 < len(values); i += 2 {
		cmpMask := cmpLtInt64NEON(&values[i], threshold)
		wordIdx := i / 64
		bitIdx := uint(i % 64)

		for lane := 0; lane < 2 && i+lane < len(values); lane++ {
			if (cmpMask & (1 << uint(lane))) != 0 {
				mask[wordIdx] |= 1 << (bitIdx + uint(lane))
			}
		}
	}

	for ; i < len(values); i++ {
		if values[i] < threshold {
			wordIdx := i / 64
			bitIdx := uint(i % 64)
			mask[wordIdx] |= 1 << bitIdx
		}
	}

	return mask
}

func cmpGeInt64Impl(values []int64, threshold int64) []bool {
	if !HasNEON() || len(values) < 8 {
		return cmpGeInt64Generic(values, threshold)
	}

	results := make([]bool, len(values))
	i := 0

	for ; i+1 < len(values); i += 2 {
		mask := cmpGeInt64NEON(&values[i], threshold)
		results[i+0] = (mask & 0x1) != 0
		results[i+1] = (mask & 0x2) != 0
	}

	for ; i < len(values); i++ {
		results[i] = values[i] >= threshold
	}

	return results
}

func cmpGeInt64MaskImpl(values []int64, threshold int64) []uint64 {
	if !HasNEON() || len(values) < 8 {
		return cmpGeInt64MaskGeneric(values, threshold)
	}

	numWords := (len(values) + 63) / 64
	mask := make([]uint64, numWords)
	i := 0

	for ; i+1 < len(values); i += 2 {
		cmpMask := cmpGeInt64NEON(&values[i], threshold)
		wordIdx := i / 64
		bitIdx := uint(i % 64)

		for lane := 0; lane < 2 && i+lane < len(values); lane++ {
			if (cmpMask & (1 << uint(lane))) != 0 {
				mask[wordIdx] |= 1 << (bitIdx + uint(lane))
			}
		}
	}

	for ; i < len(values); i++ {
		if values[i] >= threshold {
			wordIdx := i / 64
			bitIdx := uint(i % 64)
			mask[wordIdx] |= 1 << bitIdx
		}
	}

	return mask
}

func cmpLeInt64Impl(values []int64, threshold int64) []bool {
	if !HasNEON() || len(values) < 8 {
		return cmpLeInt64Generic(values, threshold)
	}

	results := make([]bool, len(values))
	i := 0

	for ; i+1 < len(values); i += 2 {
		mask := cmpLeInt64NEON(&values[i], threshold)
		results[i+0] = (mask & 0x1) != 0
		results[i+1] = (mask & 0x2) != 0
	}

	for ; i < len(values); i++ {
		results[i] = values[i] <= threshold
	}

	return results
}

func cmpLeInt64MaskImpl(values []int64, threshold int64) []uint64 {
	if !HasNEON() || len(values) < 8 {
		return cmpLeInt64MaskGeneric(values, threshold)
	}

	numWords := (len(values) + 63) / 64
	mask := make([]uint64, numWords)
	i := 0

	for ; i+1 < len(values); i += 2 {
		cmpMask := cmpLeInt64NEON(&values[i], threshold)
		wordIdx := i / 64
		bitIdx := uint(i % 64)

		for lane := 0; lane < 2 && i+lane < len(values); lane++ {
			if (cmpMask & (1 << uint(lane))) != 0 {
				mask[wordIdx] |= 1 << (bitIdx + uint(lane))
			}
		}
	}

	for ; i < len(values); i++ {
		if values[i] <= threshold {
			wordIdx := i / 64
			bitIdx := uint(i % 64)
			mask[wordIdx] |= 1 << bitIdx
		}
	}

	return mask
}

func andBitmapImpl(a, b []uint64) []uint64 {
	if !HasNEON() || len(a) < 4 {
		return andBitmapGeneric(a, b)
	}

	length := len(a)
	if len(b) < length {
		length = len(b)
	}

	result := make([]uint64, length)
	andBitmapNEON(&result[0], &a[0], &b[0], length)
	return result
}

func orBitmapImpl(a, b []uint64) []uint64 {
	if !HasNEON() || len(a) < 4 {
		return orBitmapGeneric(a, b)
	}

	length := len(a)
	if len(b) < length {
		length = len(b)
	}

	result := make([]uint64, length)
	orBitmapNEON(&result[0], &a[0], &b[0], length)
	return result
}

func xorBitmapImpl(a, b []uint64) []uint64 {
	if !HasNEON() || len(a) < 4 {
		return xorBitmapGeneric(a, b)
	}

	length := len(a)
	if len(b) < length {
		length = len(b)
	}

	result := make([]uint64, length)
	xorBitmapNEON(&result[0], &a[0], &b[0], length)
	return result
}

func notBitmapImpl(a []uint64) []uint64 {
	if !HasNEON() || len(a) < 4 {
		return notBitmapGeneric(a)
	}

	result := make([]uint64, len(a))
	notBitmapNEON(&result[0], &a[0], len(a))
	return result
}

func popCountImpl(bitmap []uint64) int {
	if !HasNEON() || len(bitmap) < 4 {
		return popCountGeneric(bitmap)
	}

	return popCountNEON(&bitmap[0], len(bitmap))
}

// ============================================================================
// Phase 2: Aggregation Operations
// ============================================================================

func sumInt64Impl(values []int64) int64 {
	if !HasNEON() || len(values) < 8 {
		return sumInt64Generic(values)
	}

	return sumInt64NEON(&values[0], len(values))
}

func minInt64Impl(values []int64) int64 {
	// TODO: NEON min/max needs debugging - use generic for now
	return minInt64Generic(values)

	// if !HasNEON() || len(values) < 8 {
	// 	return minInt64Generic(values)
	// }
	// return minInt64NEON(&values[0], len(values))
}

func maxInt64Impl(values []int64) int64 {
	// TODO: NEON min/max needs debugging - use generic for now
	return maxInt64Generic(values)

	// if !HasNEON() || len(values) < 8 {
	// 	return maxInt64Generic(values)
	// }
	// return maxInt64NEON(&values[0], len(values))
}

func countNonNullImpl(values []int64, nullBitmap []uint64) int64 {
	// TODO: NEON countNonNull needs debugging - use generic for now
	return countNonNullGeneric(values, nullBitmap)

	// if len(nullBitmap) == 0 {
	// 	return int64(len(values))
	// }
	// if !HasNEON() || len(values) < 8 {
	// 	return countNonNullGeneric(values, nullBitmap)
	// }
	// return countNonNullNEON(&values[0], &nullBitmap[0], len(values))
}

// ============================================================================
// Phase 3: Hashing Operations
// ============================================================================

func hashInt64Impl(values []int64, output []uint64) {
	// TODO: NEON FNV hash needs debugging - use generic for now
	hashInt64SliceGeneric(values, output)
	return

	// if !HasNEON() || len(values) < 8 {
	// 	hashInt64SliceGeneric(values, output)
	// 	return
	// }
	// i := 0
	// // Process 2 elements at a time with NEON
	// for ; i+1 < len(values); i += 2 {
	// 	hashInt64NEON(&values[i], &output[i], 2)
	// }
	// // Handle remainder with scalar
	// for ; i < len(values); i++ {
	// 	output[i] = hashInt64Generic(values[i])
	// }
}

func crc32Int64Impl(values []int64, output []uint32) {
	// TODO: NEON CRC32C needs debugging - use generic for now
	crc32Int64SliceGeneric(values, output)
	return

	// // Always use NEON if available (hardware CRC32C)
	// if !HasNEON() {
	// 	crc32Int64SliceGeneric(values, output)
	// 	return
	// }
	// crc32Int64NEON(&values[0], &output[0], len(values))
}

func xxhash64Impl(values []int64, output []uint64) {
	if !HasNEON() || len(values) < 8 {
		xxhash64SliceGeneric(values, output)
		return
	}

	i := 0
	// Process 2 elements at a time with NEON
	for ; i+1 < len(values); i += 2 {
		xxhash64NEON(&values[i], &output[i], 2)
	}

	// Handle remainder with scalar
	for ; i < len(values); i++ {
		output[i] = xxhash64Generic(values[i])
	}
}

// ============================================================================
// Phase 4: String Operations
// ============================================================================

func strCmpImpl(a, b []byte) int {
	return strCmpGeneric(a, b)
}

func strPrefixCmpImpl(str, prefix []byte) bool {
	if !HasNEON() || len(prefix) < 16 {
		return strPrefixCmpGeneric(str, prefix)
	}

	if len(prefix) > len(str) {
		return false
	}

	result := strPrefixCmpNEON(&str[0], &prefix[0], len(str), len(prefix))
	return result == 1
}

func strEqImpl(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	if len(a) == 0 {
		return true
	}

	if !HasNEON() || len(a) < 16 {
		return strEqGeneric(a, b)
	}

	result := strEqNEON(&a[0], &b[0], len(a))
	return result == 1
}

func strToLowerImpl(s []byte) {
	if len(s) == 0 {
		return
	}

	// Note: NEON implementation disabled - Go's ARM64 assembler doesn't support
	// vector comparison instructions (VCMGE, VCMGT) needed for case conversion range checks
	strToLowerGeneric(s)
}

func strToUpperImpl(s []byte) {
	if len(s) == 0 {
		return
	}

	// Note: NEON implementation disabled - Go's ARM64 assembler doesn't support
	// vector comparison instructions (VCMGE, VCMGT) needed for case conversion range checks
	strToUpperGeneric(s)
}
