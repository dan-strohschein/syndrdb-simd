#include "textflag.h"

// func andBitmapAVX2(dst, a, b *uint64, length int)
//
// Performs bitwise AND operation on two uint64 arrays using AVX2.
// Processes 4 uint64 values (256 bits) per iteration for maximum throughput.
//
// The operation: dst[i] = a[i] & b[i] for all i in [0, length)
//
// AVX2 advantages:
// - Processes 4 uint64s per instruction vs 1 in scalar code (4x throughput)
// - Single VPAND instruction replaces 4 scalar AND operations
// - Better memory bandwidth utilization with 256-bit loads/stores
//
// Register usage:
//   SI: pointer to destination array (dst)
//   DI: pointer to source array a
//   DX: pointer to source array b
//   CX: remaining elements to process
//   YMM0: loaded values from array a
//   YMM1: loaded values from array b
//   YMM2: result of AND operation
//
TEXT ·andBitmapAVX2(SB), NOSPLIT, $0-32
    // Load arguments from Go stack frame
    MOVQ    dst+0(FP), SI           // SI = destination pointer
    MOVQ    a+8(FP), DI             // DI = first source array pointer
    MOVQ    b+16(FP), DX            // DX = second source array pointer
    MOVQ    length+24(FP), CX       // CX = number of uint64 elements
    
    // Check if we have enough elements to process with SIMD (at least 4)
    CMPQ    CX, $4
    JL      and_remainder           // If less than 4 elements, skip to scalar processing
    
and_loop:
    // Load 4 uint64 values from array a (256 bits total)
    // VMOVDQU handles unaligned memory access gracefully
    VMOVDQU (DI), Y0                // Y0 = [a[0], a[1], a[2], a[3]]
    
    // Load 4 uint64 values from array b
    VMOVDQU (DX), Y1                // Y1 = [b[0], b[1], b[2], b[3]]
    
    // Perform bitwise AND on all 4 pairs simultaneously
    // This single instruction replaces 4 scalar AND operations
    VPAND   Y0, Y1, Y2              // Y2[i] = Y0[i] & Y1[i] for i in [0,3]
    
    // Store 4 results back to destination array
    VMOVDQU Y2, (SI)                // dst[0:3] = Y2
    
    // Advance all pointers by 32 bytes (4 × 8 bytes per uint64)
    ADDQ    $32, SI                 // Move dst pointer forward
    ADDQ    $32, DI                 // Move a pointer forward
    ADDQ    $32, DX                 // Move b pointer forward
    
    // Decrease counter by 4 elements processed
    SUBQ    $4, CX
    
    // Continue if we still have at least 4 elements remaining
    CMPQ    CX, $4
    JGE     and_loop
    
and_remainder:
    // Handle remaining elements (0-3) with scalar operations
    // This ensures we process arrays that aren't multiples of 4
    TESTQ   CX, CX                  // Check if any elements remain
    JZ      and_done                // If zero, we're finished
    
and_remainder_loop:
    // Scalar processing: load one uint64 at a time
    MOVQ    (DI), AX                // AX = a[i]
    MOVQ    (DX), BX                // BX = b[i]
    ANDQ    BX, AX                  // AX = AX & BX
    MOVQ    AX, (SI)                // dst[i] = AX
    
    // Advance pointers by 8 bytes (1 uint64)
    ADDQ    $8, SI
    ADDQ    $8, DI
    ADDQ    $8, DX
    
    // Decrement counter and loop
    DECQ    CX
    JNZ     and_remainder_loop
    
and_done:
    // Clean up AVX state before returning
    VZEROUPPER
    RET

// func orBitmapAVX2(dst, a, b *uint64, length int)
//
// Performs bitwise OR operation on two uint64 arrays using AVX2.
// Implementation mirrors andBitmapAVX2 but uses VPOR instead of VPAND.
//
TEXT ·orBitmapAVX2(SB), NOSPLIT, $0-32
    MOVQ    dst+0(FP), SI
    MOVQ    a+8(FP), DI
    MOVQ    b+16(FP), DX
    MOVQ    length+24(FP), CX
    
    CMPQ    CX, $4
    JL      or_remainder
    
or_loop:
    VMOVDQU (DI), Y0
    VMOVDQU (DX), Y1
    
    // Bitwise OR on all 4 uint64 pairs
    VPOR    Y0, Y1, Y2              // Y2[i] = Y0[i] | Y1[i]
    
    VMOVDQU Y2, (SI)
    
    ADDQ    $32, SI
    ADDQ    $32, DI
    ADDQ    $32, DX
    SUBQ    $4, CX
    
    CMPQ    CX, $4
    JGE     or_loop
    
or_remainder:
    TESTQ   CX, CX
    JZ      or_done
    
or_remainder_loop:
    MOVQ    (DI), AX
    MOVQ    (DX), BX
    ORQ     BX, AX                  // Scalar OR operation
    MOVQ    AX, (SI)
    
    ADDQ    $8, SI
    ADDQ    $8, DI
    ADDQ    $8, DX
    DECQ    CX
    JNZ     or_remainder_loop
    
or_done:
    VZEROUPPER
    RET

// func xorBitmapAVX2(dst, a, b *uint64, length int)
//
// Performs bitwise XOR operation on two uint64 arrays using AVX2.
// XOR is particularly useful for bitmap diff operations.
//
TEXT ·xorBitmapAVX2(SB), NOSPLIT, $0-32
    MOVQ    dst+0(FP), SI
    MOVQ    a+8(FP), DI
    MOVQ    b+16(FP), DX
    MOVQ    length+24(FP), CX
    
    CMPQ    CX, $4
    JL      xor_remainder
    
xor_loop:
    VMOVDQU (DI), Y0
    VMOVDQU (DX), Y1
    
    // Bitwise XOR on all 4 uint64 pairs
    VPXOR   Y0, Y1, Y2              // Y2[i] = Y0[i] ^ Y1[i]
    
    VMOVDQU Y2, (SI)
    
    ADDQ    $32, SI
    ADDQ    $32, DI
    ADDQ    $32, DX
    SUBQ    $4, CX
    
    CMPQ    CX, $4
    JGE     xor_loop
    
xor_remainder:
    TESTQ   CX, CX
    JZ      xor_done
    
xor_remainder_loop:
    MOVQ    (DI), AX
    MOVQ    (DX), BX
    XORQ    BX, AX                  // Scalar XOR operation
    MOVQ    AX, (SI)
    
    ADDQ    $8, SI
    ADDQ    $8, DI
    ADDQ    $8, DX
    DECQ    CX
    JNZ     xor_remainder_loop
    
xor_done:
    VZEROUPPER
    RET

// func notBitmapAVX2(dst, src *uint64, length int)
//
// Performs bitwise NOT operation on a uint64 array using AVX2.
// The NOT operation inverts all bits: dst[i] = ^src[i]
//
// Implementation note:
// AVX2 doesn't have a direct NOT instruction, so we XOR with all 1s.
// We create an all-ones vector using VPCMPEQQ (compare equal to itself).
//
TEXT ·notBitmapAVX2(SB), NOSPLIT, $0-24
    MOVQ    dst+0(FP), SI           // SI = destination pointer
    MOVQ    src+8(FP), DI           // DI = source pointer
    MOVQ    length+16(FP), CX       // CX = number of elements
    
    // Create an all-ones vector for XOR-based NOT operation
    // VPCMPEQQ compares a register with itself, always producing all 1s
    VPCMPEQQ Y3, Y3, Y3             // Y3 = 0xFFFFFFFFFFFFFFFF (all lanes)
    
    CMPQ    CX, $4
    JL      not_remainder
    
not_loop:
    VMOVDQU (DI), Y0                // Load 4 uint64 values
    
    // XOR with all 1s to flip all bits (equivalent to NOT)
    VPXOR   Y3, Y0, Y1              // Y1[i] = Y0[i] ^ 0xFFFF... = ~Y0[i]
    
    VMOVDQU Y1, (SI)                // Store inverted values
    
    ADDQ    $32, SI
    ADDQ    $32, DI
    SUBQ    $4, CX
    
    CMPQ    CX, $4
    JGE     not_loop
    
not_remainder:
    TESTQ   CX, CX
    JZ      not_done
    
not_remainder_loop:
    MOVQ    (DI), AX
    NOTQ    AX                      // Scalar NOT operation
    MOVQ    AX, (SI)
    
    ADDQ    $8, SI
    ADDQ    $8, DI
    DECQ    CX
    JNZ     not_remainder_loop
    
not_done:
    VZEROUPPER
    RET

// func popCountAVX2(bitmap *uint64, length int) int
//
// Counts the number of set bits (1s) in a bitmap using AVX2.
// This is a critical operation for database query optimization.
//
// Strategy:
// AVX2 doesn't have a native POPCNT instruction for vectors, so we:
// 1. Use scalar POPCNT on each uint64 (available since SSE4.2)
// 2. Accumulate results in a register
//
// TODO: I could optimize this further with lookup tables or AVX512 VPOPCNT
// when AVX512VPOPCNTDQ is available. For now, this provides good performance
// by unrolling the loop and processing multiple elements between accumulations.
//
TEXT ·popCountAVX2(SB), NOSPLIT, $0-24
    MOVQ    bitmap+0(FP), SI        // SI = bitmap pointer
    MOVQ    length+8(FP), CX        // CX = number of uint64 elements
    XORQ    AX, AX                  // AX = accumulator (initialize to 0)
    
    // Process in chunks of 4 for better instruction-level parallelism
    CMPQ    CX, $4
    JL      pop_remainder
    
pop_loop:
    // Load and count 4 uint64 values
    // POPCNT instruction counts set bits in a 64-bit register
    MOVQ    (SI), BX
    POPCNTQ BX, BX                  // BX = popcount(bitmap[0])
    ADDQ    BX, AX                  // Accumulate
    
    MOVQ    8(SI), BX
    POPCNTQ BX, BX                  // BX = popcount(bitmap[1])
    ADDQ    BX, AX
    
    MOVQ    16(SI), BX
    POPCNTQ BX, BX                  // BX = popcount(bitmap[2])
    ADDQ    BX, AX
    
    MOVQ    24(SI), BX
    POPCNTQ BX, BX                  // BX = popcount(bitmap[3])
    ADDQ    BX, AX
    
    ADDQ    $32, SI                 // Advance pointer by 4 uint64s
    SUBQ    $4, CX
    
    CMPQ    CX, $4
    JGE     pop_loop
    
pop_remainder:
    TESTQ   CX, CX
    JZ      pop_done
    
pop_remainder_loop:
    MOVQ    (SI), BX
    POPCNTQ BX, BX
    ADDQ    BX, AX
    
    ADDQ    $8, SI
    DECQ    CX
    JNZ     pop_remainder_loop
    
pop_done:
    // Return accumulated count
    MOVQ    AX, ret+16(FP)
    RET
