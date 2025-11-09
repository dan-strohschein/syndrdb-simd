package syndrdbsimd

import (
	"testing"
)

// Happy path tests for bitmap operations
func TestAndBitmap_HappyPath(t *testing.T) {
	a := []uint64{0xFF00FF00FF00FF00, 0x00FF00FF00FF00FF}
	b := []uint64{0xF0F0F0F0F0F0F0F0, 0x0F0F0F0F0F0F0F0F}

	result := AndBitmap(a, b)

	expected := []uint64{0xF000F000F000F000, 0x000F000F000F000F}
	if len(result) != len(expected) {
		t.Fatalf("Expected length %d, got %d", len(expected), len(result))
	}

	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Index %d: expected 0x%x, got 0x%x", i, expected[i], result[i])
		}
	}
}

func TestOrBitmap_HappyPath(t *testing.T) {
	a := []uint64{0xFF00000000000000, 0x00000000000000FF}
	b := []uint64{0x00FF000000000000, 0x0000000000000F00}

	result := OrBitmap(a, b)

	expected := []uint64{0xFFFF000000000000, 0x0000000000000FFF}
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Index %d: expected 0x%x, got 0x%x", i, expected[i], result[i])
		}
	}
}

func TestXorBitmap_HappyPath(t *testing.T) {
	a := []uint64{0xFFFFFFFFFFFFFFFF, 0xAAAAAAAAAAAAAAAA}
	b := []uint64{0xFF00FF00FF00FF00, 0x5555555555555555}

	result := XorBitmap(a, b)

	expected := []uint64{0x00FF00FF00FF00FF, 0xFFFFFFFFFFFFFFFF}
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Index %d: expected 0x%x, got 0x%x", i, expected[i], result[i])
		}
	}
}

func TestNotBitmap_HappyPath(t *testing.T) {
	values := []uint64{0xFFFFFFFFFFFFFFFF, 0x0000000000000000, 0xAAAAAAAAAAAAAAAA}

	result := NotBitmap(values)

	expected := []uint64{0x0000000000000000, 0xFFFFFFFFFFFFFFFF, 0x5555555555555555}
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Index %d: expected 0x%x, got 0x%x", i, expected[i], result[i])
		}
	}
}

func TestPopCount_HappyPath(t *testing.T) {
	values := []uint64{
		0x0000000000000000, // 0 bits
		0xFFFFFFFFFFFFFFFF, // 64 bits
		0x00000000000000FF, // 8 bits
		0xAAAAAAAAAAAAAAAA, // 32 bits
	}

	result := PopCount(values)

	// PopCount returns total bits across all uint64s: 0 + 64 + 8 + 32 = 104
	expected := 104
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

// Edge case: empty arrays
func TestAndBitmap_EmptyArray(t *testing.T) {
	a := []uint64{}
	b := []uint64{}

	result := AndBitmap(a, b)

	if len(result) != 0 {
		t.Errorf("Expected empty result, got length %d", len(result))
	}
}

func TestNotBitmap_EmptyArray(t *testing.T) {
	values := []uint64{}

	result := NotBitmap(values)

	if len(result) != 0 {
		t.Errorf("Expected empty result, got length %d", len(result))
	}
}

func TestPopCount_EmptyArray(t *testing.T) {
	values := []uint64{}

	result := PopCount(values)

	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}
}

// Edge case: single element
func TestAndBitmap_SingleElement(t *testing.T) {
	a := []uint64{0xF0F0F0F0F0F0F0F0}
	b := []uint64{0xFF00FF00FF00FF00}

	result := AndBitmap(a, b)

	expected := uint64(0xF000F000F000F000)
	if len(result) != 1 || result[0] != expected {
		t.Errorf("Expected [0x%x], got [0x%x]", expected, result[0])
	}
}

func TestPopCount_SingleElement(t *testing.T) {
	values := []uint64{0x00000000FFFFFFFF}

	result := PopCount(values)

	if result != 32 {
		t.Errorf("Expected 32, got %d", result)
	}
}

// Edge case: odd lengths
func TestOrBitmap_OddLength(t *testing.T) {
	a := []uint64{0xFF, 0xF0, 0x0F}
	b := []uint64{0x00, 0x0F, 0xF0}

	result := OrBitmap(a, b)

	expected := []uint64{0xFF, 0xFF, 0xFF}
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Index %d: expected 0x%x, got 0x%x", i, expected[i], result[i])
		}
	}
}

func TestPopCount_OddLength(t *testing.T) {
	values := []uint64{0xFF, 0xFFFF, 0xFFFFFF}

	result := PopCount(values)

	// 8 + 16 + 24 = 48 total bits
	expected := 48
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

// Corner case: all zeros
func TestAndBitmap_AllZeros(t *testing.T) {
	a := []uint64{0, 0, 0, 0}
	b := []uint64{0xFFFFFFFFFFFFFFFF, 0xAAAAAAAAAAAAAAAA, 0x5555555555555555, 0xFF00FF00FF00FF00}

	result := AndBitmap(a, b)

	for i := range result {
		if result[i] != 0 {
			t.Errorf("Index %d: expected 0, got 0x%x", i, result[i])
		}
	}
}

func TestPopCount_AllZeros(t *testing.T) {
	values := []uint64{0, 0, 0, 0, 0}

	result := PopCount(values)

	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}
}

// Corner case: all ones
func TestOrBitmap_AllOnes(t *testing.T) {
	a := []uint64{0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF}
	b := []uint64{0x0000000000000000, 0x1234567890ABCDEF}

	result := OrBitmap(a, b)

	for i := range result {
		if result[i] != 0xFFFFFFFFFFFFFFFF {
			t.Errorf("Index %d: expected 0xFFFFFFFFFFFFFFFF, got 0x%x", i, result[i])
		}
	}
}

func TestPopCount_AllOnes(t *testing.T) {
	values := []uint64{0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF}

	result := PopCount(values)

	// 64 * 3 = 192 total bits
	if result != 192 {
		t.Errorf("Expected 192, got %d", result)
	}
}

// Corner case: XOR identity (A XOR A = 0)
func TestXorBitmap_Identity(t *testing.T) {
	values := []uint64{0xDEADBEEFCAFEBABE, 0x123456789ABCDEF0}

	result := XorBitmap(values, values)

	for i := range result {
		if result[i] != 0 {
			t.Errorf("Index %d: expected 0, got 0x%x", i, result[i])
		}
	}
}

// Corner case: NOT involution (NOT(NOT(A)) = A)
func TestNotBitmap_Involution(t *testing.T) {
	original := []uint64{0xDEADBEEFCAFEBABE, 0x123456789ABCDEF0, 0xAAAAAAAAAAAAAAAA}

	once := NotBitmap(original)
	twice := NotBitmap(once)

	for i := range original {
		if twice[i] != original[i] {
			t.Errorf("Index %d: expected 0x%x, got 0x%x", i, original[i], twice[i])
		}
	}
}

// Test SIMD threshold boundaries
func TestAndBitmap_Length8_NEONBoundary(t *testing.T) {
	a := make([]uint64, 8)
	b := make([]uint64, 8)
	for i := range a {
		a[i] = 0xFFFFFFFFFFFFFFFF
		b[i] = uint64(i)
	}

	result := AndBitmap(a, b)

	for i := range result {
		if result[i] != uint64(i) {
			t.Errorf("Index %d: expected %d, got 0x%x", i, i, result[i])
		}
	}
}

func TestOrBitmap_Length16_AVX2Boundary(t *testing.T) {
	a := make([]uint64, 16)
	b := make([]uint64, 16)
	for i := range a {
		a[i] = uint64(i) << 32
		b[i] = uint64(i)
	}

	result := OrBitmap(a, b)

	for i := range result {
		expected := (uint64(i) << 32) | uint64(i)
		if result[i] != expected {
			t.Errorf("Index %d: expected 0x%x, got 0x%x", i, expected, result[i])
		}
	}
}

// Test large arrays
func TestAndBitmap_LargeArray(t *testing.T) {
	size := 1000
	a := make([]uint64, size)
	b := make([]uint64, size)
	for i := range a {
		a[i] = 0xFFFFFFFF00000000
		b[i] = 0x00000000FFFFFFFF
	}

	result := AndBitmap(a, b)

	for i := range result {
		if result[i] != 0 {
			t.Errorf("Index %d: expected 0, got 0x%x", i, result[i])
		}
	}
}

func TestPopCount_LargeArray(t *testing.T) {
	size := 500
	values := make([]uint64, size)
	for i := range values {
		values[i] = 0xFFFFFFFFFFFFFFFF
	}

	result := PopCount(values)

	// 500 * 64 = 32000 total bits
	if result != 32000 {
		t.Errorf("Expected 32000, got %d", result)
	}
}

// Test bitmask conversion utilities
func TestBoolsToBitmask_HappyPath(t *testing.T) {
	bools := []bool{true, false, true, true, false, false, true, false}

	result := BoolsToBitmask(bools)

	// Bits: 0b01001101 = 0x4D (in first uint64)
	if len(result) != 1 {
		t.Fatalf("Expected length 1, got %d", len(result))
	}
	if result[0] != 0x4D {
		t.Errorf("Expected 0x4D, got 0x%x", result[0])
	}
}

func TestBitmaskToBools_HappyPath(t *testing.T) {
	mask := []uint64{0b10110010}
	length := 8

	result := BitmaskToBools(mask, length)

	expected := []bool{false, true, false, false, true, true, false, true}
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Index %d: expected %v, got %v", i, expected[i], result[i])
		}
	}
}

func TestBoolsToBitmask_Empty(t *testing.T) {
	bools := []bool{}

	result := BoolsToBitmask(bools)

	if len(result) != 0 {
		t.Errorf("Expected length 0, got %d", len(result))
	}
}

func TestBitmaskToBools_Empty(t *testing.T) {
	result := BitmaskToBools([]uint64{0xFF}, 0)

	if len(result) != 0 {
		t.Errorf("Expected empty result, got length %d", len(result))
	}
}

func TestBoolsToBitmask_SingleTrue(t *testing.T) {
	bools := []bool{true}

	result := BoolsToBitmask(bools)

	if len(result) != 1 || result[0] != 1 {
		t.Errorf("Expected [1], got %v", result)
	}
}

func TestBitmaskToBools_SingleBit(t *testing.T) {
	result := BitmaskToBools([]uint64{1}, 1)

	if len(result) != 1 || !result[0] {
		t.Errorf("Expected [true], got %v", result)
	}
}

func TestBoolsToBitmask_Max64Bits(t *testing.T) {
	bools := make([]bool, 64)
	for i := range bools {
		bools[i] = i%2 == 0
	}

	result := BoolsToBitmask(bools)

	// Alternating pattern: 0x5555555555555555
	if len(result) != 1 {
		t.Fatalf("Expected length 1, got %d", len(result))
	}
	if result[0] != 0x5555555555555555 {
		t.Errorf("Expected 0x5555555555555555, got 0x%x", result[0])
	}
}

func TestBitmaskToBools_Max64Bits(t *testing.T) {
	mask := []uint64{0xAAAAAAAAAAAAAAAA}

	result := BitmaskToBools(mask, 64)

	if len(result) != 64 {
		t.Fatalf("Expected length 64, got %d", len(result))
	}

	for i := range result {
		expected := i%2 == 1
		if result[i] != expected {
			t.Errorf("Index %d: expected %v, got %v", i, expected, result[i])
		}
	}
}

// Test round-trip conversions
func TestBitmaskConversion_RoundTrip(t *testing.T) {
	original := []bool{true, false, true, false, true, true, false, false, true}

	mask := BoolsToBitmask(original)
	recovered := BitmaskToBools(mask, len(original))

	for i := range original {
		if recovered[i] != original[i] {
			t.Errorf("Index %d: expected %v, got %v", i, original[i], recovered[i])
		}
	}
}
