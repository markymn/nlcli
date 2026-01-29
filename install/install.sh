#!/bin/bash

# Get the directory where the script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# Get the project root (one level up)
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "Building nlcli in $PROJECT_ROOT..."
cd "$PROJECT_ROOT" || exit
if go build -o nlcli ./cmd/nlcli; then
    echo -e "\033[0;32mBuild successful.\033[0m"
else
    echo -e "\033[0;31mBuild failed. Please check your Go installation.\033[0m"
    exit 1
fi

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

if grep -q "export PATH=.*$PROJECT_ROOT" "$SHELL_CONFIG"; then
    echo -e "\033[0;36mInfo: '$PROJECT_ROOT' is already in your PATH in $SHELL_CONFIG\033[0m"
else
    echo "" >> "$SHELL_CONFIG"
    echo "# Added by nlcli installer" >> "$SHELL_CONFIG"
    echo "export PATH=\$PATH:\"$PROJECT_ROOT\"" >> "$SHELL_CONFIG"
    echo -e "\033[0;32mSuccess: Added '$PROJECT_ROOT' to $SHELL_CONFIG\033[0m"
    echo -e "\033[0;33mPlease restart your terminal or run 'source $SHELL_CONFIG' to apply changes.\033[0m"
fi
