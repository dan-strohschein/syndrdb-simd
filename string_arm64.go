//go:build arm64
// +build arm64

package syndrdbsimd

// strEqNEON compares two strings for equality using NEON.
// Returns 1 if equal, 0 if not equal.
//
//go:noescape
func strEqNEON(a, b *byte, length int) int

// strPrefixCmpNEON checks if string starts with prefix using NEON.
// Returns 1 if str has prefix, 0 otherwise.
//
//go:noescape
func strPrefixCmpNEON(str, prefix *byte, strLen, prefixLen int) int

// Note: strToLowerNEON and strToUpperNEON are disabled
// Go's ARM64 assembler doesn't support vector comparison instructions (VCMGE, VCMGT)
// needed for case conversion range checks. Using generic implementations instead.
