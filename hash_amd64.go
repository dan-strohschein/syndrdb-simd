//go:build amd64
// +build amd64

package syndrdbsimd

// hashInt64AVX2 computes FNV-1a hashes for a slice of int64 values using AVX2.
// Processes 4 int64 values at a time.
//
//go:noescape
func hashInt64AVX2(values *int64, output *uint64, count int)

// crc32Int64AVX2 computes CRC32 checksums for a slice of int64 values using AVX2.
// Uses hardware CRC32C instruction when available.
//
//go:noescape
func crc32Int64AVX2(values *int64, output *uint32, count int)

// xxhash64AVX2 computes XXHash64 hashes for a slice of int64 values using AVX2.
// Processes 4 int64 values at a time.
//
//go:noescape
func xxhash64AVX2(values *int64, output *uint64, count int)
