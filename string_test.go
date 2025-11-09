package syndrdbsimd

import (
	"bytes"
	"strings"
	"testing"
)

// TestStrCmp tests string comparison
func TestStrCmp(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want int
	}{
		{"equal strings", "hello", "hello", 0},
		{"a < b", "abc", "abd", -1},
		{"a > b", "abd", "abc", 1},
		{"different lengths a < b", "ab", "abc", -1},
		{"different lengths a > b", "abc", "ab", 1},
		{"empty strings", "", "", 0},
		{"empty vs non-empty", "", "a", -1},
		{"non-empty vs empty", "a", "", 1},
		{"long equal strings", strings.Repeat("x", 100), strings.Repeat("x", 100), 0},
		{"long different strings", strings.Repeat("a", 100), strings.Repeat("b", 100), -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StrCmp([]byte(tt.a), []byte(tt.b))
			if (got < 0 && tt.want >= 0) || (got > 0 && tt.want <= 0) || (got == 0 && tt.want != 0) {
				t.Errorf("StrCmp(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// TestStrLen tests string length
func TestStrLen(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{"empty string", "", 0},
		{"single char", "a", 1},
		{"short string", "hello", 5},
		{"long string", strings.Repeat("x", 100), 100},
		{"unicode string", "hello 世界", 12}, // UTF-8 bytes, not rune count
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StrLen([]byte(tt.s))
			if got != tt.want {
				t.Errorf("StrLen(%q) = %d, want %d", tt.s, got, tt.want)
			}
		})
	}
}

// TestStrPrefixCmp tests prefix comparison
func TestStrPrefixCmp(t *testing.T) {
	tests := []struct {
		name   string
		str    string
		prefix string
		want   bool
	}{
		{"exact match", "hello", "hello", true},
		{"has prefix", "hello world", "hello", true},
		{"no prefix", "hello", "world", false},
		{"prefix too long", "hi", "hello", false},
		{"empty prefix", "hello", "", true},
		{"empty string empty prefix", "", "", true},
		{"empty string non-empty prefix", "", "a", false},
		{"long prefix match", strings.Repeat("a", 100) + "x", strings.Repeat("a", 100), true},
		{"long prefix no match", strings.Repeat("a", 100) + "x", strings.Repeat("b", 100), false},
		{"case sensitive", "Hello", "hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StrPrefixCmp([]byte(tt.str), []byte(tt.prefix))
			if got != tt.want {
				t.Errorf("StrPrefixCmp(%q, %q) = %v, want %v", tt.str, tt.prefix, got, tt.want)
			}
		})
	}
}

// TestStrContains tests substring search
func TestStrContains(t *testing.T) {
	tests := []struct {
		name   string
		str    string
		substr string
		want   bool
	}{
		{"contains at start", "hello world", "hello", true},
		{"contains at end", "hello world", "world", true},
		{"contains in middle", "hello world", "lo wo", true},
		{"does not contain", "hello world", "xyz", false},
		{"empty substring", "hello", "", true},
		{"empty string empty substring", "", "", true},
		{"empty string non-empty substring", "", "a", false},
		{"exact match", "hello", "hello", true},
		{"substring too long", "hi", "hello", false},
		{"long substring match", strings.Repeat("x", 50) + "abc" + strings.Repeat("y", 50), "abc", true},
		{"case sensitive", "Hello", "hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StrContains([]byte(tt.str), []byte(tt.substr))
			if got != tt.want {
				t.Errorf("StrContains(%q, %q) = %v, want %v", tt.str, tt.substr, got, tt.want)
			}
		})
	}
}

// TestStrEq tests string equality
func TestStrEq(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want bool
	}{
		{"equal strings", "hello", "hello", true},
		{"different strings", "hello", "world", false},
		{"different lengths", "hello", "hi", false},
		{"empty strings", "", "", true},
		{"empty vs non-empty", "", "a", false},
		{"long equal strings", strings.Repeat("x", 100), strings.Repeat("x", 100), true},
		{"long different strings", strings.Repeat("a", 100), strings.Repeat("b", 100), false},
		{"case sensitive", "Hello", "hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StrEq([]byte(tt.a), []byte(tt.b))
			if got != tt.want {
				t.Errorf("StrEq(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// TestStrToLower tests lowercase conversion
func TestStrToLower(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"already lowercase", "hello", "hello"},
		{"all uppercase", "HELLO", "hello"},
		{"mixed case", "HeLLo WoRLd", "hello world"},
		{"with numbers", "Hello123", "hello123"},
		{"with punctuation", "Hello, World!", "hello, world!"},
		{"empty string", "", ""},
		{"only numbers", "12345", "12345"},
		{"only symbols", "!@#$%", "!@#$%"},
		{"long string", strings.Repeat("ABC", 50), strings.Repeat("abc", 50)},
		{"non-ASCII unchanged", "Café", "café"}, // Only ASCII 'C' is converted
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := []byte(tt.input)
			StrToLower(input)
			if !bytes.Equal(input, []byte(tt.want)) {
				t.Errorf("StrToLower(%q) = %q, want %q", tt.input, string(input), tt.want)
			}
		})
	}
}

// TestStrToUpper tests uppercase conversion
func TestStrToUpper(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"already uppercase", "HELLO", "HELLO"},
		{"all lowercase", "hello", "HELLO"},
		{"mixed case", "HeLLo WoRLd", "HELLO WORLD"},
		{"with numbers", "hello123", "HELLO123"},
		{"with punctuation", "hello, world!", "HELLO, WORLD!"},
		{"empty string", "", ""},
		{"only numbers", "12345", "12345"},
		{"only symbols", "!@#$%", "!@#$%"},
		{"long string", strings.Repeat("abc", 50), strings.Repeat("ABC", 50)},
		{"non-ASCII unchanged", "café", "CAFé"}, // Only ASCII 'c' is converted
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := []byte(tt.input)
			StrToUpper(input)
			if !bytes.Equal(input, []byte(tt.want)) {
				t.Errorf("StrToUpper(%q) = %q, want %q", tt.input, string(input), tt.want)
			}
		})
	}
}

// TestStrEqIgnoreCase tests case-insensitive equality
func TestStrEqIgnoreCase(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want bool
	}{
		{"equal same case", "hello", "hello", true},
		{"equal different case", "Hello", "hello", true},
		{"equal all caps", "HELLO", "hello", true},
		{"equal mixed case", "HeLLo", "hEllO", true},
		{"different strings", "hello", "world", false},
		{"different lengths", "hello", "hi", false},
		{"empty strings", "", "", true},
		{"with numbers", "Hello123", "hello123", true},
		{"with punctuation", "Hello, World!", "hello, world!", true},
		{"long equal strings", strings.Repeat("ABC", 50), strings.Repeat("abc", 50), true},
		{"long different strings", strings.Repeat("A", 100), strings.Repeat("B", 100), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StrEqIgnoreCase([]byte(tt.a), []byte(tt.b))
			if got != tt.want {
				t.Errorf("StrEqIgnoreCase(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// Benchmark tests for string operations

func BenchmarkStrEq(b *testing.B) {
	s1 := []byte(strings.Repeat("hello world ", 10))
	s2 := []byte(strings.Repeat("hello world ", 10))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StrEq(s1, s2)
	}
}

func BenchmarkStrPrefixCmp(b *testing.B) {
	str := []byte(strings.Repeat("hello world ", 10))
	prefix := []byte("hello")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StrPrefixCmp(str, prefix)
	}
}

func BenchmarkStrContains(b *testing.B) {
	str := []byte(strings.Repeat("abcdefghij", 10))
	substr := []byte("def")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StrContains(str, substr)
	}
}

func BenchmarkStrToLower(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := []byte(strings.Repeat("HELLO WORLD ", 10))
		StrToLower(s)
	}
}

func BenchmarkStrToUpper(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := []byte(strings.Repeat("hello world ", 10))
		StrToUpper(s)
	}
}

func BenchmarkStrEqIgnoreCase(b *testing.B) {
	s1 := []byte(strings.Repeat("Hello World ", 10))
	s2 := []byte(strings.Repeat("hello world ", 10))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StrEqIgnoreCase(s1, s2)
	}
}
