package syndrdbsimd

import "math/bits"

// boolsToBitmask converts a slice of booleans to a bitmask representation.
// Each uint64 in the result holds 64 boolean values as bits.
// Bit i in result[j] corresponds to bools[j*64 + i].
func boolsToBitmask(bools []bool) []uint64 {
	// Calculate number of uint64s needed to hold all bools
	numWords := (len(bools) + 63) / 64
	mask := make([]uint64, numWords)

	for i, b := range bools {
		if b {
			wordIdx := i / 64
			bitIdx := uint(i % 64)
			mask[wordIdx] |= 1 << bitIdx
		}
	}

	return mask
}

// bitmaskToBools converts a bitmask to a slice of booleans.
// The length parameter specifies how many booleans to extract.
func bitmaskToBools(mask []uint64, length int) []bool {
	bools := make([]bool, length)

	for i := 0; i < length; i++ {
		wordIdx := i / 64
		bitIdx := uint(i % 64)
		if wordIdx < len(mask) {
			bools[i] = (mask[wordIdx] & (1 << bitIdx)) != 0
		}
	}

	return bools
}

// andBitmapGeneric performs bitwise AND on two bitmasks using scalar operations.
// Returns a new bitmask where result[i] = a[i] & b[i].
func andBitmapGeneric(a, b []uint64) []uint64 {
	// Use the shorter length
	length := len(a)
	if len(b) < length {
		length = len(b)
	}

	result := make([]uint64, length)
	for i := 0; i < length; i++ {
		result[i] = a[i] & b[i]
	}

	return result
}

// orBitmapGeneric performs bitwise OR on two bitmasks using scalar operations.
// Returns a new bitmask where result[i] = a[i] | b[i].
func orBitmapGeneric(a, b []uint64) []uint64 {
	// Use the shorter length
	length := len(a)
	if len(b) < length {
		length = len(b)
	}

	result := make([]uint64, length)
	for i := 0; i < length; i++ {
		result[i] = a[i] | b[i]
	}

	return result
}

// xorBitmapGeneric performs bitwise XOR on two bitmasks using scalar operations.
// Returns a new bitmask where result[i] = a[i] ^ b[i].
func xorBitmapGeneric(a, b []uint64) []uint64 {
	// Use the shorter length
	length := len(a)
	if len(b) < length {
		length = len(b)
	}

	result := make([]uint64, length)
	for i := 0; i < length; i++ {
		result[i] = a[i] ^ b[i]
	}

	return result
}

// notBitmapGeneric performs bitwise NOT on a bitmask using scalar operations.
// Returns a new bitmask where result[i] = ^a[i].
func notBitmapGeneric(a []uint64) []uint64 {
	result := make([]uint64, len(a))
	for i := 0; i < len(a); i++ {
		result[i] = ^a[i]
	}

	return result
}

// popCountGeneric counts the number of set bits in a bitmask using scalar operations.
// This is equivalent to counting how many true values the bitmask represents.
func popCountGeneric(bitmap []uint64) int {
	count := 0
	for _, word := range bitmap {
		count += bits.OnesCount64(word)
	}
	return count
}

// andBitmapInPlaceGeneric performs in-place bitwise AND on bitmask a with bitmask b.
// This modifies a directly: a[i] &= b[i].
// TODO: I could add SIMD versions of in-place operations for better cache performance.
func andBitmapInPlaceGeneric(a, b []uint64) {
	length := len(a)
	if len(b) < length {
		length = len(b)
	}

	for i := 0; i < length; i++ {
		a[i] &= b[i]
	}
}

// orBitmapInPlaceGeneric performs in-place bitwise OR on bitmask a with bitmask b.
// This modifies a directly: a[i] |= b[i].
func orBitmapInPlaceGeneric(a, b []uint64) {
	length := len(a)
	if len(b) < length {
		length = len(b)
	}

	for i := 0; i < length; i++ {
		a[i] |= b[i]
	}
}

// xorBitmapInPlaceGeneric performs in-place bitwise XOR on bitmask a with bitmask b.
// This modifies a directly: a[i] ^= b[i].
func xorBitmapInPlaceGeneric(a, b []uint64) {
	length := len(a)
	if len(b) < length {
		length = len(b)
	}

	for i := 0; i < length; i++ {
		a[i] ^= b[i]
	}
}

// notBitmapInPlaceGeneric performs in-place bitwise NOT on bitmask a.
// This modifies a directly: a[i] = ^a[i].
func notBitmapInPlaceGeneric(a []uint64) {
	for i := 0; i < len(a); i++ {
		a[i] = ^a[i]
	}
}
