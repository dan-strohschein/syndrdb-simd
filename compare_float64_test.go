package syndrdbsimd

import (
	"math"
	"testing"
)

// TestCmpGtFloat64 tests greater-than comparisons for float64 values
func TestCmpGtFloat64(t *testing.T) {
	tests := []struct {
		name      string
		values    []float64
		threshold float64
		expected  []bool
	}{
		{
			name:      "empty slice",
			values:    []float64{},
			threshold: 5.0,
			expected:  []bool{},
		},
		{
			name:      "single element - true",
			values:    []float64{10.0},
			threshold: 5.0,
			expected:  []bool{true},
		},
		{
			name:      "single element - false",
			values:    []float64{3.0},
			threshold: 5.0,
			expected:  []bool{false},
		},
		{
			name:      "small array - mixed results",
			values:    []float64{1.0, 5.0, 10.0, 3.0, 8.0},
			threshold: 5.0,
			expected:  []bool{false, false, true, false, true},
		},
		{
			name:      "SIMD threshold boundary - 15 elements",
			values:    []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			threshold: 8.0,
			expected:  []bool{false, false, false, false, false, false, false, false, true, true, true, true, true, true, true},
		},
		{
			name:      "SIMD threshold boundary - 16 elements",
			values:    []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			threshold: 8.0,
			expected:  []bool{false, false, false, false, false, false, false, false, true, true, true, true, true, true, true, true},
		},
		{
			name:      "large array - 100 elements",
			values:    makeFloat64Range(1.0, 100.0, 1.0),
			threshold: 50.0,
			expected:  makeBoolRange(100, func(i int) bool { return float64(i+1) > 50.0 }),
		},
		{
			name:      "negative values",
			values:    []float64{-10.0, -5.0, 0.0, 5.0, 10.0},
			threshold: 0.0,
			expected:  []bool{false, false, false, true, true},
		},
		{
			name:      "NaN comparisons - always false",
			values:    []float64{math.NaN(), 5.0, math.NaN(), 10.0},
			threshold: 5.0,
			expected:  []bool{false, false, false, true},
		},
		{
			name:      "NaN threshold - all false",
			values:    []float64{1.0, 5.0, 10.0, math.NaN()},
			threshold: math.NaN(),
			expected:  []bool{false, false, false, false},
		},
		{
			name:      "infinity values",
			values:    []float64{math.Inf(-1), -1000.0, 0.0, 1000.0, math.Inf(1)},
			threshold: 0.0,
			expected:  []bool{false, false, false, true, true},
		},
		{
			name:      "infinity threshold",
			values:    []float64{-1000.0, 0.0, 1000.0, math.Inf(1)},
			threshold: math.Inf(1),
			expected:  []bool{false, false, false, false},
		},
		{
			name:      "zero comparisons",
			values:    []float64{-0.0, 0.0, 1e-308, -1e-308},
			threshold: 0.0,
			expected:  []bool{false, false, true, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CmpGtFloat64(tt.values, tt.threshold)
			if len(result) != len(tt.expected) {
				t.Fatalf("length mismatch: got %d, want %d", len(result), len(tt.expected))
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("index %d: got %v, want %v (value=%v, threshold=%v)",
						i, result[i], tt.expected[i], tt.values[i], tt.threshold)
				}
			}
		})
	}
}

// TestCmpGtFloat64Mask tests the bitmask variant of greater-than comparison
func TestCmpGtFloat64Mask(t *testing.T) {
	tests := []struct {
		name      string
		values    []float64
		threshold float64
	}{
		{
			name:      "empty slice",
			values:    []float64{},
			threshold: 5.0,
		},
		{
			name:      "64 elements - full word",
			values:    makeFloat64Range(1.0, 64.0, 1.0),
			threshold: 32.0,
		},
		{
			name:      "65 elements - multiple words",
			values:    makeFloat64Range(1.0, 65.0, 1.0),
			threshold: 32.0,
		},
		{
			name:      "128 elements - two full words",
			values:    makeFloat64Range(1.0, 128.0, 1.0),
			threshold: 64.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			boolResult := CmpGtFloat64(tt.values, tt.threshold)
			maskResult := CmpGtFloat64Mask(tt.values, tt.threshold)

			// Convert mask back to bools and compare
			convertedBools := BitmaskToBools(maskResult, len(tt.values))
			if len(convertedBools) != len(boolResult) {
				t.Fatalf("length mismatch: got %d, want %d", len(convertedBools), len(boolResult))
			}
			for i := range convertedBools {
				if convertedBools[i] != boolResult[i] {
					t.Errorf("index %d: mask gave %v, bool gave %v", i, convertedBools[i], boolResult[i])
				}
			}
		})
	}
}

// TestCmpGeFloat64 tests greater-or-equal comparisons
func TestCmpGeFloat64(t *testing.T) {
	tests := []struct {
		name      string
		values    []float64
		threshold float64
		expected  []bool
	}{
		{
			name:      "equality boundary",
			values:    []float64{4.0, 5.0, 6.0},
			threshold: 5.0,
			expected:  []bool{false, true, true},
		},
		{
			name:      "NaN - always false",
			values:    []float64{math.NaN(), 5.0, 10.0},
			threshold: 5.0,
			expected:  []bool{false, true, true},
		},
		{
			name:      "large array",
			values:    makeFloat64Range(1.0, 100.0, 1.0),
			threshold: 50.0,
			expected:  makeBoolRange(100, func(i int) bool { return float64(i+1) >= 50.0 }),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CmpGeFloat64(tt.values, tt.threshold)
			if len(result) != len(tt.expected) {
				t.Fatalf("length mismatch: got %d, want %d", len(result), len(tt.expected))
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("index %d: got %v, want %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

// TestCmpLtFloat64 tests less-than comparisons
func TestCmpLtFloat64(t *testing.T) {
	tests := []struct {
		name      string
		values    []float64
		threshold float64
		expected  []bool
	}{
		{
			name:      "basic comparison",
			values:    []float64{1.0, 5.0, 10.0, 3.0, 8.0},
			threshold: 5.0,
			expected:  []bool{true, false, false, true, false},
		},
		{
			name:      "NaN - always false",
			values:    []float64{math.NaN(), 5.0, 1.0},
			threshold: 5.0,
			expected:  []bool{false, false, true},
		},
		{
			name:      "infinity comparisons",
			values:    []float64{math.Inf(-1), -1000.0, 0.0, 1000.0, math.Inf(1)},
			threshold: 0.0,
			expected:  []bool{true, true, false, false, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CmpLtFloat64(tt.values, tt.threshold)
			if len(result) != len(tt.expected) {
				t.Fatalf("length mismatch: got %d, want %d", len(result), len(tt.expected))
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("index %d: got %v, want %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

// TestCmpLeFloat64 tests less-or-equal comparisons
func TestCmpLeFloat64(t *testing.T) {
	tests := []struct {
		name      string
		values    []float64
		threshold float64
		expected  []bool
	}{
		{
			name:      "equality boundary",
			values:    []float64{4.0, 5.0, 6.0},
			threshold: 5.0,
			expected:  []bool{true, true, false},
		},
		{
			name:      "NaN - always false",
			values:    []float64{1.0, math.NaN(), 10.0},
			threshold: 5.0,
			expected:  []bool{true, false, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CmpLeFloat64(tt.values, tt.threshold)
			if len(result) != len(tt.expected) {
				t.Fatalf("length mismatch: got %d, want %d", len(result), len(tt.expected))
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("index %d: got %v, want %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

// TestCmpEqFloat64 tests equality comparisons
func TestCmpEqFloat64(t *testing.T) {
	tests := []struct {
		name      string
		values    []float64
		threshold float64
		expected  []bool
	}{
		{
			name:      "basic equality",
			values:    []float64{1.0, 5.0, 5.0, 10.0},
			threshold: 5.0,
			expected:  []bool{false, true, true, false},
		},
		{
			name:      "NaN equality - always false (even NaN == NaN)",
			values:    []float64{math.NaN(), 5.0, math.NaN()},
			threshold: math.NaN(),
			expected:  []bool{false, false, false},
		},
		{
			name:      "zero equality",
			values:    []float64{-0.0, 0.0, 1e-308},
			threshold: 0.0,
			expected:  []bool{true, true, false},
		},
		{
			name:      "infinity equality",
			values:    []float64{math.Inf(1), math.Inf(-1), 1e308},
			threshold: math.Inf(1),
			expected:  []bool{true, false, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CmpEqFloat64(tt.values, tt.threshold)
			if len(result) != len(tt.expected) {
				t.Fatalf("length mismatch: got %d, want %d", len(result), len(tt.expected))
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("index %d: got %v, want %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

// TestCmpNeFloat64 tests inequality comparisons
func TestCmpNeFloat64(t *testing.T) {
	tests := []struct {
		name      string
		values    []float64
		threshold float64
		expected  []bool
	}{
		{
			name:      "basic inequality",
			values:    []float64{1.0, 5.0, 5.0, 10.0},
			threshold: 5.0,
			expected:  []bool{true, false, false, true},
		},
		{
			name:      "NaN inequality - always true (even NaN != NaN)",
			values:    []float64{math.NaN(), 5.0, math.NaN()},
			threshold: math.NaN(),
			expected:  []bool{true, true, true},
		},
		{
			name:      "NaN vs value - always true",
			values:    []float64{math.NaN(), 5.0, 10.0},
			threshold: 5.0,
			expected:  []bool{true, false, true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CmpNeFloat64(tt.values, tt.threshold)
			if len(result) != len(tt.expected) {
				t.Fatalf("length mismatch: got %d, want %d", len(result), len(tt.expected))
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("index %d: got %v, want %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

// Helper functions

// makeFloat64Range creates a slice of float64 values from start to end with a given step
func makeFloat64Range(start, end, step float64) []float64 {
	count := int((end-start)/step) + 1
	result := make([]float64, count)
	for i := 0; i < count; i++ {
		result[i] = start + float64(i)*step
	}
	return result
}

// makeBoolRange creates a slice of booleans using a function
func makeBoolRange(count int, fn func(i int) bool) []bool {
	result := make([]bool, count)
	for i := 0; i < count; i++ {
		result[i] = fn(i)
	}
	return result
}
