package syndrdbsimd

// cmpGtFloat64Generic performs element-wise greater-than comparison using scalar operations.
// Returns a slice of booleans where true indicates values[i] > threshold.
// NaN comparisons always return false per IEEE 754.
func cmpGtFloat64Generic(values []float64, threshold float64) []bool {
	results := make([]bool, len(values))
	for i, v := range values {
		results[i] = v > threshold
	}
	return results
}

// cmpGtFloat64MaskGeneric performs element-wise greater-than comparison using scalar operations.
// Returns a bitmask where bit i is set if values[i] > threshold.
func cmpGtFloat64MaskGeneric(values []float64, threshold float64) []uint64 {
	bools := cmpGtFloat64Generic(values, threshold)
	return boolsToBitmask(bools)
}

// cmpGeFloat64Generic performs element-wise greater-than-or-equal comparison using scalar operations.
// Returns a slice of booleans where true indicates values[i] >= threshold.
// NaN comparisons always return false per IEEE 754.
func cmpGeFloat64Generic(values []float64, threshold float64) []bool {
	results := make([]bool, len(values))
	for i, v := range values {
		results[i] = v >= threshold
	}
	return results
}

// cmpGeFloat64MaskGeneric performs element-wise greater-than-or-equal comparison using scalar operations.
// Returns a bitmask where bit i is set if values[i] >= threshold.
func cmpGeFloat64MaskGeneric(values []float64, threshold float64) []uint64 {
	bools := cmpGeFloat64Generic(values, threshold)
	return boolsToBitmask(bools)
}

// cmpLtFloat64Generic performs element-wise less-than comparison using scalar operations.
// Returns a slice of booleans where true indicates values[i] < threshold.
// NaN comparisons always return false per IEEE 754.
func cmpLtFloat64Generic(values []float64, threshold float64) []bool {
	results := make([]bool, len(values))
	for i, v := range values {
		results[i] = v < threshold
	}
	return results
}

// cmpLtFloat64MaskGeneric performs element-wise less-than comparison using scalar operations.
// Returns a bitmask where bit i is set if values[i] < threshold.
func cmpLtFloat64MaskGeneric(values []float64, threshold float64) []uint64 {
	bools := cmpLtFloat64Generic(values, threshold)
	return boolsToBitmask(bools)
}

// cmpLeFloat64Generic performs element-wise less-than-or-equal comparison using scalar operations.
// Returns a slice of booleans where true indicates values[i] <= threshold.
// NaN comparisons always return false per IEEE 754.
func cmpLeFloat64Generic(values []float64, threshold float64) []bool {
	results := make([]bool, len(values))
	for i, v := range values {
		results[i] = v <= threshold
	}
	return results
}

// cmpLeFloat64MaskGeneric performs element-wise less-than-or-equal comparison using scalar operations.
// Returns a bitmask where bit i is set if values[i] <= threshold.
func cmpLeFloat64MaskGeneric(values []float64, threshold float64) []uint64 {
	bools := cmpLeFloat64Generic(values, threshold)
	return boolsToBitmask(bools)
}

// cmpEqFloat64Generic performs element-wise equality comparison using scalar operations.
// Returns a slice of booleans where true indicates values[i] == threshold.
// NaN comparisons always return false per IEEE 754.
func cmpEqFloat64Generic(values []float64, threshold float64) []bool {
	results := make([]bool, len(values))
	for i, v := range values {
		results[i] = v == threshold
	}
	return results
}

// cmpEqFloat64MaskGeneric performs element-wise equality comparison using scalar operations.
// Returns a bitmask where bit i is set if values[i] == threshold.
func cmpEqFloat64MaskGeneric(values []float64, threshold float64) []uint64 {
	bools := cmpEqFloat64Generic(values, threshold)
	return boolsToBitmask(bools)
}

// cmpNeFloat64Generic performs element-wise inequality comparison using scalar operations.
// Returns a slice of booleans where true indicates values[i] != threshold.
// NaN != x is true for all x (including NaN) per IEEE 754.
func cmpNeFloat64Generic(values []float64, threshold float64) []bool {
	results := make([]bool, len(values))
	for i, v := range values {
		results[i] = v != threshold
	}
	return results
}

// cmpNeFloat64MaskGeneric performs element-wise inequality comparison using scalar operations.
// Returns a bitmask where bit i is set if values[i] != threshold.
func cmpNeFloat64MaskGeneric(values []float64, threshold float64) []uint64 {
	bools := cmpNeFloat64Generic(values, threshold)
	return boolsToBitmask(bools)
}
