#!/bin/bash

DIST_DIR="dist"
rm -rf "$DIST_DIR"
mkdir -p "$DIST_DIR"

PLATFORMS=(
    "windows/amd64/.exe"
    "windows/arm64/.exe"
    "linux/amd64/"
    "linux/arm64/"
    "darwin/amd64/"
    "darwin/arm64/"
)

for P in "${PLATFORMS[@]}"; do
    IFS="/" read -r OS ARCH EXT <<< "$P"
    NAME="nlcli-$OS-$ARCH$EXT"
    echo "Building $NAME..."
    GOOS=$OS GOARCH=$ARCH go build -o "$DIST_DIR/$NAME" ./cmd/nlcli
done

echo "Build complete. Binaries are in '$DIST_DIR/'"
