#!/bin/bash
set -e

echo "Installing nlcli..."

# 1. Check Prerequisites
if ! command -v git &> /dev/null; then
    echo "Error: 'git' is not installed."
    exit 1
fi
if ! command -v go &> /dev/null; then
    echo "Error: 'go' is not installed."
    exit 1
fi

# 2. Setup Directories
INSTALL_DIR="$HOME/.nlcli"
SRC_DIR="$INSTALL_DIR/src"
BIN_DIR="$INSTALL_DIR/bin"

mkdir -p "$BIN_DIR"

# 3. Clone Repository
rm -rf "$SRC_DIR"
echo "Cloning repository..."
git clone --depth 1 https://github.com/markymn/nlcli.git "$SRC_DIR"

# 4. Build
echo "Building nlcli..."
cd "$SRC_DIR"
go build -o nlcli ./cmd/nlcli

# 5. Install Binary
mv nlcli "$BIN_DIR/nlcli"
echo "Installed to $BIN_DIR/nlcli"

# 6. Update PATH
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

# Cleanup
cd "$HOME"
rm -rf "$SRC_DIR"

echo "Done!"
