// +build amd64

package syndrdbsimd

// strEqAVX2 compares two strings for equality using AVX2.
// Returns 1 if equal, 0 if not equal.
//
//go:noescape
func strEqAVX2(a, b *byte, length int) int

// strPrefixCmpAVX2 checks if string starts with prefix using AVX2.
// Returns 1 if str has prefix, 0 otherwise.
//
//go:noescape
func strPrefixCmpAVX2(str, prefix *byte, strLen, prefixLen int) int

// strToLowerAVX2 converts ASCII string to lowercase using AVX2.
// Modifies the string in-place.
//
//go:noescape
func strToLowerAVX2(s *byte, length int)

// strToUpperAVX2 converts ASCII string to uppercase using AVX2.
// Modifies the string in-place.
//
//go:noescape
func strToUpperAVX2(s *byte, length int)
