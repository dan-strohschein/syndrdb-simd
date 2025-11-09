//go:build amd64

package syndrdbsimd

// andBitmapAVX2 performs bitwise AND on bitmaps using AVX2.
// Processes 4 uint64 values (256 bits) at a time.
//
//go:noescape
func andBitmapAVX2(dst, a, b *uint64, length int)

// orBitmapAVX2 performs bitwise OR on bitmaps using AVX2.
// Processes 4 uint64 values (256 bits) at a time.
//
//go:noescape
func orBitmapAVX2(dst, a, b *uint64, length int)

// xorBitmapAVX2 performs bitwise XOR on bitmaps using AVX2.
// Processes 4 uint64 values (256 bits) at a time.
//
//go:noescape
func xorBitmapAVX2(dst, a, b *uint64, length int)

// notBitmapAVX2 performs bitwise NOT on a bitmap using AVX2.
// Processes 4 uint64 values (256 bits) at a time.
//
//go:noescape
func notBitmapAVX2(dst, src *uint64, length int)

// popCountAVX2 counts set bits in a bitmap using AVX2.
// Returns the total number of 1 bits across all uint64 values.
//
//go:noescape
func popCountAVX2(bitmap *uint64, length int) int
