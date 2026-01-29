#!/bin/bash
set -e

echo "Installing nlcli..."

# 1. Detect OS and Architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Error: Unsupported architecture $ARCH"; exit 1 ;;
esac

BINARY_NAME="nlcli-$OS-$ARCH"
DOWNLOAD_URL="https://github.com/markymn/nlcli/releases/latest/download/$BINARY_NAME"

# 2. Setup Directories
INSTALL_DIR="$HOME/.nlcli"
BIN_DIR="$INSTALL_DIR/bin"

mkdir -p "$BIN_DIR"

# 3. Download Binary
echo "Downloading $BINARY_NAME..."
if ! curl -L -o "$BIN_DIR/nlcli" "$DOWNLOAD_URL"; then
    echo "Error: Failed to download binary from GitHub. Ensure a release exists."
    exit 1
fi

chmod +x "$BIN_DIR/nlcli"
echo "Installed to $BIN_DIR/nlcli"

# 4. Update PATH
SHELL_CONFIG=""
if [ -f "$HOME/.zshrc" ]; then
    SHELL_CONFIG="$HOME/.zshrc"
elif [ -f "$HOME/.bashrc" ]; then
    SHELL_CONFIG="$HOME/.bashrc"
elif [ -f "$HOME/.bash_profile" ]; then
    SHELL_CONFIG="$HOME/.bash_profile"
else
    SHELL_CONFIG="$HOME/.bashrc"
    touch "$SHELL_CONFIG"
fi

if grep -q "export PATH=.*$BIN_DIR" "$SHELL_CONFIG"; then
    echo "Info: '$BIN_DIR' is already in your PATH."
else
    echo "" >> "$SHELL_CONFIG"
    echo "# Added by nlcli installer" >> "$SHELL_CONFIG"
    echo "export PATH=\$PATH:\"$BIN_DIR\"" >> "$SHELL_CONFIG"
    echo "Success: Added '$BIN_DIR' to $SHELL_CONFIG"
    echo "Please restart your terminal or run 'source $SHELL_CONFIG' to use 'nlcli'."
fi

echo "Done!"
