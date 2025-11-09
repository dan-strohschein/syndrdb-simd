package syndrdbsimd

// cmpEqInt64Generic performs element-wise equality comparison using scalar operations.
// Returns a slice of booleans where true indicates values[i] == threshold.
func cmpEqInt64Generic(values []int64, threshold int64) []bool {
	results := make([]bool, len(values))
	for i, v := range values {
		results[i] = v == threshold
	}
	return results
}

// cmpEqInt64MaskGeneric performs element-wise equality comparison using scalar operations.
// Returns a bitmask where bit i is set if values[i] == threshold.
// Each uint64 in the result holds 64 comparison results.
func cmpEqInt64MaskGeneric(values []int64, threshold int64) []uint64 {
	// TODO: I could optimize this further by processing 64 values at a time directly into uint64
	bools := cmpEqInt64Generic(values, threshold)
	return boolsToBitmask(bools)
}

// cmpNeInt64Generic performs element-wise inequality comparison using scalar operations.
// Returns a slice of booleans where true indicates values[i] != threshold.
func cmpNeInt64Generic(values []int64, threshold int64) []bool {
	results := make([]bool, len(values))
	for i, v := range values {
		results[i] = v != threshold
	}
	return results
}

// cmpNeInt64MaskGeneric performs element-wise inequality comparison using scalar operations.
// Returns a bitmask where bit i is set if values[i] != threshold.
func cmpNeInt64MaskGeneric(values []int64, threshold int64) []uint64 {
	bools := cmpNeInt64Generic(values, threshold)
	return boolsToBitmask(bools)
}

// cmpGtInt64Generic performs element-wise greater-than comparison using scalar operations.
// Returns a slice of booleans where true indicates values[i] > threshold.
func cmpGtInt64Generic(values []int64, threshold int64) []bool {
	results := make([]bool, len(values))
	for i, v := range values {
		results[i] = v > threshold
	}
	return results
}

// cmpGtInt64MaskGeneric performs element-wise greater-than comparison using scalar operations.
// Returns a bitmask where bit i is set if values[i] > threshold.
func cmpGtInt64MaskGeneric(values []int64, threshold int64) []uint64 {
	bools := cmpGtInt64Generic(values, threshold)
	return boolsToBitmask(bools)
}

// cmpLtInt64Generic performs element-wise less-than comparison using scalar operations.
// Returns a slice of booleans where true indicates values[i] < threshold.
func cmpLtInt64Generic(values []int64, threshold int64) []bool {
	results := make([]bool, len(values))
	for i, v := range values {
		results[i] = v < threshold
	}
	return results
}

// cmpLtInt64MaskGeneric performs element-wise less-than comparison using scalar operations.
// Returns a bitmask where bit i is set if values[i] < threshold.
func cmpLtInt64MaskGeneric(values []int64, threshold int64) []uint64 {
	bools := cmpLtInt64Generic(values, threshold)
	return boolsToBitmask(bools)
}

// cmpGeInt64Generic performs element-wise greater-than-or-equal comparison using scalar operations.
// Returns a slice of booleans where true indicates values[i] >= threshold.
func cmpGeInt64Generic(values []int64, threshold int64) []bool {
	results := make([]bool, len(values))
	for i, v := range values {
		results[i] = v >= threshold
	}
	return results
}

// cmpGeInt64MaskGeneric performs element-wise greater-than-or-equal comparison using scalar operations.
// Returns a bitmask where bit i is set if values[i] >= threshold.
func cmpGeInt64MaskGeneric(values []int64, threshold int64) []uint64 {
	bools := cmpGeInt64Generic(values, threshold)
	return boolsToBitmask(bools)
}

// cmpLeInt64Generic performs element-wise less-than-or-equal comparison using scalar operations.
// Returns a slice of booleans where true indicates values[i] <= threshold.
func cmpLeInt64Generic(values []int64, threshold int64) []bool {
	results := make([]bool, len(values))
	for i, v := range values {
		results[i] = v <= threshold
	}
	return results
}

// cmpLeInt64MaskGeneric performs element-wise less-than-or-equal comparison using scalar operations.
// Returns a bitmask where bit i is set if values[i] <= threshold.
func cmpLeInt64MaskGeneric(values []int64, threshold int64) []uint64 {
	bools := cmpLeInt64Generic(values, threshold)
	return boolsToBitmask(bools)
}
