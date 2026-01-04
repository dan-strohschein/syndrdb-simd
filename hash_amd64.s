// +build amd64

#include "textflag.h"

// func hashInt64AVX2(values *int64, output *uint64, count int)
// Computes FNV-1a hash for int64 values using AVX2
TEXT ·hashInt64AVX2(SB), NOSPLIT, $0-24
	MOVQ values+0(FP), SI    // SI = &values[0]
	MOVQ output+8(FP), DI    // DI = &output[0]
	MOVQ count+16(FP), CX    // CX = count

	// FNV-1a constants
	// offset64 = 14695981039346656037 = 0xCBF29CE484222325
	// prime64  = 1099511628211 = 0x100000001B3
	MOVQ $0xCBF29CE484222325, AX   // FNV offset
	MOVQ $0x100000001B3, BX        // FNV prime
	MOVQ $0xFF, DX                 // Byte mask
	
	// Broadcast constants to YMM registers
	VPBROADCASTQ AX, Y0   // Y0 = [offset, offset, offset, offset]
	VPBROADCASTQ BX, Y1   // Y1 = [prime, prime, prime, prime]
	VPBROADCASTQ DX, Y5   // Y5 = [0xFF, 0xFF, 0xFF, 0xFF] byte mask

	// Process 4 elements at a time
	SHRQ $2, CX          // CX = count / 4
	JZ remainder

loop:
	// Load 4 int64 values
	VMOVDQU (SI), Y2     // Y2 = [v0, v1, v2, v3]
	
	// Initialize hash accumulators with FNV offset
	VMOVDQA Y0, Y3       // Y3 = current hash values
	
	// Process each byte of the int64 (8 bytes)
	// Byte 0
	VPAND Y5, Y2, Y4
	VPXOR Y4, Y3, Y3
	VPMULLQ Y1, Y3, Y3
	VPSRLQ $8, Y2, Y2
	
	// Byte 1
	VPAND Y5, Y2, Y4
	VPXOR Y4, Y3, Y3
	VPMULLQ Y1, Y3, Y3
	VPSRLQ $8, Y2, Y2
	
	// Byte 2
	VPAND Y5, Y2, Y4
	VPXOR Y4, Y3, Y3
	VPMULLQ Y1, Y3, Y3
	VPSRLQ $8, Y2, Y2
	
	// Byte 3
	VPAND Y5, Y2, Y4
	VPXOR Y4, Y3, Y3
	VPMULLQ Y1, Y3, Y3
	VPSRLQ $8, Y2, Y2
	
	// Byte 4
	VPAND Y5, Y2, Y4
	VPXOR Y4, Y3, Y3
	VPMULLQ Y1, Y3, Y3
	VPSRLQ $8, Y2, Y2
	
	// Byte 5
	VPAND Y5, Y2, Y4
	VPXOR Y4, Y3, Y3
	VPMULLQ Y1, Y3, Y3
	VPSRLQ $8, Y2, Y2
	
	// Byte 6
	VPAND Y5, Y2, Y4
	VPXOR Y4, Y3, Y3
	VPMULLQ Y1, Y3, Y3
	VPSRLQ $8, Y2, Y2
	
	// Byte 7
	VPAND Y5, Y2, Y4
	VPXOR Y4, Y3, Y3
	VPMULLQ Y1, Y3, Y3
	
	// Store results
	VMOVDQU Y3, (DI)
	
	ADDQ $32, SI         // Advance values pointer (4 * 8 bytes)
	ADDQ $32, DI         // Advance output pointer (4 * 8 bytes)
	DECQ CX
	JNZ loop

remainder:
	VZEROUPPER
	RET

// func crc32Int64AVX2(values *int64, output *uint32, count int)
// Computes CRC32C checksums for int64 values using hardware CRC32C
TEXT ·crc32Int64AVX2(SB), NOSPLIT, $0-24
	MOVQ values+0(FP), SI    // SI = &values[0]
	MOVQ output+8(FP), DI    // DI = &output[0]
	MOVQ count+16(FP), CX    // CX = count

loop_crc:
	TESTQ CX, CX
	JZ done_crc
	
	// Load int64 value
	MOVQ (SI), AX
	
	// Compute CRC32C using hardware instruction
	XORQ DX, DX           // Initialize CRC to 0
	CRC32Q AX, DX         // CRC32C of 8 bytes
	
	// Store result as uint32
	MOVL DX, (DI)
	
	ADDQ $8, SI           // Advance values pointer
	ADDQ $4, DI           // Advance output pointer
	DECQ CX
	JNZ loop_crc

done_crc:
	RET

// func xxhash64AVX2(values *int64, output *uint64, count int)
// Computes XXHash64 for int64 values using AVX2
TEXT ·xxhash64AVX2(SB), NOSPLIT, $0-24
	MOVQ values+0(FP), SI    // SI = &values[0]
	MOVQ output+8(FP), DI    // DI = &output[0]
	MOVQ count+16(FP), CX    // CX = count

	// XXHash64 constants
	// prime64_1 = 11400714785074694791 = 0x9E3779B185EBCA87
	// prime64_2 = 14029467366897019727 = 0xC2B2AE3D27D4EB4F
	// prime64_3 = 1609587929392839161  = 0x165667B19E3779F9
	// prime64_4 = 9650029242287828579  = 0x85EBCA77C2B2AE63
	// prime64_5 = 2870177450012600261  = 0x27D4EB2F165667C5
	
	MOVQ $0x9E3779B185EBCA87, R8   // prime64_1
	MOVQ $0xC2B2AE3D27D4EB4F, R9   // prime64_2
	MOVQ $0x165667B19E3779F9, R10  // prime64_3
	MOVQ $0x85EBCA77C2B2AE63, R11  // prime64_4
	MOVQ $0x27D4EB2F165667C5, R12  // prime64_5
	
	// Broadcast constants to YMM
	VPBROADCASTQ R9, Y0   // Y0 = prime64_2
	VPBROADCASTQ R8, Y1   // Y1 = prime64_1
	VPBROADCASTQ R11, Y2  // Y2 = prime64_4
	
	// Process 4 elements at a time
	SHRQ $2, CX          // CX = count / 4
	JZ remainder_xxh

loop_xxh:
	// Load 4 int64 values
	VMOVDQU (SI), Y3     // Y3 = [v0, v1, v2, v3]
	
	// h64 = prime64_5 + 8 for each
	VPBROADCASTQ R12, Y4
	MOVQ $8, R13
	VPBROADCASTQ R13, Y6
	VPADDQ Y6, Y4, Y4    // h64 = prime64_5 + 8
	
	// k1 = value * prime64_2
	VPMULLQ Y0, Y3, Y5   // Y5 = k1 = value * prime64_2
	
	// k1 = rotl64(k1, 31)
	VPSLLQ $31, Y5, Y6
	VPSRLQ $33, Y5, Y5
	VPOR Y6, Y5, Y5      // Y5 = rotl64(k1, 31)
	
	// k1 *= prime64_1
	VPMULLQ Y1, Y5, Y5   // Y5 = k1 * prime64_1
	
	// h64 ^= k1
	VPXOR Y5, Y4, Y4
	
	// h64 = rotl64(h64, 27) * prime64_1 + prime64_4
	VPSLLQ $27, Y4, Y6
	VPSRLQ $37, Y4, Y4
	VPOR Y6, Y4, Y4      // Y4 = rotl64(h64, 27)
	VPMULLQ Y1, Y4, Y4   // Y4 *= prime64_1
	VPADDQ Y2, Y4, Y4    // Y4 += prime64_4
	
	// Finalization mix
	// h64 ^= h64 >> 33
	VPSRLQ $33, Y4, Y5
	VPXOR Y5, Y4, Y4
	
	// h64 *= prime64_2
	VPMULLQ Y0, Y4, Y4
	
	// h64 ^= h64 >> 29
	VPSRLQ $29, Y4, Y5
	VPXOR Y5, Y4, Y4
	
	// h64 *= prime64_3
	VPBROADCASTQ R10, Y5
	VPMULLQ Y5, Y4, Y4
	
	// h64 ^= h64 >> 32
	VPSRLQ $32, Y4, Y5
	VPXOR Y5, Y4, Y4
	
	// Store results
	VMOVDQU Y4, (DI)
	
	ADDQ $32, SI         // Advance values pointer
	ADDQ $32, DI         // Advance output pointer
	DECQ CX
	JNZ loop_xxh

remainder_xxh:
	VZEROUPPER
	RET
