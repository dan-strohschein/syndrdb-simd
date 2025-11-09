# Build with AVX2 support
go build -tags=avx2

# Build with NEON support (ARM)
GOARCH=arm64 go build

# Run benchmarks
go test -bench=. -benchmem ./src/internal/simd/

# Check assembly output
go tool compile -S simd_amd64.s