package syndrdbsimd

import (
	"math"
	"testing"
)

// ============================================================================
// SumInt64 Tests
// ============================================================================

// Happy path tests
func TestSumInt64_HappyPath(t *testing.T) {
	values := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	result := SumInt64(values)

	expected := int64(55) // 1+2+3+...+10 = 55
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

func TestSumInt64_LargeArray(t *testing.T) {
	values := make([]int64, 1000)
	for i := range values {
		values[i] = int64(i + 1)
	}

	result := SumInt64(values)

	// Sum of 1+2+...+1000 = 1000*1001/2 = 500500
	expected := int64(500500)
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

// Edge case: empty array
func TestSumInt64_EmptyArray(t *testing.T) {
	values := []int64{}

	result := SumInt64(values)

	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}
}

// Edge case: single element
func TestSumInt64_SingleElement(t *testing.T) {
	values := []int64{42}

	result := SumInt64(values)

	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
}

// Edge case: odd length
func TestSumInt64_OddLength(t *testing.T) {
	values := []int64{1, 2, 3, 4, 5, 6, 7}

	result := SumInt64(values)

	expected := int64(28) // 1+2+3+4+5+6+7 = 28
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

// Corner case: negative numbers
func TestSumInt64_NegativeNumbers(t *testing.T) {
	values := []int64{-5, -3, -1, 0, 1, 3, 5}

	result := SumInt64(values)

	expected := int64(0) // Sum cancels out
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

// Corner case: all zeros
func TestSumInt64_AllZeros(t *testing.T) {
	values := make([]int64, 100)

	result := SumInt64(values)

	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}
}

// SIMD boundary tests
func TestSumInt64_Length16_AVX2Boundary(t *testing.T) {
	values := make([]int64, 16)
	for i := range values {
		values[i] = int64(i + 1)
	}

	result := SumInt64(values)

	expected := int64(136) // 1+2+...+16 = 136
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

func TestSumInt64_Length8_NEONBoundary(t *testing.T) {
	values := []int64{1, 2, 3, 4, 5, 6, 7, 8}

	result := SumInt64(values)

	expected := int64(36) // 1+2+...+8 = 36
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

// ============================================================================
// MinInt64 Tests
// ============================================================================

func TestMinInt64_HappyPath(t *testing.T) {
	values := []int64{10, 3, 7, 1, 9, 2, 8, 4, 6, 5}

	result := MinInt64(values)

	if result != 1 {
		t.Errorf("Expected 1, got %d", result)
	}
}

func TestMinInt64_AllSame(t *testing.T) {
	values := []int64{7, 7, 7, 7, 7}

	result := MinInt64(values)

	if result != 7 {
		t.Errorf("Expected 7, got %d", result)
	}
}

func TestMinInt64_EmptyArray(t *testing.T) {
	values := []int64{}

	result := MinInt64(values)

	if result != math.MaxInt64 {
		t.Errorf("Expected math.MaxInt64 (%d), got %d", int64(math.MaxInt64), result)
	}
}

func TestMinInt64_SingleElement(t *testing.T) {
	values := []int64{42}

	result := MinInt64(values)

	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
}

func TestMinInt64_NegativeNumbers(t *testing.T) {
	values := []int64{-5, -100, -3, -1, -50}

	result := MinInt64(values)

	if result != -100 {
		t.Errorf("Expected -100, got %d", result)
	}
}

func TestMinInt64_MinAtEnd(t *testing.T) {
	values := []int64{10, 20, 30, 40, 50, 60, 70, 5}

	result := MinInt64(values)

	if result != 5 {
		t.Errorf("Expected 5, got %d", result)
	}
}

func TestMinInt64_MinAtStart(t *testing.T) {
	values := []int64{1, 10, 20, 30, 40, 50, 60, 70}

	result := MinInt64(values)

	if result != 1 {
		t.Errorf("Expected 1, got %d", result)
	}
}

func TestMinInt64_LargeArray(t *testing.T) {
	values := make([]int64, 1000)
	for i := range values {
		values[i] = int64(i + 1)
	}
	values[500] = -100 // Insert a min value

	result := MinInt64(values)

	if result != -100 {
		t.Errorf("Expected -100, got %d", result)
	}
}

// ============================================================================
// MaxInt64 Tests
// ============================================================================

func TestMaxInt64_HappyPath(t *testing.T) {
	values := []int64{10, 3, 7, 100, 9, 2, 8, 4, 6, 5}

	result := MaxInt64(values)

	if result != 100 {
		t.Errorf("Expected 100, got %d", result)
	}
}

func TestMaxInt64_AllSame(t *testing.T) {
	values := []int64{-5, -5, -5, -5, -5}

	result := MaxInt64(values)

	if result != -5 {
		t.Errorf("Expected -5, got %d", result)
	}
}

func TestMaxInt64_EmptyArray(t *testing.T) {
	values := []int64{}

	result := MaxInt64(values)

	if result != math.MinInt64 {
		t.Errorf("Expected math.MinInt64 (%d), got %d", int64(math.MinInt64), result)
	}
}

func TestMaxInt64_SingleElement(t *testing.T) {
	values := []int64{-99}

	result := MaxInt64(values)

	if result != -99 {
		t.Errorf("Expected -99, got %d", result)
	}
}

func TestMaxInt64_NegativeNumbers(t *testing.T) {
	values := []int64{-5, -100, -3, -1, -50}

	result := MaxInt64(values)

	if result != -1 {
		t.Errorf("Expected -1, got %d", result)
	}
}

func TestMaxInt64_MaxAtEnd(t *testing.T) {
	values := []int64{10, 20, 30, 40, 50, 60, 70, 500}

	result := MaxInt64(values)

	if result != 500 {
		t.Errorf("Expected 500, got %d", result)
	}
}

func TestMaxInt64_MaxAtStart(t *testing.T) {
	values := []int64{999, 10, 20, 30, 40, 50, 60, 70}

	result := MaxInt64(values)

	if result != 999 {
		t.Errorf("Expected 999, got %d", result)
	}
}

func TestMaxInt64_LargeArray(t *testing.T) {
	values := make([]int64, 1000)
	for i := range values {
		values[i] = int64(i + 1)
	}
	values[500] = 10000 // Insert a max value

	result := MaxInt64(values)

	if result != 10000 {
		t.Errorf("Expected 10000, got %d", result)
	}
}

// ============================================================================
// CountNonNull Tests
// ============================================================================

func TestCountNonNull_AllNonNull(t *testing.T) {
	values := []int64{1, 2, 3, 4, 5, 6, 7, 8}
	nullBitmap := []uint64{} // No nulls

	result := CountNonNull(values, nullBitmap)

	if result != 8 {
		t.Errorf("Expected 8, got %d", result)
	}
}

func TestCountNonNull_SomeNulls(t *testing.T) {
	values := []int64{1, 2, 3, 4, 5, 6, 7, 8}
	// Bits: 0b01010101 = positions 0, 2, 4, 6 are null
	nullBitmap := []uint64{0b01010101}

	result := CountNonNull(values, nullBitmap)

	expected := int64(4) // Positions 1, 3, 5, 7 are not null
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

func TestCountNonNull_AllNull(t *testing.T) {
	values := []int64{1, 2, 3, 4, 5, 6, 7, 8}
	nullBitmap := []uint64{0xFF} // First 8 bits set = all null

	result := CountNonNull(values, nullBitmap)

	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}
}

func TestCountNonNull_EmptyArray(t *testing.T) {
	values := []int64{}
	nullBitmap := []uint64{}

	result := CountNonNull(values, nullBitmap)

	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}
}

func TestCountNonNull_SingleNonNull(t *testing.T) {
	values := []int64{42}
	nullBitmap := []uint64{0} // Bit 0 = 0 = not null

	result := CountNonNull(values, nullBitmap)

	if result != 1 {
		t.Errorf("Expected 1, got %d", result)
	}
}

func TestCountNonNull_SingleNull(t *testing.T) {
	values := []int64{42}
	nullBitmap := []uint64{1} // Bit 0 = 1 = null

	result := CountNonNull(values, nullBitmap)

	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}
}

func TestCountNonNull_LargeArray(t *testing.T) {
	values := make([]int64, 100)
	// Set every other bit as null
	nullBitmap := []uint64{
		0xAAAAAAAAAAAAAAAA, // 64 bits: alternating 1010...
		0xAAAAAAAAAAAAAAAA, // Next 36 bits (only 36 used)
	}

	result := CountNonNull(values, nullBitmap)

	expected := int64(50) // Half are non-null
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

func TestCountNonNull_PartialBitmap(t *testing.T) {
	values := make([]int64, 100)
	// Only first 64 bits defined, rest are implicitly not null
	nullBitmap := []uint64{0xFFFFFFFFFFFFFFFF} // First 64 are null

	result := CountNonNull(values, nullBitmap)

	expected := int64(36) // Last 36 elements are non-null
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

// ============================================================================
// AvgInt64 Tests
// ============================================================================

func TestAvgInt64_HappyPath(t *testing.T) {
	values := []int64{1, 2, 3, 4, 5}

	result := AvgInt64(values)

	expected := 3.0 // (1+2+3+4+5)/5 = 15/5 = 3.0
	if result != expected {
		t.Errorf("Expected %f, got %f", expected, result)
	}
}

func TestAvgInt64_EmptyArray(t *testing.T) {
	values := []int64{}

	result := AvgInt64(values)

	if result != 0.0 {
		t.Errorf("Expected 0.0, got %f", result)
	}
}

func TestAvgInt64_SingleElement(t *testing.T) {
	values := []int64{42}

	result := AvgInt64(values)

	if result != 42.0 {
		t.Errorf("Expected 42.0, got %f", result)
	}
}

func TestAvgInt64_NegativeNumbers(t *testing.T) {
	values := []int64{-10, -5, 0, 5, 10}

	result := AvgInt64(values)

	expected := 0.0 // Sum = 0, avg = 0
	if result != expected {
		t.Errorf("Expected %f, got %f", expected, result)
	}
}

func TestAvgInt64_LargeNumbers(t *testing.T) {
	values := []int64{1000, 2000, 3000, 4000, 5000}

	result := AvgInt64(values)

	expected := 3000.0 // (1000+2000+3000+4000+5000)/5 = 15000/5 = 3000
	if result != expected {
		t.Errorf("Expected %f, got %f", expected, result)
	}
}

// ============================================================================
// Integration Tests (combining multiple operations)
// ============================================================================

func TestAggregations_Combined(t *testing.T) {
	values := []int64{5, 2, 9, 1, 7, 3, 8, 4, 6}

	sum := SumInt64(values)
	min := MinInt64(values)
	max := MaxInt64(values)
	avg := AvgInt64(values)

	if sum != 45 {
		t.Errorf("Sum: expected 45, got %d", sum)
	}
	if min != 1 {
		t.Errorf("Min: expected 1, got %d", min)
	}
	if max != 9 {
		t.Errorf("Max: expected 9, got %d", max)
	}
	if avg != 5.0 {
		t.Errorf("Avg: expected 5.0, got %f", avg)
	}
}

func TestAggregations_WithNulls(t *testing.T) {
	values := []int64{10, 20, 30, 40, 50, 60, 70, 80}
	// Nulls at positions 1, 3, 5 (bits 1, 3, 5 set)
	nullBitmap := []uint64{0b00101010}

	count := CountNonNull(values, nullBitmap)

	// Positions 0, 2, 4, 6, 7 are non-null = 5 values
	expected := int64(5)
	if count != expected {
		t.Errorf("Expected %d non-null, got %d", expected, count)
	}
}
