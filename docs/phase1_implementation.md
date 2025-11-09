# Phase 1 Implementation: SIMD Library for SyndrDB

## Overview

This document describes the Phase 1 implementation of the SIMD (Single Instruction Multiple Data) library for SyndrDB. The library provides hardware-accelerated comparison and bitmap operations for int64 values, targeting both x86-64 (AVX2) and ARM64 (NEON) architectures with generic fallbacks.

## Architecture

The library uses a multi-layered architecture to support different CPU instruction sets:

```
User Code
    ↓
Public API (api.go)
    ↓
Architecture Routing (impl_*.go)
    ↓  ↓  ↓
   AVX2  NEON  Generic
   Assembly Fallback
```

### File Organization

- **api.go**: Public exported functions
- **impl_amd64.go**: AMD64/x86-64 routing logic
- **impl_arm64.go**: ARM64 routing logic  
- **impl_generic.go**: Generic fallback routing
- **compare_*.go**: Comparison function stubs per architecture
- **compare_*.s**: Assembly implementations (AMD64, ARM64)
- **bitmap_*.go**: Bitmap function stubs per architecture
- **bitmap_*.s**: Assembly implementations (AMD64, ARM64)
- **cpu_detection_*.go**: Runtime CPU feature detection per architecture

## Implemented Operations

### Comparison Operations

All comparison operations take a slice of int64 values and a threshold value, returning either a slice of booleans or a bitmask ([]uint64).

#### Bool Variants ([]bool return type)
- **CmpEqInt64**: Tests if values[i] == threshold
- **CmpNeInt64**: Tests if values[i] != threshold
- **CmpGtInt64**: Tests if values[i] > threshold
- **CmpGeInt64**: Tests if values[i] >= threshold
- **CmpLtInt64**: Tests if values[i] < threshold
- **CmpLeInt64**: Tests if values[i] <= threshold

#### Mask Variants ([]uint64 return type)
- **CmpEqInt64Mask**: Returns bitmask for equality
- **CmpNeInt64Mask**: Returns bitmask for inequality
- **CmpGtInt64Mask**: Returns bitmask for greater-than
- **CmpGeInt64Mask**: Returns bitmask for greater-or-equal
- **CmpLtInt64Mask**: Returns bitmask for less-than
- **CmpLeInt64Mask**: Returns bitmask for less-or-equal

### Bitmap Operations

All bitmap operations work on slices of uint64 values representing packed bitmaps.

- **AndBitmap**: Bitwise AND of two bitmaps
- **OrBitmap**: Bitwise OR of two bitmaps
- **XorBitmap**: Bitwise XOR of two bitmaps
- **NotBitmap**: Bitwise NOT of a bitmap
- **PopCount**: Count total set bits across all uint64s

### Utility Functions

- **BoolsToBitmask**: Convert []bool to packed []uint64 bitmask
- **BitmaskToBools**: Convert packed []uint64 bitmask to []bool

## SIMD Implementation Details

### AMD64 (x86-64) - AVX2

**Register Width**: 256-bit (YMM registers)  
**Elements per Operation**: 4 × int64 (or 4 × uint64 for bitmaps)  
**Minimum Array Size**: 16 elements (SIMD kicks in at 16+)

#### Instructions Used
- `VPBROADCASTQ`: Broadcast int64 threshold to all lanes
- `VMOVDQU`: Unaligned vector loads/stores
- `VPCMPGTQ`: Compare packed int64 for greater-than
- `VPCMPEQQ`: Compare packed int64 for equality
- `VMOVMSKPD`: Extract comparison result bits to scalar register
- `VPAND/VPOR/VPXOR`: Bitmap AND/OR/XOR operations
- `POPCNTQ`: Population count (bit count)
- `VZEROUPPER`: Clean up YMM registers before return

#### Performance
Approximately **4× throughput** vs scalar operations for large arrays.

### ARM64 - NEON

**Register Width**: 128-bit (V registers)  
**Elements per Operation**: 2 × int64 (scalar CMP approach) or 2 × uint64 for bitmaps  
**Minimum Array Size**: 8 elements (SIMD kicks in at 8+)

#### Instructions Used

**Comparisons** (scalar with assembly loop overhead reduction):
- `MOVD`: Load int64 values
- `CMP`: Scalar comparison
- `CSET`: Conditional set based on comparison flags
- `LSL/ORR`: Bit manipulation to build result mask

**Bitmap Operations**:
- `VLD1/VST1`: Vector loads/stores
- `VAND/VORR/VEOR`: Vector AND/OR/XOR
- `VCNT`: Vector population count (byte-level)
- `VMOV`: Extract bytes for accumulation

#### Performance Notes
ARM64 comparisons use scalar operations within tight assembly loops for overhead reduction. Go's ARM64 assembler doesn't support direct vector comparison mnemonics like VCMGT/VCMGE, requiring the CMP+CSET pattern.

Bitmap operations achieve **~2× throughput** vs pure Go loops.

### Generic Fallback

Used when:
- CPU doesn't support AVX2 (AMD64) or NEON (ARM64)
- Array size is below SIMD threshold
- Running on non-AMD64/ARM64 architecture

Pure Go implementation with simple for-loops.

## Usage Examples

### Comparison Operations

```go
package main

import (
    "fmt"
    simd "github.com/dan-strohschein/syndrdb-simd"
)

func main() {
    values := []int64{10, 20, 30, 40, 50, 60, 70, 80}
    threshold := int64(45)
    
    // Bool variant
    results := simd.CmpGtInt64(values, threshold)
    fmt.Println(results) // [false false false false true true true true]
    
    // Mask variant - more compact
    mask := simd.CmpGtInt64Mask(values, threshold)
    fmt.Printf("Mask: 0x%x\n", mask[0]) // Mask: 0xf0 (bits 4-7 set)
}
```

### Bitmap Operations

```go
package main

import (
    "fmt"
    simd "github.com/dan-strohschein/syndrdb-simd"
)

func main() {
    a := []uint64{0xFF00FF00FF00FF00, 0x00FF00FF00FF00FF}
    b := []uint64{0xF0F0F0F0F0F0F0F0, 0x0F0F0F0F0F0F0F0F}
    
    // Bitwise AND
    result := simd.AndBitmap(a, b)
    fmt.Printf("AND result: %x\n", result)
    
    // Count set bits
    count := simd.PopCount(a)
    fmt.Printf("Bit count: %d\n", count)
}
```

### Bitmask Conversion

```go
package main

import (
    "fmt"
    simd "github.com/dan-strohschein/syndrdb-simd"
)

func main() {
    bools := []bool{true, false, true, true, false, false, true, false}
    
    // Convert to compact bitmask
    mask := simd.BoolsToBitmask(bools)
    fmt.Printf("Mask: 0x%x\n", mask[0]) // Mask: 0x4d
    
    // Convert back
    recovered := simd.BitmaskToBools(mask, len(bools))
    fmt.Println(recovered) // [true false true true false false true false]
}
```

## CPU Feature Detection

The library automatically detects CPU capabilities at runtime:

### AMD64
- **AVX2**: 256-bit SIMD (required for acceleration)
- **AVX512**: Detected but not yet used (future enhancement)
- **SSE4.2**: Detected for potential future use

### ARM64
- **NEON**: 128-bit SIMD (standard on all ARM64 CPUs)
- **SVE**: Scalable Vector Extension (detected but not yet used)

Detection happens once at package initialization via `golang.org/x/sys/cpu`.

## Performance Characteristics

### SIMD Activation Thresholds

To avoid overhead of SIMD setup for small arrays:
- **AMD64 (AVX2)**: Arrays with 16+ elements
- **ARM64 (NEON)**: Arrays with 8+ elements
- **Below threshold**: Uses generic implementation

### Expected Speedup

For large arrays (1000+ elements):
- **AMD64 (AVX2)**: ~4× faster than generic
- **ARM64 (NEON)**: ~2× faster than generic (bitmap ops)

Actual speedup depends on:
- Array size and alignment
- CPU cache characteristics
- Memory bandwidth
- Comparison operation type

## Testing

Comprehensive test suite in `compare_test.go` and `bitmap_test.go` covering:

### Happy Path
- Normal operation with various array sizes
- All comparison types (6 operations × 2 variants)
- All bitmap operations (AND, OR, XOR, NOT, PopCount)

### Edge Cases
- Empty arrays
- Single-element arrays
- Odd-length arrays (remainder handling)

### Corner Cases
- Negative numbers
- Max/min int64 values (boundary conditions)
- All matches / no matches
- XOR identity (A ⊕ A = 0)
- NOT involution (¬¬A = A)

### SIMD Boundary Tests
- Arrays at exactly 8 elements (NEON threshold)
- Arrays at exactly 16 elements (AVX2 threshold)
- Large arrays (1000+ elements)

Run tests with:
```bash
go test -v
```

## Build and Import

### As a Library

```bash
go get github.com/dan-strohschein/syndrdb-simd
```

```go
import simd "github.com/dan-strohschein/syndrdb-simd"
```

### Building from Source

```bash
git clone https://github.com/dan-strohschein/syndrdb-simd.git
cd syndrdb-simd
go build .
go test -v
```

## Limitations and Future Work

### Current Limitations
1. **ARM64 Comparisons**: Uses scalar CMP+CSET instead of vector comparisons due to Go assembler limitations
2. **Bitmask Size**: Limited to 64 bits per uint64 (multi-uint64 support exists but adds overhead)
3. **Alignment**: No explicit alignment requirements, but aligned access may be faster

### Phase 2 Enhancements (Planned)
- Vectorized ARM64 comparisons (if Go assembler support improves)
- AVX-512 support for AMD64 (512-bit registers, 8× int64 per op)
- ARM SVE support (scalable vectors)
- Prefetching for large arrays
- Benchmarking framework

## License

[Add your license information here]

## Contributors

- Dan Strohschein (initial implementation)
