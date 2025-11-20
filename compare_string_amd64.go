//go:build amd64

package syndrdbsimd

// The adaptive threshold routing in impl_amd64.go will decide when to use these.// in a loop. Future optimization could implement bulk SIMD for multiple strings.// For now, AMD64 string comparisons use the existing single-string SIMD functions
