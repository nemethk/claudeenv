#!/bin/bash

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== claudeenv Installer ===${NC}"
echo

# Detect OS and architecture
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
    Linux*)
        OS_LOWER="linux"
        case "$ARCH" in
            x86_64)   ARCH_LOWER="amd64" ;;
            aarch64)  ARCH_LOWER="arm64" ;;
            *)        echo "Unsupported architecture: $ARCH"; exit 1 ;;
        esac
        ;;
    Darwin*)
        OS_LOWER="darwin"
        case "$ARCH" in
            x86_64)   ARCH_LOWER="amd64" ;;
            arm64)    ARCH_LOWER="arm64" ;;
            *)        echo "Unsupported architecture: $ARCH"; exit 1 ;;
        esac
        ;;
    *)
        echo "Unsupported OS: $OS (Windows is not yet supported)"
        exit 1
        ;;
esac

echo -e "${GREEN}✓${NC} Detected: $OS $ARCH"
echo

# Get latest release tag
echo "Fetching latest release..."
LATEST_TAG=$(curl -s https://api.github.com/repos/nemethk/claudeenv/releases/latest | grep -o '"tag_name": "[^"]*' | cut -d'"' -f4)

if [ -z "$LATEST_TAG" ]; then
    echo "Error: Could not determine latest release tag"
    exit 1
fi

echo -e "${GREEN}✓${NC} Latest version: $LATEST_TAG"
echo

# Download archive
ARCHIVE_NAME="claudeenv_${OS_LOWER}_${ARCH_LOWER}.tar.gz"
DOWNLOAD_URL="https://github.com/nemethk/claudeenv/releases/download/${LATEST_TAG}/${ARCHIVE_NAME}"
TEMP_DIR=$(mktemp -d)
TEMP_ARCHIVE="$TEMP_DIR/$ARCHIVE_NAME"

echo "Downloading from: $DOWNLOAD_URL"
curl -fsSL -o "$TEMP_ARCHIVE" "$DOWNLOAD_URL"

if [ ! -s "$TEMP_ARCHIVE" ]; then
    echo "Error: Download failed or file is empty"
    rm -rf "$TEMP_DIR"
    exit 1
fi

echo -e "${GREEN}✓${NC} Downloaded successfully"
echo

# Extract binary
echo "Extracting binary..."
tar -xzf "$TEMP_ARCHIVE" -C "$TEMP_DIR"
BINARY="$TEMP_DIR/claudeenv"

if [ ! -f "$BINARY" ]; then
    echo "Error: Binary not found in archive"
    rm -rf "$TEMP_DIR"
    exit 1
fi

echo -e "${GREEN}✓${NC} Extracted successfully"
echo

# Install to /usr/local/bin
echo "Installing to /usr/local/bin/claudeenv..."
sudo install -m 0755 "$BINARY" /usr/local/bin/claudeenv
rm -rf "$TEMP_DIR"

echo -e "${GREEN}✓${NC} Installation complete"
echo

# Verify installation
VERSION=$(/usr/local/bin/claudeenv --version 2>&1 || true)
echo "Installed version: $VERSION"
echo

echo "Next steps:"
echo "  1. Add shell integration to ~/.zshrc or ~/.bashrc:"
echo "     eval \"\$(claudeenv init)\""
echo
echo "  2. Reload your shell:"
echo "     source ~/.zshrc  # or ~/.bashrc"
echo
echo "  3. Configure profile locations (optional):"
echo "     export CLAUDEENV_DIR_GLOBAL=~/my-profiles/claude-global"
echo "     export CLAUDEENV_DIR_PROJECT=~/my-profiles/claude-project"
echo
echo "  4. View profiles:"
echo "     claudeenv project list"
echo
echo "For more info, see: https://github.com/nemethk/claudeenv/blob/main/INSTALL.md"
