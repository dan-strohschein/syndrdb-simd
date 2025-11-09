//go:build !amd64 && !arm64

package syndrdbsimd

// HasAVX2 returns false on non-amd64 architectures.
func HasAVX2() bool {
	return false
}

// HasAVX512 returns false on non-amd64 architectures.
func HasAVX512() bool {
	return false
}

// HasSSE42 returns false on non-amd64 architectures.
func HasSSE42() bool {
	return false
}

// HasNEON returns false on non-arm64 architectures.
func HasNEON() bool {
	return false
}

// HasSVE returns false on non-arm64 architectures.
func HasSVE() bool {
	return false
}
