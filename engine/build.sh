#!/usr/bin/env bash
# build.sh - Cross-compile το Go engine για Windows, macOS, Linux
# Τρέξε από τον φάκελο engine/: ./build.sh

set -e

echo "📦 Κατεβάζω dependencies..."
go mod tidy

echo ""
echo "🔨 Compiling για όλες τις πλατφόρμες..."

BIN_DIR="../bin"
mkdir -p "$BIN_DIR"

# Linux (amd64)
echo "  → Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "$BIN_DIR/engine-linux" .

# macOS (arm64 - Apple Silicon)
echo "  → macOS (arm64 - Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o "$BIN_DIR/engine-mac-arm64" .

# macOS (amd64 - Intel)
echo "  → macOS (amd64 - Intel)..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o "$BIN_DIR/engine-mac" .

# Windows (amd64)
echo "  → Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o "$BIN_DIR/engine-windows.exe" .

echo ""
echo "✅ Build ολοκληρώθηκε! Τα binaries είναι στο φάκελο: $BIN_DIR"
ls -lh "$BIN_DIR"
