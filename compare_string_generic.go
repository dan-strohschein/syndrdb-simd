package syndrdbsimd

import (
	"bytes"
)

// cmpEqStringGeneric performs element-wise equality comparison on strings using scalar operations.
// Returns a slice of booleans where true indicates values[i] == threshold.
func cmpEqStringGeneric(values [][]byte, threshold []byte) []bool {
	results := make([]bool, len(values))
	for i, v := range values {
		results[i] = bytes.Equal(v, threshold)
	}
	return results
}

// cmpEqStringMaskGeneric performs element-wise equality comparison on strings.
// Returns a bitmask where bit i is set if values[i] == threshold.
func cmpEqStringMaskGeneric(values [][]byte, threshold []byte) []uint64 {
	bools := cmpEqStringGeneric(values, threshold)
	return boolsToBitmask(bools)
}

// cmpNeStringGeneric performs element-wise inequality comparison on strings.
// Returns a slice of booleans where true indicates values[i] != threshold.
func cmpNeStringGeneric(values [][]byte, threshold []byte) []bool {
	results := make([]bool, len(values))
	for i, v := range values {
		results[i] = !bytes.Equal(v, threshold)
	}
	return results
}

// cmpNeStringMaskGeneric performs element-wise inequality comparison on strings.
// Returns a bitmask where bit i is set if values[i] != threshold.
func cmpNeStringMaskGeneric(values [][]byte, threshold []byte) []uint64 {
	bools := cmpNeStringGeneric(values, threshold)
	return boolsToBitmask(bools)
}

// cmpHasPrefixStringGeneric checks if strings start with a prefix using scalar operations.
// Returns a slice of booleans where true indicates values[i] has the prefix.
func cmpHasPrefixStringGeneric(values [][]byte, prefix []byte) []bool {
	results := make([]bool, len(values))
	for i, v := range values {
		results[i] = bytes.HasPrefix(v, prefix)
	}
	return results
}

// cmpHasPrefixStringMaskGeneric checks prefix match and returns a bitmask.
func cmpHasPrefixStringMaskGeneric(values [][]byte, prefix []byte) []uint64 {
	bools := cmpHasPrefixStringGeneric(values, prefix)
	return boolsToBitmask(bools)
}

// cmpHasSuffixStringGeneric checks if strings end with a suffix using scalar operations.
// Returns a slice of booleans where true indicates values[i] has the suffix.
func cmpHasSuffixStringGeneric(values [][]byte, suffix []byte) []bool {
	results := make([]bool, len(values))
	for i, v := range values {
		results[i] = bytes.HasSuffix(v, suffix)
	}
	return results
}

// cmpHasSuffixStringMaskGeneric checks suffix match and returns a bitmask.
func cmpHasSuffixStringMaskGeneric(values [][]byte, suffix []byte) []uint64 {
	bools := cmpHasSuffixStringGeneric(values, suffix)
	return boolsToBitmask(bools)
}

// cmpContainsStringGeneric checks if strings contain a substring using scalar operations.
// Returns a slice of booleans where true indicates values[i] contains the substring.
func cmpContainsStringGeneric(values [][]byte, substr []byte) []bool {
	results := make([]bool, len(values))
	for i, v := range values {
		results[i] = bytes.Contains(v, substr)
	}
	return results
}

// cmpContainsStringMaskGeneric checks substring match and returns a bitmask.
func cmpContainsStringMaskGeneric(values [][]byte, substr []byte) []uint64 {
	bools := cmpContainsStringGeneric(values, substr)
	return boolsToBitmask(bools)
}

// cmpEqStringIgnoreCaseGeneric performs case-insensitive equality comparison (ASCII only).
// Non-ASCII bytes (>127) are compared as-is.
func cmpEqStringIgnoreCaseGeneric(values [][]byte, threshold []byte) []bool {
	results := make([]bool, len(values))
	
	// Convert threshold to lowercase once
	thresholdLower := make([]byte, len(threshold))
	copy(thresholdLower, threshold)
	for i := 0; i < len(thresholdLower); i++ {
		if thresholdLower[i] >= 'A' && thresholdLower[i] <= 'Z' {
			thresholdLower[i] += 32
		}
	}
	
	// Compare each value (converted to lowercase) with lowercase threshold
	for i, v := range values {
		if len(v) != len(thresholdLower) {
			results[i] = false
			continue
		}
		
		match := true
		for j := 0; j < len(v); j++ {
			vChar := v[j]
			if vChar >= 'A' && vChar <= 'Z' {
				vChar += 32
			}
			if vChar != thresholdLower[j] {
				match = false
				break
			}
		}
		results[i] = match
	}
	
	return results
}

// cmpEqStringIgnoreCaseMaskGeneric performs case-insensitive equality and returns bitmask.
func cmpEqStringIgnoreCaseMaskGeneric(values [][]byte, threshold []byte) []uint64 {
	bools := cmpEqStringIgnoreCaseGeneric(values, threshold)
	return boolsToBitmask(bools)
}

// matchWildcard matches a single string against a pattern with % and _ wildcards.
// % matches zero or more characters, _ matches exactly one character.
func matchWildcard(str, pattern []byte) bool {
	sLen, pLen := len(str), len(pattern)
	si, pi := 0, 0
	starIdx, matchIdx := -1, 0
	
	for si < sLen {
		if pi < pLen {
			if pattern[pi] == '%' {
				// Found %, mark position and try to match rest
				starIdx = pi
				matchIdx = si
				pi++
				continue
			} else if pattern[pi] == '_' || pattern[pi] == str[si] {
				// _ matches any char, or exact match
				si++
				pi++
				continue
			}
		}
		
		// No match, backtrack to last %
		if starIdx != -1 {
			pi = starIdx + 1
			matchIdx++
			si = matchIdx
			continue
		}
		
		return false
	}
	
	// Consume trailing % in pattern
	for pi < pLen && pattern[pi] == '%' {
		pi++
	}
	
	return pi == pLen
}

// cmpMatchWildcardGeneric matches strings against a wildcard pattern.
func cmpMatchWildcardGeneric(values [][]byte, pattern []byte) []bool {
	results := make([]bool, len(values))
	for i, v := range values {
		results[i] = matchWildcard(v, pattern)
	}
	return results
}

// cmpMatchWildcardMaskGeneric matches wildcard pattern and returns bitmask.
func cmpMatchWildcardMaskGeneric(values [][]byte, pattern []byte) []uint64 {
	bools := cmpMatchWildcardGeneric(values, pattern)
	return boolsToBitmask(bools)
}
