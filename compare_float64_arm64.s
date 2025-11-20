#include "textflag.h"

// func cmpGtFloat64NEON(values *float64, threshold float64) uint64
//
// ARM64 NEON compares 2 float64 values using scalar comparisons.
// NEON has limited 64-bit float SIMD support, so we use efficient scalar operations.
//
TEXT ·cmpGtFloat64NEON(SB), NOSPLIT, $0-24
    MOVD    values+0(FP), R0        // R0 = pointer to values array
    FMOVD   threshold+8(FP), F0     // F0 = threshold value
    MOVD    $0, R3                  // R3 = result mask
    
    // Load and compare first value
    FMOVD   0(R0), F1               // F1 = values[0]
    FCMPD   F0, F1                  // Compare F1 with F0
    CSET    GT, R5                  // R5 = 1 if F1 > F0, else 0
    ORR     R5, R3, R3              // Set bit 0
    
    // Load and compare second value
    FMOVD   8(R0), F1               // F1 = values[1]
    FCMPD   F0, F1
    CSET    GT, R5
    LSL     $1, R5, R5              // Shift to bit 1
    ORR     R5, R3, R3              // Set bit 1
    
    MOVD    R3, ret+16(FP)
    RET

// func cmpGeFloat64NEON(values *float64, threshold float64) uint64
TEXT ·cmpGeFloat64NEON(SB), NOSPLIT, $0-24
    MOVD    values+0(FP), R0
    FMOVD   threshold+8(FP), F0
    MOVD    $0, R3
    
    FMOVD   0(R0), F1
    FCMPD   F0, F1
    CSET    GE, R5
    ORR     R5, R3, R3
    
    FMOVD   8(R0), F1
    FCMPD   F0, F1
    CSET    GE, R5
    LSL     $1, R5, R5
    ORR     R5, R3, R3
    
    MOVD    R3, ret+16(FP)
    RET

// func cmpLtFloat64NEON(values *float64, threshold float64) uint64
TEXT ·cmpLtFloat64NEON(SB), NOSPLIT, $0-24
    MOVD    values+0(FP), R0
    FMOVD   threshold+8(FP), F0
    MOVD    $0, R3
    
    FMOVD   0(R0), F1
    FCMPD   F0, F1
    CSET    LT, R5
    ORR     R5, R3, R3
    
    FMOVD   8(R0), F1
    FCMPD   F0, F1
    CSET    LT, R5
    LSL     $1, R5, R5
    ORR     R5, R3, R3
    
    MOVD    R3, ret+16(FP)
    RET

// func cmpLeFloat64NEON(values *float64, threshold float64) uint64
TEXT ·cmpLeFloat64NEON(SB), NOSPLIT, $0-24
    MOVD    values+0(FP), R0
    FMOVD   threshold+8(FP), F0
    MOVD    $0, R3
    
    FMOVD   0(R0), F1
    FCMPD   F0, F1
    CSET    LE, R5
    ORR     R5, R3, R3
    
    FMOVD   8(R0), F1
    FCMPD   F0, F1
    CSET    LE, R5
    LSL     $1, R5, R5
    ORR     R5, R3, R3
    
    MOVD    R3, ret+16(FP)
    RET

// func cmpEqFloat64NEON(values *float64, threshold float64) uint64
TEXT ·cmpEqFloat64NEON(SB), NOSPLIT, $0-24
    MOVD    values+0(FP), R0
    FMOVD   threshold+8(FP), F0
    MOVD    $0, R3
    
    FMOVD   0(R0), F1
    FCMPD   F0, F1
    CSET    EQ, R5
    ORR     R5, R3, R3
    
    FMOVD   8(R0), F1
    FCMPD   F0, F1
    CSET    EQ, R5
    LSL     $1, R5, R5
    ORR     R5, R3, R3
    
    MOVD    R3, ret+16(FP)
    RET

// func cmpNeFloat64NEON(values *float64, threshold float64) uint64
TEXT ·cmpNeFloat64NEON(SB), NOSPLIT, $0-24
    MOVD    values+0(FP), R0
    FMOVD   threshold+8(FP), F0
    MOVD    $0, R3
    
    FMOVD   0(R0), F1
    FCMPD   F0, F1
    CSET    NE, R5
    ORR     R5, R3, R3
    
    FMOVD   8(R0), F1
    FCMPD   F0, F1
    CSET    NE, R5
    LSL     $1, R5, R5
    ORR     R5, R3, R3
    
    MOVD    R3, ret+16(FP)
    RET
