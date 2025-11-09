package syndrdbsimd

import (
	"hash/crc32"
)

// hashInt64Generic computes a hash of an int64 value using FNV-1a algorithm.
// This is a simple, fast hash suitable for hash table operations.
func hashInt64Generic(value int64) uint64 {
	// FNV-1a hash constants
	const (
		offset64 = 14695981039346656037
		prime64  = 1099511628211
	)

	hash := uint64(offset64)
	bytes := uint64(value)

	// Process 8 bytes
	for i := 0; i < 8; i++ {
		hash ^= bytes & 0xFF
		hash *= prime64
		bytes >>= 8
	}

	return hash
}

// hashInt64SliceGeneric computes hashes for a slice of int64 values.
// Results are written to the output slice which must be the same length as values.
func hashInt64SliceGeneric(values []int64, output []uint64) {
	for i := range values {
		output[i] = hashInt64Generic(values[i])
	}
}

// crc32Generic computes the CRC32 checksum of a byte slice using the IEEE polynomial.
func crc32Generic(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}

// crc32Int64Generic computes the CRC32 checksum of an int64 value.
func crc32Int64Generic(value int64) uint32 {
	bytes := [8]byte{
		byte(value),
		byte(value >> 8),
		byte(value >> 16),
		byte(value >> 24),
		byte(value >> 32),
		byte(value >> 40),
		byte(value >> 48),
		byte(value >> 56),
	}
	return crc32.ChecksumIEEE(bytes[:])
}

// crc32Int64SliceGeneric computes CRC32 checksums for a slice of int64 values.
// Results are written to the output slice which must be the same length as values.
func crc32Int64SliceGeneric(values []int64, output []uint32) {
	for i := range values {
		output[i] = crc32Int64Generic(values[i])
	}
}

// xxhash64Generic implements the XXHash64 algorithm for a single int64 value.
// XXHash64 is a fast, high-quality non-cryptographic hash.
func xxhash64Generic(value int64) uint64 {
	// XXHash64 constants
	const (
		prime64_1 uint64 = 11400714785074694791
		prime64_2 uint64 = 14029467366897019727
		prime64_3 uint64 = 1609587929392839161
		prime64_4 uint64 = 9650029242287828579
		prime64_5 uint64 = 2870177450012600261
	)

	h64 := uint64(prime64_5 + 8) // seed=0, len=8

	// Process the 8 bytes
	k1 := uint64(value)
	k1 *= prime64_2
	k1 = rotl64(k1, 31)
	k1 *= prime64_1
	h64 ^= k1
	h64 = rotl64(h64, 27)*prime64_1 + prime64_4

	// Finalization
	h64 ^= h64 >> 33
	h64 *= prime64_2
	h64 ^= h64 >> 29
	h64 *= prime64_3
	h64 ^= h64 >> 32

	return h64
}

// xxhash64SliceGeneric computes XXHash64 hashes for a slice of int64 values.
// Results are written to the output slice which must be the same length as values.
func xxhash64SliceGeneric(values []int64, output []uint64) {
	for i := range values {
		output[i] = xxhash64Generic(values[i])
	}
}

// xxhash64BytesGeneric implements XXHash64 for arbitrary byte slices.
func xxhash64BytesGeneric(data []byte) uint64 {
	const (
		prime64_1 uint64 = 11400714785074694791
		prime64_2 uint64 = 14029467366897019727
		prime64_3 uint64 = 1609587929392839161
		prime64_4 uint64 = 9650029242287828579
		prime64_5 uint64 = 2870177450012600261
	)

	length := uint64(len(data))
	var h64 uint64

	if length >= 32 {
		var v1, v2, v3, v4 uint64
		v1 = prime64_1
		v1 += prime64_2
		v2 = prime64_2
		v3 = 0
		v4 = 0
		v4 -= prime64_1

		for len(data) >= 32 {
			v1 = round64(v1, u64(data[0:8]))
			v2 = round64(v2, u64(data[8:16]))
			v3 = round64(v3, u64(data[16:24]))
			v4 = round64(v4, u64(data[24:32]))
			data = data[32:]
		}

		h64 = rotl64(v1, 1) + rotl64(v2, 7) + rotl64(v3, 12) + rotl64(v4, 18)
		h64 = mergeRound64(h64, v1)
		h64 = mergeRound64(h64, v2)
		h64 = mergeRound64(h64, v3)
		h64 = mergeRound64(h64, v4)
	} else {
		h64 = prime64_5
	}

	h64 += length

	// Process remaining bytes
	for len(data) >= 8 {
		k1 := u64(data[0:8])
		k1 *= prime64_2
		k1 = rotl64(k1, 31)
		k1 *= prime64_1
		h64 ^= k1
		h64 = rotl64(h64, 27)*prime64_1 + prime64_4
		data = data[8:]
	}

	if len(data) >= 4 {
		h64 ^= uint64(u32(data[0:4])) * prime64_1
		h64 = rotl64(h64, 23)*prime64_2 + prime64_3
		data = data[4:]
	}

	for len(data) > 0 {
		h64 ^= uint64(data[0]) * prime64_5
		h64 = rotl64(h64, 11) * prime64_1
		data = data[1:]
	}

	// Finalization
	h64 ^= h64 >> 33
	h64 *= prime64_2
	h64 ^= h64 >> 29
	h64 *= prime64_3
	h64 ^= h64 >> 32

	return h64
}

// Helper functions for XXHash64

func rotl64(x uint64, r uint8) uint64 {
	return (x << r) | (x >> (64 - r))
}

func round64(acc, input uint64) uint64 {
	const prime64_1 uint64 = 11400714785074694791
	const prime64_2 uint64 = 14029467366897019727

	acc += input * prime64_2
	acc = rotl64(acc, 31)
	acc *= prime64_1
	return acc
}

func mergeRound64(acc, val uint64) uint64 {
	const prime64_1 uint64 = 11400714785074694791
	const prime64_2 uint64 = 14029467366897019727
	const prime64_4 uint64 = 9650029242287828579

	val = round64(0, val)
	acc ^= val
	acc = acc*prime64_1 + prime64_4
	return acc
}

func u64(b []byte) uint64 {
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func u32(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}
