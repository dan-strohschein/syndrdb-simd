//go:build arm64

package syndrdbsimd

// For now, ARM64 string comparisons use the existing single-string functions
// in a loop. Future optimization could implement bulk NEON for multiple strings.
// The adaptive threshold routing in impl_arm64.go will decide when to use these.
