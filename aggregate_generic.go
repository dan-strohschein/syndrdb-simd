package syndrdbsimd

import "math"

// sumInt64Generic computes the sum of int64 values using scalar operations.
func sumInt64Generic(values []int64) int64 {
	sum := int64(0)
	for _, v := range values {
		sum += v
	}
	return sum
}

// minInt64Generic finds the minimum int64 value using scalar operations.
// Returns math.MaxInt64 if the slice is empty.
func minInt64Generic(values []int64) int64 {
	if len(values) == 0 {
		return math.MaxInt64
	}
	
	min := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] < min {
			min = values[i]
		}
	}
	return min
}

// maxInt64Generic finds the maximum int64 value using scalar operations.
// Returns math.MinInt64 if the slice is empty.
func maxInt64Generic(values []int64) int64 {
	if len(values) == 0 {
		return math.MinInt64
	}
	
	max := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] > max {
			max = values[i]
		}
	}
	return max
}

// countNonNullGeneric counts non-null values in a slice.
// In Go, we simulate null with a separate bitmap where bit i indicates if values[i] is null.
// Returns the count of values where the corresponding bit in nullBitmap is 0 (not null).
func countNonNullGeneric(values []int64, nullBitmap []uint64) int64 {
	if len(values) == 0 {
		return 0
	}
	
	count := int64(0)
	for i := range values {
		// Check if bit i in nullBitmap is 0 (not null)
		wordIdx := i / 64
		bitIdx := uint(i % 64)
		
		if wordIdx >= len(nullBitmap) {
			// No null bitmap provided for this index, count as non-null
			count++
		} else {
			isNull := (nullBitmap[wordIdx] & (1 << bitIdx)) != 0
			if !isNull {
				count++
			}
		}
	}
	return count
}

// avgInt64Generic computes the average of int64 values.
// Returns 0 if the slice is empty.
func avgInt64Generic(values []int64) float64 {
	if len(values) == 0 {
		return 0
	}
	
	sum := sumInt64Generic(values)
	return float64(sum) / float64(len(values))
}
