// +build arm64

package syndrdbsimd

// hashInt64NEON computes FNV-1a hashes for a slice of int64 values using NEON.
// Processes 2 int64 values at a time.
//
//go:noescape
func hashInt64NEON(values *int64, output *uint64, count int)

// crc32Int64NEON computes CRC32 checksums for a slice of int64 values using ARM64 CRC32.
// Uses hardware CRC32C instruction when available.
//
//go:noescape
func crc32Int64NEON(values *int64, output *uint32, count int)

// xxhash64NEON computes XXHash64 hashes for a slice of int64 values using NEON.
// Processes 2 int64 values at a time.
//
//go:noescape
func xxhash64NEON(values *int64, output *uint64, count int)
