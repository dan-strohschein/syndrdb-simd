package syndrdbsimd

import (
	"testing"
)

// Happy path tests - normal operation with various data sizes
func TestCmpEqInt64_HappyPath(t *testing.T) {
	values := []int64{1, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70, 75}
	threshold := int64(35)

	result := CmpEqInt64(values, threshold)

	expected := []bool{false, false, false, false, false, false, false, true, false, false, false, false, false, false, false, false}
	if len(result) != len(expected) {
		t.Fatalf("Expected length %d, got %d", len(expected), len(result))
	}

	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Index %d: expected %v, got %v", i, expected[i], result[i])
		}
	}
}

func TestCmpGtInt64_HappyPath(t *testing.T) {
	values := []int64{10, 20, 30, 40, 50, 60, 70, 80}
	threshold := int64(45)

	result := CmpGtInt64(values, threshold)

	expected := []bool{false, false, false, false, true, true, true, true}
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Index %d: expected %v, got %v", i, expected[i], result[i])
		}
	}
}

func TestCmpLtInt64_HappyPath(t *testing.T) {
	values := []int64{5, 15, 25, 35, 45, 55, 65, 75}
	threshold := int64(40)

	result := CmpLtInt64(values, threshold)

	expected := []bool{true, true, true, true, false, false, false, false}
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Index %d: expected %v, got %v", i, expected[i], result[i])
		}
	}
}

func TestCmpGeInt64_HappyPath(t *testing.T) {
	values := []int64{10, 20, 30, 40, 50, 60}
	threshold := int64(30)

	result := CmpGeInt64(values, threshold)

	expected := []bool{false, false, true, true, true, true}
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Index %d: expected %v, got %v", i, expected[i], result[i])
		}
	}
}

func TestCmpLeInt64_HappyPath(t *testing.T) {
	values := []int64{10, 20, 30, 40, 50, 60}
	threshold := int64(30)

	result := CmpLeInt64(values, threshold)

	expected := []bool{true, true, true, false, false, false}
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Index %d: expected %v, got %v", i, expected[i], result[i])
		}
	}
}

func TestCmpNeInt64_HappyPath(t *testing.T) {
	values := []int64{10, 20, 30, 40, 50, 60}
	threshold := int64(30)

	result := CmpNeInt64(values, threshold)

	expected := []bool{true, true, false, true, true, true}
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Index %d: expected %v, got %v", i, expected[i], result[i])
		}
	}
}

// Mask variant tests
func TestCmpEqInt64Mask_HappyPath(t *testing.T) {
	values := []int64{1, 2, 3, 4, 5, 6, 7, 8}
	threshold := int64(5)

	result := CmpEqInt64Mask(values, threshold)

	// Only index 4 should be set (value 5 == 5)
	if len(result) < 1 {
		t.Fatalf("Expected at least 1 uint64 in result")
	}
	expected := uint64(1 << 4)
	if result[0] != expected {
		t.Errorf("Expected 0x%x, got 0x%x", expected, result[0])
	}
}

func TestCmpGtInt64Mask_HappyPath(t *testing.T) {
	values := []int64{1, 2, 3, 4, 5, 6, 7, 8}
	threshold := int64(5)

	result := CmpGtInt64Mask(values, threshold)

	// Indices 5, 6, 7 should be set (values 6, 7, 8 > 5)
	if len(result) < 1 {
		t.Fatalf("Expected at least 1 uint64 in result")
	}
	expected := uint64((1 << 5) | (1 << 6) | (1 << 7))
	if result[0] != expected {
		t.Errorf("Expected 0x%x, got 0x%x", expected, result[0])
	}
}

// Edge case: empty array
func TestCmpEqInt64_EmptyArray(t *testing.T) {
	values := []int64{}
	threshold := int64(10)

	result := CmpEqInt64(values, threshold)

	if len(result) != 0 {
		t.Errorf("Expected empty result, got length %d", len(result))
	}
}

func TestCmpEqInt64Mask_EmptyArray(t *testing.T) {
	values := []int64{}
	threshold := int64(10)

	result := CmpEqInt64Mask(values, threshold)

	if len(result) != 0 {
		t.Errorf("Expected empty result, got length %d", len(result))
	}
}

// Edge case: single element
func TestCmpEqInt64_SingleElement(t *testing.T) {
	values := []int64{42}
	threshold := int64(42)

	result := CmpEqInt64(values, threshold)

	if len(result) != 1 || !result[0] {
		t.Errorf("Expected [true], got %v", result)
	}
}

func TestCmpGtInt64_SingleElementMatch(t *testing.T) {
	values := []int64{100}
	threshold := int64(50)

	result := CmpGtInt64(values, threshold)

	if len(result) != 1 || !result[0] {
		t.Errorf("Expected [true], got %v", result)
	}
}

func TestCmpGtInt64_SingleElementNoMatch(t *testing.T) {
	values := []int64{25}
	threshold := int64(50)

	result := CmpGtInt64(values, threshold)

	if len(result) != 1 || result[0] {
		t.Errorf("Expected [false], got %v", result)
	}
}

// Edge case: odd lengths (test SIMD remainder handling)
func TestCmpEqInt64_OddLength(t *testing.T) {
	values := []int64{1, 2, 3, 4, 5, 6, 7}
	threshold := int64(4)

	result := CmpEqInt64(values, threshold)

	expected := []bool{false, false, false, true, false, false, false}
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Index %d: expected %v, got %v", i, expected[i], result[i])
		}
	}
}

func TestCmpGtInt64_LengthThree(t *testing.T) {
	values := []int64{10, 20, 30}
	threshold := int64(15)

	result := CmpGtInt64(values, threshold)

	expected := []bool{false, true, true}
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Index %d: expected %v, got %v", i, expected[i], result[i])
		}
	}
}

// Corner case: negative numbers
func TestCmpGtInt64_NegativeNumbers(t *testing.T) {
	values := []int64{-50, -25, 0, 25, 50}
	threshold := int64(0)

	result := CmpGtInt64(values, threshold)

	expected := []bool{false, false, false, true, true}
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Index %d: expected %v, got %v", i, expected[i], result[i])
		}
	}
}

func TestCmpLtInt64_NegativeNumbers(t *testing.T) {
	values := []int64{-100, -50, -25, 0, 25}
	threshold := int64(-30)

	result := CmpLtInt64(values, threshold)

	expected := []bool{true, true, false, false, false}
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Index %d: expected %v, got %v", i, expected[i], result[i])
		}
	}
}

// Corner case: max/min int64 values
func TestCmpEqInt64_MaxInt64(t *testing.T) {
	const maxInt64 = int64(9223372036854775807)
	values := []int64{maxInt64 - 1, maxInt64, maxInt64 - 2}
	threshold := maxInt64

	result := CmpEqInt64(values, threshold)

	expected := []bool{false, true, false}
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Index %d: expected %v, got %v", i, expected[i], result[i])
		}
	}
}

func TestCmpEqInt64_MinInt64(t *testing.T) {
	const minInt64 = int64(-9223372036854775808)
	values := []int64{minInt64 + 1, minInt64, minInt64 + 2}
	threshold := minInt64

	result := CmpEqInt64(values, threshold)

	expected := []bool{false, true, false}
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Index %d: expected %v, got %v", i, expected[i], result[i])
		}
	}
}

// Corner case: all matches
func TestCmpEqInt64_AllMatch(t *testing.T) {
	values := []int64{7, 7, 7, 7, 7, 7, 7, 7}
	threshold := int64(7)

	result := CmpEqInt64(values, threshold)

	for i := range result {
		if !result[i] {
			t.Errorf("Index %d: expected true, got false", i)
		}
	}
}

func TestCmpGtInt64Mask_AllMatch(t *testing.T) {
	values := []int64{10, 20, 30, 40, 50, 60, 70, 80}
	threshold := int64(5)

	result := CmpGtInt64Mask(values, threshold)

	if len(result) < 1 {
		t.Fatalf("Expected at least 1 uint64 in result")
	}
	expected := uint64(0xFF) // All 8 bits set
	if result[0] != expected {
		t.Errorf("Expected 0x%x, got 0x%x", expected, result[0])
	}
}

// Corner case: no matches
func TestCmpEqInt64_NoMatch(t *testing.T) {
	values := []int64{1, 2, 3, 4, 5, 6, 7, 8}
	threshold := int64(100)

	result := CmpEqInt64(values, threshold)

	for i := range result {
		if result[i] {
			t.Errorf("Index %d: expected false, got true", i)
		}
	}
}

func TestCmpLtInt64Mask_NoMatch(t *testing.T) {
	values := []int64{100, 200, 300, 400}
	threshold := int64(50)

	result := CmpLtInt64Mask(values, threshold)

	// Should have at least one uint64, but all bits should be 0
	if len(result) == 0 {
		t.Errorf("Expected non-empty result")
	}
	for i, mask := range result {
		if mask != 0 {
			t.Errorf("Index %d: expected 0, got 0x%x", i, mask)
		}
	}
}

// Test SIMD threshold boundaries (AVX2 kicks in at 16+, NEON at 8+)
func TestCmpEqInt64_Length16_AVX2Boundary(t *testing.T) {
	values := make([]int64, 16)
	for i := range values {
		values[i] = int64(i * 10)
	}
	// Override to ensure only indices 8 and 15 match
	values[10] = 999 // Was 100, change to avoid false positive
	values[8] = 100
	values[15] = 100
	threshold := int64(100)

	result := CmpEqInt64(values, threshold)

	for i := range result {
		expected := i == 8 || i == 15
		if result[i] != expected {
			t.Errorf("Index %d: expected %v, got %v (value=%d)", i, expected, result[i], values[i])
		}
	}
}

func TestCmpGtInt64_Length8_NEONBoundary(t *testing.T) {
	values := []int64{5, 10, 15, 20, 25, 30, 35, 40}
	threshold := int64(22)

	result := CmpGtInt64(values, threshold)

	expected := []bool{false, false, false, false, true, true, true, true}
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Index %d: expected %v, got %v", i, expected[i], result[i])
		}
	}
}

// Test large arrays to ensure SIMD path is exercised
func TestCmpEqInt64_LargeArray(t *testing.T) {
	size := 1000
	values := make([]int64, size)
	for i := range values {
		values[i] = int64(i)
	}
	threshold := int64(500)

	result := CmpEqInt64(values, threshold)

	if len(result) != size {
		t.Fatalf("Expected length %d, got %d", size, len(result))
	}

	matchCount := 0
	for i := range result {
		if result[i] {
			matchCount++
			if values[i] != threshold {
				t.Errorf("False positive at index %d: value=%d", i, values[i])
			}
		}
	}

	if matchCount != 1 {
		t.Errorf("Expected 1 match, got %d", matchCount)
	}
}

func TestCmpGtInt64Mask_LargeArray(t *testing.T) {
	size := 100
	values := make([]int64, size)
	for i := range values {
		values[i] = int64(i)
	}
	threshold := int64(50)

	result := CmpGtInt64Mask(values, threshold)

	// Count set bits across all uint64s in result
	setBits := 0
	for _, mask := range result {
		for i := 0; i < 64; i++ {
			if (mask & (1 << i)) != 0 {
				setBits++
			}
		}
	}

	expectedCount := 49 // Values 51-99 are > 50
	if setBits != expectedCount {
		t.Errorf("Expected %d set bits, got %d", expectedCount, setBits)
	}
}
