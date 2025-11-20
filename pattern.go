package syndrdbsimd

import (
	"fmt"
	"strings"
)

// CompiledPattern holds a pre-analyzed SQL LIKE pattern for efficient matching.
// SyndrDB pre-parses and classifies patterns, so this struct assumes validated input.
type CompiledPattern struct {
	// Type indicates the pattern category for optimized routing
	Type PatternType

	// Segments holds the literal parts of the pattern (between wildcards)
	// For PatternExact/Prefix/Suffix/Contains: single segment
	// For PatternWildcard: multiple segments separated by % or _
	Segments [][]byte

	// WildcardPositions tracks the positions of _ (single char wildcard)
	// Only used for PatternWildcard type
	WildcardPositions []int

	// HasWildcard indicates if the pattern contains % or _ wildcards
	HasWildcard bool

	// OriginalPattern is the original pattern string (for debugging/logging)
	OriginalPattern string
}

// CompilePattern analyzes and compiles a SQL LIKE pattern for efficient matching.
//
// SyndrDB pre-validates and classifies patterns, so this function assumes:
//   - Patterns are non-empty
//   - Patterns contain at least one literal character (not only wildcards)
//   - PatternType correctly describes the pattern structure
//
// Returns an error if the pattern is malformed (should never happen with SyndrDB).
//
// TODO: Make error verbosity configurable in future.
func CompilePattern(patternType PatternType, pattern string) (*CompiledPattern, error) {
	if pattern == "" {
		return nil, fmt.Errorf("pattern is empty")
	}

	if patternType < PatternExact || patternType > PatternWildcard {
		return nil, fmt.Errorf("invalid pattern type: %d", int(patternType))
	}

	compiled := &CompiledPattern{
		Type:            patternType,
		OriginalPattern: pattern,
		HasWildcard:     strings.ContainsAny(pattern, "%_"),
	}

	// Parse pattern based on type
	switch patternType {
	case PatternExact:
		// Exact match - entire pattern is a literal
		compiled.Segments = [][]byte{stringToBytes(pattern)}

	case PatternPrefix:
		// Prefix match - strip trailing %
		literal := strings.TrimSuffix(pattern, "%")
		if literal == "" {
			return nil, fmt.Errorf("prefix pattern has no literal part: %q", pattern)
		}
		compiled.Segments = [][]byte{stringToBytes(literal)}

	case PatternSuffix:
		// Suffix match - strip leading %
		literal := strings.TrimPrefix(pattern, "%")
		if literal == "" {
			return nil, fmt.Errorf("suffix pattern has no literal part: %q", pattern)
		}
		compiled.Segments = [][]byte{stringToBytes(literal)}

	case PatternContains:
		// Contains match - strip leading and trailing %
		literal := strings.TrimPrefix(strings.TrimSuffix(pattern, "%"), "%")
		if literal == "" {
			return nil, fmt.Errorf("contains pattern has no literal part: %q", pattern)
		}
		compiled.Segments = [][]byte{stringToBytes(literal)}

	case PatternWildcard:
		// Complex pattern with % and/or _ wildcards
		// Split by % to get literal segments
		segments := strings.Split(pattern, "%")
		
		// Filter out empty segments and convert to []byte
		for _, seg := range segments {
			if seg != "" {
				compiled.Segments = append(compiled.Segments, stringToBytes(seg))
			}
		}

		if len(compiled.Segments) == 0 {
			return nil, fmt.Errorf("wildcard pattern has no literal segments: %q", pattern)
		}

		// Track positions of _ wildcards within segments
		// Note: Full _ position tracking across entire pattern would be complex
		// For now, we handle _ during matching rather than pre-parsing positions
		// This works well with SyndrDB's constraint of max 3 % characters
		for _, seg := range compiled.Segments {
			if strings.Contains(string(seg), "_") {
				// Mark that we have _ wildcards
				// Detailed position tracking can be added if needed
				break
			}
		}

	default:
		return nil, fmt.Errorf("unsupported pattern type: %v", patternType)
	}

	return compiled, nil
}

// DetectPatternType analyzes a SQL LIKE pattern and returns its type.
// This is a helper for cases where the pattern type is not known in advance.
func DetectPatternType(pattern string) PatternType {
	if pattern == "" {
		return PatternExact
	}

	hasPercent := strings.Contains(pattern, "%")
	hasUnderscore := strings.Contains(pattern, "_")

	// No wildcards at all
	if !hasPercent && !hasUnderscore {
		return PatternExact
	}

	// Complex wildcard pattern
	if hasUnderscore {
		return PatternWildcard
	}

	// Only % wildcards - check pattern structure
	if strings.HasPrefix(pattern, "%") && strings.HasSuffix(pattern, "%") {
		// %...% pattern
		if strings.Count(pattern, "%") == 2 {
			return PatternContains
		}
		return PatternWildcard
	}

	if strings.HasPrefix(pattern, "%") {
		// %... pattern (suffix match)
		if strings.Count(pattern, "%") == 1 {
			return PatternSuffix
		}
		return PatternWildcard
	}

	if strings.HasSuffix(pattern, "%") {
		// ...% pattern (prefix match)
		if strings.Count(pattern, "%") == 1 {
			return PatternPrefix
		}
		return PatternWildcard
	}

	// % in the middle somewhere - complex pattern
	return PatternWildcard
}

// CompilePatternAuto analyzes and compiles a SQL LIKE pattern with auto-detection.
// This is a convenience wrapper around CompilePattern that auto-detects the pattern type.
func CompilePatternAuto(pattern string) (*CompiledPattern, error) {
	patternType := DetectPatternType(pattern)
	return CompilePattern(patternType, pattern)
}
