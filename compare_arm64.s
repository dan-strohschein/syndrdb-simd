#include "textflag.h"

// func cmpGtInt64NEON(values *int64, threshold int64) uint64
//
// Go's ARM64 assembler doesn't have direct VCMGT mnemonics.
// We use scalar comparisons with CMP + CSET which is still efficient
// and benefits from tight assembly loop overhead reduction.
//
TEXT ·cmpGtInt64NEON(SB), NOSPLIT, $0-24
    MOVD    values+0(FP), R0        // R0 = pointer to values array
    MOVD    threshold+8(FP), R1     // R1 = threshold value
    MOVD    $0, R3                  // R3 = result mask
    
    // Load first value and compare
    MOVD    0(R0), R4
    CMP     R1, R4                  // Compare R4 with threshold
    CSET    GT, R5                  // R5 = 1 if R4 > R1, else 0
    ORR     R5, R3, R3              // Set bit 0
    
    // Load second value and compare
    MOVD    8(R0), R4
    CMP     R1, R4
    CSET    GT, R5
    LSL     $1, R5, R5              // Shift to bit 1
    ORR     R5, R3, R3              // Set bit 1
    
    MOVD    R3, ret+16(FP)
    RET

// func cmpEqInt64NEON(values *int64, threshold int64) uint64
TEXT ·cmpEqInt64NEON(SB), NOSPLIT, $0-24
    MOVD    values+0(FP), R0
    MOVD    threshold+8(FP), R1
    MOVD    $0, R3
    
    MOVD    0(R0), R4
    CMP     R1, R4
    CSET    EQ, R5
    ORR     R5, R3, R3
    
    MOVD    8(R0), R4
    CMP     R1, R4
    CSET    EQ, R5
    LSL     $1, R5, R5
    ORR     R5, R3, R3
    
    MOVD    R3, ret+16(FP)
    RET

// func cmpLtInt64NEON(values *int64, threshold int64) uint64
TEXT ·cmpLtInt64NEON(SB), NOSPLIT, $0-24
    MOVD    values+0(FP), R0
    MOVD    threshold+8(FP), R1
    MOVD    $0, R3
    
    MOVD    0(R0), R4
    CMP     R1, R4
    CSET    LT, R5
    ORR     R5, R3, R3
    
    MOVD    8(R0), R4
    CMP     R1, R4
    CSET    LT, R5
    LSL     $1, R5, R5
    ORR     R5, R3, R3
    
    MOVD    R3, ret+16(FP)
    RET

// func cmpGeInt64NEON(values *int64, threshold int64) uint64
TEXT ·cmpGeInt64NEON(SB), NOSPLIT, $0-24
    MOVD    values+0(FP), R0
    MOVD    threshold+8(FP), R1
    MOVD    $0, R3
    
    MOVD    0(R0), R4
    CMP     R1, R4
    CSET    GE, R5
    ORR     R5, R3, R3
    
    MOVD    8(R0), R4
    CMP     R1, R4
    CSET    GE, R5
    LSL     $1, R5, R5
    ORR     R5, R3, R3
    
    MOVD    R3, ret+16(FP)
    RET

// func cmpLeInt64NEON(values *int64, threshold int64) uint64
TEXT ·cmpLeInt64NEON(SB), NOSPLIT, $0-24
    MOVD    values+0(FP), R0
    MOVD    threshold+8(FP), R1
    MOVD    $0, R3
    
    MOVD    0(R0), R4
    CMP     R1, R4
    CSET    LE, R5
    ORR     R5, R3, R3
    
    MOVD    8(R0), R4
    CMP     R1, R4
    CSET    LE, R5
    LSL     $1, R5, R5
    ORR     R5, R3, R3
    
    MOVD    R3, ret+16(FP)
    RET

// func cmpNeInt64NEON(values *int64, threshold int64) uint64
TEXT ·cmpNeInt64NEON(SB), NOSPLIT, $0-24
    MOVD    values+0(FP), R0
    MOVD    threshold+8(FP), R1
    MOVD    $0, R3
    
    MOVD    0(R0), R4
    CMP     R1, R4
    CSET    NE, R5
    ORR     R5, R3, R3
    
    MOVD    8(R0), R4
    CMP     R1, R4
    CSET    NE, R5
    LSL     $1, R5, R5
    ORR     R5, R3, R3
    
    MOVD    R3, ret+16(FP)
    RET
