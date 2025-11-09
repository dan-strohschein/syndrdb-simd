//go:build arm64

package syndrdbsimd

// andBitmapNEON performs bitwise AND on bitmaps using NEON.
// Processes 2 uint64 values (128 bits) at a time.
//
//go:noescape
func andBitmapNEON(dst, a, b *uint64, length int)

// orBitmapNEON performs bitwise OR on bitmaps using NEON.
// Processes 2 uint64 values (128 bits) at a time.
//
//go:noescape
func orBitmapNEON(dst, a, b *uint64, length int)

// xorBitmapNEON performs bitwise XOR on bitmaps using NEON.
// Processes 2 uint64 values (128 bits) at a time.
//
//go:noescape
func xorBitmapNEON(dst, a, b *uint64, length int)

// notBitmapNEON performs bitwise NOT on a bitmap using NEON.
// Processes 2 uint64 values (128 bits) at a time.
//
//go:noescape
func notBitmapNEON(dst, src *uint64, length int)

// popCountNEON counts set bits in a bitmap using NEON.
// Returns the total number of 1 bits across all uint64 values.
//
//go:noescape
func popCountNEON(bitmap *uint64, length int) int
