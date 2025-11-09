#include "textflag.h"

// func cmpGtInt64AVX2(values *int64, threshold int64) uint64
//
// This function compares 4 int64 values against a threshold using AVX2 SIMD instructions.
// It leverages 256-bit YMM registers to process 4×64-bit integers in parallel.
//
// AVX2 provides significant speedup (4x theoretical) over scalar comparisons because:
// - Single instruction processes 4 comparisons simultaneously
// - Reduces loop overhead and branch mispredictions
// - Better utilization of CPU execution units
//
// Register usage:
//   YMM0: Broadcasted threshold value [threshold, threshold, threshold, threshold]
//   YMM1: Input values from memory [values[0], values[1], values[2], values[3]]
//   YMM2: Comparison result mask [0xFFFF... if true, 0x0000... if false per lane]
//   AX:   Final 4-bit result mask
//
TEXT ·cmpGtInt64AVX2(SB), NOSPLIT, $0-24
    // Load function arguments from stack frame
    // Go calling convention: arguments at positive offsets from FP (frame pointer)
    MOVQ    values+0(FP), SI        // SI = pointer to values array
    MOVQ    threshold+8(FP), AX     // AX = threshold value to compare against
    
    // Broadcast the threshold value to all 4 lanes of a 256-bit YMM register
    // VPBROADCASTQ takes a 64-bit value and replicates it across all lanes
    // This creates: YMM0 = [threshold, threshold, threshold, threshold]
    VPBROADCASTQ AX, Y0
    
    // Load 4 consecutive int64 values from memory into YMM1
    // VMOVDQU performs an unaligned load (works regardless of memory alignment)
    // Loads 256 bits = 4 × 64-bit integers
    VMOVDQU (SI), Y1                // Y1 = [values[0], values[1], values[2], values[3]]
    
    // Perform packed comparison: Y1 > Y0
    // VPCMPGTQ compares 4 pairs of 64-bit signed integers for "greater than"
    // Result: Each 64-bit lane becomes 0xFFFFFFFFFFFFFFFF if true, 0x0000000000000000 if false
    VPCMPGTQ Y0, Y1, Y2             // Y2[i] = (Y1[i] > Y0[i]) ? 0xFFFF... : 0x0000...
    
    // Convert vector comparison results to a compact bitmask
    // VMOVMSKPD extracts the sign bit from each 64-bit double-precision lane
    // Since our comparison results are all 1s (negative) or all 0s (positive),
    // the sign bit perfectly represents the comparison result
    // Result: 4-bit mask where bit i = 1 if lane i comparison was true
    VMOVMSKPD Y2, AX                // AX = compact 4-bit mask (bits 0-3 valid)
    
    // Clean up AVX state
    // VZEROUPPER clears the upper 128 bits of all YMM registers
    // This is CRITICAL to avoid severe performance penalties when mixing AVX and SSE code
    // Without this, subsequent SSE instructions can stall for dozens of cycles
    VZEROUPPER
    
    // Return the result mask
    // Go expects return values at specific stack offsets
    MOVQ    AX, ret+16(FP)
    RET

// func cmpEqInt64AVX2(values *int64, threshold int64) uint64
//
// Performs equality comparison using AVX2. Similar to cmpGtInt64AVX2 but uses VPCMPEQQ.
// VPCMPEQQ compares for equality instead of greater-than.
//
TEXT ·cmpEqInt64AVX2(SB), NOSPLIT, $0-24
    MOVQ    values+0(FP), SI
    MOVQ    threshold+8(FP), AX
    
    // Broadcast threshold to all lanes
    VPBROADCASTQ AX, Y0
    
    // Load 4 int64 values
    VMOVDQU (SI), Y1
    
    // Compare for equality: Y1 == Y0
    // VPCMPEQQ sets each 64-bit lane to 0xFFFF... if equal, 0x0000... otherwise
    VPCMPEQQ Y0, Y1, Y2
    
    // Extract sign bits to create bitmask
    VMOVMSKPD Y2, AX
    
    // Clean up AVX state to prevent performance penalties
    VZEROUPPER
    
    MOVQ    AX, ret+16(FP)
    RET

// func cmpLtInt64AVX2(values *int64, threshold int64) uint64
//
// Performs less-than comparison. Since AVX2 doesn't have a direct "less than" instruction
// for 64-bit integers, we compare threshold > values (which is equivalent).
//
TEXT ·cmpLtInt64AVX2(SB), NOSPLIT, $0-24
    MOVQ    values+0(FP), SI
    MOVQ    threshold+8(FP), AX
    
    VPBROADCASTQ AX, Y0
    VMOVDQU (SI), Y1
    
    // Compare: Y0 > Y1 (equivalent to Y1 < Y0)
    // We swap the operands to achieve "less than" semantics
    VPCMPGTQ Y1, Y0, Y2             // Y2[i] = (Y0[i] > Y1[i]) ≡ (Y1[i] < Y0[i])
    
    VMOVMSKPD Y2, AX
    VZEROUPPER
    
    MOVQ    AX, ret+16(FP)
    RET

// func cmpGeInt64AVX2(values *int64, threshold int64) uint64
//
// Performs greater-than-or-equal comparison.
// Implemented as: (values > threshold) OR (values == threshold)
// This requires two comparisons and a bitwise OR.
//
TEXT ·cmpGeInt64AVX2(SB), NOSPLIT, $0-24
    MOVQ    values+0(FP), SI
    MOVQ    threshold+8(FP), AX
    
    VPBROADCASTQ AX, Y0
    VMOVDQU (SI), Y1
    
    // First comparison: values > threshold
    VPCMPGTQ Y0, Y1, Y2             // Y2 = (Y1 > Y0)
    
    // Second comparison: values == threshold
    VPCMPEQQ Y0, Y1, Y3             // Y3 = (Y1 == Y0)
    
    // Combine with OR: (Y1 > Y0) OR (Y1 == Y0) ≡ (Y1 >= Y0)
    VPOR Y3, Y2, Y2                 // Y2 = Y2 | Y3
    
    VMOVMSKPD Y2, AX
    VZEROUPPER
    
    MOVQ    AX, ret+16(FP)
    RET

// func cmpLeInt64AVX2(values *int64, threshold int64) uint64
//
// Performs less-than-or-equal comparison.
// Implemented as: (values < threshold) OR (values == threshold)
// Uses the same technique as cmpGeInt64AVX2 but with operands swapped.
//
TEXT ·cmpLeInt64AVX2(SB), NOSPLIT, $0-24
    MOVQ    values+0(FP), SI
    MOVQ    threshold+8(FP), AX
    
    VPBROADCASTQ AX, Y0
    VMOVDQU (SI), Y1
    
    // First comparison: threshold > values (i.e., values < threshold)
    VPCMPGTQ Y1, Y0, Y2             // Y2 = (Y0 > Y1) ≡ (Y1 < Y0)
    
    // Second comparison: values == threshold
    VPCMPEQQ Y0, Y1, Y3             // Y3 = (Y1 == Y0)
    
    // Combine with OR
    VPOR Y3, Y2, Y2
    
    VMOVMSKPD Y2, AX
    VZEROUPPER
    
    MOVQ    AX, ret+16(FP)
    RET

// func cmpNeInt64AVX2(values *int64, threshold int64) uint64
//
// Performs inequality (not-equal) comparison.
// Implemented as: NOT(values == threshold)
// We compare for equality and then invert the result mask.
//
TEXT ·cmpNeInt64AVX2(SB), NOSPLIT, $0-24
    MOVQ    values+0(FP), SI
    MOVQ    threshold+8(FP), AX
    
    VPBROADCASTQ AX, Y0
    VMOVDQU (SI), Y1
    
    // Compare for equality
    VPCMPEQQ Y0, Y1, Y2             // Y2 = (Y1 == Y0)
    
    // Extract bitmask
    VMOVMSKPD Y2, AX
    
    // Invert the 4-bit result (XOR with 0xF to flip bits 0-3)
    // Since we only use the bottom 4 bits, XOR with 0xF inverts them
    XORQ $0xF, AX                   // AX = ~AX (for bottom 4 bits)
    
    VZEROUPPER
    
    MOVQ    AX, ret+16(FP)
    RET
