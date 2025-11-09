# 1. Verify Docker supports x86-64 emulation
docker run --rm --platform linux/amd64 alpine uname -m
# Output should be: x86_64

# 2. Create a Dockerfile for AVX2 testing
cat > Dockerfile.avx2 << 'EOF'
FROM golang:1.22-bullseye

# Install build tools
RUN apt-get update && apt-get install -y \
    build-essential \
    gcc \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy your code
COPY . .

# Build Go with AVX2
ENV GOARCH=amd64
ENV GOOS=linux

# Verify CPU features (QEMU will emulate AVX2)
RUN cat /proc/cpuinfo | grep -i avx2 || echo "AVX2 not detected (expected in QEMU)"

# Build your SIMD library
RUN go build -o syndrdb .

CMD ["/bin/bash"]
EOF

# 3. Build the Docker image
docker build --platform linux/amd64 -t syndrdb-avx2 -f Dockerfile.avx2 .

# 4. Run tests in the container
docker run --rm --platform linux/amd64 syndrdb-avx2 go test -v ./src/internal/simd/

# 5. Run benchmarks
docker run --rm --platform linux/amd64 syndrdb-avx2 go test -bench=. ./src/internal/simd/