#!/bin/bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration
REPO="softkr/claude-provider"
BINARY_NAME="claude-switch"
INSTALL_DIR="/usr/local/bin"

echo -e "${CYAN}🤖 Claude Code API Switcher Installer${NC}"
echo ""

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
    linux*)
        OS="linux"
        ;;
    darwin*)
        OS="darwin"
        ;;
    *)
        echo -e "${RED}❌ Unsupported operating system: $OS${NC}"
        exit 1
        ;;
esac

# Detect Architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64|amd64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}❌ Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

echo -e "${CYAN}📋 System Information${NC}"
echo "   OS: $OS"
echo "   Arch: $ARCH"
echo ""

# Get latest release version
echo -e "${CYAN}🔍 Fetching latest release...${NC}"
LATEST_VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo -e "${YELLOW}⚠️  Could not fetch latest version, using 'latest'${NC}"
    LATEST_VERSION="latest"
fi

echo "   Version: $LATEST_VERSION"
echo ""

# Construct download URL
BINARY_FILENAME="${BINARY_NAME}-${OS}-${ARCH}"
if [ "$OS" = "windows" ]; then
    BINARY_FILENAME="${BINARY_FILENAME}.exe"
fi

DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/${BINARY_FILENAME}"

# Create temporary directory
TMP_DIR=$(mktemp -d)
TMP_FILE="${TMP_DIR}/${BINARY_NAME}"

cleanup() {
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT

# Download binary
echo -e "${CYAN}📥 Downloading ${BINARY_FILENAME}...${NC}"
if ! curl -fsSL "$DOWNLOAD_URL" -o "$TMP_FILE"; then
    echo -e "${RED}❌ Failed to download binary${NC}"
    echo -e "${YELLOW}   URL: ${DOWNLOAD_URL}${NC}"
    echo ""
    echo -e "${YELLOW}💡 Tip: Make sure the release exists with the correct binary name.${NC}"
    echo -e "${YELLOW}   Expected: ${BINARY_FILENAME}${NC}"
    exit 1
fi

# Make executable
chmod +x "$TMP_FILE"

# Verify binary
echo -e "${CYAN}✅ Verifying binary...${NC}"
if ! "$TMP_FILE" --version > /dev/null 2>&1; then
    echo -e "${RED}❌ Binary verification failed${NC}"
    exit 1
fi

VERSION_OUTPUT=$("$TMP_FILE" --version)
echo "   $VERSION_OUTPUT"
echo ""

# Install binary
INSTALL_PATH="${INSTALL_DIR}/${BINARY_NAME}"

echo -e "${CYAN}📦 Installing to ${INSTALL_PATH}...${NC}"

# Check if we need sudo
if [ -w "$INSTALL_DIR" ]; then
    cp "$TMP_FILE" "$INSTALL_PATH"
    chmod +x "$INSTALL_PATH"
else
    echo -e "${YELLOW}⚠️  Need sudo permission to install to ${INSTALL_DIR}${NC}"
    sudo cp "$TMP_FILE" "$INSTALL_PATH"
    sudo chmod +x "$INSTALL_PATH"
fi

echo -e "${GREEN}✅ Binary installed successfully${NC}"
echo ""

# Setup shell aliases
setup_aliases() {
    local shell_rc="$1"
    local is_fish="$2"

    if [ ! -f "$shell_rc" ]; then
        return 1
    fi

    # Check if aliases already exist
    if grep -q "Claude Code API Switcher" "$shell_rc" 2>/dev/null; then
        echo -e "${YELLOW}⚠️  Aliases already exist in ${shell_rc}${NC}"
        return 0
    fi

    # Add aliases
    if [ "$is_fish" = "true" ]; then
        cat >> "$shell_rc" << EOF

# Claude Code API Switcher
alias claude-switch '${INSTALL_PATH}'
alias claude-anthropic '${INSTALL_PATH} --anthropic'
alias claude-zai '${INSTALL_PATH} --zai'
alias claude-status '${INSTALL_PATH} --status'
EOF
    else
        cat >> "$shell_rc" << EOF

# Claude Code API Switcher
alias claude-switch='${INSTALL_PATH}'
alias claude-anthropic='${INSTALL_PATH} --anthropic'
alias claude-zai='${INSTALL_PATH} --zai'
alias claude-status='${INSTALL_PATH} --status'
EOF
    fi

    echo -e "${GREEN}✅ Aliases added to ${shell_rc}${NC}"
    return 0
}

echo -e "${CYAN}🔧 Setting up shell aliases...${NC}"

SHELL_NAME=$(basename "$SHELL")
ALIASES_INSTALLED=false

case "$SHELL_NAME" in
    zsh)
        if setup_aliases "$HOME/.zshrc" "false"; then
            ALIASES_INSTALLED=true
        fi
        ;;
    bash)
        if setup_aliases "$HOME/.bashrc" "false"; then
            ALIASES_INSTALLED=true
        elif setup_aliases "$HOME/.bash_profile" "false"; then
            ALIASES_INSTALLED=true
        fi
        ;;
    fish)
        if setup_aliases "$HOME/.config/fish/config.fish" "true"; then
            ALIASES_INSTALLED=true
        fi
        ;;
    *)
        # Try common shell configs
        if [ -f "$HOME/.zshrc" ]; then
            setup_aliases "$HOME/.zshrc" "false" && ALIASES_INSTALLED=true
        elif [ -f "$HOME/.bashrc" ]; then
            setup_aliases "$HOME/.bashrc" "false" && ALIASES_INSTALLED=true
        fi
        ;;
esac

echo ""
echo -e "${GREEN}🎉 Installation complete!${NC}"
echo ""
echo -e "${CYAN}Available commands:${NC}"
echo "  claude-switch --anthropic  # Switch to Anthropic API"
echo "  claude-switch --zai        # Switch to Z.AI API"
echo "  claude-switch --status     # Check current config"
echo "  claude-anthropic           # Quick switch to Anthropic"
echo "  claude-zai                 # Quick switch to Z.AI"
echo "  claude-status              # Quick status check"
echo ""

if [ "$ALIASES_INSTALLED" = true ]; then
    echo -e "${CYAN}Reload your shell to use aliases:${NC}"
    case "$SHELL_NAME" in
        zsh)
            echo "  source ~/.zshrc"
            ;;
        bash)
            echo "  source ~/.bashrc"
            ;;
        fish)
            echo "  source ~/.config/fish/config.fish"
            ;;
        *)
            echo "  source your shell config file"
            ;;
    esac
else
    echo -e "${YELLOW}Note: Could not detect shell config. You can run:${NC}"
    echo "  ${INSTALL_PATH} --install"
    echo "to manually set up aliases."
fi

echo ""
echo -e "${CYAN}Get started:${NC}"
echo "  claude-switch --status"
echo ""
