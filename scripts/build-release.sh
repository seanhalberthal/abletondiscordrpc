#!/bin/bash

# Build script for creating release binaries
set -e

VERSION=${1:-"latest"}
OUTPUT_DIR="releases"

echo "Building Ableton Discord RPC v$VERSION"

# Create releases directory
mkdir -p $OUTPUT_DIR

# Build for macOS Intel
echo "Building for macOS Intel..."
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o $OUTPUT_DIR/ableton-discord-rpc-intel ./

# Build for macOS Apple Silicon
echo "Building for macOS Apple Silicon..."
GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o $OUTPUT_DIR/ableton-discord-rpc-arm64 ./

# Create universal binary
echo "Creating universal binary..."
lipo -create -output $OUTPUT_DIR/ableton-discord-rpc $OUTPUT_DIR/ableton-discord-rpc-intel $OUTPUT_DIR/ableton-discord-rpc-arm64

# Create checksums
echo "Creating checksums..."
cd $OUTPUT_DIR
shasum -a 256 ableton-discord-rpc* > checksums.txt
cd ..

echo "✅ Release binaries created in $OUTPUT_DIR/"
echo "📁 Files:"
ls -la $OUTPUT_DIR/

echo ""
echo "🚀 Upload ableton-discord-rpc to GitHub releases as:"
echo "   - ableton-discord-rpc (universal binary)"
echo "   - checksums.txt"