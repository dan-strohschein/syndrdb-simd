//go:build arm64

package syndrdbsimd

import (
	"golang.org/x/sys/cpu"
)

var (
	hasNEON bool
	hasSVE  bool
)

func init() {
	// NEON is always available on ARM64
	hasNEON = true
	hasSVE = cpu.ARM64.HasSVE
}

// HasNEON returns true if the CPU supports NEON instructions.
// On ARM64, this is always true.
func HasNEON() bool {
	return hasNEON
}

// HasSVE returns true if the CPU supports Scalable Vector Extension.
func HasSVE() bool {
	return hasSVE
}
