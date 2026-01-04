#include "textflag.h"

// func sumInt64AVX2(values *int64, length int) int64
//
// Computes the sum of int64 values using AVX2 SIMD.
// AVX2 allows processing 4 int64 values (256 bits) per iteration.
//
// Algorithm:
// 1. Initialize 4 accumulators in a YMM register (all zeros)
// 2. Process array in chunks of 4, adding to accumulators
// 3. Handle remainder with scalar loop
// 4. Horizontally sum the 4 accumulators
// 5. Return total sum
//
TEXT ·sumInt64AVX2(SB), NOSPLIT, $0-24
	MOVQ    values+0(FP), SI       // SI = pointer to values array
	MOVQ    length+8(FP), CX       // CX = length
	
	// Check if length is 0
	TESTQ   CX, CX
	JZ      return_zero
	
	// Initialize accumulator: 4 int64 sums in YMM0
	VPXOR   Y0, Y0, Y0            // Y0 = [0, 0, 0, 0]
	
	// Calculate number of full 4-element chunks
	MOVQ    CX, DX
	SHRQ    $2, DX                // DX = length / 4
	JZ      remainder             // If no full chunks, go to remainder
	
loop:
	// Load 4 int64 values into YMM1
	VMOVDQU  0(SI), Y1             // Y1 = [values[0], values[1], values[2], values[3]]
	
	// Add to accumulator
	VPADDQ   Y1, Y0, Y0            // Y0 += Y1
	
	ADDQ     $32, SI               // Move pointer forward (4 * 8 bytes)
	DECQ     DX
	JNZ      loop
	
remainder:
	// Handle remaining elements (0-3)
	MOVQ     CX, DX
	ANDQ     $3, DX                // DX = length % 4
	JZ       horizontal_sum        // No remainder
	
	XORQ     R8, R8                // R8 = scalar sum for remainder
	
remainder_loop:
	ADDQ     0(SI), R8             // R8 += *SI
	ADDQ     $8, SI
	DECQ     DX
	JNZ      remainder_loop
	
horizontal_sum:
	// Horizontal sum of YMM0: sum all 4 int64 values
	// Extract high 128 bits to XMM1
	VEXTRACTI128 $1, Y0, X1         // X1 = Y0[127:64]
	
	// Add high and low 128-bit halves
	VPADDQ   X1, X0, X0            // X0 = X0 + X1 (now 2 int64 sums)
	
	// Horizontal add within 128 bits
	VPSHUFD  $0xEE, X0, X1         // X1 = X0 with high 64 bits duplicated
	VPADDQ   X1, X0, X0            // X0[0] = sum of both elements
	
	// Extract result to scalar register
	VMOVQ    X0, AX                // AX = final sum
	
	// Add remainder sum if any
	ADDQ     R8, AX
	
	// Clean up YMM registers
	VZEROUPPER
	
	MOVQ     AX, ret+16(FP)
	RET
	
return_zero:
	MOVQ     $0, ret+16(FP)
	RET

// func minInt64AVX2(values *int64, length int) int64
//
// Finds the minimum int64 value using AVX2 SIMD.
// AVX2 doesn't have a direct VPMINSQ (min signed 64-bit) until AVX-512.
// We use VPCMPGTQ + VPBLENDVB to simulate min operation.
//
TEXT ·minInt64AVX2(SB), NOSPLIT, $0-24
	MOVQ    values+0(FP), SI
	MOVQ    length+8(FP), CX
	
	// Check if length is 0
	TESTQ   CX, CX
	JZ      return_max_int64
	
	// Initialize min vector with first 4 values (or broadcast first if length < 4)
	CMPQ    CX, $4
	JGE     init_full
	
	// Length < 4: broadcast first value
	VPBROADCASTQ 0(SI), Y0          // Y0 = [values[0], values[0], values[0], values[0]]
	JMP     remainder_min
	
init_full:
	VMOVDQU  0(SI), Y0              // Y0 = [values[0], values[1], values[2], values[3]]
	ADDQ     $32, SI
	SUBQ     $4, CX
	
	// Calculate number of remaining full chunks
	MOVQ     CX, DX
	SHRQ     $2, DX
	JZ       remainder_min
	
loop_min:
	VMOVDQU  0(SI), Y1              // Load next 4 values
	
	// Compare: Y2 = (Y0 > Y1) ? 0xFFFF... : 0
	VPCMPGTQ Y1, Y0, Y2
	
	// Blend: if Y0 > Y1, select Y1 (smaller), else select Y0
	VPBLENDVB Y2, Y1, Y0, Y0
	
	ADDQ     $32, SI
	DECQ     DX
	JNZ      loop_min
	
remainder_min:
	// First, find min across the 4 lanes we have
	// Extract high 128 bits
	VEXTRACTI128 $1, Y0, X1
	
	// Min of high and low halves: if X0 > X1, select X1
	VPCMPGTQ X1, X0, X2
	VPBLENDVB X2, X1, X0, X0
	
	// Min within 128 bits (2 int64 values)
	VPSHUFD  $0xEE, X0, X1
	VPCMPGTQ X1, X0, X2
	VPBLENDVB X2, X1, X0, X0
	
	// Extract current min to R8
	VMOVQ    X0, R8
	
	// Handle remaining elements
	MOVQ     CX, DX
	ANDQ     $3, DX
	JZ       return_min_result
	
remainder_min_loop:
	MOVQ     0(SI), R9
	CMPQ     R9, R8
	CMOVQLT  R9, R8                // R8 = min(R8, R9)
	ADDQ     $8, SI
	DECQ     DX
	JNZ      remainder_min_loop
	
return_min_result:
	MOVQ     R8, AX
	VZEROUPPER
	MOVQ     AX, ret+16(FP)
	RET
	
horizontal_min:
	// Find minimum across 4 lanes of Y0 (no remainder case)
	// Extract high 128 bits
	VEXTRACTI128 $1, Y0, X1
	
	// Min of high and low halves: if X0 > X1, select X1
	VPCMPGTQ X1, X0, X2
	VPBLENDVB X2, X1, X0, X0
	
	// Min within 128 bits (2 int64 values)
	VPSHUFD  $0xEE, X0, X1
	VPCMPGTQ X1, X0, X2
	VPBLENDVB X2, X1, X0, X0
	
	// Extract result
	VMOVQ    X0, AX
	
cleanup_min:
	VZEROUPPER
	MOVQ     AX, ret+16(FP)
	RET
	
return_max_int64:
	MOVQ     $0x7FFFFFFFFFFFFFFF, AX
	MOVQ     AX, ret+16(FP)  // math.MaxInt64
	RET

// func maxInt64AVX2(values *int64, length int) int64
//
// Finds the maximum int64 value using AVX2 SIMD.
// Uses VPCMPGTQ + VPBLENDVB to simulate max operation.
//
TEXT ·maxInt64AVX2(SB), NOSPLIT, $0-24
	MOVQ    values+0(FP), SI
	MOVQ    length+8(FP), CX
	
	// Check if length is 0
	TESTQ   CX, CX
	JZ      return_min_int64
	
	// Initialize max vector with first 4 values (or broadcast first if length < 4)
	CMPQ    CX, $4
	JGE     init_full_max
	
	// Length < 4: broadcast first value
	VPBROADCASTQ 0(SI), Y0
	JMP     remainder_max
	
init_full_max:
	VMOVDQU  0(SI), Y0
	ADDQ     $32, SI
	SUBQ     $4, CX
	
	// Calculate number of remaining full chunks
	MOVQ     CX, DX
	SHRQ     $2, DX
	JZ       remainder_max
	
loop_max:
	VMOVDQU  0(SI), Y1
	
	// Compare: Y2 = (Y1 > Y0) ? 0xFFFF... : 0
	VPCMPGTQ Y0, Y1, Y2
	
	// Blend: if Y1 > Y0, select Y1 (larger), else select Y0
	VPBLENDVB Y2, Y1, Y0, Y0
	
	ADDQ     $32, SI
	DECQ     DX
	JNZ      loop_max
	
remainder_max:
	// First, find max across the 4 lanes we have
	VEXTRACTI128 $1, Y0, X1
	
	// Max of high and low halves: if X1 > X0, select X1
	VPCMPGTQ X0, X1, X2
	VPBLENDVB X2, X1, X0, X0
	
	// Max within 128 bits
	VPSHUFD  $0xEE, X0, X1
	VPCMPGTQ X0, X1, X2
	VPBLENDVB X2, X1, X0, X0
	
	// Extract current max to R8
	VMOVQ    X0, R8
	
	// Handle remaining elements
	MOVQ     CX, DX
	ANDQ     $3, DX
	JZ       return_max_result
	
remainder_max_loop:
	MOVQ     0(SI), R9
	CMPQ     R9, R8
	CMOVQGT  R9, R8                // R8 = max(R8, R9)
	ADDQ     $8, SI
	DECQ     DX
	JNZ      remainder_max_loop
	
return_max_result:
	MOVQ     R8, AX
	VZEROUPPER
	MOVQ     AX, ret+16(FP)
	RET
	
horizontal_max:
	// Find maximum across 4 lanes (no remainder case)
	VEXTRACTI128 $1, Y0, X1
	
	// Max of high and low halves: if X1 > X0, select X1
	VPCMPGTQ X0, X1, X2
	VPBLENDVB X2, X1, X0, X0
	
	// Max within 128 bits
	VPSHUFD  $0xEE, X0, X1
	VPCMPGTQ X0, X1, X2
	VPBLENDVB X2, X1, X0, X0
	
	// Extract result
	VMOVQ    X0, AX
	
cleanup_max:
	VZEROUPPER
	MOVQ     AX, ret+16(FP)
	RET
	
return_min_int64:
	MOVQ     $0x8000000000000000, AX
	MOVQ     AX, ret+16(FP)  // math.MinInt64
	RET

// func countNonNullAVX2(values *int64, nullBitmap *uint64, length int) int64
//
// Counts non-null values using AVX2.
// For each bit in nullBitmap: 0 = not null (count it), 1 = null (skip it).
//
TEXT ·countNonNullAVX2(SB), NOSPLIT, $0-32
	MOVQ    values+0(FP), SI
	MOVQ    nullBitmap+8(FP), DI
	MOVQ    length+16(FP), CX
	
	// Check if length is 0
	TESTQ   CX, CX
	JZ      return_zero_count
	
	XORQ     AX, AX                // AX = count
	XORQ     R8, R8                // R8 = current bit index
	
	// Check if nullBitmap is provided
	TESTQ    DI, DI
	JZ       count_all             // No null bitmap, count all
	
count_loop:
	// Get the uint64 word containing the bit
	MOVQ     R8, R9
	SHRQ     $6, R9                // R9 = bit_index / 64 (word index)
	MOVQ     0(DI)(R9*8), R10      // R10 = nullBitmap[word_index]
	
	// Get bit position within the word
	MOVQ     R8, R11
	ANDQ     $63, R11              // R11 = bit_index % 64
	
	// Check if bit is set (null)
	BTQ      R11, R10              // Test bit
	JC       skip_null             // If carry (bit=1), it's null, skip
	
	// Not null, increment count
	INCQ     AX
	
skip_null:
	INCQ     R8
	CMPQ     R8, CX
	JL       count_loop
	
	JMP      return_count
	
count_all:
	// No null bitmap provided, count = length
	MOVQ     CX, AX
	
return_count:
	MOVQ     AX, ret+24(FP)
	RET
	
return_zero_count:
	MOVQ     $0, ret+24(FP)
	RET
