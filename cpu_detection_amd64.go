//go:build amd64

package syndrdbsimd

import (
	"golang.org/x/sys/cpu"
)

var (
	hasAVX2   bool
	hasAVX512 bool
	hasSSE42  bool
)

func init() {
	hasAVX2 = cpu.X86.HasAVX2
	hasAVX512 = cpu.X86.HasAVX512F
	hasSSE42 = cpu.X86.HasSSE42
}

// HasAVX2 returns true if the CPU supports AVX2 instructions.
func HasAVX2() bool {
	return hasAVX2
}

// HasAVX512 returns true if the CPU supports AVX-512 Foundation instructions.
func HasAVX512() bool {
	return hasAVX512
}

// HasSSE42 returns true if the CPU supports SSE4.2 instructions.
func HasSSE42() bool {
	return hasSSE42
}
