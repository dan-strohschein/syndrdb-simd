package syndrdbsimd

import (
	"unsafe"
)

// stringToBytes converts a string to []byte without allocation.
// This is safe because Go strings are immutable and we only use read-only access.
// Uses unsafe.Slice and unsafe.StringData which are standard as of Go 1.20+.
func stringToBytes(s string) []byte {
	if len(s) == 0 {
		return []byte{}
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// stringsToBytes converts []string to [][]byte without allocation.
// Uses unsafe conversion for zero-copy transformation.
func stringsToBytes(strs []string) [][]byte {
	if len(strs) == 0 {
		return [][]byte{}
	}
	
	result := make([][]byte, len(strs))
	for i, s := range strs {
		result[i] = stringToBytes(s)
	}
	return result
}
