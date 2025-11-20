Based on my batch evaluator implementation, I need these additional SIMD functions:

## Float64 Comparisons
```go
func CmpGtFloat64(values []float64, threshold float64) []bool
func CmpGeFloat64(values []float64, threshold float64) []bool
func CmpLtFloat64(values []float64, threshold float64) []bool
func CmpLeFloat64(values []float64, threshold float64) []bool
func CmpEqFloat64(values []float64, threshold float64) []bool
func CmpNeFloat64(values []float64, threshold float64) []bool
```

## String Comparisons
```go
func CmpEqString(values []string, threshold string) []bool
func CmpNeString(values []string, threshold string) []bool
```

## Optional but Useful - Bulk String Comparisons
```go
// For case-insensitive equality (useful for WHERE Name ILIKE "john")
func CmpEqStringIgnoreCase(values []string, threshold string) []bool

// For prefix matching (useful for WHERE Name LIKE "John%")
func CmpHasPrefixString(values []string, prefix string) []bool

// For suffix matching (useful for WHERE Email LIKE "%@gmail.com")
func CmpHasSuffixString(values []string, suffix string) []bool

// For substring matching (useful for WHERE Description LIKE "%error%")
func CmpContainsString(values []string, substr string) []bool
```

The **critical ones** are:
1. **Float64 comparisons** (all 6 operators) - needed for numeric WHERE clauses
2. **String equality/inequality** - needed for string WHERE clauses

The optional string matching functions would enable SIMD-accelerated LIKE queries, which would be a huge win but aren't blocking for Phase 3.

All functions should follow the same pattern as the existing int64 functions:
- Take a slice of values and a single threshold
- Return a bool slice with the same length
- Handle AVX2/NEON with scalar fallback
- Use the same efficient vectorization approach