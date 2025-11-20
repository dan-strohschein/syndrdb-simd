package syndrdbsimd

import (
	"strings"
	"testing"
)

// TestCmpEqString tests string equality comparisons
func TestCmpEqString(t *testing.T) {
	tests := []struct {
		name      string
		values    []string
		threshold string
		expected  []bool
	}{
		{
			name:      "empty slice",
			values:    []string{},
			threshold: "test",
			expected:  []bool{},
		},
		{
			name:      "single match",
			values:    []string{"test"},
			threshold: "test",
			expected:  []bool{true},
		},
		{
			name:      "single no match",
			values:    []string{"other"},
			threshold: "test",
			expected:  []bool{false},
		},
		{
			name:      "mixed results",
			values:    []string{"hello", "test", "world", "test", "foo"},
			threshold: "test",
			expected:  []bool{false, true, false, true, false},
		},
		{
			name:      "empty string matches",
			values:    []string{"", "test", ""},
			threshold: "",
			expected:  []bool{true, false, true},
		},
		{
			name:      "case sensitive",
			values:    []string{"Test", "test", "TEST"},
			threshold: "test",
			expected:  []bool{false, true, false},
		},
		{
			name:      "SIMD threshold - 15 strings",
			values:    repeatString("test", 15),
			threshold: "test",
			expected:  repeatBool(true, 15),
		},
		{
			name:      "SIMD threshold - 16 strings",
			values:    repeatString("test", 16),
			threshold: "test",
			expected:  repeatBool(true, 16),
		},
		{
			name:      "large strings - 100 bytes each",
			values:    []string{strings.Repeat("a", 100), strings.Repeat("b", 100), strings.Repeat("a", 100)},
			threshold: strings.Repeat("a", 100),
			expected:  []bool{true, false, true},
		},
		{
			name:      "Unicode strings",
			values:    []string{"hello", "世界", "test", "世界"},
			threshold: "世界",
			expected:  []bool{false, true, false, true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CmpEqString(tt.values, tt.threshold)
			if len(result) != len(tt.expected) {
				t.Fatalf("length mismatch: got %d, want %d", len(result), len(tt.expected))
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("index %d: got %v, want %v (value=%q, threshold=%q)",
						i, result[i], tt.expected[i], tt.values[i], tt.threshold)
				}
			}
		})
	}
}

// TestCmpEqStringMask tests bitmask variant of string equality
func TestCmpEqStringMask(t *testing.T) {
	tests := []struct {
		name      string
		values    []string
		threshold string
	}{
		{
			name:      "64 strings - full word",
			values:    repeatString("test", 64),
			threshold: "test",
		},
		{
			name:      "65 strings - multiple words",
			values:    append(repeatString("test", 64), "other"),
			threshold: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			boolResult := CmpEqString(tt.values, tt.threshold)
			maskResult := CmpEqStringMask(tt.values, tt.threshold)

			convertedBools := BitmaskToBools(maskResult, len(tt.values))
			if len(convertedBools) != len(boolResult) {
				t.Fatalf("length mismatch: got %d, want %d", len(convertedBools), len(boolResult))
			}
			for i := range convertedBools {
				if convertedBools[i] != boolResult[i] {
					t.Errorf("index %d: mask gave %v, bool gave %v", i, convertedBools[i], boolResult[i])
				}
			}
		})
	}
}

// TestCmpNeString tests string inequality comparisons
func TestCmpNeString(t *testing.T) {
	tests := []struct {
		name      string
		values    []string
		threshold string
		expected  []bool
	}{
		{
			name:      "mixed results",
			values:    []string{"hello", "test", "world", "test"},
			threshold: "test",
			expected:  []bool{true, false, true, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CmpNeString(tt.values, tt.threshold)
			if len(result) != len(tt.expected) {
				t.Fatalf("length mismatch: got %d, want %d", len(result), len(tt.expected))
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("index %d: got %v, want %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

// TestCmpHasPrefixString tests prefix matching
func TestCmpHasPrefixString(t *testing.T) {
	tests := []struct {
		name     string
		values   []string
		prefix   string
		expected []bool
	}{
		{
			name:     "basic prefix",
			values:   []string{"hello", "hell", "help", "world"},
			prefix:   "hel",
			expected: []bool{true, true, true, false},
		},
		{
			name:     "empty prefix - all match",
			values:   []string{"hello", "world", "test"},
			prefix:   "",
			expected: []bool{true, true, true},
		},
		{
			name:     "prefix longer than string",
			values:   []string{"hi", "hello"},
			prefix:   "hello world",
			expected: []bool{false, false},
		},
		{
			name:     "exact match",
			values:   []string{"test", "testing"},
			prefix:   "test",
			expected: []bool{true, true},
		},
		{
			name:     "case sensitive",
			values:   []string{"Hello", "hello"},
			prefix:   "hel",
			expected: []bool{false, true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CmpHasPrefixString(tt.values, tt.prefix)
			if len(result) != len(tt.expected) {
				t.Fatalf("length mismatch: got %d, want %d", len(result), len(tt.expected))
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("index %d: got %v, want %v (value=%q, prefix=%q)",
						i, result[i], tt.expected[i], tt.values[i], tt.prefix)
				}
			}
		})
	}
}

// TestCmpHasSuffixString tests suffix matching
func TestCmpHasSuffixString(t *testing.T) {
	tests := []struct {
		name     string
		values   []string
		suffix   string
		expected []bool
	}{
		{
			name:     "basic suffix",
			values:   []string{"hello", "world", "test"},
			suffix:   "llo",
			expected: []bool{true, false, false},
		},
		{
			name:     "empty suffix - all match",
			values:   []string{"hello", "world"},
			suffix:   "",
			expected: []bool{true, true},
		},
		{
			name:     "suffix longer than string",
			values:   []string{"hi", "hello"},
			suffix:   "hello world",
			expected: []bool{false, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CmpHasSuffixString(tt.values, tt.suffix)
			if len(result) != len(tt.expected) {
				t.Fatalf("length mismatch: got %d, want %d", len(result), len(tt.expected))
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("index %d: got %v, want %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

// TestCmpContainsString tests substring matching
func TestCmpContainsString(t *testing.T) {
	tests := []struct {
		name     string
		values   []string
		substr   string
		expected []bool
	}{
		{
			name:     "basic contains",
			values:   []string{"hello world", "test", "world map"},
			substr:   "world",
			expected: []bool{true, false, true},
		},
		{
			name:     "empty substring - all match",
			values:   []string{"hello", "world"},
			substr:   "",
			expected: []bool{true, true},
		},
		{
			name:     "substring not found",
			values:   []string{"hello", "world"},
			substr:   "xyz",
			expected: []bool{false, false},
		},
		{
			name:     "substring at beginning",
			values:   []string{"hello world"},
			substr:   "hello",
			expected: []bool{true},
		},
		{
			name:     "substring at end",
			values:   []string{"hello world"},
			substr:   "world",
			expected: []bool{true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CmpContainsString(tt.values, tt.substr)
			if len(result) != len(tt.expected) {
				t.Fatalf("length mismatch: got %d, want %d", len(result), len(tt.expected))
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("index %d: got %v, want %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

// TestCmpEqStringIgnoreCase tests case-insensitive equality
func TestCmpEqStringIgnoreCase(t *testing.T) {
	tests := []struct {
		name      string
		values    []string
		threshold string
		expected  []bool
	}{
		{
			name:      "case insensitive match",
			values:    []string{"HELLO", "hello", "HeLLo", "world"},
			threshold: "hello",
			expected:  []bool{true, true, true, false},
		},
		{
			name:      "mixed case ASCII",
			values:    []string{"TeSt", "TEST", "test"},
			threshold: "test",
			expected:  []bool{true, true, true},
		},
		{
			name:      "non-ASCII - exact match only",
			values:    []string{"CAFÉ", "café", "Café"},
			threshold: "café",
			expected:  []bool{false, true, true}, // Café matches because ASCII 'C' lowercases to 'c', é is exact byte match
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CmpEqStringIgnoreCase(tt.values, tt.threshold)
			if len(result) != len(tt.expected) {
				t.Fatalf("length mismatch: got %d, want %d", len(result), len(tt.expected))
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("index %d: got %v, want %v (value=%q, threshold=%q)",
						i, result[i], tt.expected[i], tt.values[i], tt.threshold)
				}
			}
		})
	}
}

// TestCmpLikeString tests SQL LIKE pattern matching
func TestCmpLikeString(t *testing.T) {
	tests := []struct {
		name     string
		values   []string
		pattern  string
		expected []bool
	}{
		{
			name:     "exact pattern - no wildcards",
			values:   []string{"hello", "hello", "world"},
			pattern:  "hello",
			expected: []bool{true, true, false},
		},
		{
			name:     "prefix pattern - hello%",
			values:   []string{"hello", "hello world", "hi there"},
			pattern:  "hello%",
			expected: []bool{true, true, false},
		},
		{
			name:     "suffix pattern - %world",
			values:   []string{"hello world", "world", "hello"},
			pattern:  "%world",
			expected: []bool{true, true, false},
		},
		{
			name:     "contains pattern - %test%",
			values:   []string{"this is a test", "testing", "hello"},
			pattern:  "%test%",
			expected: []bool{true, true, false},
		},
		{
			name:     "underscore wildcard - h_llo",
			values:   []string{"hello", "hallo", "hxllo", "hllo"},
			pattern:  "h_llo",
			expected: []bool{true, true, true, false},
		},
		{
			name:     "complex pattern - a%b%c",
			values:   []string{"abc", "aXbYc", "aXXbYYc", "ac", "ab"},
			pattern:  "a%b%c",
			expected: []bool{true, true, true, false, false},
		},
		{
			name:     "pattern with multiple underscores - __st",
			values:   []string{"test", "best", "st", "aest"},
			pattern:  "__st",
			expected: []bool{true, true, false, true},
		},
		{
			name:     "mixed wildcards - %t_st%",
			values:   []string{"test", "testing", "latest result", "st"},
			pattern:  "%t_st%",
			expected: []bool{true, true, true, false}, // "test"→"t"+"e"+"st", "testing"→"t"+"e"+"st"+"ing", "latest result"→"late"+"s"+"t "+"result"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CmpLikeString(tt.values, tt.pattern)
			if len(result) != len(tt.expected) {
				t.Fatalf("length mismatch: got %d, want %d", len(result), len(tt.expected))
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("index %d: got %v, want %v (value=%q, pattern=%q)",
						i, result[i], tt.expected[i], tt.values[i], tt.pattern)
				}
			}
		})
	}
}

// TestCmpLikeStringCompiled tests pre-compiled pattern matching
func TestCmpLikeStringCompiled(t *testing.T) {
	values := []string{"hello", "hello world", "world", "test"}

	// Test pattern reuse
	pattern, err := CompilePatternAuto("hello%")
	if err != nil {
		t.Fatalf("CompilePatternAuto failed: %v", err)
	}

	result1 := CmpLikeStringCompiled(values, pattern)
	result2 := CmpLikeStringCompiled(values, pattern)

	// Results should be identical
	for i := range result1 {
		if result1[i] != result2[i] {
			t.Errorf("inconsistent results at index %d: %v vs %v", i, result1[i], result2[i])
		}
	}

	// Verify correctness
	expected := []bool{true, true, false, false}
	for i := range result1 {
		if result1[i] != expected[i] {
			t.Errorf("index %d: got %v, want %v", i, result1[i], expected[i])
		}
	}
}

// TestDetectPatternType tests pattern type detection
func TestDetectPatternType(t *testing.T) {
	tests := []struct {
		pattern      string
		expectedType PatternType
	}{
		{"hello", PatternExact},
		{"hello%", PatternPrefix},
		{"%world", PatternSuffix},
		{"%test%", PatternContains},
		{"h_llo", PatternWildcard},
		{"a%b%c", PatternWildcard},
		{"%a%b%", PatternWildcard},
		{"", PatternExact},
	}

	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			result := DetectPatternType(tt.pattern)
			if result != tt.expectedType {
				t.Errorf("got %v, want %v", result, tt.expectedType)
			}
		})
	}
}

// TestCompilePattern tests pattern compilation
func TestCompilePattern(t *testing.T) {
	tests := []struct {
		name        string
		patternType PatternType
		pattern     string
		wantErr     bool
	}{
		{
			name:        "valid exact pattern",
			patternType: PatternExact,
			pattern:     "hello",
			wantErr:     false,
		},
		{
			name:        "valid prefix pattern",
			patternType: PatternPrefix,
			pattern:     "hello%",
			wantErr:     false,
		},
		{
			name:        "valid suffix pattern",
			patternType: PatternSuffix,
			pattern:     "%world",
			wantErr:     false,
		},
		{
			name:        "valid contains pattern",
			patternType: PatternContains,
			pattern:     "%test%",
			wantErr:     false,
		},
		{
			name:        "valid wildcard pattern",
			patternType: PatternWildcard,
			pattern:     "a%b%c",
			wantErr:     false,
		},
		{
			name:        "empty pattern",
			patternType: PatternExact,
			pattern:     "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CompilePattern(tt.patternType, tt.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompilePattern() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result == nil {
				t.Error("CompilePattern() returned nil result without error")
			}
			if !tt.wantErr && result.Type != tt.patternType {
				t.Errorf("CompilePattern() type = %v, want %v", result.Type, tt.patternType)
			}
		})
	}
}

// TestBufferPoolIntegration tests buffer pool with string operations
func TestBufferPoolIntegration(t *testing.T) {
	// Create large strings that might use pooled buffers
	largeStrings := make([]string, 100)
	for i := range largeStrings {
		largeStrings[i] = strings.Repeat("test", 100) // 400 bytes each
	}

	// Run multiple operations to exercise the pool
	for i := 0; i < 10; i++ {
		_ = CmpEqString(largeStrings, "test")
		_ = CmpHasPrefixString(largeStrings, "test")
		_ = CmpContainsString(largeStrings, "test")
	}

	// Get pool stats
	stats := GetBufferPoolStats()
	
	// Verify stats are being tracked (uint64 so always >= 0)
	// Just verify we can call the function without panic
	_ = stats.PoolHits
	_ = stats.PoolMisses
	_ = stats.HeapFallbacks
	
	// Check time window is reasonable
	if stats.WindowEndTime.Before(stats.WindowStartTime) {
		t.Error("WindowEndTime should be after WindowStartTime")
	}
}

// Helper functions

// repeatString creates a slice of n identical strings
func repeatString(s string, n int) []string {
	result := make([]string, n)
	for i := range result {
		result[i] = s
	}
	return result
}

// repeatBool creates a slice of n identical booleans
func repeatBool(b bool, n int) []bool {
	result := make([]bool, n)
	for i := range result {
		result[i] = b
	}
	return result
}
