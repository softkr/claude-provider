#!/bin/bash
# Installation script for Go-based Claude Code API Switcher

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
RED='\033[0;31m'
NC='\033[0m'

# Configuration
BINARY_NAME="claude-switch"
INSTALL_DIR="/usr/local/bin"
REPO_URL="https://github.com/softkr/claude-provider"
RELEASE_VERSION="v2.0.0"

echo -e "${CYAN}🤖 Claude Code API Switcher v2.0.0 - Installation${NC}"
echo ""

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        echo -e "${RED}❌ Go is not installed${NC}"
        echo ""
        echo "Please install Go first:"
        echo "  • macOS: brew install go"
        echo "  • Ubuntu/Debian: sudo apt install golang-go"
        echo "  • Or download from: https://golang.org/dl/"
        echo ""
        exit 1
    fi

    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    echo -e "${GREEN}✅ Go found: $GO_VERSION${NC}"
}

# Detect platform
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case $OS in
        darwin)
            PLATFORM_OS="darwin"
            ;;
        linux)
            PLATFORM_OS="linux"
            ;;
        windows)
            PLATFORM_OS="windows"
            ;;
        *)
            echo -e "${RED}❌ Unsupported OS: $OS${NC}"
            exit 1
            ;;
    esac

    case $ARCH in
        x86_64|amd64)
            PLATFORM_ARCH="amd64"
            ;;
        arm64|aarch64)
            PLATFORM_ARCH="arm64"
            ;;
        *)
            echo -e "${RED}❌ Unsupported architecture: $ARCH${NC}"
            exit 1
            ;;
    esac

    echo -e "${GREEN}✅ Platform detected: $PLATFORM_OS-$PLATFORM_ARCH${NC}"
}

# Installation methods
install_from_source() {
    echo -e "${CYAN}📦 Installing from source...${NC}"

    # Create temporary directory
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"

    # Clone repository
    echo "📥 Cloning repository..."
    git clone "$REPO_URL" .

    # Build
    echo "🔨 Building..."
    make build

    # Install
    echo "📥 Installing to $INSTALL_DIR..."
    sudo cp build/$BINARY_NAME $INSTALL_DIR/
    sudo chmod +x $INSTALL_DIR/$BINARY_NAME

    # Cleanup
    cd /
    rm -rf "$TEMP_DIR"

    echo -e "${GREEN}✅ Installed from source${NC}"
}

install_prebuilt() {
    echo -e "${CYAN}📦 Installing pre-built binary...${NC}"

    BINARY="$BINARY_NAME-$PLATFORM_OS-$PLATFORM_ARCH"
    if [ "$PLATFORM_OS" = "windows" ]; then
        BINARY="$BINARY.exe"
    fi

    DOWNLOAD_URL="$REPO_URL/releases/download/$RELEASE_VERSION/$BINARY"

    echo "📥 Downloading: $DOWNLOAD_URL"

    # Download
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"

    if command -v curl &> /dev/null; then
        curl -L -o "$BINARY_NAME" "$DOWNLOAD_URL"
    elif command -v wget &> /dev/null; then
        wget -O "$BINARY_NAME" "$DOWNLOAD_URL"
    else
        echo -e "${RED}❌ Neither curl nor wget found${NC}"
        exit 1
    fi

    # Install
    echo "📥 Installing to $INSTALL_DIR..."
    sudo cp "$BINARY_NAME" $INSTALL_DIR/
    sudo chmod +x $INSTALL_DIR/$BINARY_NAME

    # Cleanup
    cd /
    rm -rf "$TEMP_DIR"

    echo -e "${GREEN}✅ Pre-built binary installed${NC}"
}

setup_aliases() {
    echo -e "${CYAN}🔧 Setting up shell aliases...${NC}"

    # Create alias block
    ALIAS_BLOCK="# Claude Code API Switcher v2.0.0
alias claude-switch='$INSTALL_DIR/$BINARY_NAME'
alias claude-anthropic='$INSTALL_DIR/$BINARY_NAME -a'
alias claude-zai='$INSTALL_DIR/$BINARY_NAME -z'
alias claude-status='$INSTALL_DIR/$BINARY_NAME -s'"

    # Setup for multiple shell configs
    SHELL_CONFIGS=()

    # Check which shell configs exist
    [ -f "$HOME/.zshrc" ] && SHELL_CONFIGS+=("$HOME/.zshrc")
    [ -f "$HOME/.bashrc" ] && SHELL_CONFIGS+=("$HOME/.bashrc")
    [ -f "$HOME/.bash_profile" ] && SHELL_CONFIGS+=("$HOME/.bash_profile")

    # If no config found, create for current shell
    if [ ${#SHELL_CONFIGS[@]} -eq 0 ]; then
        if [[ "$SHELL" == *"zsh"* ]]; then
            SHELL_CONFIGS+=("$HOME/.zshrc")
        else
            SHELL_CONFIGS+=("$HOME/.bashrc")
        fi
    fi

    for SHELL_RC in "${SHELL_CONFIGS[@]}"; do
        # Remove old aliases if exist
        if grep -q "Claude Code API Switcher" "$SHELL_RC" 2>/dev/null; then
            # Create temp file without old aliases
            grep -v "Claude Code API Switcher\|claude-switch\|claude-anthropic\|claude-zai\|claude-status" "$SHELL_RC" > "$SHELL_RC.tmp"
            mv "$SHELL_RC.tmp" "$SHELL_RC"
            echo -e "${YELLOW}🗑️  Removed old aliases from $SHELL_RC${NC}"
        fi

        # Add new aliases
        echo "" >> "$SHELL_RC"
        echo "$ALIAS_BLOCK" >> "$SHELL_RC"
        echo -e "${GREEN}✅ Aliases added to $SHELL_RC${NC}"
    done
}

# Installation flow
main() {
    check_go
    detect_platform

    echo ""
    echo -e "${CYAN}Choose installation method:${NC}"
    echo "1) Install from source (requires Git and Make)"
    echo "2) Install pre-built binary (recommended)"
    echo "3) Cancel"
    echo ""
    read -p "Enter choice [1-3]: " choice

    case $choice in
        1)
            if ! command -v git &> /dev/null; then
                echo -e "${RED}❌ Git is required for source installation${NC}"
                exit 1
            fi
            install_from_source
            ;;
        2)
            install_prebuilt
            ;;
        3)
            echo -e "${YELLOW}Installation cancelled${NC}"
            exit 0
            ;;
        *)
            echo -e "${RED}❌ Invalid choice${NC}"
            exit 1
            ;;
    esac

    setup_aliases

    echo ""
    echo -e "${GREEN}🎉 Installation complete!${NC}"
    echo ""
    echo -e "${CYAN}Next steps:${NC}"
    echo "1. Reload your shell:"
    echo "   source $SHELL_RC"
    echo ""
    echo "2. Usage:"
    echo "   claude-switch -z  # Switch to Z.AI (backup web token)"
    echo "   claude-switch -a  # Switch to Anthropic (restore web token)"
    echo "   claude-switch -s  # Check current config"
    echo ""
    echo "   Quick aliases:"
    echo "   claude-zai        # Same as claude-switch -z"
    echo "   claude-anthropic  # Same as claude-switch -a"
    echo "   claude-status     # Same as claude-switch -s"
    echo ""
    echo "3. Run Claude Code:"
    echo "   claude"
    echo ""
    echo -e "${YELLOW}Note: Anthropic uses web login token (backed up automatically)${NC}"
    echo -e "${YELLOW}      Z.AI uses API key (prompted on first use)${NC}"
}

# Check if running with options
case "${1:-}" in
    --source)
        check_go
        detect_platform
        install_from_source
        setup_aliases
        echo -e "${GREEN}✅ Installation from source complete!${NC}"
        ;;
    --prebuilt)
        detect_platform
        install_prebuilt
        setup_aliases
        echo -e "${GREEN}✅ Pre-built installation complete!${NC}"
        ;;
    *)
        main
        ;;
esac