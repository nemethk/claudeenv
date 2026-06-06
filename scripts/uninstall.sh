#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=== claudeenv Uninstaller ==="
echo

# Remove binary
BINARY_PATH="/usr/local/bin/claudeenv"
if [ -f "$BINARY_PATH" ]; then
    sudo rm "$BINARY_PATH"
    echo -e "${GREEN}✓${NC} Removed $BINARY_PATH"
else
    echo -e "${YELLOW}✓${NC} Binary already removed ($BINARY_PATH not found)"
fi

echo

# Ask about removing profiles
echo "Remove profiles? (optional)"
echo "This will delete:"
echo "  - ~/claudeenv/ (default location)"
echo "  - CLAUDEENV_DIR_GLOBAL and CLAUDEENV_DIR_PROJECT (if set)"
read -p "Remove profiles? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    rm -rf ~/claudeenv
    echo -e "${GREEN}✓${NC} Removed ~/claudeenv/"

    # Check for custom profile dirs in env
    if [ -n "$CLAUDEENV_DIR_GLOBAL" ] && [ -d "$CLAUDEENV_DIR_GLOBAL" ]; then
        read -p "Remove CLAUDEENV_DIR_GLOBAL ($CLAUDEENV_DIR_GLOBAL)? (y/n) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            rm -rf "$CLAUDEENV_DIR_GLOBAL"
            echo -e "${GREEN}✓${NC} Removed $CLAUDEENV_DIR_GLOBAL"
        fi
    fi

    if [ -n "$CLAUDEENV_DIR_PROJECT" ] && [ -d "$CLAUDEENV_DIR_PROJECT" ]; then
        read -p "Remove CLAUDEENV_DIR_PROJECT ($CLAUDEENV_DIR_PROJECT)? (y/n) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            rm -rf "$CLAUDEENV_DIR_PROJECT"
            echo -e "${GREEN}✓${NC} Removed $CLAUDEENV_DIR_PROJECT"
        fi
    fi
else
    echo -e "${YELLOW}✓${NC} Skipped profile removal"
fi

echo

# Remove shell integration
remove_from_shell() {
    local shell_rc="$1"

    if [ ! -f "$shell_rc" ]; then
        return
    fi

    if grep -q 'eval.*claudeenv init' "$shell_rc"; then
        # Create backup
        cp "$shell_rc" "$shell_rc.backup"

        # Remove the line
        grep -v 'eval.*claudeenv init' "$shell_rc.backup" > "$shell_rc.tmp"
        mv "$shell_rc.tmp" "$shell_rc"

        echo -e "${GREEN}✓${NC} Removed claudeenv from $shell_rc"
        echo "  (backup saved as $shell_rc.backup)"
    fi
}

remove_from_shell ~/.zshrc
remove_from_shell ~/.bashrc

echo

echo -e "${GREEN}✓${NC} Uninstall complete"
echo
echo "Next steps:"
echo "  1. Reload your shell: source ~/.zshrc  # or ~/.bashrc"
echo "  2. Remove any .claudeenv files from your projects (optional)"
echo "     find ~ -name '.claudeenv' -type f"
