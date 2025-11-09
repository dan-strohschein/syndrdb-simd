// +build arm64

#include "textflag.h"

// func strEqNEON(a, b *byte, length int) int
// Compares two strings for equality using NEON
TEXT ·strEqNEON(SB), NOSPLIT, $0-32
	MOVD a+0(FP), R0         // R0 = &a[0]
	MOVD b+8(FP), R1         // R1 = &b[0]
	MOVD length+16(FP), R2   // R2 = length
	MOVD $0, R3              // R3 = result (0 = not equal)

	// Handle length < 16
	CMP $16, R2
	BLT scalar_cmp

loop_16:
	CMP $16, R2
	BLT remainder_cmp
	
	// Load 16 bytes from each string
	VLD1 (R0), [V0.B16]
	VLD1 (R1), [V1.B16]
	
	// Compare
	VCMEQ V0.B16, V1.B16, V2.B16
	
	// Check if all bytes matched
	VMOV V2.D[0], R4
	VMOV V2.D[1], R5
	AND R5, R4, R4
	CMP $-1, R4
	BNE not_equal
	
	ADD $16, R0, R0
	ADD $16, R1, R1
	SUB $16, R2, R2
	B loop_16

remainder_cmp:
	// Handle remaining bytes (< 16)
	CBZ R2, equal

scalar_cmp:
	CBZ R2, equal
	
	MOVB (R0), R4
	MOVB (R1), R5
	CMP R5, R4
	BNE not_equal
	
	ADD $1, R0, R0
	ADD $1, R1, R1
	SUB $1, R2, R2
	B scalar_cmp

equal:
	MOVD $1, R3
	MOVD R3, ret+24(FP)
	RET

not_equal:
	MOVD $0, R3
	MOVD R3, ret+24(FP)
	RET

// func strPrefixCmpNEON(str, prefix *byte, strLen, prefixLen int) int
// Checks if string starts with prefix using NEON
TEXT ·strPrefixCmpNEON(SB), NOSPLIT, $0-40
	MOVD str+0(FP), R0           // R0 = &str[0]
	MOVD prefix+8(FP), R1        // R1 = &prefix[0]
	MOVD strLen+16(FP), R2       // R2 = strLen
	MOVD prefixLen+24(FP), R3    // R3 = prefixLen
	MOVD $0, R4                  // R4 = result

	// If prefix is longer than string, can't match
	CMP R2, R3
	BGT no_match
	
	// If prefix is empty, always matches
	CBZ R3, has_prefix
	
	MOVD R3, R5                  // R5 = bytes to compare

	// Handle prefixLen < 16
	CMP $16, R5
	BLT scalar_prefix

loop_prefix_16:
	CMP $16, R5
	BLT remainder_prefix
	
	// Load 16 bytes from each
	VLD1 (R0), [V0.B16]
	VLD1 (R1), [V1.B16]
	
	// Compare
	VCMEQ V0.B16, V1.B16, V2.B16
	
	// Check if all bytes matched
	VMOV V2.D[0], R6
	VMOV V2.D[1], R7
	AND R7, R6, R6
	CMP $-1, R6
	BNE no_match
	
	ADD $16, R0, R0
	ADD $16, R1, R1
	SUB $16, R5, R5
	B loop_prefix_16

remainder_prefix:
	CBZ R5, has_prefix

scalar_prefix:
	CBZ R5, has_prefix
	
	MOVB (R0), R6
	MOVB (R1), R7
	CMP R7, R6
	BNE no_match
	
	ADD $1, R0, R0
	ADD $1, R1, R1
	SUB $1, R5, R5
	B scalar_prefix

has_prefix:
	MOVD $1, R4
	MOVD R4, ret+32(FP)
	RET

no_match:
	MOVD $0, R4
	MOVD R4, ret+32(FP)
	RET

// Note: strToLowerNEON and strToUpperNEON are disabled
// Go's ARM64 assembler doesn't support vector comparison instructions (VCMGE, VCMGT)
// needed for case conversion range checks. Using generic implementations instead.
