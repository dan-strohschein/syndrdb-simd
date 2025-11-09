package syndrdbsimd

import (
	"bytes"
)

// strCmpGeneric compares two byte slices lexicographically.
// Returns:
//
//	-1 if a < b
//	 0 if a == b
//	+1 if a > b
func strCmpGeneric(a, b []byte) int {
	return bytes.Compare(a, b)
}

// strCmpSliceGeneric compares multiple pairs of strings.
// Both slices must have the same length.
// Results are written to output: -1, 0, or 1 for each pair.
func strCmpSliceGeneric(a, b [][]byte, output []int8) {
	for i := range a {
		output[i] = int8(bytes.Compare(a[i], b[i]))
	}
}

// strLenGeneric returns the length of a byte slice.
func strLenGeneric(s []byte) int {
	return len(s)
}

// strLenSliceGeneric computes lengths for multiple strings.
// Results are written to output.
func strLenSliceGeneric(strings [][]byte, output []int) {
	for i, s := range strings {
		output[i] = len(s)
	}
}

// strPrefixCmpGeneric checks if 'str' starts with 'prefix'.
// Returns true if str has prefix, false otherwise.
func strPrefixCmpGeneric(str, prefix []byte) bool {
	return bytes.HasPrefix(str, prefix)
}

// strPrefixCmpSliceGeneric checks prefix match for multiple strings.
// All strings in 'strings' are compared against the same 'prefix'.
// Results are written to output as booleans.
func strPrefixCmpSliceGeneric(strings [][]byte, prefix []byte, output []bool) {
	for i, s := range strings {
		output[i] = bytes.HasPrefix(s, prefix)
	}
}

// strContainsGeneric checks if 'str' contains 'substr'.
// Returns true if substr is found in str, false otherwise.
func strContainsGeneric(str, substr []byte) bool {
	return bytes.Contains(str, substr)
}

// strContainsSliceGeneric checks substring match for multiple strings.
// All strings in 'strings' are checked for containing 'substr'.
// Results are written to output as booleans.
func strContainsSliceGeneric(strings [][]byte, substr []byte, output []bool) {
	for i, s := range strings {
		output[i] = bytes.Contains(s, substr)
	}
}

// strEqGeneric checks if two byte slices are equal.
// Returns true if equal, false otherwise.
func strEqGeneric(a, b []byte) bool {
	return bytes.Equal(a, b)
}

// strEqSliceGeneric checks equality for multiple string pairs.
// Both slices must have the same length.
// Results are written to output as booleans.
func strEqSliceGeneric(a, b [][]byte, output []bool) {
	for i := range a {
		output[i] = bytes.Equal(a[i], b[i])
	}
}

// strToLowerGeneric converts a byte slice to lowercase (ASCII only).
// The input slice is modified in-place.
func strToLowerGeneric(s []byte) {
	for i := 0; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			s[i] += 32 // Convert to lowercase
		}
	}
}

// strToUpperGeneric converts a byte slice to uppercase (ASCII only).
// The input slice is modified in-place.
func strToUpperGeneric(s []byte) {
	for i := 0; i < len(s); i++ {
		if s[i] >= 'a' && s[i] <= 'z' {
			s[i] -= 32 // Convert to uppercase
		}
	}
}

// strEqIgnoreCaseGeneric checks if two byte slices are equal, ignoring case (ASCII only).
func strEqIgnoreCaseGeneric(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca := a[i]
		cb := b[i]
		// Convert to lowercase for comparison
		if ca >= 'A' && ca <= 'Z' {
			ca += 32
		}
		if cb >= 'A' && cb <= 'Z' {
			cb += 32
		}
		if ca != cb {
			return false
		}
	}
	return true
}
