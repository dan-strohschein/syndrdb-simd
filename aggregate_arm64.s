#include "textflag.h"

// func sumInt64NEON(values *int64, length int) int64
//
// Computes the sum of int64 values using NEON SIMD.
// NEON processes 2 int64 values (128 bits) per iteration.
//
TEXT ·sumInt64NEON(SB), NOSPLIT, $0-24
	MOVD    values+0(FP), R0       // R0 = pointer to values array
	MOVD    length+8(FP), R1       // R1 = length
	
	// Check if length is 0
	CBZ     R1, return_zero
	
	// Initialize accumulator V0 = [0, 0]
	VEOR    V0.B16, V0.B16, V0.B16
	
	// Calculate number of full 2-element chunks
	MOVD    R1, R2
	LSR     $1, R2                 // R2 = length / 2
	CBZ     R2, remainder
	
loop:
	// Load 2 int64 values into V1
	VLD1    (R0), [V1.D2]          // V1 = [values[0], values[1]]
	
	// Add to accumulator
	VADD    V1.D2, V0.D2, V0.D2    // V0 += V1
	
	ADD     $16, R0                // Move pointer forward (2 * 8 bytes)
	SUB     $1, R2
	CBNZ    R2, loop
	
remainder:
	// Handle remaining element (0 or 1)
	AND     $1, R1, R2
	CBZ     R2, horizontal_sum
	
	// Load last element
	MOVD    (R0), R3
	
	// Extract V0[0] and add remainder
	VMOV    V0.D[0], R4
	ADD     R3, R4
	VMOV    R4, V0.D[0]
	
horizontal_sum:
	// Horizontal sum: add V0[0] + V0[1]
	VMOV    V0.D[0], R4
	VMOV    V0.D[1], R5
	ADD     R5, R4
	
	MOVD    R4, ret+16(FP)
	RET
	
return_zero:
	MOVD    $0, ret+16(FP)
	RET

// func minInt64NEON(values *int64, length int) int64
//
// Finds the minimum int64 value using NEON SIMD.
// NEON has SMIN instruction for minimum signed values.
//
TEXT ·minInt64NEON(SB), NOSPLIT, $0-24
	MOVD    values+0(FP), R0
	MOVD    length+8(FP), R1
	
	// Check if length is 0
	CBZ     R1, return_max_int64
	
	// Check if length < 2
	CMP     $2, R1
	BLT     single_element_min
	
	// Initialize min vector with first 2 values
	VLD1    (R0), [V0.D2]
	ADD     $16, R0
	SUB     $2, R1
	
	// Calculate number of remaining full chunks
	MOVD    R1, R2
	LSR     $1, R2
	CBZ     R2, remainder_min
	
loop_min:
	VLD1    (R0), [V1.D2]
	
	// Manual min using CMP for each lane
	// Compare V0[0] with V1[0]
	VMOV    V0.D[0], R4
	VMOV    V1.D[0], R5
	CMP     R5, R4
	CSEL    LT, R5, R4, R4
	VMOV    R4, V0.D[0]
	
	// Compare V0[1] with V1[1]
	VMOV    V0.D[1], R4
	VMOV    V1.D[1], R5
	CMP     R5, R4
	CSEL    LT, R5, R4, R4
	VMOV    R4, V0.D[1]
	
	ADD     $16, R0
	SUB     $1, R2
	CBNZ    R2, loop_min
	
remainder_min:
	// Handle remaining element
	AND     $1, R1, R2
	CBZ     R2, horizontal_min
	
	// Load last element
	MOVD    (R0), R3
	
	// Compare with V0[0]
	VMOV    V0.D[0], R4
	CMP     R3, R4
	CSEL    LT, R3, R4, R4
	VMOV    R4, V0.D[0]
	
horizontal_min:
	// Find min of V0[0] and V0[1]
	VMOV    V0.D[0], R4
	VMOV    V0.D[1], R5
	CMP     R5, R4
	CSEL    LT, R5, R4, R4
	
	MOVD    R4, ret+16(FP)
	RET
	
single_element_min:
	MOVD    (R0), R4
	MOVD    R4, ret+16(FP)
	RET
	
return_max_int64:
	// Load math.MaxInt64 using two 32-bit loads
	MOVD    $0x7FFFFFFF, R4
	LSL     $32, R4, R4
	MOVD    $0xFFFFFFFF, R5
	ORR     R5, R4, R4
	MOVD    R4, ret+16(FP)
	RET

// func maxInt64NEON(values *int64, length int) int64
//
// Finds the maximum int64 value using NEON SIMD.
// NEON has SMAX instruction for maximum signed values.
//
TEXT ·maxInt64NEON(SB), NOSPLIT, $0-24
	MOVD    values+0(FP), R0
	MOVD    length+8(FP), R1
	
	// Check if length is 0
	CBZ     R1, return_min_int64
	
	// Check if length < 2
	CMP     $2, R1
	BLT     single_element_max
	
	// Initialize max vector with first 2 values
	VLD1    (R0), [V0.D2]
	ADD     $16, R0
	SUB     $2, R1
	
	// Calculate number of remaining full chunks
	MOVD    R1, R2
	LSR     $1, R2
	CBZ     R2, remainder_max
	
loop_max:
	VLD1    (R0), [V1.D2]
	
	// Manual max using CMP for each lane
	// Compare V0[0] with V1[0]
	VMOV    V0.D[0], R4
	VMOV    V1.D[0], R5
	CMP     R5, R4
	CSEL    GT, R5, R4, R4
	VMOV    R4, V0.D[0]
	
	// Compare V0[1] with V1[1]
	VMOV    V0.D[1], R4
	VMOV    V1.D[1], R5
	CMP     R5, R4
	CSEL    GT, R5, R4, R4
	VMOV    R4, V0.D[1]
	
	ADD     $16, R0
	SUB     $1, R2
	CBNZ    R2, loop_max
	
remainder_max:
	// Handle remaining element
	AND     $1, R1, R2
	CBZ     R2, horizontal_max
	
	// Load last element
	MOVD    (R0), R3
	
	// Compare with V0[0]
	VMOV    V0.D[0], R4
	CMP     R3, R4
	CSEL    GT, R3, R4, R4
	VMOV    R4, V0.D[0]
	
horizontal_max:
	// Find max of V0[0] and V0[1]
	VMOV    V0.D[0], R4
	VMOV    V0.D[1], R5
	CMP     R5, R4
	CSEL    GT, R5, R4, R4
	
	MOVD    R4, ret+16(FP)
	RET
	
single_element_max:
	MOVD    (R0), R4
	MOVD    R4, ret+16(FP)
	RET
	
return_min_int64:
	// Load math.MinInt64 = -9223372036854775808 = 1 << 63
	MOVD    $1, R4
	LSL     $63, R4, R4
	MOVD    R4, ret+16(FP)
	RET

// func countNonNullNEON(values *int64, nullBitmap *uint64, length int) int64
//
// Counts non-null values.
// For each bit in nullBitmap: 0 = not null, 1 = null.
//
TEXT ·countNonNullNEON(SB), NOSPLIT, $0-32
	MOVD    values+0(FP), R0
	MOVD    nullBitmap+8(FP), R1
	MOVD    length+16(FP), R2
	
	// Check if length is 0
	CBZ     R2, return_zero_count
	
	MOVD    $0, R3                 // R3 = count
	MOVD    $0, R4                 // R4 = current bit index
	
	// Check if nullBitmap is provided
	CBZ     R1, count_all
	
count_loop:
	// Get word index: R5 = bit_index / 64
	MOVD    R4, R5
	LSR     $6, R5
	
	// Load nullBitmap[word_index]
	MOVD    (R1)(R5<<3), R6
	
	// Get bit position: R7 = bit_index % 64
	AND     $63, R4, R7
	
	// Check if bit is set
	LSR     R7, R6, R8
	AND     $1, R8, R8
	
	// If bit is 0 (not null), increment count
	CMP     $0, R8
	CINC    EQ, R3, R3
	
	ADD     $1, R4
	CMP     R4, R2
	BLT     count_loop
	
	JMP     return_count
	
count_all:
	// No null bitmap, count = length
	MOVD    R2, R3
	
return_count:
	MOVD    R3, ret+24(FP)
	RET
	
return_zero_count:
	MOVD    $0, ret+24(FP)
	RET
