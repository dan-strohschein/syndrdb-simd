// +build amd64

#include "textflag.h"

// func strEqAVX2(a, b *byte, length int) int
// Compares two strings for equality using AVX2
TEXT ·strEqAVX2(SB), NOSPLIT, $0-32
	MOVQ a+0(FP), SI         // SI = &a[0]
	MOVQ b+8(FP), DI         // DI = &b[0]
	MOVQ length+16(FP), CX   // CX = length
	XORQ AX, AX              // AX = result (0 = not equal, 1 = equal)

	// Handle length < 32
	CMPQ CX, $32
	JL scalar_cmp

loop_32:
	CMPQ CX, $32
	JL remainder_cmp
	
	// Load 32 bytes from each string
	VMOVDQU (SI), Y0
	VMOVDQU (DI), Y1
	
	// Compare
	VPCMPEQB Y0, Y1, Y2
	VPMOVMSKB Y2, DX
	
	// Check if all bytes matched (mask should be 0xFFFFFFFF)
	CMPL DX, $-1
	JNE not_equal
	
	ADDQ $32, SI
	ADDQ $32, DI
	SUBQ $32, CX
	JMP loop_32

remainder_cmp:
	// Handle remaining bytes (< 32)
	TESTQ CX, CX
	JZ equal

scalar_cmp:
	TESTQ CX, CX
	JZ equal
	
	MOVB (SI), AL
	CMPB AL, (DI)
	JNE not_equal
	
	INCQ SI
	INCQ DI
	DECQ CX
	JMP scalar_cmp

equal:
	MOVQ $1, AX
	MOVQ AX, ret+24(FP)
	VZEROUPPER
	RET

not_equal:
	XORQ AX, AX
	MOVQ AX, ret+24(FP)
	VZEROUPPER
	RET

// func strPrefixCmpAVX2(str, prefix *byte, strLen, prefixLen int) int
// Checks if string starts with prefix using AVX2
TEXT ·strPrefixCmpAVX2(SB), NOSPLIT, $0-40
	MOVQ str+0(FP), SI           // SI = &str[0]
	MOVQ prefix+8(FP), DI        // DI = &prefix[0]
	MOVQ strLen+16(FP), R8       // R8 = strLen
	MOVQ prefixLen+24(FP), R9    // R9 = prefixLen
	XORQ AX, AX                  // AX = result

	// If prefix is longer than string, can't match
	CMPQ R9, R8
	JG no_match
	
	// If prefix is empty, always matches
	TESTQ R9, R9
	JZ has_prefix
	
	MOVQ R9, CX                  // CX = bytes to compare

	// Handle prefixLen < 32
	CMPQ CX, $32
	JL scalar_prefix

loop_prefix_32:
	CMPQ CX, $32
	JL remainder_prefix
	
	// Load 32 bytes from each
	VMOVDQU (SI), Y0
	VMOVDQU (DI), Y1
	
	// Compare
	VPCMPEQB Y0, Y1, Y2
	VPMOVMSKB Y2, DX
	
	// Check if all bytes matched
	CMPL DX, $-1
	JNE no_match
	
	ADDQ $32, SI
	ADDQ $32, DI
	SUBQ $32, CX
	JMP loop_prefix_32

remainder_prefix:
	TESTQ CX, CX
	JZ has_prefix

scalar_prefix:
	TESTQ CX, CX
	JZ has_prefix
	
	MOVB (SI), AL
	CMPB AL, (DI)
	JNE no_match
	
	INCQ SI
	INCQ DI
	DECQ CX
	JMP scalar_prefix

has_prefix:
	MOVQ $1, AX
	MOVQ AX, ret+32(FP)
	VZEROUPPER
	RET

no_match:
	XORQ AX, AX
	MOVQ AX, ret+32(FP)
	VZEROUPPER
	RET

// func strToLowerAVX2(s *byte, length int)
// Converts ASCII string to lowercase using AVX2
TEXT ·strToLowerAVX2(SB), NOSPLIT, $0-16
	MOVQ s+0(FP), SI         // SI = &s[0]
	MOVQ length+8(FP), CX    // CX = length

	// Prepare constants
	VPCMPEQB Y5, Y5, Y5      // Y5 = all 1s
	VPSRLW $8, Y5, Y5
	VPSLLDQ $1, Y5, Y5       // Y5 = 0xFF00FF00... (mask for odd bytes)
	
	// ASCII 'A' = 65, 'Z' = 90
	// To convert to lowercase: add 32 if in range [A-Z]
	VPBROADCASTB $65, Y6     // Y6 = 'A' repeated
	VPBROADCASTB $90, Y7     // Y7 = 'Z' repeated
	VPBROADCASTB $32, Y8     // Y8 = 32 repeated

loop_lower_32:
	CMPQ CX, $32
	JL remainder_lower
	
	// Load 32 bytes
	VMOVDQU (SI), Y0
	
	// Check if >= 'A' and <= 'Z'
	VPCMPGTB Y6, Y0, Y1      // Y1 = (char > 'A' - 1) ? 0xFF : 0
	VPCMPGTB Y0, Y7, Y2      // Y2 = (char < 'Z' + 1) ? 0xFF : 0
	VPANDN Y1, Y2, Y3        // Y3 = in range [A-Z]
	VPAND Y8, Y3, Y4         // Y4 = 32 where uppercase, 0 elsewhere
	VPADDB Y4, Y0, Y0        // Add 32 to uppercase chars
	
	// Store result
	VMOVDQU Y0, (SI)
	
	ADDQ $32, SI
	SUBQ $32, CX
	JMP loop_lower_32

remainder_lower:
	TESTQ CX, CX
	JZ done_lower
	
	// Handle remainder byte-by-byte
scalar_lower:
	MOVB (SI), AL
	CMPB AL, $65
	JL skip_lower
	CMPB AL, $90
	JG skip_lower
	ADDB $32, AL
	MOVB AL, (SI)

skip_lower:
	INCQ SI
	DECQ CX
	JNZ scalar_lower

done_lower:
	VZEROUPPER
	RET

// func strToUpperAVX2(s *byte, length int)
// Converts ASCII string to uppercase using AVX2
TEXT ·strToUpperAVX2(SB), NOSPLIT, $0-16
	MOVQ s+0(FP), SI         // SI = &s[0]
	MOVQ length+8(FP), CX    // CX = length

	// ASCII 'a' = 97, 'z' = 122
	// To convert to uppercase: subtract 32 if in range [a-z]
	VPBROADCASTB $97, Y6     // Y6 = 'a' repeated
	VPBROADCASTB $122, Y7    // Y7 = 'z' repeated
	VPBROADCASTB $32, Y8     // Y8 = 32 repeated

loop_upper_32:
	CMPQ CX, $32
	JL remainder_upper
	
	// Load 32 bytes
	VMOVDQU (SI), Y0
	
	// Check if >= 'a' and <= 'z'
	VPCMPGTB Y6, Y0, Y1      // Y1 = (char > 'a' - 1) ? 0xFF : 0
	VPCMPGTB Y0, Y7, Y2      // Y2 = (char < 'z' + 1) ? 0xFF : 0
	VPANDN Y1, Y2, Y3        // Y3 = in range [a-z]
	VPAND Y8, Y3, Y4         // Y4 = 32 where lowercase, 0 elsewhere
	VPSUBB Y4, Y0, Y0        // Subtract 32 from lowercase chars
	
	// Store result
	VMOVDQU Y0, (SI)
	
	ADDQ $32, SI
	SUBQ $32, CX
	JMP loop_upper_32

remainder_upper:
	TESTQ CX, CX
	JZ done_upper
	
	// Handle remainder byte-by-byte
scalar_upper:
	MOVB (SI), AL
	CMPB AL, $97
	JL skip_upper
	CMPB AL, $122
	JG skip_upper
	SUBB $32, AL
	MOVB AL, (SI)

skip_upper:
	INCQ SI
	DECQ CX
	JNZ scalar_upper

done_upper:
	VZEROUPPER
	RET
