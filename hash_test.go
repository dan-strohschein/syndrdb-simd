package syndrdbsimd

import (
	"testing"
)

// ============================================================================
// HashInt64 Tests
// ============================================================================

func TestHashInt64_Empty(t *testing.T) {
	values := []int64{}
	output := []uint64{}
	HashInt64(values, output)
	// Should handle empty gracefully
}

func TestHashInt64_Single(t *testing.T) {
	values := []int64{42}
	output := make([]uint64, 1)
	HashInt64(values, output)

	expected := hashInt64Generic(42)
	if output[0] != expected {
		t.Errorf("HashInt64([42]) = %d, expected %d", output[0], expected)
	}
}

func TestHashInt64_Multiple(t *testing.T) {
	values := []int64{1, 2, 3, 4, 5}
	output := make([]uint64, len(values))
	HashInt64(values, output)

	// Verify each hash
	for i, v := range values {
		expected := hashInt64Generic(v)
		if output[i] != expected {
			t.Errorf("HashInt64: output[%d] = %d, expected %d for value %d", i, output[i], expected, v)
		}
	}
}

func TestHashInt64_LargeArray(t *testing.T) {
	// Test SIMD path with 100 elements
	values := make([]int64, 100)
	for i := range values {
		values[i] = int64(i * 12345)
	}

	output := make([]uint64, len(values))
	HashInt64(values, output)

	// Verify correctness against generic implementation
	expected := make([]uint64, len(values))
	hashInt64SliceGeneric(values, expected)

	for i := range output {
		if output[i] != expected[i] {
			t.Errorf("HashInt64: output[%d] = %d, expected %d", i, output[i], expected[i])
		}
	}
}

func TestHashInt64_OddLength(t *testing.T) {
	// Odd length to test remainder handling
	values := []int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110}
	output := make([]uint64, len(values))
	HashInt64(values, output)

	expected := make([]uint64, len(values))
	hashInt64SliceGeneric(values, expected)

	for i := range output {
		if output[i] != expected[i] {
			t.Errorf("HashInt64: output[%d] = %d, expected %d", i, output[i], expected[i])
		}
	}
}

func TestHashInt64_NegativeValues(t *testing.T) {
	values := []int64{-100, -50, 0, 50, 100}
	output := make([]uint64, len(values))
	HashInt64(values, output)

	expected := make([]uint64, len(values))
	hashInt64SliceGeneric(values, expected)

	for i := range output {
		if output[i] != expected[i] {
			t.Errorf("HashInt64: output[%d] = %d, expected %d", i, output[i], expected[i])
		}
	}
}

func TestHashInt64_Deterministic(t *testing.T) {
	// Same values should produce same hashes
	values := []int64{42, 42, 42}
	output := make([]uint64, len(values))
	HashInt64(values, output)

	if output[0] != output[1] || output[1] != output[2] {
		t.Errorf("HashInt64: same values produced different hashes: %d, %d, %d", output[0], output[1], output[2])
	}
}

func TestHashInt64_Distribution(t *testing.T) {
	// Different values should produce different hashes (collision test)
	values := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	output := make([]uint64, len(values))
	HashInt64(values, output)

	// Check for uniqueness
	seen := make(map[uint64]bool)
	for i, hash := range output {
		if seen[hash] {
			t.Errorf("HashInt64: collision detected for value %d at index %d", values[i], i)
		}
		seen[hash] = true
	}
}

// ============================================================================
// CRC32Int64 Tests
// ============================================================================

func TestCRC32Int64_Empty(t *testing.T) {
	values := []int64{}
	output := []uint32{}
	CRC32Int64(values, output)
	// Should handle empty gracefully
}

func TestCRC32Int64_Single(t *testing.T) {
	values := []int64{42}
	output := make([]uint32, 1)
	CRC32Int64(values, output)

	expected := crc32Int64Generic(42)
	if output[0] != expected {
		t.Errorf("CRC32Int64([42]) = %d, expected %d", output[0], expected)
	}
}

func TestCRC32Int64_Multiple(t *testing.T) {
	values := []int64{1, 2, 3, 4, 5}
	output := make([]uint32, len(values))
	CRC32Int64(values, output)

	// Verify each CRC
	for i, v := range values {
		expected := crc32Int64Generic(v)
		if output[i] != expected {
			t.Errorf("CRC32Int64: output[%d] = %d, expected %d for value %d", i, output[i], expected, v)
		}
	}
}

func TestCRC32Int64_LargeArray(t *testing.T) {
	values := make([]int64, 100)
	for i := range values {
		values[i] = int64(i * 9876)
	}

	output := make([]uint32, len(values))
	CRC32Int64(values, output)

	// Verify correctness
	expected := make([]uint32, len(values))
	crc32Int64SliceGeneric(values, expected)

	for i := range output {
		if output[i] != expected[i] {
			t.Errorf("CRC32Int64: output[%d] = %d, expected %d", i, output[i], expected[i])
		}
	}
}

func TestCRC32Int64_NegativeValues(t *testing.T) {
	values := []int64{-1000, -1, 0, 1, 1000}
	output := make([]uint32, len(values))
	CRC32Int64(values, output)

	expected := make([]uint32, len(values))
	crc32Int64SliceGeneric(values, expected)

	for i := range output {
		if output[i] != expected[i] {
			t.Errorf("CRC32Int64: output[%d] = %d, expected %d", i, output[i], expected[i])
		}
	}
}

func TestCRC32_Bytes(t *testing.T) {
	data := []byte("Hello, World!")
	crc := CRC32(data)

	expected := crc32Generic(data)
	if crc != expected {
		t.Errorf("CRC32 = %d, expected %d", crc, expected)
	}
}

// ============================================================================
// XXHash64 Tests
// ============================================================================

func TestXXHash64_Empty(t *testing.T) {
	values := []int64{}
	output := []uint64{}
	XXHash64(values, output)
	// Should handle empty gracefully
}

func TestXXHash64_Single(t *testing.T) {
	values := []int64{42}
	output := make([]uint64, 1)
	XXHash64(values, output)

	expected := xxhash64Generic(42)
	if output[0] != expected {
		t.Errorf("XXHash64([42]) = %d, expected %d", output[0], expected)
	}
}

func TestXXHash64_Multiple(t *testing.T) {
	values := []int64{1, 2, 3, 4, 5}
	output := make([]uint64, len(values))
	XXHash64(values, output)

	// Verify each hash
	for i, v := range values {
		expected := xxhash64Generic(v)
		if output[i] != expected {
			t.Errorf("XXHash64: output[%d] = %d, expected %d for value %d", i, output[i], expected, v)
		}
	}
}

func TestXXHash64_LargeArray(t *testing.T) {
	values := make([]int64, 100)
	for i := range values {
		values[i] = int64(i * 54321)
	}

	output := make([]uint64, len(values))
	XXHash64(values, output)

	// Verify correctness
	expected := make([]uint64, len(values))
	xxhash64SliceGeneric(values, expected)

	for i := range output {
		if output[i] != expected[i] {
			t.Errorf("XXHash64: output[%d] = %d, expected %d", i, output[i], expected[i])
		}
	}
}

func TestXXHash64_OddLength(t *testing.T) {
	// Odd length to test remainder handling
	values := []int64{10, 20, 30, 40, 50, 60, 70}
	output := make([]uint64, len(values))
	XXHash64(values, output)

	expected := make([]uint64, len(values))
	xxhash64SliceGeneric(values, expected)

	for i := range output {
		if output[i] != expected[i] {
			t.Errorf("XXHash64: output[%d] = %d, expected %d", i, output[i], expected[i])
		}
	}
}

func TestXXHash64_NegativeValues(t *testing.T) {
	values := []int64{-999, -1, 0, 1, 999}
	output := make([]uint64, len(values))
	XXHash64(values, output)

	expected := make([]uint64, len(values))
	xxhash64SliceGeneric(values, expected)

	for i := range output {
		if output[i] != expected[i] {
			t.Errorf("XXHash64: output[%d] = %d, expected %d", i, output[i], expected[i])
		}
	}
}

func TestXXHash64_Deterministic(t *testing.T) {
	// Same values should produce same hashes
	values := []int64{123, 123, 123}
	output := make([]uint64, len(values))
	XXHash64(values, output)

	if output[0] != output[1] || output[1] != output[2] {
		t.Errorf("XXHash64: same values produced different hashes: %d, %d, %d", output[0], output[1], output[2])
	}
}

func TestXXHash64_Distribution(t *testing.T) {
	// Different values should produce different hashes
	values := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	output := make([]uint64, len(values))
	XXHash64(values, output)

	// Check for uniqueness
	seen := make(map[uint64]bool)
	for i, hash := range output {
		if seen[hash] {
			t.Errorf("XXHash64: collision detected for value %d at index %d", values[i], i)
		}
		seen[hash] = true
	}
}

func TestXXHash64Bytes(t *testing.T) {
	data := []byte("The quick brown fox jumps over the lazy dog")
	hash := XXHash64Bytes(data)

	expected := xxhash64BytesGeneric(data)
	if hash != expected {
		t.Errorf("XXHash64Bytes = %d, expected %d", hash, expected)
	}
}

func TestXXHash64Bytes_Empty(t *testing.T) {
	data := []byte{}
	hash := XXHash64Bytes(data)

	expected := xxhash64BytesGeneric(data)
	if hash != expected {
		t.Errorf("XXHash64Bytes(empty) = %d, expected %d", hash, expected)
	}
}

func TestXXHash64Bytes_LongString(t *testing.T) {
	// Create a long string (> 32 bytes to test the full loop)
	data := make([]byte, 128)
	for i := range data {
		data[i] = byte(i % 256)
	}

	hash := XXHash64Bytes(data)
	expected := xxhash64BytesGeneric(data)

	if hash != expected {
		t.Errorf("XXHash64Bytes(long) = %d, expected %d", hash, expected)
	}
}

// ============================================================================
// SIMD Boundary Tests
// ============================================================================

func TestHashInt64_SIMDBoundary_8(t *testing.T) {
	// 8 elements - boundary for NEON
	values := make([]int64, 8)
	for i := range values {
		values[i] = int64(i + 1)
	}

	output := make([]uint64, len(values))
	HashInt64(values, output)

	expected := make([]uint64, len(values))
	hashInt64SliceGeneric(values, expected)

	for i := range output {
		if output[i] != expected[i] {
			t.Errorf("HashInt64 (len=8): output[%d] = %d, expected %d", i, output[i], expected[i])
		}
	}
}

func TestHashInt64_SIMDBoundary_16(t *testing.T) {
	// 16 elements - boundary for AVX2
	values := make([]int64, 16)
	for i := range values {
		values[i] = int64(i * 100)
	}

	output := make([]uint64, len(values))
	HashInt64(values, output)

	expected := make([]uint64, len(values))
	hashInt64SliceGeneric(values, expected)

	for i := range output {
		if output[i] != expected[i] {
			t.Errorf("HashInt64 (len=16): output[%d] = %d, expected %d", i, output[i], expected[i])
		}
	}
}

func TestXXHash64_SIMDBoundary_Large(t *testing.T) {
	// 1000 elements to thoroughly test SIMD
	values := make([]int64, 1000)
	for i := range values {
		values[i] = int64(i*i - 500)
	}

	output := make([]uint64, len(values))
	XXHash64(values, output)

	expected := make([]uint64, len(values))
	xxhash64SliceGeneric(values, expected)

	for i := range output {
		if output[i] != expected[i] {
			t.Errorf("XXHash64 (len=1000): output[%d] = %d, expected %d for value %d",
				i, output[i], expected[i], values[i])
			break // Only show first error to avoid spam
		}
	}
}

// ============================================================================
// Integration Tests
// ============================================================================

func TestHashFunctions_Integration(t *testing.T) {
	// Test that all three hash functions work on the same data
	values := []int64{100, 200, 300, 400, 500}

	// HashInt64
	hashOutput := make([]uint64, len(values))
	HashInt64(values, hashOutput)

	// CRC32Int64
	crcOutput := make([]uint32, len(values))
	CRC32Int64(values, crcOutput)

	// XXHash64
	xxOutput := make([]uint64, len(values))
	XXHash64(values, xxOutput)

	// All should produce valid results (non-zero for these values)
	for i := range values {
		if hashOutput[i] == 0 {
			t.Errorf("HashInt64[%d] produced zero hash", i)
		}
		if crcOutput[i] == 0 {
			t.Errorf("CRC32Int64[%d] produced zero CRC", i)
		}
		if xxOutput[i] == 0 {
			t.Errorf("XXHash64[%d] produced zero hash", i)
		}
	}
}
