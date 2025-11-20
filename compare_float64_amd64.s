#include "textflag.h"

// func cmpGtFloat64AVX2(values *float64, threshold float64) uint64
//
// This function compares 4 float64 values against a threshold using AVX2 SIMD instructions.
// It uses VCMPPD with the greater-than predicate (0x1E) to compare 4×64-bit floats in parallel.
//
// AVX2 comparison predicates for VCMPPD:
//   0x00: EQ (equal)
//   0x01: LT (less than)
//   0x02: LE (less than or equal)
//   0x11: LT (unordered, non-signaling)
//   0x12: LE (unordered, non-signaling)
//   0x1D: GE (greater than or equal, ordered, non-signaling)
//   0x1E: GT (greater than, ordered, non-signaling)
//   0x04: NEQ (not equal)
//
// Using ordered comparisons (0x1D, 0x1E) ensures NaN comparisons return false per IEEE 754.
//
TEXT ·cmpGtFloat64AVX2(SB), NOSPLIT, $0-24
    MOVQ    values+0(FP), SI        // SI = pointer to values array
    MOVQ    threshold+8(FP), X0     // X0 = threshold (lower 64 bits)
    
    // Broadcast the threshold to all 4 lanes of YMM0
    VBROADCASTSD X0, Y0             // Y0 = [threshold, threshold, threshold, threshold]
    
    // Load 4 float64 values from memory
    VMOVUPD (SI), Y1                // Y1 = [values[0], values[1], values[2], values[3]]
    
    // Compare: Y1 > Y0 (ordered, non-signaling - NaN returns false)
    // VCMPPD with immediate 0x1E (GT ordered)
    VCMPPD  $0x1E, Y0, Y1, Y2       // Y2[i] = (Y1[i] > Y0[i]) ? 0xFFFF... : 0x0000...
    
    // Extract bitmask from comparison results
    VMOVMSKPD Y2, AX                // AX = 4-bit mask
    
    // Clean up AVX state
    VZEROUPPER
    
    MOVQ    AX, ret+16(FP)
    RET

// func cmpGeFloat64AVX2(values *float64, threshold float64) uint64
TEXT ·cmpGeFloat64AVX2(SB), NOSPLIT, $0-24
    MOVQ    values+0(FP), SI
    MOVQ    threshold+8(FP), X0
    
    VBROADCASTSD X0, Y0
    VMOVUPD (SI), Y1
    
    // Compare: Y1 >= Y0 (ordered, non-signaling)
    VCMPPD  $0x1D, Y0, Y1, Y2
    
    VMOVMSKPD Y2, AX
    VZEROUPPER
    
    MOVQ    AX, ret+16(FP)
    RET

// func cmpLtFloat64AVX2(values *float64, threshold float64) uint64
TEXT ·cmpLtFloat64AVX2(SB), NOSPLIT, $0-24
    MOVQ    values+0(FP), SI
    MOVQ    threshold+8(FP), X0
    
    VBROADCASTSD X0, Y0
    VMOVUPD (SI), Y1
    
    // Compare: Y1 < Y0 (ordered, non-signaling)
    VCMPPD  $0x11, Y0, Y1, Y2
    
    VMOVMSKPD Y2, AX
    VZEROUPPER
    
    MOVQ    AX, ret+16(FP)
    RET

// func cmpLeFloat64AVX2(values *float64, threshold float64) uint64
TEXT ·cmpLeFloat64AVX2(SB), NOSPLIT, $0-24
    MOVQ    values+0(FP), SI
    MOVQ    threshold+8(FP), X0
    
    VBROADCASTSD X0, Y0
    VMOVUPD (SI), Y1
    
    // Compare: Y1 <= Y0 (ordered, non-signaling)
    VCMPPD  $0x12, Y0, Y1, Y2
    
    VMOVMSKPD Y2, AX
    VZEROUPPER
    
    MOVQ    AX, ret+16(FP)
    RET

// func cmpEqFloat64AVX2(values *float64, threshold float64) uint64
TEXT ·cmpEqFloat64AVX2(SB), NOSPLIT, $0-24
    MOVQ    values+0(FP), SI
    MOVQ    threshold+8(FP), X0
    
    VBROADCASTSD X0, Y0
    VMOVUPD (SI), Y1
    
    // Compare: Y1 == Y0 (ordered, quiet)
    VCMPPD  $0x00, Y0, Y1, Y2
    
    VMOVMSKPD Y2, AX
    VZEROUPPER
    
    MOVQ    AX, ret+16(FP)
    RET

// func cmpNeFloat64AVX2(values *float64, threshold float64) uint64
TEXT ·cmpNeFloat64AVX2(SB), NOSPLIT, $0-24
    MOVQ    values+0(FP), SI
    MOVQ    threshold+8(FP), X0
    
    VBROADCASTSD X0, Y0
    VMOVUPD (SI), Y1
    
    // Compare: Y1 != Y0 (unordered - NaN != x returns true)
    VCMPPD  $0x04, Y0, Y1, Y2
    
    VMOVMSKPD Y2, AX
    VZEROUPPER
    
    MOVQ    AX, ret+16(FP)
    RET
