#include "textflag.h"

// func andBitmapNEON(dst, a, b *uint64, length int)
//
// Performs bitwise AND operation on two uint64 arrays using NEON.
// Processes 2 uint64 values (128 bits) per iteration.
//
TEXT ·andBitmapNEON(SB), NOSPLIT, $0-32
    MOVD    dst+0(FP), R0           // R0 = destination pointer
    MOVD    a+8(FP), R1             // R1 = first source array pointer
    MOVD    b+16(FP), R2            // R2 = second source array pointer
    MOVD    length+24(FP), R3       // R3 = number of uint64 elements
    
    CMP     $2, R3
    BLT     and_remainder
    
and_loop:
    // Load 2 uint64 values from array a
    VLD1    (R1), [V0.D2]           // V0 = [a[0], a[1]]
    
    // Load 2 uint64 values from array b
    VLD1    (R2), [V1.D2]           // V1 = [b[0], b[1]]
    
    // Bitwise AND on both uint64 pairs
    VAND    V0.B16, V1.B16, V2.B16  // V2 = V0 & V1
    
    // Store 2 results back to destination
    VST1    [V2.D2], (R0)           // dst[0:1] = V2
    
    // Advance all pointers by 16 bytes (2 × 8 bytes per uint64)
    ADD     $16, R0
    ADD     $16, R1
    ADD     $16, R2
    SUB     $2, R3
    
    CMP     $2, R3
    BGE     and_loop
    
and_remainder:
    // Handle remaining elements (0-1) with scalar operations
    CBZ     R3, and_done
    
and_remainder_loop:
    MOVD    (R1), R4                // R4 = a[i]
    MOVD    (R2), R5                // R5 = b[i]
    AND     R5, R4, R4              // R4 = R4 & R5
    MOVD    R4, (R0)                // dst[i] = R4
    
    ADD     $8, R0
    ADD     $8, R1
    ADD     $8, R2
    SUB     $1, R3
    CBNZ    R3, and_remainder_loop
    
and_done:
    RET

// func orBitmapNEON(dst, a, b *uint64, length int)
TEXT ·orBitmapNEON(SB), NOSPLIT, $0-32
    MOVD    dst+0(FP), R0
    MOVD    a+8(FP), R1
    MOVD    b+16(FP), R2
    MOVD    length+24(FP), R3
    
    CMP     $2, R3
    BLT     or_remainder
    
or_loop:
    VLD1    (R1), [V0.D2]
    VLD1    (R2), [V1.D2]
    
    // Bitwise OR
    VORR    V0.B16, V1.B16, V2.B16  // V2 = V0 | V1
    
    VST1    [V2.D2], (R0)
    
    ADD     $16, R0
    ADD     $16, R1
    ADD     $16, R2
    SUB     $2, R3
    
    CMP     $2, R3
    BGE     or_loop
    
or_remainder:
    CBZ     R3, or_done
    
or_remainder_loop:
    MOVD    (R1), R4
    MOVD    (R2), R5
    ORR     R5, R4, R4              // Scalar OR
    MOVD    R4, (R0)
    
    ADD     $8, R0
    ADD     $8, R1
    ADD     $8, R2
    SUB     $1, R3
    CBNZ    R3, or_remainder_loop
    
or_done:
    RET

// func xorBitmapNEON(dst, a, b *uint64, length int)
TEXT ·xorBitmapNEON(SB), NOSPLIT, $0-32
    MOVD    dst+0(FP), R0
    MOVD    a+8(FP), R1
    MOVD    b+16(FP), R2
    MOVD    length+24(FP), R3
    
    CMP     $2, R3
    BLT     xor_remainder
    
xor_loop:
    VLD1    (R1), [V0.D2]
    VLD1    (R2), [V1.D2]
    
    // Bitwise XOR
    VEOR    V0.B16, V1.B16, V2.B16  // V2 = V0 ^ V1
    
    VST1    [V2.D2], (R0)
    
    ADD     $16, R0
    ADD     $16, R1
    ADD     $16, R2
    SUB     $2, R3
    
    CMP     $2, R3
    BGE     xor_loop
    
xor_remainder:
    CBZ     R3, xor_done
    
xor_remainder_loop:
    MOVD    (R1), R4
    MOVD    (R2), R5
    EOR     R5, R4, R4              // Scalar XOR
    MOVD    R4, (R0)
    
    ADD     $8, R0
    ADD     $8, R1
    ADD     $8, R2
    SUB     $1, R3
    CBNZ    R3, xor_remainder_loop
    
xor_done:
    RET

// func notBitmapNEON(dst, src *uint64, length int)
TEXT ·notBitmapNEON(SB), NOSPLIT, $0-24
    MOVD    dst+0(FP), R0           // R0 = destination pointer
    MOVD    src+8(FP), R1           // R1 = source pointer
    MOVD    length+16(FP), R3       // R3 = number of elements
    
    CMP     $2, R3
    BLT     not_remainder
    
not_loop:
    VLD1    (R1), [V0.D2]           // Load 2 uint64 values
    
    // For NOT, we need all-ones to XOR with
    // Create using VMOVI (move immediate to vector)
    // VMOVI can create various patterns - use #0xFF
    // Alternatively, for maximum compatibility, just use scalar NOT in loop
    // Let's process with scalar operations for maximum compatibility
    
    // Extract first uint64, NOT it, store it
    VMOV    V0.D[0], R4
    MVN     R4, R4
    VMOV    R4, V1.D[0]
    
    // Extract second uint64, NOT it, store it
    VMOV    V0.D[1], R4
    MVN     R4, R4
    VMOV    R4, V1.D[1]
    
    VST1    [V1.D2], (R0)           // Store inverted values
    
    ADD     $16, R0
    ADD     $16, R1
    SUB     $2, R3
    
    CMP     $2, R3
    BGE     not_loop
    
not_remainder:
    CBZ     R3, not_done
    
not_remainder_loop:
    MOVD    (R1), R4
    MVN     R4, R4                  // Scalar NOT
    MOVD    R4, (R0)
    
    ADD     $8, R0
    ADD     $8, R1
    SUB     $1, R3
    CBNZ    R3, not_remainder_loop
    
not_done:
    RET

// func popCountNEON(bitmap *uint64, length int) int
//
// Counts set bits in a bitmap using NEON CNT instruction.
// CNT counts bits per byte, so we need to sum the results.
//
TEXT ·popCountNEON(SB), NOSPLIT, $0-24
    MOVD    bitmap+0(FP), R0        // R0 = bitmap pointer
    MOVD    length+8(FP), R3        // R3 = number of uint64 elements
    MOVD    $0, R4                  // R4 = accumulator
    
    CMP     $2, R3
    BLT     pop_remainder
    
pop_loop:
    // Load 2 uint64 values (128 bits)
    VLD1    (R0), [V0.D2]           // V0 = [bitmap[0], bitmap[1]]
    
    // Count bits in each byte
    // CNT (or VCNT) counts population per byte
    VCNT    V0.B16, V1.B16          // V1[i] = popcount(V0_byte[i]) for each of 16 bytes
    
    // Now we need to sum all 16 bytes in V1
    // Extract each byte and add to accumulator
    // This is verbose but works with Go's assembler
    VMOV    V1.B[0], R5
    ADD     R5, R4
    VMOV    V1.B[1], R5
    ADD     R5, R4
    VMOV    V1.B[2], R5
    ADD     R5, R4
    VMOV    V1.B[3], R5
    ADD     R5, R4
    VMOV    V1.B[4], R5
    ADD     R5, R4
    VMOV    V1.B[5], R5
    ADD     R5, R4
    VMOV    V1.B[6], R5
    ADD     R5, R4
    VMOV    V1.B[7], R5
    ADD     R5, R4
    VMOV    V1.B[8], R5
    ADD     R5, R4
    VMOV    V1.B[9], R5
    ADD     R5, R4
    VMOV    V1.B[10], R5
    ADD     R5, R4
    VMOV    V1.B[11], R5
    ADD     R5, R4
    VMOV    V1.B[12], R5
    ADD     R5, R4
    VMOV    V1.B[13], R5
    ADD     R5, R4
    VMOV    V1.B[14], R5
    ADD     R5, R4
    VMOV    V1.B[15], R5
    ADD     R5, R4
    
    ADD     $16, R0
    SUB     $2, R3
    
    CMP     $2, R3
    BGE     pop_loop
    
pop_remainder:
    CBZ     R3, pop_done
    
pop_remainder_loop:
    MOVD    (R0), R5                // Load one uint64
    
    // Use VCNT by moving to vector register
    VMOV    R5, V0.D[0]
    VCNT    V0.B8, V1.B8            // Count bits in 8 bytes
    
    // Extract and sum 8 bytes
    VMOV    V1.B[0], R6
    ADD     R6, R4
    VMOV    V1.B[1], R6
    ADD     R6, R4
    VMOV    V1.B[2], R6
    ADD     R6, R4
    VMOV    V1.B[3], R6
    ADD     R6, R4
    VMOV    V1.B[4], R6
    ADD     R6, R4
    VMOV    V1.B[5], R6
    ADD     R6, R4
    VMOV    V1.B[6], R6
    ADD     R6, R4
    VMOV    V1.B[7], R6
    ADD     R6, R4
    
    ADD     $8, R0
    SUB     $1, R3
    CBNZ    R3, pop_remainder_loop
    
pop_done:
    // Return accumulated count
    MOVD    R4, ret+16(FP)
    RET
