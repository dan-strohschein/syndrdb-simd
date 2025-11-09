// +build arm64

#include "textflag.h"

// func hashInt64NEON(values *int64, output *uint64, count int)
// Computes FNV-1a hash for int64 values using NEON
TEXT ·hashInt64NEON(SB), NOSPLIT, $0-24
	MOVD values+0(FP), R0    // R0 = &values[0]
	MOVD output+8(FP), R1    // R1 = &output[0]
	MOVD count+16(FP), R2    // R2 = count

	// FNV-1a constants
	// offset64 = 14695981039346656037 = 0xCBF29CE484222325
	// prime64  = 1099511628211 = 0x100000001B3
	
	// Load FNV offset constant
	MOVD $0x4222325, R3
	MOVD $0xCBF29CE4, R4
	LSL $32, R4, R4
	ORR R4, R3, R3           // R3 = FNV offset
	
	// Load FNV prime constant
	MOVD $0x1B3, R4
	MOVD $0x100000001, R5
	LSL $32, R5, R5
	ORR R5, R4, R4           // R4 = FNV prime

	// Process 2 elements at a time
	LSR $1, R2, R2           // R2 = count / 2
	CBZ R2, remainder

loop:
	// Load 2 int64 values
	VLD1 (R0), [V2.D2]       // V2 = [v0, v1]
	
	// Initialize hash with offset - use scalar approach
	// Process each value individually using scalar operations
	VMOV V2.D[0], R5         // R5 = v0
	VMOV V2.D[1], R6         // R6 = v1
	
	// Hash v0
	MOVD R3, R7              // R7 = FNV offset
	MOVD $8, R8              // Byte counter
hash_loop_0:
	AND $0xFF, R5, R9        // Extract byte
	EOR R9, R7, R7           // XOR byte
	MUL R4, R7, R7           // Multiply by prime
	LSR $8, R5, R5           // Shift value right
	SUBS $1, R8, R8
	BNE hash_loop_0
	
	// Store result 0
	MOVD R7, (R1)
	
	// Hash v1
	MOVD R3, R7              // R7 = FNV offset
	MOVD $8, R8              // Byte counter
hash_loop_1:
	AND $0xFF, R6, R9        // Extract byte
	EOR R9, R7, R7           // XOR byte
	MUL R4, R7, R7           // Multiply by prime
	LSR $8, R6, R6           // Shift value right
	SUBS $1, R8, R8
	BNE hash_loop_1
	
	// Store result 1
	MOVD R7, 8(R1)
	
	ADD $16, R0, R0          // Advance values pointer (2 * 8 bytes)
	ADD $16, R1, R1          // Advance output pointer (2 * 8 bytes)
	SUBS $1, R2, R2
	BNE loop

remainder:
	RET

// func crc32Int64NEON(values *int64, output *uint32, count int)
// Computes CRC32C checksums for int64 values using ARM64 CRC32
TEXT ·crc32Int64NEON(SB), NOSPLIT, $0-24
	MOVD values+0(FP), R0    // R0 = &values[0]
	MOVD output+8(FP), R1    // R1 = &output[0]
	MOVD count+16(FP), R2    // R2 = count

loop_crc:
	CBZ R2, done_crc
	
	// Load int64 value
	MOVD (R0), R3
	
	// Compute CRC32C using hardware instruction
	MOVD $0, R4              // Initialize CRC to 0
	CRC32CX R3, R4           // CRC32C of 8 bytes
	
	// Store result as uint32
	MOVW R4, (R1)
	
	ADD $8, R0, R0           // Advance values pointer
	ADD $4, R1, R1           // Advance output pointer
	SUBS $1, R2, R2
	BNE loop_crc

done_crc:
	RET

// func xxhash64NEON(values *int64, output *uint64, count int)
// Computes XXHash64 for int64 values using NEON
TEXT ·xxhash64NEON(SB), NOSPLIT, $0-24
	MOVD values+0(FP), R0    // R0 = &values[0]
	MOVD output+8(FP), R1    // R1 = &output[0]
	MOVD count+16(FP), R2    // R2 = count

	// XXHash64 constants (loaded in parts due to ARM64 limitations)
	// prime64_1 = 0x9E3779B185EBCA87
	MOVD $0x85EBCA87, R3
	MOVD $0x9E3779B1, R4
	LSL $32, R4, R4
	ORR R4, R3, R3           // R3 = prime64_1
	
	// prime64_2 = 0xC2B2AE3D27D4EB4F
	MOVD $0x27D4EB4F, R4
	MOVD $0xC2B2AE3D, R5
	LSL $32, R5, R5
	ORR R5, R4, R4           // R4 = prime64_2
	
	// prime64_3 = 0x165667B19E3779F9
	MOVD $0x9E3779F9, R5
	MOVD $0x165667B1, R6
	LSL $32, R6, R6
	ORR R6, R5, R5           // R5 = prime64_3
	
	// prime64_4 = 0x85EBCA77C2B2AE63
	MOVD $0xC2B2AE63, R6
	MOVD $0x85EBCA77, R7
	LSL $32, R7, R7
	ORR R7, R6, R6           // R6 = prime64_4
	
	// prime64_5 = 0x27D4EB2F165667C5
	MOVD $0x165667C5, R7
	MOVD $0x27D4EB2F, R8
	LSL $32, R8, R8
	ORR R8, R7, R7           // R7 = prime64_5

	// Process 2 elements at a time
	LSR $1, R2, R2           // R2 = count / 2
	CBZ R2, remainder_xxh

loop_xxh:
	// Load 2 int64 values
	VLD1 (R0), [V0.D2]       // V0 = [v0, v1]
	
	// Process each value individually
	VMOV V0.D[0], R8         // R8 = v0
	VMOV V0.D[1], R9         // R9 = v1
	
	// Hash v0
	MOVD R7, R10             // R10 = h64 = prime64_5
	ADD $8, R10, R10         // h64 += 8
	
	MUL R4, R8, R11          // R11 = k1 = v0 * prime64_2
	ROR $33, R11, R11        // k1 = rotl64(k1, 31) = ror(k1, 33)
	MUL R3, R11, R11         // k1 *= prime64_1
	EOR R11, R10, R10        // h64 ^= k1
	
	ROR $37, R10, R10        // h64 = rotl64(h64, 27) = ror(h64, 37)
	MUL R3, R10, R10         // h64 *= prime64_1
	ADD R6, R10, R10         // h64 += prime64_4
	
	// Finalization
	LSR $33, R10, R11
	EOR R11, R10, R10        // h64 ^= h64 >> 33
	MUL R4, R10, R10         // h64 *= prime64_2
	LSR $29, R10, R11
	EOR R11, R10, R10        // h64 ^= h64 >> 29
	MUL R5, R10, R10         // h64 *= prime64_3
	LSR $32, R10, R11
	EOR R11, R10, R10        // h64 ^= h64 >> 32
	
	// Store result 0
	MOVD R10, (R1)
	
	// Hash v1
	MOVD R7, R10             // R10 = h64 = prime64_5
	ADD $8, R10, R10         // h64 += 8
	
	MUL R4, R9, R11          // R11 = k1 = v1 * prime64_2
	ROR $33, R11, R11        // k1 = rotl64(k1, 31) = ror(k1, 33)
	MUL R3, R11, R11         // k1 *= prime64_1
	EOR R11, R10, R10        // h64 ^= k1
	
	ROR $37, R10, R10        // h64 = rotl64(h64, 27) = ror(h64, 37)
	MUL R3, R10, R10         // h64 *= prime64_1
	ADD R6, R10, R10         // h64 += prime64_4
	
	// Finalization
	LSR $33, R10, R11
	EOR R11, R10, R10        // h64 ^= h64 >> 33
	MUL R4, R10, R10         // h64 *= prime64_2
	LSR $29, R10, R11
	EOR R11, R10, R10        // h64 ^= h64 >> 29
	MUL R5, R10, R10         // h64 *= prime64_3
	LSR $32, R10, R11
	EOR R11, R10, R10        // h64 ^= h64 >> 32
	
	// Store result 1
	MOVD R10, 8(R1)
	
	ADD $16, R0, R0          // Advance values pointer
	ADD $16, R1, R1          // Advance output pointer
	SUBS $1, R2, R2
	BNE loop_xxh

remainder_xxh:
	RET
