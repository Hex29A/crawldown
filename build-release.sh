#!/bin/bash
set -e

VERSION=${1:-"dev"}
OUTPUT_DIR="./dist"

echo "Building CrawlDown v${VERSION}"
echo "================================"

# Clean previous builds
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

# Build flags:
# -s: omit symbol table
# -w: omit DWARF symbol table
# -X: inject version info
LDFLAGS="-s -w -X main.version=${VERSION}"

echo ""
echo "Building binaries..."
echo ""

# Linux AMD64
echo "📦 Building Linux (amd64)..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$LDFLAGS" -o "$OUTPUT_DIR/crawldown-linux-amd64" ./cmd/crawldown

# Linux ARM64
echo "📦 Building Linux (arm64)..."
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="$LDFLAGS" -o "$OUTPUT_DIR/crawldown-linux-arm64" ./cmd/crawldown

# macOS Intel
echo "📦 Building macOS (amd64)..."
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="$LDFLAGS" -o "$OUTPUT_DIR/crawldown-darwin-amd64" ./cmd/crawldown

# macOS Apple Silicon
echo "📦 Building macOS (arm64)..."
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="$LDFLAGS" -o "$OUTPUT_DIR/crawldown-darwin-arm64" ./cmd/crawldown

# Windows AMD64
echo "📦 Building Windows (amd64)..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="$LDFLAGS" -o "$OUTPUT_DIR/crawldown-windows-amd64.exe" ./cmd/crawldown

echo ""
echo "✓ Build complete!"
echo ""
echo "Binaries in $OUTPUT_DIR:"
ls -lh "$OUTPUT_DIR"
echo ""

# Calculate sizes before compression
echo "Binary sizes (before compression):"
du -h "$OUTPUT_DIR"/* | sort -h
echo ""

# Compress with UPX if available
if command -v upx &> /dev/null; then
    echo "🗜️  Compressing binaries with UPX..."
    echo ""
    
    for binary in "$OUTPUT_DIR"/*; do
        echo "  Compressing $(basename "$binary")..."
        upx --best --lzma "$binary"
    done
    
    echo ""
    echo "✓ Compression complete!"
    echo ""
    echo "Binary sizes (after UPX):"
    du -h "$OUTPUT_DIR"/* | sort -h
    echo ""
else
    echo "⚠️  UPX not found. Install with:"
    echo "    apt install upx-ucl  (Ubuntu/Debian)"
    echo "    brew install upx     (macOS)"
    echo ""
    echo "  Binaries are already stripped with -s -w flags"
    echo ""
fi

echo "To upload to GitHub:"
echo "1. Create release at https://github.com/Hex29A/crawldown/releases/new"
echo "2. Tag: v${VERSION}"
echo "3. Upload all files from ${OUTPUT_DIR}/"
