//go:build amd64

package syndrdbsimd

// cmpGtInt64AVX2 compares 4 int64 values against a threshold using AVX2.
// Returns a bitmask where bit i is set if values[i] > threshold.
// This function processes exactly 4 int64 values (32 bytes).
//
//go:noescape
func cmpGtInt64AVX2(values *int64, threshold int64) uint64

// cmpEqInt64AVX2 compares 4 int64 values for equality against a threshold using AVX2.
// Returns a bitmask where bit i is set if values[i] == threshold.
//
//go:noescape
func cmpEqInt64AVX2(values *int64, threshold int64) uint64

// cmpLtInt64AVX2 compares 4 int64 values against a threshold using AVX2.
// Returns a bitmask where bit i is set if values[i] < threshold.
//
//go:noescape
func cmpLtInt64AVX2(values *int64, threshold int64) uint64

// cmpGeInt64AVX2 compares 4 int64 values against a threshold using AVX2.
// Returns a bitmask where bit i is set if values[i] >= threshold.
//
//go:noescape
func cmpGeInt64AVX2(values *int64, threshold int64) uint64

// cmpLeInt64AVX2 compares 4 int64 values against a threshold using AVX2.
// Returns a bitmask where bit i is set if values[i] <= threshold.
//
//go:noescape
func cmpLeInt64AVX2(values *int64, threshold int64) uint64

// cmpNeInt64AVX2 compares 4 int64 values for inequality against a threshold using AVX2.
// Returns a bitmask where bit i is set if values[i] != threshold.
//
//go:noescape
func cmpNeInt64AVX2(values *int64, threshold int64) uint64
